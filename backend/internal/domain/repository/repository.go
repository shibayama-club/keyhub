package repository

//go:generate go run go.uber.org/mock/mockgen@latest -source=$GOFILE -destination=mock/mock_repository.go -package=mock

import "context"

type Repository interface {
	Transaction

	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}

type Transaction interface {
	UserRepository
	TenantRepository
	ConsoleSessionRepository
}
