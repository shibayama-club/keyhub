package sqlc

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	"github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
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

func (r *SqlcRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx repository.Transaction) error) (err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	committed := false
	defer func() {
		if committed {
			return
		}
		if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			if err != nil {
				err = errors.Wrapf(err, "failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	if txErr := fn(ctx, &SqlcTransaction{queries: r.queries.WithTx(tx)}); txErr != nil {
		err = errors.Wrap(txErr, "failed to Queries")
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to Commit")
	}
	committed = true

	return nil
}

func (r *SqlcRepository) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}
