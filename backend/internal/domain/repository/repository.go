package repository

import "context"

type Repository interface {
	Transaction

	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}

type Transaction interface {
	UserRepository
}