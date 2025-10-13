package model

import (
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type TenantID uuid.UUID

func (id TenantID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id TenantID) String() string {
	return uuid.UUID(id).String()
}

type TenantName string

func (n TenantName) String() string {
	return string(n)
}

func (n TenantName) Validate() error {
	if n == "" {
		return errors.WithHint(
			errors.New("tenant name is required"),
			"Tenant Nameは必須です。",
		)
	}

	if utf8.RuneCountInString(string(n)) > 30 {
		return errors.WithHint(
			errors.New("Please enter a tenantname within 30 characters"),
			"テナント名は30文字以内で入力してください。",
		)
	}
	return nil
}

func NewTenantName(value string) (TenantName, error) {
	n := TenantName(value)
	if err := n.Validate(); err != nil {
		return "", err
	}
	return n, nil
}

type TenantDescription string

func (d TenantDescription) String() string {
	return string(d)
}

func (d TenantDescription) Validate() error {
	if utf8.RuneCountInString(string(d)) > 500 {
		return errors.WithHint(
			errors.New("Please enter a tenant description within 500 characters"),
			"テナントの説明は500文字以内で入力してください。",
		)
	}
	return nil
}

func NewTenantDescription(value string) (TenantDescription, error) {
	d := TenantDescription(value)
	if err := d.Validate(); err != nil {
		return "", err
	}
	return d, nil
}

type TenantType string

const (
	TenantTypeDepartment TenantType = "department"
	TenantTypeLaboratory TenantType = "laboratory"
	TenantTypeProject    TenantType = "project"
	TenantTypeTeam       TenantType = "team"
)

func (t TenantType) String() string {
	return string(t)
}

func (t TenantType) Validate() error {
	switch t {
	case TenantTypeDepartment, TenantTypeLaboratory, TenantTypeProject, TenantTypeTeam:
		return nil
	default:
		return errors.WithHint(
			errors.Newf("invalid tenant type: %s", t),
			"テナントタイプは department, laboratory, division のいずれかである必要があります。",
		)
	}
}

func NewTenantType(value string) (TenantType, error) {
	t := TenantType(value)
	if err := t.Validate(); err != nil {
		return "", err
	}
	return t, nil
}

type Tenant struct {
	ID          TenantID
	Name        TenantName
	Description TenantDescription
	Type        TenantType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var isSlugFormatRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func IsSlugFormat(str string) bool {
	return isSlugFormatRegex.MatchString(str)
}
