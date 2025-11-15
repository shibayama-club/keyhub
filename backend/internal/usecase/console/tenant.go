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
	if input.Name == "" {
		return "", errors.WithHint(
			errors.Mark(errors.New("name is required"), domainerrors.ErrValidation),
			"テナント名を入力してください。",
		)
	}

	tenant, err := model.NewTenant(
		input.OrganizationID,
		input.Name,
		input.Description,
		input.TenantType,
	)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create tenant")
	}

	joinCodeEntity, err := model.NewTenantJoinCodeEntity(
		tenant.ID,
		input.JoinCode,
		input.JoinCodeExpiry,
		input.JoinCodeMaxUse,
	)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create tenant join code entity")
	}

	var createdTenant model.Tenant
	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		createdTenant, err = tx.CreateTenant(ctx, repository.CreateTenantArg{
			ID:             tenant.ID,
			OrganizationID: tenant.OrganizationID,
			Name:           tenant.Name,
			Description:    tenant.Description,
			Type:           tenant.Type,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create tenant in repository")
		}

		_, err = tx.CreateTenantJoinCode(ctx, repository.CreateTenantJoinCodeArg{
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

	return createdTenant.ID.String(), nil
}

func (u *UseCase) GetAllTenants(ctx context.Context, organizationID model.OrganizationID) ([]model.Tenant, error) {
	tenants, err := u.repo.GetAllTenants(ctx, organizationID)
	if err != nil {
		return nil, errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to get all tenants from repository")
	}

	return tenants, nil
}
