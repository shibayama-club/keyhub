package repository

import (
	"context"
	"time"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateRoomAssignmentArg struct {
	ID         model.RoomAssignmentID
	TenantID   model.TenantID
	RoomID     model.RoomID
	AssignedAt time.Time
	ExpiresAt  *time.Time
}

type RoomAssignmentRepository interface {
	CreateRoomAssignment(ctx context.Context, arg CreateRoomAssignmentArg) error
}
