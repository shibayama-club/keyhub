package sqlc

import (
	"context"

	"github.com/samber/lo"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcTenant(tenant sqlcgen.Tenant) (model.Tenant, error) {
	return model.Tenant{
		ID:             model.TenantID(tenant.ID),
		OrganizationID: model.OrganizationID(tenant.OrganizationID),
		Name:           model.TenantName(tenant.Name),
		Description:    model.TenantDescription(tenant.Description),
		Type:           model.TenantType(tenant.TenantType),
		CreatedAt:      tenant.CreatedAt.Time,
		UpdatedAt:      tenant.UpdatedAt.Time,
	}, nil
}

func (t *SqlcTransaction) CreateTenant(ctx context.Context, arg repository.CreateTenantArg) (model.Tenant, error) {
	sqlcTenantRow, err := t.queries.CreateTenant(ctx, sqlcgen.CreateTenantParams{
		ID:             arg.ID.UUID(),
		OrganizationID: arg.OrganizationID.UUID(),
		Name:           arg.Name.String(),
		Description:    arg.Description.String(),
		TenantType:     arg.Type.String(),
	})
	if err != nil {
		return model.Tenant{}, err
	}
	return parseSqlcTenant(sqlcTenantRow.Tenant)
}

func (t *SqlcTransaction) GetTenant(ctx context.Context, id model.TenantID) (model.Tenant, error) {
	row, err := t.queries.GetTenant(ctx, id.UUID())
	if err != nil {
		return model.Tenant{}, err
	}
	return parseSqlcTenant(row.Tenant)
}

func (t *SqlcTransaction) GetAllTenants(ctx context.Context, organizationID model.OrganizationID) ([]model.Tenant, error) {
	rows, err := t.queries.GetAllTenants(ctx, organizationID.UUID())
	if err != nil {
		return nil, err
	}

	tenants := lo.Map(rows, func(row sqlcgen.GetAllTenantsRow, _ int) model.Tenant {
		tenant, _ := parseSqlcTenant(row.Tenant)
		return tenant
	})

	return tenants, nil
}
