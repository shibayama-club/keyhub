package model

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type AppSessionID string

func (id AppSessionID) String() string {
	return string(id)
}

func (id AppSessionID) Validate() error {
	if id == "" {
		return errors.WithHint(
			errors.New("session ID is required"),
			"セッションIDは必須です。",
		)
	}
	return nil
}

func NewAppSessionID(value string) (AppSessionID, error) {
	id := AppSessionID(value)
	if err := id.Validate(); err != nil {
		return "", err
	}
	return id, nil
}

type AppSession struct {
	SessionID          AppSessionID
	UserID             UserID
	ActiveMembershipID *uuid.UUID //ここは外部キーで現在このキーがあるテーブル自体が存在しないので一時的にプリミティブ型
	CreatedAt          time.Time
	ExpiresAt          time.Time
	CSRFToken          *string
	Revoked            bool
}

func (s AppSession) String() string {
	return string(s.SessionID)
}

func (s AppSession) Validate() error {
	if err := s.SessionID.Validate(); err != nil {
		return err
	}

	// UserID は uuid.UUID 型なので、ゼロ値かどうかをチェック
	if s.UserID == UserID(uuid.Nil) {
		return errors.WithHint(
			errors.New("user ID is required"),
			"ユーザーIDは必須です。",
		)
	}

	if s.CreatedAt.IsZero() {
		return errors.WithHint(
			errors.New("created_at is required"),
			"作成日時は必須です。",
		)
	}

	if s.ExpiresAt.IsZero() {
		return errors.WithHint(
			errors.New("expires_at is required"),
			"有効期限は必須です。",
		)
	}

	if s.ExpiresAt.Before(s.CreatedAt) {
		return errors.WithHint(
			errors.New("expires_at must be after created_at"),
			"有効期限は作成日時より後である必要があります。",
		)
	}

	return nil
}

// IsExpired はセッションが期限切れかどうかを確認する
func (s AppSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid はセッションが有効かどうか（無効化されておらず、期限切れでない）を確認する
func (s AppSession) IsValid() bool {
	return !s.Revoked && !s.IsExpired()
}

func NewAppSession(
	sessionID AppSessionID,
	userID UserID,
	expiresAt time.Time,
) (AppSession, error) {
	session := AppSession{
		SessionID: sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	if err := session.Validate(); err != nil {
		return AppSession{}, err
	}

	return session, nil
}
