package model

import (
	"regexp"
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

type UserEmail string

func (e UserEmail) String() string {
	return string(e)
}

func (e UserEmail) Validate() error {
	if e == "" {
		return errors.WithHint(
			errors.New("email is required"),
			"Emailは必須です。",
		)
	}
	
	if !IsEmailFormat(string(e)) {
		return errors.WithHint(
			errors.New("Please enter a validate email address"),
			"Emailの正しい形式で入力してください。",
		)
	}
	return nil
}

func NewUserEmail(value string) (UserEmail, error) {
	e := UserEmail(value)
	if err := e.Validate(); err != nil {
		return "", err
	}
	return UserEmail(value), nil
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
	if utf8.RuneCountInString(string(n)) > 30 {
		return errors.WithHint(
			errors.New("Please enter a username within 30 characters"),
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

type UserIcon string

func (i UserIcon) Validate() error {
	if i == "" {
		return nil
	}
	
	if !IsIconURLFormat(string(i)) {
		return errors.WithHint(
			errors.New("icon must be a valid HTTPS URL"),
			"アイコンURLは有効なHTTPS形式である必要があります。",
		)
	}

	if utf8.RuneCountInString(string(i)) > 2048 {
		return errors.WithHint(
			errors.New("icon URL is too long"),
			"アイコンURLが長すぎます。",
		)
	}

	return nil
}

func NewUserIcon(value string) (UserIcon, error) {
	i := UserIcon(value)
	if err := i.Validate(); err != nil {
		return "", err
	}
	return UserIcon(value), nil
}

type User struct {
	Id UserID
	Email UserEmail
	Name UserName
	Icon UserIcon
	CreatedAt time.Time
	UpdatedAt time.Time
}

var isEmailFormatRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func IsEmailFormat(str string) bool {
	return isEmailFormatRegex.MatchString(str)
}

var isIconURLFormatRegex = regexp.MustCompile(`^https://[^\s]+$`)

func IsIconURLFormat(str string) bool {
	return isIconURLFormatRegex.MatchString(str)
}
