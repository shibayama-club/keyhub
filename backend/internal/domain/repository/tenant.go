package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateTenantArg struct {
	ID             model.TenantID
	OrganizationID model.OrganizationID
	Name           model.TenantName
	Description    model.TenantDescription
	Type           model.TenantType
}

type TenantWithMemberCount struct {
	Tenant      model.Tenant
	MemberCount int32
}
type UpdateTenantArg struct {
	ID          model.TenantID
	Name        model.TenantName
	Description model.TenantDescription
	Type        model.TenantType
}
type TenantWithJoinCode struct {
	Tenant   model.Tenant
	JoinCode model.TenantJoinCodeEntity
}

type TenantRepository interface {
	CreateTenant(ctx context.Context, arg CreateTenantArg) (model.Tenant, error)
	GetAllTenants(ctx context.Context, organizationID model.OrganizationID) ([]model.Tenant, error)
	GetTenantsByUserID(ctx context.Context, userID model.UserID) ([]TenantWithMemberCount, error)
	GetTenantByID(ctx context.Context, id model.TenantID) (TenantWithJoinCode, error)
	UpdateTenant(ctx context.Context, arg UpdateTenantArg) error
}
