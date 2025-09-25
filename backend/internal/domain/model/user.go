package user

import (
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type UserID uuid.UUID

func (id UserID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id UserID) String() string {
	return uuid.UUID(id).String()
}

type UserName string

func (n UserName) String() string {
	return string(n)
}

func (n UserName) Validate() error {
	if n == "" {
		return errors.WithHint(
			errors.New("username is required"),
			"Nameは必須です。",
		)
	}
	if utf8.RuneCountInString(string(n)) > 10 {
		return errors.WithHint(
			errors.New("Please enter a username within 10 characters"),
			"ユーザーネームは10字以内で入力してください。",
		)
	}
	return nil
}

func NewUserName(value string) (UserName, error) {
	e := UserName(value)
	if err := e.Validate(); err != nil {
		return "", err
	}
	return UserName(value), nil
}

type User struct {
	Id UserID
	Name UserName
	CreatedAt time.Time
	UpdatedAt time.Time
}