package v1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/shibayama-club/keyhub/internal/domain"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	appv1 "github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/app/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) GetMe(
	ctx context.Context,
	req *connect.Request[appv1.GetMeRequest],
) (*connect.Response[appv1.GetMeResponse], error) {
	userID, ok := domain.Value[model.UserID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user not authenticated"))
	}

	user, err := h.useCase.GetUserByID(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err, "failed to get user"))
	}

	return connect.NewResponse(&appv1.GetMeResponse{
		User: &appv1.User{
			Id:        user.UserId.String(),
			Email:     user.Email.String(),
			Name:      user.Name.String(),
			Icon:      user.Icon.String(),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}), nil
}

func (h *Handler) Logout(
	ctx context.Context,
	req *connect.Request[appv1.LogoutRequest],
) (*connect.Response[appv1.LogoutResponse], error) {
	sessionID, ok := domain.Value[model.AppSessionID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("session not found"))
	}

	if err := h.useCase.Logout(ctx, sessionID.String()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err, "failed to logout"))
	}

	return connect.NewResponse(&appv1.LogoutResponse{
		Success: true,
	}), nil
}
