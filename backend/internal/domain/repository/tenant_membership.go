package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type TenantMembershipRepository interface {
	CreateTenantMembership(ctx context.Context, membership model.TenantMembership) error
	IncrementJoinCodeUsedCount(ctx context.Context, code model.TenantJoinCode) error
	GetTenantMembershipByTenantAndUser(ctx context.Context, tenantID model.TenantID, userID model.UserID) (model.TenantMembership, error)
}
