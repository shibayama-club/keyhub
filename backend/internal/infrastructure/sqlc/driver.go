package sqlc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/shibayama-club/keyhub/cmd/config"
)

func NewPool(ctx context.Context, cf config.DBConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s",
		cf.Host, cf.Port, cf.User, cf.Password, cf.Database)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse pgx config")
	}

	if cf.Debug {
		config.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   newLogger(),
			LogLevel: tracelog.LogLevelTrace,
		}
	}

	config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		if _, err := conn.Exec(ctx, "RESET ALL"); err != nil {
			slog.ErrorContext(ctx, "BeforeAcquire: failed to RESET context", slog.String("error", err.Error()))
			return false
		}
		return true
	}

	config.AfterRelease = func(conn *pgx.Conn) bool {
		if _, err := conn.Exec(context.Background(), "RESET ALL"); err != nil {
			slog.ErrorContext(ctx, "AfterRelease: failed to RESET context", slog.String("error", err.Error()))
			return false
		}

		return true
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pgx pool")
	}

	return pool, nil
}

type logger struct {
	l *slog.Logger
}

func newLogger() *logger {
	return &logger{
		l: slog.Default(),
	}
}

func (l *logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	attrs := make([]any, 0, 2*len(data))
	for key, value := range data {
		attrs = append(attrs, slog.Any(key, value))
	}

	switch level {
	case tracelog.LogLevelTrace:
		l.l.With("PGX_LOG_LEVEL", level).Log(ctx, slog.LevelDebug, msg, attrs...)
	case tracelog.LogLevelDebug:
		l.l.Log(ctx, slog.LevelDebug, msg, attrs...)
	case tracelog.LogLevelInfo:
		l.l.Log(ctx, slog.LevelInfo, msg, attrs...)
	case tracelog.LogLevelWarn:
		l.l.Log(ctx, slog.LevelWarn, msg, attrs...)
	case tracelog.LogLevelError:
		l.l.Log(ctx, slog.LevelError, msg, attrs...)
	default:
		l.l.With("INVALID_PGX_LOG_LEVEL", level).Log(ctx, slog.LevelError, msg, attrs...)
	}
}
