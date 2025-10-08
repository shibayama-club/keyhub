package sqlc

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlc "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSslcTenant(tenant sqlc.Tenant) (model.Tenant, error) {
	var tenantSlug *model.TenantSlug
	if tenant.Slug != nil {
		slug, err := model.NewTenantSlug(*tenant.Slug)
		if err != nil {
			return model.Tenant{}, err
		}
		tenantSlug = &slug
	}
	return model.Tenant{
		TenantId:           model.TenantID(tenant.ID),
		TenantName:         model.TenantName(tenant.Name),
		TenantSlug:         tenantSlug,
		TenantPasswordHash: model.TenantPasswordHash(tenant.PasswordHash),
		CreatedAt:          tenant.CreatedAt.Time,
		UpdatedAt:          tenant.UpdatedAt.Time,
	}, nil
}

func (t *SqlcTransaction) InsertTenant(ctx context.Context, arg repository.InsertTenantArg) (model.Tenant, error) {
	var slugPtr *string
	if slug := arg.Slug.String(); slug != "" {
		s := slug
		slugPtr = &s
	}
	sqlcTenant, err := t.queries.InsertTenant(ctx, sqlc.InsertTenantParams{
		ID:           arg.Id.UUID(),
		Name:         arg.Name.String(),
		Slug:         slugPtr,
		PasswordHash: arg.PasswordHash.String(),
	})
	if err != nil {
		return model.Tenant{}, err
	}
	return parseSslcTenant(sqlcTenant)
}
