package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateRoomArg struct {
	ID             model.RoomID
	OrganizationID model.OrganizationID
	Name           model.RoomName
	BuildingName   model.BuildingName
	FloorNumber    model.FloorNumber
	Type           model.RoomType
	Description    model.RoomDescription
}

type RoomRepository interface {
	CreateRoom(ctx context.Context, arg CreateRoomArg) error
	GetRoomByID(ctx context.Context, id model.RoomID) (model.Room, error)
	GetAllRooms(ctx context.Context) ([]model.Room, error)
	GetRoomsByTenant(ctx context.Context, tenantID model.TenantID) ([]model.Room, error)
}
