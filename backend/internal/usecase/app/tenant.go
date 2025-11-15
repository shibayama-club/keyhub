package app

import (
	"context"

	"github.com/cockroachdb/errors"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/usecase/app/dto"
)

func (u *UseCase) GetTenantByJoinCode(ctx context.Context, joinCode string) (dto.GetTenantByJoinCodeOutput, error) {
	code := model.TenantJoinCode(joinCode)

	tenant, err := u.repo.GetTenantByJoinCode(ctx, code)
	if err != nil {
		return dto.GetTenantByJoinCodeOutput{}, errors.Wrap(errors.Mark(err, domainerrors.ErrNotFound), "tenant not found")
	}

	return dto.GetTenantByJoinCodeOutput{
		ID:          tenant.ID.String(),
		Name:        tenant.Name.String(),
		Description: tenant.Description.String(),
		TenantType:  tenant.Type.String(),
	}, nil
}
