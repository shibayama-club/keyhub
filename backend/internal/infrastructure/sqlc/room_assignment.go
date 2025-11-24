package sqlc

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
	"github.com/shibayama-club/keyhub/internal/util"
)

func parseSqlcRoomAssignment(ra sqlcgen.RoomAssignment) (model.RoomAssignment, error) {
	return model.RoomAssignment{
		ID:         model.RoomAssignmentID(ra.ID),
		TenantID:   model.TenantID(ra.TenantID),
		RoomID:     model.RoomID(ra.RoomID),
		AssignedAt: ra.AssignedAt.Time,
		ExpiresAt:  util.PgTimestamptzToGoTime(ra.ExpiresAt),
		CreatedAt:  ra.CreatedAt.Time,
		UpdatedAt:  ra.UpdatedAt.Time,
	}, nil
}

func (t *SqlcTransaction) CreateRoomAssignment(ctx context.Context, arg repository.CreateRoomAssignmentArg) error {
	return t.queries.CreateRoomAssignment(ctx, sqlcgen.CreateRoomAssignmentParams{
		ID:         arg.ID.UUID(),
		TenantID:   arg.TenantID.UUID(),
		RoomID:     arg.RoomID.UUID(),
		AssignedAt: util.GoTimeToPgTimestamptz(&arg.AssignedAt),
		ExpiresAt:  util.GoTimeToPgTimestamptz(arg.ExpiresAt),
	})
}
