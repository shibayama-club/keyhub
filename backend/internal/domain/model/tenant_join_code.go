package model

import (
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type TenantJoinCodeID uuid.UUID

func (id TenantJoinCodeID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id TenantJoinCodeID) String() string {
	return uuid.UUID(id).String()
}

type TenantJoinCode string

func (c TenantJoinCode) String() string {
	return string(c)
}

func (c TenantJoinCode) Validate() error {
	if c == "" {
		return errors.WithHint(
			errors.New("tenant join code is required"),
			"テナント参加コードは必須です。",
		)
	}

	length := utf8.RuneCountInString(string(c))
	if length < 6 || length > 20 {
		return errors.WithHint(
			errors.New("tenant join code must be between 6 and 20 characters"),
			"テナント参加コードは6文字以上20文字以下で入力してください。",
		)
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, string(c))
	if err != nil {
		return errors.Wrap(err, "failed to validate tenant join code format")
	}
	if !matched {
		return errors.WithHint(
			errors.New("tenant join code must contain only alphanumeric characters"),
			"テナント参加コードは英数字のみで入力してください。",
		)
	}

	return nil
}

func NewTenantJoinCode(value string) (TenantJoinCode, error) {
	c := TenantJoinCode(value)
	if err := c.Validate(); err != nil {
		return "", err
	}
	return c, nil
}

type TenantJoinCodeMaxUses int32

func (m TenantJoinCodeMaxUses) Int32() int32 {
	return int32(m)
}

func (m TenantJoinCodeMaxUses) Validate() error {
	if m < 0 {
		return errors.WithHint(
			errors.New("max uses cannot be negative"),
			"最大使用回数は0以上である必要があります。",
		)
	}
	return nil
}

func NewTenantJoinCodeMaxUses(value int32) (TenantJoinCodeMaxUses, error) {
	m := TenantJoinCodeMaxUses(value)
	if err := m.Validate(); err != nil {
		return 0, err
	}
	return m, nil
}

type TenantJoinCodeExpiresAt *time.Time

func ValidateTenantJoinCodeExpiresAt(expiresAt TenantJoinCodeExpiresAt) error {
	if expiresAt != nil && (*expiresAt).Before(time.Now()) {
		return errors.WithHint(
			errors.New("expiration time must be in the future"),
			"有効期限は未来の日時である必要があります。",
		)
	}
	return nil
}

func NewTenantJoinCodeExpiresAt(value *time.Time) (TenantJoinCodeExpiresAt, error) {
	expiresAt := TenantJoinCodeExpiresAt(value)
	if err := ValidateTenantJoinCodeExpiresAt(expiresAt); err != nil {
		return nil, err
	}
	return expiresAt, nil
}

type TenantJoinCodeEntity struct {
	ID        TenantJoinCodeID
	TenantID  TenantID
	Code      TenantJoinCode
	ExpiresAt TenantJoinCodeExpiresAt
	MaxUses   TenantJoinCodeMaxUses
	UsedCount int
	CreatedAt time.Time
}

func (t TenantJoinCodeEntity) Validate() error {
	if err := t.Code.Validate(); err != nil {
		return err
	}

	if err := t.MaxUses.Validate(); err != nil {
		return err
	}

	if err := ValidateTenantJoinCodeExpiresAt(t.ExpiresAt); err != nil {
		return err
	}

	if t.UsedCount < 0 {
		return errors.WithHint(
			errors.New("used count cannot be negative"),
			"使用回数は0以上である必要があります。",
		)
	}

	if t.UsedCount > int(t.MaxUses.Int32()) && t.MaxUses.Int32() > 0 {
		return errors.WithHint(
			errors.New("used count cannot exceed max uses"),
			"使用回数が最大使用回数を超えています。",
		)
	}

	if t.CreatedAt.IsZero() {
		return errors.WithHint(
			errors.New("created_at is required"),
			"作成日時は必須です。",
		)
	}

	return nil
}

func (t TenantJoinCodeEntity) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*t.ExpiresAt)
}

func (t TenantJoinCodeEntity) IsUsable() bool {
	if t.IsExpired() {
		return false
	}
	if t.MaxUses.Int32() > 0 && t.UsedCount >= int(t.MaxUses.Int32()) {
		return false
	}
	return true
}

func NewTenantJoinCodeEntity(
	tenantID TenantID,
	code TenantJoinCode,
	expiresAt TenantJoinCodeExpiresAt,
	maxUses TenantJoinCodeMaxUses,
) (TenantJoinCodeEntity, error) {
	entity := TenantJoinCodeEntity{
		ID:        TenantJoinCodeID(uuid.New()),
		TenantID:  tenantID,
		Code:      code,
		ExpiresAt: expiresAt,
		MaxUses:   maxUses,
		UsedCount: 0,
		CreatedAt: time.Now(),
	}

	if err := entity.Validate(); err != nil {
		return TenantJoinCodeEntity{}, err
	}

	return entity, nil
}
