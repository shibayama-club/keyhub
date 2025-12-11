package iface

//go:generate go run go.uber.org/mock/mockgen@latest -source=$GOFILE -destination=../mock/mock_usecase.go -package=mock

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/usecase/console/dto"
)

type IUseCase interface {
	LoginWithOrgId(ctx context.Context, orgID, orgKey string) (string, int64, error)
	Logout(ctx context.Context, sessionID string) error
	ValidateSession(ctx context.Context, token string) (model.ConsoleSession, error)
	CreateTenant(ctx context.Context, input dto.CreateTenantInput) (string, error)
	GetAllTenants(ctx context.Context) ([]model.Tenant, error)
	GetTenantById(ctx context.Context, tenantId model.TenantID) (dto.GetTenantByIdOutput, error)
	UpdateTenant(ctx context.Context, input dto.UpdateTenantInput) error
	CreateRoom(ctx context.Context, input dto.CreateRoomInput) (string, error)
	GetAllRooms(ctx context.Context) ([]model.Room, error)
	AssignRoomToTenant(ctx context.Context, input dto.AssignRoomToTenantInput) (string, error)
	CreateKey(ctx context.Context, input dto.CreateKeyInput) (string, error)
	GetKeysByRoom(ctx context.Context, roomID model.RoomID) ([]model.Key, error)
	DeleteTenantById(ctx context.Context, tenantId model.TenantID) error
}
