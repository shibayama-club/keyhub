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
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
)

func ServeApp() *cobra.Command {
	cmd := &cobra.Command{
		Use: "app",
		PreRunE: config.ParseConfig[config.Config],
		RunE: runApp,
	}
	flags := cmd.Flags()
	flags.Int("port", 8080, "Port number to listen")

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

func SetupApp(ctx context.Context, config config.Config) (*echo.Echo, error) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: config.Sentry.DSN,
		Environment: config.Env,
		TracesSampleRate: 1.0,
		EnableTracing: true,
		AttachStacktrace: true,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to initialize Sentry")
	}

	e := echo.New()
	e.Use(
		middleware.Recover(),
		slogecho.New(slog.Default()),
	)

	return e, nil
}