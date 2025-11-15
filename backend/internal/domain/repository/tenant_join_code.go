package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateTenantJoinCodeArg struct {
	ID        model.TenantJoinCodeID
	TenantID  model.TenantID
	Code      model.TenantJoinCode
	ExpiresAt model.TenantJoinCodeExpiresAt
	MaxUses   model.TenantJoinCodeMaxUses
	UsedCount int
}

type TenantJoinCodeRepository interface {
	CreateTenantJoinCode(ctx context.Context, arg CreateTenantJoinCodeArg) (model.TenantJoinCodeEntity, error)
}
