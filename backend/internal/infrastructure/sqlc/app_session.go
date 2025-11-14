package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcSession(session sqlcgen.Session) (model.AppSession, error) {
	return model.AppSession{
		SessionID:          model.AppSessionID(session.SessionID),
		UserID:             model.UserID(session.UserID),
		ActiveMembershipID: session.ActiveMembershipID,
		CreatedAt:          session.CreatedAt.Time,
		ExpiresAt:          session.ExpiresAt.Time,
		CSRFToken:          session.CsrfToken,
		Revoked:            session.Revoked,
	}, nil
}

func (t *SqlcTransaction) CreateAppSession(ctx context.Context, arg repository.CreateAppSessionArg) (model.AppSession, error) {
	sqlcSessionRow, err := t.queries.CreateAppSession(ctx, sqlcgen.CreateAppSessionParams{
		SessionID:          arg.SessionID.String(),
		UserID:             arg.UserID.UUID(),
		ActiveMembershipID: nil, // Optional field
		ExpiresAt: pgtype.Timestamptz{
			Time:  arg.ExpiresAt,
			Valid: true,
		},
		CsrfToken: nil, // Optional field
	})
	if err != nil {
		return model.AppSession{}, err
	}
	return parseSqlcSession(sqlcSessionRow.Session)
}

func (t *SqlcTransaction) GetAppSession(ctx context.Context, sessionID model.AppSessionID) (model.AppSession, error) {
	sqlcSessionRow, err := t.queries.GetAppSession(ctx, sessionID.String())
	if err != nil {
		return model.AppSession{}, err
	}

	return parseSqlcSession(sqlcSessionRow.Session)
}

func (t *SqlcTransaction) RevokeAppSession(ctx context.Context, sessionID model.AppSessionID) error {
	return t.queries.RevokeAppSession(ctx, sessionID.String())
}
