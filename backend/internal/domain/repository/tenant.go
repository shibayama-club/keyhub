package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type InsertTenantArg struct {
	Id          model.TenantID
	Name        model.TenantName
	Description model.TenantDescription
	Type        model.TenantType
}

type TenantRepository interface {
	InsertTenant(ctx context.Context, arg InsertTenantArg) (model.Tenant, error)
}
