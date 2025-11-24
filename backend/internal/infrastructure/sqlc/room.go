package sqlc

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcRoom(room sqlcgen.Room) (model.Room, error) {
	return model.Room{
		ID:             model.RoomID(room.ID),
		OrganizationID: model.OrganizationID(room.OrganizationID),
		Name:           model.RoomName(room.Name),
		BuildingName:   model.BuildingName(room.BuildingName),
		FloorNumber:    model.FloorNumber(room.FloorNumber),
		Type:           model.RoomType(room.RoomType),
		Description:    model.RoomDescription(room.Description),
		CreatedAt:      room.CreatedAt.Time,
		UpdatedAt:      room.UpdatedAt.Time,
	}, nil
}

func (t *SqlcTransaction) CreateRoom(ctx context.Context, arg repository.CreateRoomArg) error {
	return t.queries.CreateRoom(ctx, sqlcgen.CreateRoomParams{
		ID:             arg.ID.UUID(),
		OrganizationID: arg.OrganizationID.UUID(),
		Name:           arg.Name.String(),
		BuildingName:   arg.BuildingName.String(),
		FloorNumber:    arg.FloorNumber.String(),
		RoomType:       arg.Type.String(),
		Description:    arg.Description.String(),
	})
}

func (t *SqlcTransaction) GetRoomByID(ctx context.Context, id model.RoomID) (model.Room, error) {
	row, err := t.queries.GetRoomById(ctx, id.UUID())
	if err != nil {
		return model.Room{}, err
	}
	return parseSqlcRoom(row.Room)
}
