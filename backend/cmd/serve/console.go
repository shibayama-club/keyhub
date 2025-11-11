package serve

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/shibayama-club/keyhub/cmd/config"
	"github.com/shibayama-club/keyhub/internal/domain/healthcheck"
	consoleauth "github.com/shibayama-club/keyhub/internal/infrastructure/auth/console"
	"github.com/shibayama-club/keyhub/internal/infrastructure/sqlc"
	consolev1 "github.com/shibayama-club/keyhub/internal/interface/console/v1"
	"github.com/shibayama-club/keyhub/internal/interface/console/v1/interceptor"
	"github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/console/v1/consolev1connect"
	"github.com/shibayama-club/keyhub/internal/interface/health"
	"github.com/shibayama-club/keyhub/internal/usecase/console"
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
	flags.Int("port", 8081, "Listen port")

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

func SetupConsole(ctx context.Context, cfg config.Config) (*echo.Echo, error) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.Sentry.DSN,
		Environment:      cfg.Env,
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
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{cfg.FrontendURL.Console},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			AllowHeaders: []string{"*"},
		}),
	)

	healthCheckers := make([]healthcheck.HealthChecker, 0)

	pool, err := sqlc.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create postgres pool")
	}
	repo := sqlc.NewRepository(pool)
	healthCheckers = append(healthCheckers, healthcheck.NewHealthCheckFunc("repository", repo.Ping))

	jwtSecret := cfg.Console.JWTSecret
	if jwtSecret == "" {
		jwtSecret = console.DEFAULT_JWT_SECRET
	}
	consoleAuth, err := consoleauth.NewAuthService(jwtSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create console auth service")
	}

	consoleUseCase, err := console.NewUseCase(ctx, repo, cfg, consoleAuth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create console use case")
	}

	sentryInterceptor := interceptor.NewSentryInterceptor()
	authInterceptor := interceptor.NewAuthInterceptor(consoleUseCase)

	consoleHandler, err := consolev1.NewHandler(consoleUseCase, jwtSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create console handler")
	}

	// ConsoleAuthServiceをConnectRPCに登録
	authPath, authHandler := consolev1connect.NewConsoleAuthServiceHandler(
		consoleHandler,
		connect.WithInterceptors(sentryInterceptor, authInterceptor),
	)
	e.Any(authPath+"*", echo.WrapHandler(authHandler))

	// ConsoleServiceをConnectRPCに登録
	servicePath, serviceHandler := consolev1connect.NewConsoleServiceHandler(
		consoleHandler,
		connect.WithInterceptors(sentryInterceptor, authInterceptor),
	)
	e.Any(servicePath+"*", echo.WrapHandler(serviceHandler))

	healthHandler := health.NewHealthCheck(healthCheckers...)
	e.GET("/keyhub.console.v1.HealthService/Check", healthHandler.Check)

	return e, nil
}
