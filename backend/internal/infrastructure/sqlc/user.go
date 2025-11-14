package sqlc

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcUser(user sqlcgen.User) (model.User, error) {
	return model.User{
		UserId:    model.UserID(user.ID),
		Email:     model.UserEmail(user.Email),
		Name:      model.UserName(user.Name),
		Icon:      model.UserIcon(user.Icon),
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (t *SqlcTransaction) GetUser(ctx context.Context, userID model.UserID) (model.User, error) {
	sqlcUserRow, err := t.queries.GetUser(ctx, userID.UUID())
	if err != nil {
		return model.User{}, err
	}

	return parseSqlcUser(sqlcUserRow.User)
}

func (t *SqlcTransaction) UpsertUser(ctx context.Context, arg repository.UpsertUserArg) (model.User, error) {
	sqlcUserRow, err := t.queries.UpsertUser(ctx, sqlcgen.UpsertUserParams{
		Email: arg.Email.String(),
		Name:  arg.Name.String(),
		Icon:  arg.Icon.String(),
	})
	if err != nil {
		return model.User{}, err
	}
	return parseSqlcUser(sqlcUserRow.User)
}

func (t *SqlcTransaction) GetUserByProviderIdentity(ctx context.Context, provider, providerSub string) (model.User, error) {
	sqlcUserRow, err := t.queries.GetUserByProviderIdentity(ctx, sqlcgen.GetUserByProviderIdentityParams{
		Provider:    provider,
		ProviderSub: providerSub,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return parseSqlcUser(sqlcUserRow.User)
}

func (t *SqlcTransaction) UpsertUserIdentity(ctx context.Context, arg repository.UpsertUserIdentityArg) error {
	return t.queries.UpsertUserIdentity(ctx, sqlcgen.UpsertUserIdentityParams{
		UserID:      arg.UserID.UUID(),
		Provider:    arg.Provider,
		ProviderSub: arg.ProviderSub,
	})
}
