package serve

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/shibayama-club/keyhub/cmd/config"
	"github.com/shibayama-club/keyhub/internal/domain/healthcheck"
	"github.com/shibayama-club/keyhub/internal/infrastructure/auth/google"
	"github.com/shibayama-club/keyhub/internal/infrastructure/sqlc"
	appv1 "github.com/shibayama-club/keyhub/internal/interface/app/v1"
	"github.com/shibayama-club/keyhub/internal/interface/app/v1/interceptor"
	"github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/app/v1/appv1connect"
	"github.com/shibayama-club/keyhub/internal/interface/health"
	"github.com/shibayama-club/keyhub/internal/usecase/app"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
)

func ServeApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app",
		PreRunE: config.ParseConfig[config.Config],
		RunE:    runApp,
	}
	flags := cmd.Flags()
	flags.Int("port", 8080, "Listen Port")

	config.ConfigFlags(flags)

	return cmd
}

func runApp(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cfg, ok := ctx.Value(cmd).(config.Config)
	if !ok {
		return errors.New("failed to get config")
	}
	e, err := SetupApp(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to setup app api server")
	}
	slog.Info("starting app api server")
	return e.StartH2CServer(fmt.Sprintf(":%d", cfg.Port), &http2.Server{})
}

func SetupApp(ctx context.Context, cfg config.Config) (*echo.Echo, error) {
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
			AllowOrigins: []string{cfg.FrontendURL.App},
			AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
			AllowHeaders: []string{
				"Accept",
				"Accept-Encoding",
				"Content-Type",
				"Connect-Protocol-Version",
				"Connect-Timeout-Ms",
				"Grpc-Timeout",
				"X-Grpc-Web",
				"X-User-Agent",
				"Cookie",
			},
			ExposeHeaders:    []string{"Content-Length", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           3600,
		}),
	)

	healthCheckers := make([]healthcheck.HealthChecker, 0)

	pool, err := sqlc.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create postgres pool")
	}
	repo := sqlc.NewRepository(pool)
	healthCheckers = append(healthCheckers, healthcheck.NewHealthCheckFunc("repository", repo.Ping))

	oauthService, err := google.NewOAuthService(google.OAuthConfig{
		ClientID:     cfg.Auth.Google.ClientID,
		ClientSecret: cfg.Auth.Google.ClientSecret,
		RedirectURI:  cfg.Auth.Google.RedirectURI,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create OAuth service")
	}

	appUseCase, err := app.NewUseCase(ctx, repo, cfg, oauthService)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create app use case")
	}

	appHandler := appv1.NewHandler(appUseCase, cfg.Env, cfg.FrontendURL.App)

	e.GET("/auth/google/login", appHandler.GoogleLogin)
	e.GET("/auth/google/callback", appHandler.GoogleCallback)

	authInterceptor := interceptor.NewAuthInterceptor(appUseCase)

	authPath, authHandler := appv1connect.NewAuthServiceHandler(
		appHandler,
		connect.WithInterceptors(authInterceptor),
	)
	e.Any(authPath+"*", echo.WrapHandler(authHandler))

	tenantPath, tenantHandler := appv1connect.NewTenantServiceHandler(
		appHandler,
		connect.WithInterceptors(authInterceptor),
	)
	e.Any(tenantPath+"*", echo.WrapHandler(tenantHandler))

	roomPath, roomHandler := appv1connect.NewRoomServiceHandler(
		appHandler,
		connect.WithInterceptors(authInterceptor),
	)
	e.Any(roomPath+"*", echo.WrapHandler(roomHandler))

	healthHandler := health.NewHealthCheck(healthCheckers...)
	e.GET("/keyhub.app.v1.HealthService/Check", healthHandler.Check)

	return e, nil
}
