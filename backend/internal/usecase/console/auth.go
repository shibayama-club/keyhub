package console

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
)

func (u *UseCase) LoginWithOrgId(ctx context.Context, orgID, orgKey string) (string, int64, error) {
	expectedOrgID := u.config.Console.OrganizationId
	expectedOrgKey := u.config.Console.OrganizationKey

	if expectedOrgID == "" {
		expectedOrgID = DEFAULT_ORGANIZATION_ID
	}
	if expectedOrgKey == "" {
		expectedOrgKey = DEFAULT_ORGANIZATION_KEY
	}

	if orgID != expectedOrgID || orgKey != expectedOrgKey {
		return "", 0, errors.WithHint(
			errors.New("invalid organization credentials"),
			"組織IDまたはキーが正しくありません。",
		)
	}

	sessionBytes := make([]byte, 32)
	if _, err := rand.Read(sessionBytes); err != nil {
		return "", 0, errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to generate session ID")
	}
	sessionIDStr := "console_sess_" + hex.EncodeToString(sessionBytes)

	sessionID, err := model.NewConsoleSessionID(sessionIDStr)
	if err != nil {
		return "", 0, errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create session ID")
	}

	orgUUID := uuid.MustParse(orgID)
	organizationID, err := model.NewOrganizationID(orgUUID)
	if err != nil {
		return "", 0, errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create organization ID")
	}

	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		_, err := tx.CreateSession(ctx, repository.CreateConsoleSessionArg{
			SessionID:      sessionID,
			OrganizationID: organizationID,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create console session")
		}
		return nil
	})
	if err != nil {
		return "", 0, err
	}

	// JWTトークンの有効期限
	expiresIn := 24 * time.Hour
	token, err := u.authService.GenerateToken(orgID, sessionIDStr, expiresIn)
	if err != nil {
		return "", 0, errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to generate JWT token")
	}

	return token, int64(expiresIn.Seconds()), nil
}

func (u *UseCase) Logout(ctx context.Context, sessionID string) error {
	sid := model.ConsoleSessionID(sessionID)

	err := u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		err := tx.DeleteSession(ctx, sid)
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to delete session")
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) ValidateSession(ctx context.Context, token string) (model.ConsoleSession, error) {
	claims, err := u.authService.ValidateToken(token)
	if err != nil {
		return model.ConsoleSession{}, errors.Wrap(errors.Mark(err, domainerrors.ErrUnAuthorized), "failed to validate token")
	}

	sid, err := model.NewConsoleSessionID(claims.Sid)
	if err != nil {
		return model.ConsoleSession{}, errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid session ID in token")
	}

	session, err := u.repo.GetSession(ctx, sid)
	if err != nil {
		return model.ConsoleSession{}, errors.Wrap(errors.Mark(err, domainerrors.ErrNotFound), "failed to get session from database")
	}

	if session.OrganizationID.String() != claims.Org {
		return model.ConsoleSession{}, errors.WithHint(
			errors.New("organization mismatch"),
			"トークンの組織IDがセッションと一致しません。",
		)
	}

	return session, nil
}
