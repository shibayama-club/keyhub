package serve

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/shibayama-club/keyhub/cmd/config"
	"github.com/shibayama-club/keyhub/internal/domain/healthcheck"
	"github.com/shibayama-club/keyhub/internal/infrastructure/sqlc"
	"github.com/shibayama-club/keyhub/internal/interface/health"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
)

func ServeConsole() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "console",
		PreRunE: config.ParseConfig[config.Config],
		RunE:    runConsole,
	}
	flags := cmd.Flags()
	flags.Int("port", 8081, "port number to listen")

	config.ConfigFlags(flags)

	return cmd
}

func runConsole(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cfg, ok := ctx.Value(cmd).(config.Config)
	if !ok {
		return errors.New("failed to get config")
	}

	e, err := SetupConsole(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to setup console api server")
	}
	slog.Info("starting console api server")
	return e.StartH2CServer(fmt.Sprintf(":%d", cfg.Port), &http2.Server{})
}

func SetupConsole(ctx context.Context, config config.Config) (*echo.Echo, error) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.Sentry.DSN,
		Environment:      config.Env,
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		AttachStacktrace: true,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to initialize Sentry")
	}

	e := echo.New()
	e.Use(
		middleware.Recover(),
		slogecho.New(slog.Default()),
	)

	healthCheckers := make([]healthcheck.HealthChecker, 0)

	pool, err := sqlc.NewPool(ctx, config.Postgres)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create postgres pool")
	}
	repo := sqlc.NewRepository(pool)
	healthCheckers = append(healthCheckers, healthcheck.NewHealthCheckFunc("repository", repo.Ping))

	healthHandler := health.NewHealthCheck(healthCheckers...)
	e.GET("/keyhub.console.v1.HealthService/Check", healthHandler.Check)

	return e, nil
}
