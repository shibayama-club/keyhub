package dto

import "time"

type GetTenantByJoinCodeOutput struct {
	ID          string
	Name        string
	Description string
	TenantType  string
}

type TenantOutput struct {
	ID             string
	OrganizationID string
	Name           string
	Description    string
	TenantType     string
	MemberCount    int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type GetMyTenantsOutput struct {
	Tenants []TenantOutput
}
