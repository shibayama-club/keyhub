package sqlc

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlcRepository struct {
	pool *pgxpool.Pool
	*SqlcTransaction
}

type SqlcTransaction struct {
	queries *gen.Queries
}

func NewRepository(pool *pgxpool.Pool) *SqlcRepository {
	return &SqlcRepository{
		pool: pool,
		SqlcTransaction: &SqlcTransaction{
			queries: gen.New(pool),
		},
	}
}

func (r *SqlcRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx repository.Transaction) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx, &SqlcTransaction{queries: r.queries.WithTx(tx)}); err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return errors.Wrapf(err, "failed to Queries and Rollback: %v", txErr)
		}
		return errors.Wrap(err, "failed to Queries")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to Commit")
	}
	return nil
}

func (r *SqlcRepository) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}