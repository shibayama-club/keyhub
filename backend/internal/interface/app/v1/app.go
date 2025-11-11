package v1

import (
	"log/slog"

	"github.com/shibayama-club/keyhub/internal/usecase/app/iface"
)

type Handler struct {
	l           *slog.Logger
	useCase     iface.IUseCase
	env         string
	frontendURL string
}

func NewHandler(useCase iface.IUseCase, env, frontendURL string) *Handler {
	return &Handler{
		l:           slog.Default(),
		useCase:     useCase,
		env:         env,
		frontendURL: frontendURL,
	}
}
