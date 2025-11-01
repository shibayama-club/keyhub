package console

import (
	"context"

	"github.com/shibayama-club/keyhub/cmd/config"
	"github.com/shibayama-club/keyhub/internal/domain/authenticator"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	"github.com/shibayama-club/keyhub/internal/usecase/console/iface"
)

// 開発用環境変数
const (
	DEFAULT_ORGANIZATION_ID  = "550e8400-e29b-41d4-a716-446655440000"
	DEFAULT_ORGANIZATION_KEY = "org_key_example_12345"
	DEFAULT_JWT_SECRET       = "your-secret-jwt-key-change-in-production"
)

type UseCase struct {
	repo        repository.Repository
	config      config.Config
	authService authenticator.ConsoleAuthenticator
}

var _ iface.IUseCase = (*UseCase)(nil)

func NewUseCase(
	ctx context.Context,
	repo repository.Repository,
	cf config.Config,
	auth authenticator.ConsoleAuthenticator,
) (iface.IUseCase, error) {
	return &UseCase{
		repo:        repo,
		config:      cf,
		authService: auth,
	}, nil
}
