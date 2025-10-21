package iface

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/usecase/console/dto"
)

type IUseCase interface {
	LoginWithOrgId(ctx context.Context, orgID, orgKey string) (string, int64, error)
	Logout(ctx context.Context, sessionID string) error
	ValidateSession(ctx context.Context, token string) (model.ConsoleSession, error)
	CreateTenant(ctx context.Context, input dto.CreateTenantInput) (dto.CreateTenantOutput, error)
}
