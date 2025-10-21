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
	TenantTypeUnspecified TenantType = "TENANT_TYPE_UNSPECIFIED"
	TenantTypeTeam        TenantType = "TENANT_TYPE_TEAM"
	TenantTypeDepartment  TenantType = "TENANT_TYPE_DEPARTMENT"
	TenantTypeProject     TenantType = "TENANT_TYPE_PROJECT"
	TenantTypeLaboratory  TenantType = "TENANT_TYPE_LABORATORY"
)

func (t TenantType) String() string {
	return string(t)
}

func (t TenantType) Validate() error {
	switch t {
	case TenantTypeTeam, TenantTypeDepartment, TenantTypeProject, TenantTypeLaboratory:
		return nil
	case TenantTypeUnspecified:
		return errors.WithHint(
			errors.New("tenant type must be specified"),
			"テナントタイプを指定してください。",
		)
	default:
		return errors.WithHintf(
			errors.New("invalid tenant type"),
			"無効なテナントタイプです: %s", t,
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
	ID             TenantID
	OrganizationID OrganizationID
	Name           TenantName
	Description    TenantDescription
	Type           TenantType
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (t Tenant) Validate() error {
	if err := t.OrganizationID.Validate(); err != nil {
		return err
	}

	if err := t.Name.Validate(); err != nil {
		return err
	}

	if err := t.Description.Validate(); err != nil {
		return err
	}

	if err := t.Type.Validate(); err != nil {
		return err
	}

	if t.CreatedAt.IsZero() {
		return errors.WithHint(
			errors.New("created_at is required"),
			"作成日時は必須です。",
		)
	}

	if t.UpdatedAt.IsZero() {
		return errors.WithHint(
			errors.New("updated_at is required"),
			"更新日時は必須です。",
		)
	}

	return nil
}

func NewTenant(
	organizationID OrganizationID,
	name TenantName,
	description TenantDescription,
	tenantType TenantType,
) (Tenant, error) {
	now := time.Now()
	tenant := Tenant{
		ID:             TenantID(uuid.New()),
		OrganizationID: organizationID,
		Name:           name,
		Description:    description,
		Type:           tenantType,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := tenant.Validate(); err != nil {
		return Tenant{}, err
	}

	return tenant, nil
}

var isSlugFormatRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func IsSlugFormat(str string) bool {
	return isSlugFormatRegex.MatchString(str)
}
