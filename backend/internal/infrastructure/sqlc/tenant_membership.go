package sqlc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcTenantMembership(row sqlcgen.TenantMembership) (model.TenantMembership, error) {
	return model.TenantMembership{
		ID:        model.TenantMembershipID(row.ID),
		TenantID:  model.TenantID(row.TenantID),
		UserID:    model.UserID(row.UserID),
		Role:      model.TenantMembershipRole(row.Role),
		CreatedAt: row.CreatedAt.Time,
		LeftAt:    timestamptzPtrValue(row.LeftAt),
	}, nil
}

func timestamptzPtrValue(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func (t *SqlcTransaction) CreateTenantMembership(ctx context.Context, membership model.TenantMembership) error {
	return t.queries.CreateTenantMembership(ctx, sqlcgen.CreateTenantMembershipParams{
		ID:       membership.ID.UUID(),
		TenantID: membership.TenantID.UUID(),
		UserID:   membership.UserID.UUID(),
		Role:     membership.Role.String(),
	})
}

func (t *SqlcTransaction) IncrementJoinCodeUsedCount(ctx context.Context, code model.TenantJoinCode) error {
	return t.queries.IncrementJoinCodeUsedCount(ctx, code.String())
}

func (t *SqlcTransaction) GetTenantMembershipByTenantAndUser(ctx context.Context, tenantID model.TenantID, userID model.UserID) (model.TenantMembership, error) {
	sqlcRow, err := t.queries.GetTenantMembershipByTenantAndUser(ctx, sqlcgen.GetTenantMembershipByTenantAndUserParams{
		TenantID: tenantID.UUID(),
		UserID:   userID.UUID(),
	})
	if err != nil {
		return model.TenantMembership{}, err
	}
	return parseSqlcTenantMembership(sqlcRow.TenantMembership)
}

func (t *SqlcTransaction) ClearActiveMembershipByTenantID(ctx context.Context, tenantId model.TenantID) error {
	err := t.queries.ClearActiveMembershipByTenantID(ctx, tenantId.UUID())
	if err != nil {
		return err
	}
	return nil
}
