package v1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/interface/console/v1/interceptor"
	consolev1 "github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/console/v1"
	"github.com/shibayama-club/keyhub/internal/usecase/console/dto"
)

func convertTenantType(protoType consolev1.TenantType) (string, error) {
	switch protoType {
	case consolev1.TenantType_TENANT_TYPE_TEAM:
		return model.TenantTypeTeam.String(), nil
	case consolev1.TenantType_TENANT_TYPE_DEPARTMENT:
		return model.TenantTypeDepartment.String(), nil
	case consolev1.TenantType_TENANT_TYPE_PROJECT:
		return model.TenantTypeProject.String(), nil
	case consolev1.TenantType_TENANT_TYPE_LABORATORY:
		return model.TenantTypeLaboratory.String(), nil
	default:
		return "", errors.New("invalid tenant type")
	}
}

func (h *Handler) CreateTenant(
	ctx context.Context,
	req *connect.Request[consolev1.CreateTenantRequest],
) (*connect.Response[consolev1.CreateTenantResponse], error) {
	session, ok := ctx.Value(interceptor.ConsoleSessionKey).(model.ConsoleSession)
	if !ok {
		h.l.Warn("session not found in context")
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("session not found"))
	}

	tenantTypeStr, err := convertTenantType(req.Msg.TenantType)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	output, err := h.useCase.CreateTenant(ctx, dto.CreateTenantInput{
		OrganizationID: session.OrganizationID,
		Name:           req.Msg.Name,
		Description:    req.Msg.Description,
		TenantType:     tenantTypeStr,
	})
	if err != nil {
		h.l.Error("failed to create tenant", "error", err)
		return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err, "failed to create tenant"))
	}

	return connect.NewResponse(&consolev1.CreateTenantResponse{
		Id: output.ID,
	}), nil
}
