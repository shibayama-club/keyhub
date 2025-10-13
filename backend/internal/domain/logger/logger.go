package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/shibayama-club/keyhub/internal/domain"
	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type (
	contextualLoggingHandler struct {
		slog.Handler
		debug bool
	}

	RequestID string
)

func SetupLogger(debug bool) {
	slog.SetDefault(slog.New(contextualLoggingHandler{createHandler(debug), debug}))
}

func createHandler(debug bool) slog.Handler {
	if debug {
		return tint.NewHandler(os.Stdout, &tint.Options{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
}

func (h contextualLoggingHandler) Handle(ctx context.Context, record slog.Record) error {
	requestID, ok := domain.Value[RequestID](ctx)
	if ok {
		record.Add(slog.String("request_id", string(requestID)))
	}

	// app
	userID, ok := domain.Value[model.UserID](ctx)
	if ok {
		record.Add(slog.String("user_id", userID.String()))
	}
	tenantID, ok := domain.Value[model.TenantID](ctx)
	if ok {
		record.Add(slog.String("tenant_id", tenantID.String()))
	}

	// console
	organizationID, ok := domain.Value[model.OrganizationID](ctx)
	if ok {
		record.Add(slog.String("organization_id", organizationID.String()))
	}

	return h.Handler.Handle(ctx, record)
}
