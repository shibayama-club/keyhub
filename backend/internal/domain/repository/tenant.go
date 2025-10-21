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

type TenantRepository interface {
	CreateTenant(ctx context.Context, arg CreateTenantArg) (model.Tenant, error)
	GetTenant(ctx context.Context, id model.TenantID) (model.Tenant, error)
	GetTenantsByOrganization(ctx context.Context, organizationID model.OrganizationID) ([]model.Tenant, error)
}
