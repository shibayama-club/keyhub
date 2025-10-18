package sqlc

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcTenant(tenant sqlcgen.Tenant) (model.Tenant, error) {

	return model.Tenant{
		ID:          model.TenantID(tenant.ID),
		Name:        model.TenantName(tenant.Name),
		Description: model.TenantDescription(tenant.Description),
		Type:        model.TenantType(tenant.TenantType),
		CreatedAt:   tenant.CreatedAt.Time,
		UpdatedAt:   tenant.UpdatedAt.Time,
	}, nil
}

func (t *SqlcTransaction) InsertTenant(ctx context.Context, arg repository.InsertTenantArg) (model.Tenant, error) {
	sqlcTenant, err := t.queries.InsertTenant(ctx, sqlcgen.InsertTenantParams{
		ID:           arg.Id.UUID(),
		Name:         arg.Name.String(),
		Slug:         arg.Description.String(),
		PasswordHash: arg.Type.String(),
	})
	if err != nil {
		return model.Tenant{}, err
	}
	return parseSqlcTenant(sqlcTenant)
}
