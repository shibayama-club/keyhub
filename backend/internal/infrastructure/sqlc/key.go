package sqlc

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func (t *SqlcTransaction) CreateKey(ctx context.Context, arg repository.CreateKeyArg) error {
	return t.queries.CreateKey(ctx, sqlcgen.CreateKeyParams{
		ID:             arg.ID.UUID(),
		RoomID:         arg.RoomID.UUID(),
		OrganizationID: arg.OrganizationID.UUID(),
		KeyNumber:      arg.KeyNumber.String(),
		Status:         arg.Status.String(),
	})
}
