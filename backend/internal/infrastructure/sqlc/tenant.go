package sqlc

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlc "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSslcTenant(tenant sqlc.Tenant)(model.Tenant, error){
	return model.Tenant{
		TenantId: model.TenantID(tenant.ID),
		TenantName: model.TenantName(tenant.Name),
		TenantSlug: model.TenantSlug(tenant.Slug),
		TenantPasswordHash: model.TenantPasswordHash(tenant.PasswordHash),
		CreatedAt: tenant.CreatedAt.Time,
		UpdatedAt: tenant.UpdatedAt.Time,
	},nil
}

func (t *SqlcTransaction) InsertTenant(ctx context.Context, arg repository.InsertTenantArg)(model.Tenant, error){
	sqlcTenant, err := t.queries.InsertTenant(ctx, sqlc.InsertTenantParams{
		ID: arg.Id.UUID(),
		Name: arg.Name.String(),
		Slug: arg.Slug.String(),
		PasswordHash: arg.PasswordHash.String(),
	})
	if err != nil{
		return model.Tenant{}, err
	}
	return parseSslcTenant(sqlcTenant)
}

