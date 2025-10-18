package model

import (
	"time"

	"github.com/cockroachdb/errors"
)

type ConsoleSessionID string

func (id ConsoleSessionID) String() string {
	return string(id)
}

func (id ConsoleSessionID) Validate() error {
	if id == "" {
		return errors.WithHint(
			errors.New("session ID is required"),
			"セッションIDは必須です。",
		)
	}
	return nil
}

func NewConsoleSessionID(value string) (ConsoleSessionID, error) {
	id := ConsoleSessionID(value)
	if err := id.Validate(); err != nil {
		return "", err
	}
	return id, nil
}

type ConsoleSession struct {
	SessionID      ConsoleSessionID
	OrganizationID OrganizationID
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

func (cs ConsoleSession) String() string {
	return string(cs.SessionID)
}

func (cs ConsoleSession) Validate() error {
	if err := cs.SessionID.Validate(); err != nil {
		return err
	}

	if err := cs.OrganizationID.Validate(); err != nil {
		return err
	}

	if cs.CreatedAt.IsZero() {
		return errors.WithHint(
			errors.New("created_at is required"),
			"作成日時は必須です。",
		)
	}

	if cs.ExpiresAt.IsZero() {
		return errors.WithHint(
			errors.New("expires_at is required"),
			"有効期限は必須です。",
		)
	}

	if cs.ExpiresAt.Before(cs.CreatedAt) {
		return errors.WithHint(
			errors.New("expires_at must be after created_at"),
			"有効期限は作成日時より後である必要があります。",
		)
	}

	return nil
}

func (cs ConsoleSession) IsExpired() bool {
	return time.Now().After(cs.ExpiresAt)
}

func NewConsoleSession(sessionID ConsoleSessionID, organizationID OrganizationID, expiresAt time.Time) (ConsoleSession, error) {
	session := ConsoleSession{
		SessionID:      sessionID,
		OrganizationID: organizationID,
		CreatedAt:      time.Now(),
		ExpiresAt:      expiresAt,
	}

	if err := session.Validate(); err != nil {
		return ConsoleSession{}, err
	}

	return session, nil
}
