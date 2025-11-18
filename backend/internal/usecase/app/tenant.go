package app

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
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

func (u *UseCase) JoinTenant(ctx context.Context, userID model.UserID, joinCode string) error {
	code := model.TenantJoinCode(joinCode)

	err := u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		tenant, err := tx.GetTenantByJoinCode(ctx, code)
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrNotFound), "tenant not found")
		}

		_, err = tx.GetTenantMembershipByTenantAndUser(ctx, tenant.ID, userID)
		if err == nil {
			return errors.Mark(errors.New("user is already a member of this tenant"), domainerrors.ErrAlreadyExists)
		}

		membership := model.TenantMembership{
			ID:       model.TenantMembershipID(uuid.New()),
			TenantID: tenant.ID,
			UserID:   userID,
			Role:     model.TenantMembershipRoleMember,
		}

		_, err = tx.CreateTenantMembership(ctx, membership)
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create tenant membership in repository")
		}

		err = tx.IncrementJoinCodeUsedCount(ctx, code)
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to increment join code used count in repository")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
