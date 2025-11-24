package sqlc

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
	"github.com/shibayama-club/keyhub/internal/util"
)

func (t *SqlcTransaction) CreateRoomAssignment(ctx context.Context, arg repository.CreateRoomAssignmentArg) error {
	return t.queries.CreateRoomAssignment(ctx, sqlcgen.CreateRoomAssignmentParams{
		ID:         arg.ID.UUID(),
		TenantID:   arg.TenantID.UUID(),
		RoomID:     arg.RoomID.UUID(),
		AssignedAt: util.GoTimeToPgTimestamptz(&arg.AssignedAt),
		ExpiresAt:  util.GoTimeToPgTimestamptz(arg.ExpiresAt),
	})
}
