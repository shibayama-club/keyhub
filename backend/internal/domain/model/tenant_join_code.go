package model

import (
	"time"

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

type TenantJoinCodeMaxUses int

func (m TenantJoinCodeMaxUses) Int() int {
	return int(m)
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

func NewTenantJoinCodeMaxUses(value int) (TenantJoinCodeMaxUses, error) {
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
	ID         TenantJoinCodeID
	TenantID   TenantID
	Code       TenantJoinCode
	ExpiresAt  TenantJoinCodeExpiresAt
	MaxUses    TenantJoinCodeMaxUses
	UsedCount  int
	CreatedAt  time.Time
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

	if t.UsedCount > t.MaxUses.Int() && t.MaxUses.Int() > 0 {
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
	if t.MaxUses.Int() > 0 && t.UsedCount >= t.MaxUses.Int() {
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
