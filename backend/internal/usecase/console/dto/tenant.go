package dto

import (
	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateTenantInput struct {
	OrganizationID model.OrganizationID
	Name           model.TenantName
	Description    model.TenantDescription
	TenantType     model.TenantType
	JoinCode       model.TenantJoinCode
	JoinCodeExpiry model.TenantJoinCodeExpiresAt
	JoinCodeMaxUse model.TenantJoinCodeMaxUses
}
