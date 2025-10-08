package sqlc

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlc "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcUser(user sqlc.User) (model.User, error) {
	return model.User{
		UserId:    model.UserID(user.ID),
		Email:     model.UserEmail(user.Email),
		Name:      model.UserName(user.Name),
		Icon:      model.UserIcon(user.Icon),
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (t *SqlcTransaction) InsertUser(ctx context.Context, arg repository.InsertUserArg) (model.User, error) {
	sqlcUser, err := t.queries.InsertUser(ctx, sqlc.InsertUserParams{
		ID:    arg.ID.UUID(),
		Email: arg.Email.String(),
		Name:  arg.Name.String(),
		Icon:  arg.Icon.String(),
	})
	if err != nil {
		return model.User{}, err
	}

	return parseSqlcUser(sqlcUser)
}
