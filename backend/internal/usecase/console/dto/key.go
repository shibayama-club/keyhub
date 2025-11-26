package dto

import (
	"github.com/shibayama-club/keyhub/internal/domain/model"
)

type CreateKeyInput struct {
	RoomID         model.RoomID
	OrganizationID model.OrganizationID
	KeyNumber      string
}
