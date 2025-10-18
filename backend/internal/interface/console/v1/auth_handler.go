package v1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	consolev1 "github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/console/v1"
)

func (h *Handler) LoginWithOrgId(
	ctx context.Context,
	req *connect.Request[consolev1.LoginWithOrgIdRequest],
) (*connect.Response[consolev1.LoginWithOrgIdResponse], error) {
	token, expiresIn, err := h.useCase.LoginWithOrgId(ctx, req.Msg.OrganizationId, req.Msg.OrganizationKey)
	if err != nil {
		h.l.Error("failed to login with org id", "error", err)
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	return connect.NewResponse(&consolev1.LoginWithOrgIdResponse{
		SessionToken: token,
		ExpiresIn:    expiresIn,
	}), nil
}

func (h *Handler) Logout(
	ctx context.Context,
	req *connect.Request[consolev1.LogoutRequest],
) (*connect.Response[consolev1.LogoutResponse], error) {
	authHeader := req.Header().Get("Authorization")
	if authHeader == "" {
		h.l.Warn("missing authorization header")
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authorization header is required"))
	}

	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		h.l.Warn("invalid token", "error", err)
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.Wrap(err, "invalid token"))
	}

	if err := h.useCase.Logout(ctx, claims.Sid); err != nil {
		h.l.Error("failed to logout", "error", err)
		return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err, "failed to logout"))
	}

	return connect.NewResponse(&consolev1.LogoutResponse{
		Success: true,
	}), nil
}
