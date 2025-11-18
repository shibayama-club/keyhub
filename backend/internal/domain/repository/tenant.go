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

type TenantRepository interface {
	CreateTenant(ctx context.Context, arg CreateTenantArg) (model.Tenant, error)
	GetTenant(ctx context.Context, id model.TenantID) (model.Tenant, error)
	GetAllTenants(ctx context.Context, organizationID model.OrganizationID) ([]model.Tenant, error)
	GetTenantsByUserID(ctx context.Context, userID model.UserID) ([]TenantWithMemberCount, error)
}
