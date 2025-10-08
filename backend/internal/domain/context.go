package domain

import "context"

type (
	ctxKey[T any] struct{}
)

func WithValue[T any](ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, ctxKey[T]{}, val)
}

func RemoveValue[T any](ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey[T]{}, nil)
}

func Value[T any](ctx context.Context) (T, bool) {
	value, ok := ctx.Value(ctxKey[T]{}).(T)
	return value, ok
}
