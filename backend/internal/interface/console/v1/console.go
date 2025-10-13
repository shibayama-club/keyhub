package v1

import (
	"log/slog"

	"github.com/shibayama-club/keyhub/internal/usecase/console"
)

type Handler struct {
	l       *slog.Logger
	useCase console.IUseCase
}

func NewHandler(useCase console.IUseCase) *Handler {
	return &Handler{
		l:       slog.Default(),
		useCase: useCase,
	}
}
