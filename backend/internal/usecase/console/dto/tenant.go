package dto

import "github.com/shibayama-club/keyhub/internal/domain/model"

type CreateTenantInput struct {
	OrganizationID model.OrganizationID
	Name           string
	Description    string
	TenantType     string
}