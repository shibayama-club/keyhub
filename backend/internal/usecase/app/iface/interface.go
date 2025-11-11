package iface

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type IUseCase interface {
	StartGoogleLogin(ctx context.Context) (authURL string, err error)
	GoogleCallback(ctx context.Context, code, state string) (sessionID string, err error)
	GetMe(ctx context.Context, sessionID string) (*model.User, error)
	GetUserByID(ctx context.Context, userID model.UserID) (*model.User, error)
	Logout(ctx context.Context, sessionID string) error
}
