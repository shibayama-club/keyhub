package v1

import (
	"log/slog"

	"github.com/shibayama-club/keyhub/internal/usecase/app"
)

type Handler struct {
	l       *slog.Logger
	useCase app.IUseCase
}

func NewHandler(useCase app.IUseCase) *Handler {
	return &Handler{
		l:       slog.Default(),
		useCase: useCase,
	}
}
