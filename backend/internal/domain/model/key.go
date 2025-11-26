package model

import (
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type KeyID uuid.UUID

func (id KeyID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id KeyID) String() string {
	return uuid.UUID(id).String()
}

func ParseKeyID(value string) (KeyID, error) {
	u, err := uuid.Parse(value)
	if err != nil {
		return KeyID{}, errors.WithHint(
			errors.Wrap(err, "failed to parse key ID"),
			"鍵IDの形式が正しくありません。",
		)
	}
	return KeyID(u), nil
}

type KeyNumber string

func (n KeyNumber) String() string {
	return string(n)
}

func (n KeyNumber) Validate() error {
	if n == "" {
		return errors.WithHint(
			errors.New("key number is required"),
			"鍵番号は必須です。",
		)
	}

	if utf8.RuneCountInString(string(n)) > 10 {
		return errors.WithHint(
			errors.New("key number must be within 10 characters"),
			"鍵番号は10文字以内で入力してください。",
		)
	}
	return nil
}

func NewKeyNumber(value string) (KeyNumber, error) {
	n := KeyNumber(value)
	if err := n.Validate(); err != nil {
		return "", err
	}
	return n, nil
}

type KeyStatus string

const (
	KeyStatusAvailable KeyStatus = "available"
	KeyStatusInUse     KeyStatus = "in_use"
	KeyStatusLost      KeyStatus = "lost"
	KeyStatusDamaged   KeyStatus = "damaged"
)

func (s KeyStatus) String() string {
	return string(s)
}

func (s KeyStatus) Validate() error {
	switch s {
	case KeyStatusAvailable, KeyStatusInUse, KeyStatusLost, KeyStatusDamaged:
		return nil
	default:
		return errors.WithHintf(
			errors.New("invalid key status"),
			"無効な鍵のステータスです: %s", s,
		)
	}
}

func NewKeyStatus(value string) (KeyStatus, error) {
	s := KeyStatus(value)
	if err := s.Validate(); err != nil {
		return "", err
	}
	return s, nil
}

type Key struct {
	ID             KeyID
	RoomID         RoomID
	OrganizationID OrganizationID
	KeyNumber      KeyNumber
	Status         KeyStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (k Key) Validate() error {
	if err := k.OrganizationID.Validate(); err != nil {
		return err
	}

	if err := k.KeyNumber.Validate(); err != nil {
		return err
	}

	if err := k.Status.Validate(); err != nil {
		return err
	}

	if k.CreatedAt.IsZero() {
		return errors.WithHint(
			errors.New("created_at is required"),
			"作成日時は必須です。",
		)
	}

	if k.UpdatedAt.IsZero() {
		return errors.WithHint(
			errors.New("updated_at is required"),
			"更新日時は必須です。",
		)
	}

	return nil
}

func NewKey(
	roomID RoomID,
	organizationID OrganizationID,
	keyNumber KeyNumber,
) (Key, error) {
	now := time.Now()
	key := Key{
		ID:             KeyID(uuid.New()),
		RoomID:         roomID,
		OrganizationID: organizationID,
		KeyNumber:      keyNumber,
		Status:         KeyStatusAvailable,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := key.Validate(); err != nil {
		return Key{}, err
	}

	return key, nil
}
