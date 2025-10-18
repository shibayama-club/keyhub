package v1

import (
	"log/slog"

	"github.com/cockroachdb/errors"
	authConsole "github.com/shibayama-club/keyhub/internal/infrastructure/auth/console"
	"github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/console/v1/consolev1connect"
	"github.com/shibayama-club/keyhub/internal/usecase/console"
)

type Handler struct {
	consolev1connect.UnimplementedConsoleAuthServiceHandler
	l           *slog.Logger
	useCase     console.IUseCase
	authService *authConsole.AuthService
}

func NewHandler(useCase console.IUseCase, jwtSecret string) (*Handler, error) {
	authService, err := authConsole.NewAuthService(jwtSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth service")
	}

	return &Handler{
		l:           slog.Default(),
		useCase:     useCase,
		authService: authService,
	}, nil
}
