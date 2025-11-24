package sqlc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
	"github.com/shibayama-club/keyhub/internal/util"
)

func parseSqlcTenantJoinCode(tjc sqlcgen.TenantJoinCode) (model.TenantJoinCodeEntity, error) {
	var expiresAt *time.Time
	if tjc.ExpiresAt.Valid {
		t := tjc.ExpiresAt.Time
		expiresAt = &t
	}

	expiresAtVO, err := model.NewTenantJoinCodeExpiresAt(expiresAt)
	if err != nil {
		return model.TenantJoinCodeEntity{}, err
	}

	return model.TenantJoinCodeEntity{
		ID:        model.TenantJoinCodeID(tjc.ID),
		TenantID:  model.TenantID(tjc.TenantID),
		Code:      model.TenantJoinCode(tjc.Code),
		ExpiresAt: expiresAtVO,
		MaxUses:   model.TenantJoinCodeMaxUses(tjc.MaxUses),
		UsedCount: int(tjc.UsedCount),
		CreatedAt: tjc.CreatedAt.Time,
	}, nil
}

func (t *SqlcTransaction) CreateTenantJoinCode(ctx context.Context, arg repository.CreateTenantJoinCodeArg) error {
	var expiresAt pgtype.Timestamptz
	if arg.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{
			Time:  *arg.ExpiresAt,
			Valid: true,
		}
	}

	return t.queries.CreateTenantJoinCode(ctx, sqlcgen.CreateTenantJoinCodeParams{
		ID:        arg.ID.UUID(),
		TenantID:  arg.TenantID.UUID(),
		Code:      arg.Code.String(),
		ExpiresAt: expiresAt,
		MaxUses:   arg.MaxUses.Int32(),
		UsedCount: int32(arg.UsedCount),
	})
}

func (t *SqlcTransaction) GetTenantByJoinCode(ctx context.Context, code model.TenantJoinCode) (model.Tenant, error) {
	sqlcRow, err := t.queries.GetTenantByJoinCode(ctx, code.String())
	if err != nil {
		return model.Tenant{}, err
	}
	return parseSqlcTenant(sqlcRow.Tenant)
}

func (t *SqlcTransaction) UpdateTenantJoinCodeByTenantId(ctx context.Context, arg repository.UpdateTenantJoinCodeArg) error {
	expiresAt := util.PaeseTimeToPgtypeTimestamptz(arg.ExpiresAt)
	err := t.queries.UpdateTenantJoinCodeByTenantId(ctx, sqlcgen.UpdateTenantJoinCodeByTenantIdParams{
		TenantID:  arg.TenantID.UUID(),
		Code:      arg.Code.String(),
		ExpiresAt: expiresAt,
		MaxUses:   arg.MaxUses.Int32(),
	})
	if err != nil {
		return err
	}
	return nil
}
