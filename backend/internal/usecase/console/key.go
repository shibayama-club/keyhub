package console

import (
	"context"

	"github.com/cockroachdb/errors"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	"github.com/shibayama-club/keyhub/internal/usecase/console/dto"
)

func (u *UseCase) CreateKey(ctx context.Context, input dto.CreateKeyInput) (string, error) {
	// Verify room exists
	_, err := u.repo.GetRoomByID(ctx, input.RoomID)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrNotFound), "room not found")
	}

	keyNumber, err := model.NewKeyNumber(input.KeyNumber)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid key number")
	}

	key, err := model.NewKey(
		input.RoomID,
		input.OrganizationID,
		keyNumber,
	)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create key")
	}

	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		err = tx.CreateKey(ctx, repository.CreateKeyArg{
			ID:             key.ID,
			RoomID:         key.RoomID,
			OrganizationID: key.OrganizationID,
			KeyNumber:      key.KeyNumber,
			Status:         key.Status,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create key in repository")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return key.ID.String(), nil
}
