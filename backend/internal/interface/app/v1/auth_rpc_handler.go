package v1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	appv1 "github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/app/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetMe retrieves the current user information
func (h *Handler) GetMe(
	ctx context.Context,
	req *connect.Request[appv1.GetMeRequest],
) (*connect.Response[appv1.GetMeResponse], error) {
	// Get session ID from cookie
	sessionID := req.Header().Get("Cookie")
	if sessionID == "" {
		h.l.Warn("missing session cookie")
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("session cookie is required"))
	}

	// Extract session ID from cookie
	// TODO: Properly parse cookie
	// For now, this is a simplified version

	user, err := h.useCase.GetMe(ctx, sessionID)
	if err != nil {
		h.l.Error("failed to get user", "error", err)
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.Wrap(err, "failed to get user"))
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

// Logout logs out the user
func (h *Handler) Logout(
	ctx context.Context,
	req *connect.Request[appv1.LogoutRequest],
) (*connect.Response[appv1.LogoutResponse], error) {
	// Get session ID from cookie
	sessionID := req.Header().Get("Cookie")
	if sessionID == "" {
		h.l.Warn("missing session cookie")
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("session cookie is required"))
	}

	// Extract session ID from cookie
	// TODO: Properly parse cookie

	if err := h.useCase.Logout(ctx, sessionID); err != nil {
		h.l.Error("failed to logout", "error", err)
		return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err, "failed to logout"))
	}

	return connect.NewResponse(&appv1.LogoutResponse{
		Success: true,
	}), nil
}
