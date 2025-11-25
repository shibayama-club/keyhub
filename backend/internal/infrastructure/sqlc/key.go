package sqlc

import (
	"context"

	"github.com/samber/lo"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcKey(key sqlcgen.Key) (model.Key, error) {
	return model.Key{
		ID:             model.KeyID(key.ID),
		RoomID:         model.RoomID(key.RoomID),
		OrganizationID: model.OrganizationID(key.OrganizationID),
		KeyNumber:      model.KeyNumber(key.KeyNumber),
		Status:         model.KeyStatus(key.Status),
		CreatedAt:      key.CreatedAt.Time,
		UpdatedAt:      key.UpdatedAt.Time,
	}, nil
}

func (t *SqlcTransaction) CreateKey(ctx context.Context, arg repository.CreateKeyArg) error {
	return t.queries.CreateKey(ctx, sqlcgen.CreateKeyParams{
		ID:             arg.ID.UUID(),
		RoomID:         arg.RoomID.UUID(),
		OrganizationID: arg.OrganizationID.UUID(),
		KeyNumber:      arg.KeyNumber.String(),
		Status:         arg.Status.String(),
	})
}

func (t *SqlcTransaction) GetKeysByRoom(ctx context.Context, roomID model.RoomID) ([]model.Key, error) {
	rows, err := t.queries.GetKeysByRoom(ctx, roomID.UUID())
	if err != nil {
		return nil, err
	}

	keys := lo.Map(rows, func(row sqlcgen.GetKeysByRoomRow, _ int) model.Key {
		key, _ := parseSqlcKey(row.Key)
		return key
	})

	return keys, nil
}
