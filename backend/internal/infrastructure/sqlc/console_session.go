package sqlc

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	sqlcgen "github.com/shibayama-club/keyhub/internal/infrastructure/sqlc/gen"
)

func parseSqlcConsoleSession(consoleSession sqlcgen.ConsoleSession) (model.ConsoleSession, error) {
	return model.ConsoleSession{
		SessionID:      model.ConsoleSessionID(consoleSession.SessionID),
		OrganizationID: model.OrganizationID(consoleSession.OrganizationID),
		CreatedAt:      consoleSession.CreatedAt.Time,
		ExpiresAt:      consoleSession.ExpiresAt.Time,
	}, nil
}

func (t *SqlcTransaction) CreateSession(ctx context.Context, arg repository.CreateConsoleSessionArg) (model.ConsoleSession, error) {
	params := sqlcgen.CreateConsoleSessionParams{
		SessionID:      arg.SessionID.String(),
		OrganizationID: arg.OrganizationID.UUID(),
	}
	sqlcConsoleSession, err := t.queries.CreateConsoleSession(ctx, params)
	if err != nil {
		return model.ConsoleSession{}, err
	}
	return parseSqlcConsoleSession(sqlcConsoleSession)
}

func (t *SqlcTransaction) GetSession(ctx context.Context, sessionID model.ConsoleSessionID) (model.ConsoleSession, error) {
	sqlcConsoleSession, err := t.queries.GetConsoleSession(ctx, sessionID.String())
	if err != nil {
		return model.ConsoleSession{}, err
	}
	return parseSqlcConsoleSession(sqlcConsoleSession)
}

func (t *SqlcTransaction) DeleteSession(ctx context.Context, sessionID model.ConsoleSessionID) error {
	return t.queries.DeleteConsoleSession(ctx, sessionID.String())
}
