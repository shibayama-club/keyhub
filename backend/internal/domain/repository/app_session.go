package repository

import (
	"context"
	"time"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type AppSessionRepository interface {
	CreateAppSession(ctx context.Context, arg CreateAppSessionArg) (model.AppSession, error)
	GetAppSession(ctx context.Context, sessionID model.AppSessionID) (model.AppSession, error)
	RevokeAppSession(ctx context.Context, sessionID model.AppSessionID) error
}

type CreateAppSessionArg struct {
	SessionID model.AppSessionID
	UserID    model.UserID
	ExpiresAt time.Time
}
