package v1

import (
	"context"

	"connectrpc.com/connect"
	appv1 "github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/app/v1"
)

func (h *Handler) GetTenantByJoinCode(ctx context.Context, req *connect.Request[appv1.GetTenantByJoinCodeRequest]) (*connect.Response[appv1.GetTenantByJoinCodeResponse], error) {
	output, err := h.useCase.GetTenantByJoinCode(ctx, req.Msg.JoinCode)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&appv1.GetTenantByJoinCodeResponse{
		Id:          output.ID,
		Name:        output.Name,
		Description: output.Description,
		TenantType:  convertStringToTenantTypeProto(output.TenantType),
	}), nil
}

func convertStringToTenantTypeProto(tenantType string) appv1.TenantType {
	switch tenantType {
	case "TENANT_TYPE_TEAM":
		return appv1.TenantType_TENANT_TYPE_TEAM
	case "TENANT_TYPE_DEPARTMENT":
		return appv1.TenantType_TENANT_TYPE_DEPARTMENT
	case "TENANT_TYPE_PROJECT":
		return appv1.TenantType_TENANT_TYPE_PROJECT
	case "TENANT_TYPE_LABORATORY":
		return appv1.TenantType_TENANT_TYPE_LABORATORY
	default:
		return appv1.TenantType_TENANT_TYPE_UNSPECIFIED
	}
}
