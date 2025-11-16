package model

import (
	"time"

	"github.com/google/uuid"
)

type TenantMembershipID string

func (id TenantMembershipID) UUID() uuid.UUID {
	return uuid.MustParse(string(id))
}

func (id TenantMembershipID) String() string {
	return string(id)
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
