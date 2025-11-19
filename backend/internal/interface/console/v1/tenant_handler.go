package v1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/samber/lo"
	"github.com/shibayama-club/keyhub/internal/domain"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	consolev1 "github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/console/v1"
	"github.com/shibayama-club/keyhub/internal/usecase/console/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// protobuf特有の型からmodelの型に変更
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
	orgID, ok := domain.Value[model.OrganizationID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.WithMessage(domainerrors.ErrNotFound, "organization not found"))
	}

	tenantTypeStr, err := convertTenantType(req.Msg.TenantType)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	input := dto.CreateTenantInput{
		OrganizationID: orgID,
		Name:           req.Msg.Name,
		Description:    req.Msg.Description,
		TenantType:     tenantTypeStr,
		JoinCode:       req.Msg.JoinCode,
		JoinCodeMaxUse: req.Msg.JoinCodeMaxUse,
	}

	if req.Msg.JoinCodeExpiry != nil {
		expiryTime := req.Msg.JoinCodeExpiry.AsTime()
		input.JoinCodeExpiry = &expiryTime
	}

	tenantID, err := h.useCase.CreateTenant(ctx, input)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&consolev1.CreateTenantResponse{
		Id: tenantID,
	}), nil
}

// modelの型からprotobuf特有の型に変更(convertTenantType関数の逆)
func convertModelTenantTypeToProto(modelType model.TenantType) consolev1.TenantType {
	switch modelType.String() {
	case model.TenantTypeTeam.String():
		return consolev1.TenantType_TENANT_TYPE_TEAM
	case model.TenantTypeDepartment.String():
		return consolev1.TenantType_TENANT_TYPE_DEPARTMENT
	case model.TenantTypeProject.String():
		return consolev1.TenantType_TENANT_TYPE_PROJECT
	case model.TenantTypeLaboratory.String():
		return consolev1.TenantType_TENANT_TYPE_LABORATORY
	default:
		return consolev1.TenantType_TENANT_TYPE_UNSPECIFIED
	}
}

func convertModelTenantToProto(tenant model.Tenant) *consolev1.Tenant {
	return &consolev1.Tenant{
		Id:             tenant.ID.String(),
		OrganizationId: tenant.OrganizationID.String(),
		Name:           tenant.Name.String(),
		Description:    tenant.Description.String(),
		TenantType:     convertModelTenantTypeToProto(tenant.Type),
		CreatedAt:      timestamppb.New(tenant.CreatedAt),
		UpdatedAt:      timestamppb.New(tenant.UpdatedAt),
	}
}

func (h *Handler) GetAllTenants(
	ctx context.Context,
	req *connect.Request[consolev1.GetAllTenantsRequest],
) (*connect.Response[consolev1.GetAllTenantsResponse], error) {
	orgID, ok := domain.Value[model.OrganizationID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.WithMessage(domainerrors.ErrNotFound, "organization not found"))
	}

	tenants, err := h.useCase.GetAllTenants(ctx, orgID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoTenants := lo.Map(tenants, func(tenant model.Tenant, _ int) *consolev1.Tenant {
		return convertModelTenantToProto(tenant)
	})

	return connect.NewResponse(&consolev1.GetAllTenantsResponse{
		Tenants: protoTenants,
	}), nil
}

func (h *Handler) GetTenantById(
	ctx context.Context,
	req *connect.Request[consolev1.GetTenantByIdRequest],
) (*connect.Response[consolev1.GetTenantByIdResponse], error) {
	// QUESTION(sirasu):ParseTenantIDでいいか?
	tenantId, err := model.ParseTenantID(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	tenant, err := h.useCase.GetTenantById(ctx, tenantId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	orgID, ok := domain.Value[model.OrganizationID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.WithMessage(domainerrors.ErrNotFound, "organization not found"))
	}

	if tenant.OrganizationID != orgID {
		return nil, connect.NewError(connect.CodeNotFound, errors.WithMessage(domainerrors.ErrNotFound, "tenant not found"))
	}

	protoTenant := convertModelTenantToProto(tenant)

	return connect.NewResponse(&consolev1.GetTenantByIdResponse{
		Tenant: protoTenant,
	}), nil
}

func (h *Handler) UpdaTenant(
	ctx context.Context,
	req *connect.Request[consolev1.UpdateTenantRequest],
) (*connect.Response[consolev1.UpdateTenantResponse], error) {
	tenantId, err := model.ParseTenantID(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	tenant, err := h.useCase.GetTenantById(ctx, tenantId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	orgID, ok := domain.Value[model.OrganizationID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.WithMessage(domainerrors.ErrNotFound, "organization not found"))
	}

	if tenant.OrganizationID != orgID {
		return nil, connect.NewError(connect.CodeNotFound, errors.WithMessage(domainerrors.ErrNotFound, "tenant not found"))
	}

	tenantTypeStr, err := convertTenantType(req.Msg.TenantType)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	input := &dto.UpdateTenantInput{
		ID:          tenantId,
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
		TenantType:  tenantTypeStr,
	}

	tenantID, err := h.useCase.UpdateTenant(ctx, *input)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&consolev1.UpdateTenantResponse{
		Id: tenantID,
	}), nil

}
