package app

import (
	"context"

	"github.com/shibayama-club/keyhub/cmd/config"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
)

type IUseCase interface{}

type UseCase struct {
	repo   repository.Repository
	config config.Config
}

func NewUseCase(ctx context.Context, repo repository.Repository, cf config.Config) IUseCase {
	return &UseCase{
		repo:   repo,
		config: cf,
	}
}
