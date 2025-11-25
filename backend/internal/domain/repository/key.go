package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateKeyArg struct {
	ID             model.KeyID
	RoomID         model.RoomID
	OrganizationID model.OrganizationID
	KeyNumber      model.KeyNumber
	Status         model.KeyStatus
}

type KeyRepository interface {
	CreateKey(ctx context.Context, arg CreateKeyArg) error
	GetKeysByRoom(ctx context.Context, roomID model.RoomID) ([]model.Key, error)
}
