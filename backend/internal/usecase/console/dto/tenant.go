package dto

import (
	"time"

	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateTenantInput struct {
	OrganizationID model.OrganizationID
	Name           string
	Description    string
	TenantType     string
	JoinCode       string
	JoinCodeExpiry *time.Time
	JoinCodeMaxUse int32
}

type UpdateTenantInput struct{
	ID model.TenantID
	Name string
	Description string
	TenantType string
}
