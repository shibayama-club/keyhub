package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/cockroachdb/errors"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	"github.com/shibayama-club/keyhub/internal/infrastructure/auth/google"
)

func (u *UseCase) StartGoogleLogin(ctx context.Context) (authURL string, err error) {
	codeVerifier, err := google.GenerateCodeVerifier()
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to generate code verifier")
	}

	codeChallenge := google.GenerateCodeChallenge(codeVerifier)

	stateStr, err := google.GenerateRandomString(32)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to generate state")
	}

	stateValue, err := model.NewOAuthStateValue(stateStr)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid state value")
	}

	nonce, err := google.GenerateRandomString(32)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to generate nonce")
	}

	oauthState, err := model.NewOAuthState(stateValue, codeVerifier, nonce)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid OAuth state")
	}

	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		return tx.SaveOAuthState(ctx, oauthState)
	})
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to save OAuth state")
	}

	authURL = u.oauthService.BuildAuthURL(stateStr, nonce, codeChallenge)

	return authURL, nil
}

func (u *UseCase) GoogleCallback(ctx context.Context, code, state string) (sessionID string, err error) {
	oauthState, err := u.repo.GetOAuthState(ctx, state)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrUnAuthorized), "invalid or expired state")
	}

	if !oauthState.IsValid() {
		return "", errors.WithHint(
			errors.Mark(errors.New("OAuth state is invalid"), domainerrors.ErrUnAuthorized),
			"認証フローが無効です。最初からやり直してください。",
		)
	}

	if err := u.repo.ConsumeOAuthState(ctx, state); err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to consume OAuth state")
	}

	tokens, err := u.oauthService.ExchangeCode(ctx, code, oauthState.CodeVerifier)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to exchange code for tokens")
	}

	claims, err := u.oauthService.VerifyIDToken(ctx, tokens.IDToken, oauthState.Nonce)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrUnAuthorized), "failed to verify ID token")
	}

	email, name, picture, providerSub := claims.GetUserInfo()

	var userID model.UserID
	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		existingUser, err := tx.GetUserByProviderIdentity(ctx, "google", providerSub)
		if err == nil {
			userID = existingUser.UserId
			return nil
		}

		userEmail, err := model.NewUserEmail(email)
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid email")
		}

		userName, err := model.NewUserName(name)
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid name")
		}

		userIcon, err := model.NewUserIcon(picture)
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid icon URL")
		}

		user, err := tx.UpsertUser(ctx, repository.UpsertUserArg{
			Email: userEmail,
			Name:  userName,
			Icon:  userIcon,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create user")
		}

		userID = user.UserId

		if err := tx.UpsertUserIdentity(ctx, repository.UpsertUserIdentityArg{
			UserID:      userID,
			Provider:    "google",
			ProviderSub: providerSub,
		}); err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create user identity")
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	sessionBytes := make([]byte, 32)
	if _, err := rand.Read(sessionBytes); err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to generate session ID")
	}
	sessionIDStr := "app_sess_" + hex.EncodeToString(sessionBytes)

	appSessionID, err := model.NewAppSessionID(sessionIDStr)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create session ID")
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		err := tx.CreateAppSession(ctx, repository.CreateAppSessionArg{
			SessionID: appSessionID,
			UserID:    userID,
			ExpiresAt: expiresAt,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create session")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return sessionIDStr, nil
}

func (u *UseCase) GetMe(ctx context.Context, sessionID string) (model.User, error) {
	appSessionID, err := model.NewAppSessionID(sessionID)
	if err != nil {
		return model.User{}, errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid session ID")
	}

	session, err := u.repo.GetAppSession(ctx, appSessionID)
	if err != nil {
		return model.User{}, errors.Wrap(errors.Mark(err, domainerrors.ErrNotFound), "session not found")
	}

	if !session.IsValid() {
		return model.User{}, errors.WithHint(
			errors.Mark(errors.New("session is invalid or expired"), domainerrors.ErrUnAuthorized),
			"セッションが無効または期限切れです。再度ログインしてください。",
		)
	}

	user, err := u.repo.GetUser(ctx, session.UserID)
	if err != nil {
		return model.User{}, errors.Wrap(errors.Mark(err, domainerrors.ErrNotFound), "user not found")
	}

	return user, nil
}

func (u *UseCase) GetUserByID(ctx context.Context, userID model.UserID) (model.User, error) {
	user, err := u.repo.GetUser(ctx, userID)
	if err != nil {
		return model.User{}, errors.Wrap(errors.Mark(err, domainerrors.ErrNotFound), "user not found")
	}

	return user, nil
}

func (u *UseCase) Logout(ctx context.Context, sessionID string) error {
	appSessionID, err := model.NewAppSessionID(sessionID)
	if err != nil {
		return errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid session ID")
	}

	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		return tx.RevokeAppSession(ctx, appSessionID)
	})
	if err != nil {
		return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to revoke session")
	}

	return nil
}
