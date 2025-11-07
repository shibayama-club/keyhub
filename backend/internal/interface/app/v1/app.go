package v1

import (
	"log/slog"

	"github.com/shibayama-club/keyhub/internal/usecase/app/iface"
)

type Handler struct {
	l       *slog.Logger
	useCase iface.IUseCase
}

func NewHandler(useCase iface.IUseCase) *Handler {
	return &Handler{
		l:       slog.Default(),
		useCase: useCase,
	}
}
