package repository

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type InsertUserArg struct {
	ID    model.UserID
	Email model.UserEmail
	Name  model.UserName
	Icon    model.UserIcon
}

type UserRepository interface {
	InsertUser(ctx context.Context, arg InsertUserArg) (model.User, error)
}
