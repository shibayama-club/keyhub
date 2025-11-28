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

func (t *SqlcTransaction) CreateTenant(ctx context.Context, arg repository.CreateTenantArg) error {
	return t.queries.CreateTenant(ctx, sqlcgen.CreateTenantParams{
		ID:             arg.ID.UUID(),
		OrganizationID: arg.OrganizationID.UUID(),
		Name:           arg.Name.String(),
		Description:    arg.Description.String(),
		TenantType:     arg.Type.String(),
	})
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

func (t *SqlcTransaction) GetTenantsByUserID(ctx context.Context, userID model.UserID) ([]repository.TenantWithMemberCount, error) {
	rows, err := t.queries.GetTenantsByUserID(ctx, userID.UUID())
	if err != nil {
		return nil, err
	}

	tenants := lo.Map(rows, func(row sqlcgen.GetTenantsByUserIDRow, _ int) repository.TenantWithMemberCount {
		tenant, _ := parseSqlcTenant(row.Tenant)
		return repository.TenantWithMemberCount{
			Tenant:      tenant,
			MemberCount: row.MemberCount,
		}
	})

	return tenants, nil
}
func (t *SqlcTransaction) GetTenantByID(ctx context.Context, id model.TenantID) (repository.TenantWithJoinCode, error) {
	row, err := t.queries.GetTenantById(ctx, id.UUID())
	if err != nil {
		return repository.TenantWithJoinCode{}, err
	}
	tenantRow := row.Tenant
	joinCodeRow := row.TenantJoinCode
	tenant, err := parseSqlcTenant(tenantRow)
	if err != nil {
		return repository.TenantWithJoinCode{}, err
	}
	joinCodeEntity, err := parseSqlcTenantJoinCode(joinCodeRow)
	if err != nil {
		return repository.TenantWithJoinCode{}, err
	}
	return repository.TenantWithJoinCode{
		Tenant:   tenant,
		JoinCode: joinCodeEntity,
	}, nil
}

func (t *SqlcTransaction) UpdateTenant(ctx context.Context, arg repository.UpdateTenantArg) error {
	err := t.queries.UpdateTenant(ctx, sqlcgen.UpdateTenantParams{
		ID:          arg.ID.UUID(),
		Name:        arg.Name.String(),
		Description: arg.Description.String(),
		TenantType:  arg.Type.String(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *SqlcTransaction) DeleteTenant(ctx context.Context, id model.TenantID) error {
	err := t.queries.DeleteTenantById(ctx, id.UUID())
	if err != nil {
		return err
	}
	return nil
}
