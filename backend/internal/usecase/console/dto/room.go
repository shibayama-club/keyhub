package dto

import (
	"time"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateRoomInput struct {
	OrganizationID model.OrganizationID
	Name           string
	BuildingName   string
	FloorNumber    string
	RoomType       string
	Description    string
}

type AssignRoomToTenantInput struct {
	TenantID  model.TenantID
	RoomID    model.RoomID
	ExpiresAt *time.Time
}
