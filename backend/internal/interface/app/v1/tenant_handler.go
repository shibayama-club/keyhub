package v1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/samber/lo"
	"github.com/shibayama-club/keyhub/internal/domain"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	appv1 "github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/app/v1"
	"github.com/shibayama-club/keyhub/internal/usecase/app/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (h *Handler) JoinTenant(ctx context.Context, req *connect.Request[appv1.JoinTenantRequest]) (*connect.Response[appv1.JoinTenantResponse], error) {
	userID, ok := domain.Value[model.UserID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user not authenticated"))
	}

	err := h.useCase.JoinTenant(ctx, userID, req.Msg.JoinCode)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&appv1.JoinTenantResponse{}), nil
}

func (h *Handler) GetMyTenants(ctx context.Context, req *connect.Request[appv1.GetMyTenantsRequest]) (*connect.Response[appv1.GetMyTenantsResponse], error) {
	userID, ok := domain.Value[model.UserID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user not authenticated"))
	}

	output, err := h.useCase.GetMyTenants(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	tenants := lo.Map(output.Tenants, func(t dto.TenantOutput, _ int) *appv1.Tenant {
		return convertTenantOutputToProto(t)
	})

	return connect.NewResponse(&appv1.GetMyTenantsResponse{
		Tenants: tenants,
	}), nil
}

func convertTenantOutputToProto(tenant dto.TenantOutput) *appv1.Tenant {
	return &appv1.Tenant{
		Id:             tenant.ID,
		OrganizationId: tenant.OrganizationID,
		Name:           tenant.Name,
		Description:    tenant.Description,
		TenantType:     convertStringToTenantTypeProto(tenant.TenantType),
		MemberCount:    tenant.MemberCount,
		CreatedAt:      timestamppb.New(tenant.CreatedAt),
		UpdatedAt:      timestamppb.New(tenant.UpdatedAt),
	}
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
