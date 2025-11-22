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

type UpdateTenantJoinCodeArg struct {
	TenantID  model.TenantID
	Code      model.TenantJoinCode
	ExpiresAt model.TenantJoinCodeExpiresAt
	MaxUses   model.TenantJoinCodeMaxUses
}
type TenantJoinCodeRepository interface {
	CreateTenantJoinCode(ctx context.Context, arg CreateTenantJoinCodeArg) (model.TenantJoinCodeEntity, error)
	GetTenantByJoinCode(ctx context.Context, code model.TenantJoinCode) (model.Tenant, error)
	UpdateTenantJoinCodeByTenantId(ctx context.Context, arg UpdateTenantJoinCodeArg) error
}
