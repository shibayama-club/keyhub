package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type OAuthStateRepository interface {
	SaveOAuthState(ctx context.Context, arg SaveOAuthStateArg) error
	GetOAuthState(ctx context.Context, state string) (model.OAuthState, error)
	ConsumeOAuthState(ctx context.Context, state string) error
}

type SaveOAuthStateArg struct {
	State        string
	CodeVerifier string
	Nonce        string
}
