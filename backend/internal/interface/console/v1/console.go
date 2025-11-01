package v1

import (
	"log/slog"

	"github.com/cockroachdb/errors"
	authConsole "github.com/shibayama-club/keyhub/internal/infrastructure/auth/console"
	"github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/console/v1/consolev1connect"
	"github.com/shibayama-club/keyhub/internal/usecase/console/iface"
)

type Handler struct {
	consolev1connect.UnimplementedConsoleAuthServiceHandler
	consolev1connect.UnimplementedConsoleServiceHandler
	l           *slog.Logger
	useCase     iface.IUseCase
	authService *authConsole.AuthService
}

func NewHandler(useCase iface.IUseCase, jwtSecret string) (*Handler, error) {
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
