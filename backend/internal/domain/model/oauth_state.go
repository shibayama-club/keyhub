package model

import (
	"time"

	"github.com/cockroachdb/errors"
)

type OAuthStateValue string

func (s OAuthStateValue) String() string {
	return string(s)
}

func (s OAuthStateValue) Validate() error {
	if s == "" {
		return errors.WithHint(
			errors.New("OAuth state is required"),
			"OAuth状態は必須です。",
		)
	}
	return nil
}

func NewOAuthStateValue(value string) (OAuthStateValue, error) {
	state := OAuthStateValue(value)
	if err := state.Validate(); err != nil {
		return "", err
	}
	return state, nil
}

type OAuthState struct {
	State        OAuthStateValue
	CodeVerifier string
	Nonce        string
	CreatedAt    time.Time
	ConsumedAt   *time.Time
}

func (s OAuthState) Validate() error {
	if err := s.State.Validate(); err != nil {
		return err
	}

	if s.CodeVerifier == "" {
		return errors.WithHint(
			errors.New("code verifier is required"),
			"コード検証子は必須です。",
		)
	}

	if s.Nonce == "" {
		return errors.WithHint(
			errors.New("nonce is required"),
			"nonceは必須です。",
		)
	}

	if s.CreatedAt.IsZero() {
		return errors.WithHint(
			errors.New("created_at is required"),
			"作成日時は必須です。",
		)
	}

	return nil
}

// IsConsumed は状態が使用済みかどうかを確認する
func (s OAuthState) IsConsumed() bool {
	return s.ConsumedAt != nil
}

// IsExpired は状態が期限切れ（15分）かどうかを確認する
func (s OAuthState) IsExpired() bool {
	return time.Now().After(s.CreatedAt.Add(15 * time.Minute))
}

// IsValid は状態が有効（使用済みでなく、期限切れでない）かどうかを確認する
func (s OAuthState) IsValid() bool {
	return !s.IsConsumed() && !s.IsExpired()
}

// MarkAsConsumed はOAuth状態を使用済みとしてマークする
func (s *OAuthState) MarkAsConsumed() error {
	if s.IsConsumed() {
		return errors.WithHint(
			errors.New("OAuth state already consumed"),
			"OAuth状態は既に使用されています。",
		)
	}

	if s.IsExpired() {
		return errors.WithHint(
			errors.New("OAuth state has expired"),
			"OAuth状態の有効期限が切れています。",
		)
	}

	now := time.Now()
	s.ConsumedAt = &now
	return nil
}

func NewOAuthState(
	state OAuthStateValue,
	codeVerifier string,
	nonce string,
) (OAuthState, error) {
	oauthState := OAuthState{
		State:        state,
		CodeVerifier: codeVerifier,
		Nonce:        nonce,
		CreatedAt:    time.Now(),
		ConsumedAt:   nil,
	}

	if err := oauthState.Validate(); err != nil {
		return OAuthState{}, err
	}

	return oauthState, nil
}
