package sqlc

import (
	"context"
	"time"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcOAuthState(state sqlcgen.OauthState) (model.OAuthState, error) {
	var consumedAt *time.Time
	if state.ConsumedAt.Valid {
		t := state.ConsumedAt.Time
		consumedAt = &t
	}

	return model.OAuthState{
		State:        model.OAuthStateValue(state.State),
		CodeVerifier: state.CodeVerifier,
		Nonce:        state.Nonce,
		CreatedAt:    state.CreatedAt.Time,
		ConsumedAt:   consumedAt,
	}, nil
}

func (t *SqlcTransaction) SaveOAuthState(ctx context.Context, oauthState model.OAuthState) error {
	return t.queries.SaveOAuthState(ctx, sqlcgen.SaveOAuthStateParams{
		State:        oauthState.State.String(),
		CodeVerifier: oauthState.CodeVerifier,
		Nonce:        oauthState.Nonce,
	})
}

func (t *SqlcTransaction) GetOAuthState(ctx context.Context, state string) (model.OAuthState, error) {
	sqlcStateRow, err := t.queries.GetOAuthState(ctx, state)
	if err != nil {
		return model.OAuthState{}, err
	}

	return parseSqlcOAuthState(sqlcStateRow.OauthState)
}

func (t *SqlcTransaction) ConsumeOAuthState(ctx context.Context, state string) error {
	return t.queries.ConsumeOAuthState(ctx, state)
}
