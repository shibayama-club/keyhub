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
	"github.com/shibayama-club/keyhub/internal/domain"
	"github.com/shibayama-club/keyhub/internal/domain/model"
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
		if orgID, ok := domain.Value[model.OrganizationID](ctx); ok {
			if _, err := conn.Exec(ctx, "SELECT set_config('keyhub.organization_id', $1, false)", orgID.UUID().String()); err != nil {
				slog.ErrorContext(ctx, "BeforeAcquire: failed to set RLS organization_id", slog.String("error", err.Error()))
				return false
			}
		} else {
			if _, err := conn.Exec(ctx, "RESET keyhub.organization_id"); err != nil {
				slog.ErrorContext(ctx, "BeforeAcquire: failed to reset RLS organization_id", slog.String("error", err.Error()))
				return false
			}
		}

		if membershipID, ok := domain.Value[model.TenantMembershipID](ctx); ok {
			if _, err := conn.Exec(ctx, "SELECT set_config('keyhub.membership_id', $1, false)", membershipID.UUID().String()); err != nil {
				slog.ErrorContext(ctx, "BeforeAcquire: failed to set RLS membership_id", slog.String("error", err.Error()))
				return false
			}
		} else {
			if _, err := conn.Exec(ctx, "RESET keyhub.membership_id"); err != nil {
				slog.ErrorContext(ctx, "BeforeAcquire: failed to reset RLS membership_id", slog.String("error", err.Error()))
				return false
			}
		}

		return true
	}

	config.AfterRelease = func(conn *pgx.Conn) bool {
		if _, err := conn.Exec(context.Background(), "RESET keyhub.organization_id"); err != nil {
			slog.Error("AfterRelease: failed to reset organization_id", slog.String("error", err.Error()))
			return false
		}
		if _, err := conn.Exec(context.Background(), "RESET keyhub.membership_id"); err != nil {
			slog.Error("AfterRelease: failed to reset membership_id", slog.String("error", err.Error()))
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
