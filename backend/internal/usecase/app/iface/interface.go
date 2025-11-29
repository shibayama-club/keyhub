package iface

import (
	"context"

	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/usecase/app/dto"
)

type IUseCase interface {
	StartGoogleLogin(ctx context.Context) (authURL string, err error)
	GoogleCallback(ctx context.Context, code, state string) (sessionID string, err error)
	GetMe(ctx context.Context, sessionID string) (model.User, error)
	GetUserByID(ctx context.Context, userID model.UserID) (model.User, error)
	Logout(ctx context.Context, sessionID string) error
	GetTenantByJoinCode(ctx context.Context, joinCode string) (dto.GetTenantByJoinCodeOutput, error)
	JoinTenant(ctx context.Context, userID model.UserID, joinCode string) error
	GetMyTenants(ctx context.Context, userID model.UserID) (dto.GetMyTenantsOutput, error)
	GetRoomsByTenant(ctx context.Context, tenantID model.TenantID) ([]model.Room, error)
	GetKeysByRoom(ctx context.Context, roomID model.RoomID) ([]model.Key, error)
}
