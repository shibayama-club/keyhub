package app

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/shibayama-club/keyhub/cmd/config"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	"github.com/shibayama-club/keyhub/internal/infrastructure/auth/google"
	"github.com/shibayama-club/keyhub/internal/usecase/app/iface"
)

type UseCase struct {
	repo         repository.Repository
	config       config.Config
	oauthService *google.OAuthService
}

var _ iface.IUseCase = (*UseCase)(nil)

func NewUseCase(ctx context.Context, repo repository.Repository, cf config.Config, oauthService *google.OAuthService) (iface.IUseCase, error) {
	if oauthService == nil {
		return nil, errors.New("oauth service is required")
	}

	return &UseCase{
		repo:         repo,
		config:       cf,
		oauthService: oauthService,
	}, nil
}
