package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type InsertUserArg struct {
	ID    model.UserID
	Email model.UserEmail
	Name  model.UserName
	Icon  model.UserIcon
}

type UpsertUserArg struct {
	Email model.UserEmail
	Name  model.UserName
	Icon  model.UserIcon
}

type UpsertUserIdentityArg struct {
	UserID      model.UserID
	Provider    string
	ProviderSub string
}

type UserRepository interface {
	InsertUser(ctx context.Context, arg InsertUserArg) (model.User, error)
	GetUser(ctx context.Context, userID model.UserID) (model.User, error)
	GetUserByProviderIdentity(ctx context.Context, provider, providerSub string) (model.User, error)
	UpsertUser(ctx context.Context, arg UpsertUserArg) (model.User, error)
	UpsertUserIdentity(ctx context.Context, arg UpsertUserIdentityArg) error
}
