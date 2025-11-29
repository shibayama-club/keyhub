package app

import (
	"context"

	"github.com/cockroachdb/errors"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
)

func (u *UseCase) GetRoomsByTenant(ctx context.Context, tenantID model.TenantID) ([]model.Room, error) {
	rooms, err := u.repo.GetRoomsByTenant(ctx, tenantID)
	if err != nil {
		return nil, errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to get rooms by tenant")
	}
	return rooms, nil
}

func (u *UseCase) GetKeysByRoom(ctx context.Context, roomID model.RoomID) ([]model.Key, error) {
	keys, err := u.repo.GetKeysByRoom(ctx, roomID)
	if err != nil {
		return nil, errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to get keys by room")
	}
	return keys, nil
}
