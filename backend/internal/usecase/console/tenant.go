package console

import (
	"context"

	"github.com/cockroachdb/errors"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	"github.com/shibayama-club/keyhub/internal/usecase/console/dto"
)

func (u *UseCase) CreateTenant(ctx context.Context, input dto.CreateTenantInput) (string, error) {
	tenantName, err := model.NewTenantName(input.Name)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid tenant name")
	}

	tenantDescription, err := model.NewTenantDescription(input.Description)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid tenant description")
	}

	tenantType, err := model.NewTenantType(input.TenantType)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid tenant type")
	}

	tenant, err := model.NewTenant(
		input.OrganizationID,
		tenantName,
		tenantDescription,
		tenantType,
	)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create tenant")
	}

	joinCode, err := model.NewTenantJoinCode(input.JoinCode)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid join code")
	}

	joinCodeExpiry, err := model.NewTenantJoinCodeExpiresAt(input.JoinCodeExpiry)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid join code expiry")
	}

	joinCodeMaxUse, err := model.NewTenantJoinCodeMaxUses(input.JoinCodeMaxUse)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid join code max use")
	}

	joinCodeEntity, err := model.NewTenantJoinCodeEntity(
		tenant.ID,
		joinCode,
		joinCodeExpiry,
		joinCodeMaxUse,
	)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create tenant join code entity")
	}

	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		err = tx.CreateTenant(ctx, repository.CreateTenantArg{
			ID:             tenant.ID,
			OrganizationID: tenant.OrganizationID,
			Name:           tenant.Name,
			Description:    tenant.Description,
			Type:           tenant.Type,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create tenant in repository")
		}

		err = tx.CreateTenantJoinCode(ctx, repository.CreateTenantJoinCodeArg{
			ID:        joinCodeEntity.ID,
			TenantID:  joinCodeEntity.TenantID,
			Code:      joinCodeEntity.Code,
			ExpiresAt: joinCodeEntity.ExpiresAt,
			MaxUses:   joinCodeEntity.MaxUses,
			UsedCount: joinCodeEntity.UsedCount,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create tenant join code in repository")
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return tenant.ID.String(), nil
}

func (u *UseCase) GetAllTenants(ctx context.Context) ([]model.Tenant, error) {
	tenants, err := u.repo.GetAllTenants(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to get all tenants from repository")
	}

	return tenants, nil
}

func (u *UseCase) GetTenantById(ctx context.Context, TenantId model.TenantID) (dto.GetTenantByIdOutput, error) {
	tenantWithJoinCode, err := u.repo.GetTenantByID(ctx, TenantId)
	if err != nil {
		return dto.GetTenantByIdOutput{}, errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to get a tenant by id from repository")
	}

	return dto.GetTenantByIdOutput{
		Tenant:   tenantWithJoinCode.Tenant,
		JoinCode: tenantWithJoinCode.JoinCode,
	}, nil
}

func (u *UseCase) UpdateTenant(ctx context.Context, input dto.UpdateTenantInput) error {
	tenantName, err := model.NewTenantName(input.Name)
	if err != nil {
		return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid tenant name")
	}

	tenantDescription, err := model.NewTenantDescription(input.Description)
	if err != nil {
		return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid tenant description")
	}

	tenantType, err := model.NewTenantType(input.TenantType)
	if err != nil {
		return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid tenant type")
	}

	joinCode, err := model.NewTenantJoinCode(input.JoinCode)
	if err != nil {
		return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid join code")
	}

	joinCodeExpiry, err := model.NewTenantJoinCodeExpiresAt(input.JoinCodeExpiry)
	if err != nil {
		return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid join code expiry")
	}

	joinCodeMaxUse, err := model.NewTenantJoinCodeMaxUses(input.JoinCodeMaxUse)
	if err != nil {
		return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid join code max use")
	}

	joinCodeEntity, err := model.NewTenantJoinCodeEntity(
		input.TenantID,
		joinCode,
		joinCodeExpiry,
		joinCodeMaxUse,
	)
	if err != nil {
		return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create tenant join code entity")
	}

	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		err = tx.UpdateTenant(ctx, repository.UpdateTenantArg{
			ID:          input.TenantID,
			Name:        tenantName,
			Description: tenantDescription,
			Type:        tenantType,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to update a tenant in repository")
		}
		err = tx.UpdateTenantJoinCodeByTenantId(ctx, repository.UpdateTenantJoinCodeArg{
			TenantID:  input.TenantID,
			Code:      joinCodeEntity.Code,
			ExpiresAt: joinCodeEntity.ExpiresAt,
			MaxUses:   joinCodeEntity.MaxUses,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to update tenant join code in repository")
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
