package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateConsoleSessionArg struct {
	SessionID      model.ConsoleSessionID
	OrganizationID model.OrganizationID
}

type ConsoleSessionRepository interface {
	CreateSession(ctx context.Context, arg CreateConsoleSessionArg) (model.ConsoleSession, error)
	GetSession(ctx context.Context, sessionID model.ConsoleSessionID) (model.ConsoleSession, error)
	DeleteSession(ctx context.Context, sessionID model.ConsoleSessionID) error
}
