package model

import (
	"unicode/utf8"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type OrganizationID uuid.UUID

func (id OrganizationID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id OrganizationID) String() string {
	return uuid.UUID(id).String()
}

type OrganizationKey string

func (k OrganizationKey) String() string {
	return string(k)
}

func (k OrganizationKey) Validate() error {
	if k == "" {
		return errors.WithHint(
			errors.New("organization key is required"),
			"Organization Keyは必須です。",
		)
	}

	length := utf8.RuneCountInString(string(k))
	if length < 1 || length > 20 {
		return errors.WithHint(
			errors.New("Please enter an organization key between 1 and 200 characters"),
			"Organization Keyは1文字以上20文字以内で入力してください。",
		)
	}

	return nil
}

func NewOrganizationKey(value string) (OrganizationKey, error) {
	k := OrganizationKey(value)
	if err := k.Validate(); err != nil {
		return "", err
	}
	return k, nil
}

type Organization struct {
	ID  OrganizationID
	Key OrganizationKey
}

