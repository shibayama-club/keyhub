package model

import (
	"time"

	"github.com/google/uuid"
)

type TenantMembershipID uuid.UUID

func (id TenantMembershipID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id TenantMembershipID) String() string {
	return uuid.UUID(id).String()
}

type TenantMembershipRole string

const (
	TenantMembershipRoleAdmin  TenantMembershipRole = "admin"
	TenantMembershipRoleMember TenantMembershipRole = "member"
)

func (r TenantMembershipRole) String() string {
	return string(r)
}

type TenantMembership struct {
	ID        TenantMembershipID
	TenantID  TenantID
	UserID    UserID
	Role      TenantMembershipRole
	CreatedAt time.Time
	LeftAt    *time.Time
}
