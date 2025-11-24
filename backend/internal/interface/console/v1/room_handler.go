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
)

func convertRoomType(protoType consolev1.RoomType) (string, error) {
	switch protoType {
	case consolev1.RoomType_ROOM_TYPE_CLASSROOM:
		return model.RoomTypeClassroom.String(), nil
	case consolev1.RoomType_ROOM_TYPE_MEETING_ROOM:
		return model.RoomTypeMeetingRoom.String(), nil
	case consolev1.RoomType_ROOM_TYPE_LABORATORY:
		return model.RoomTypeLaboratory.String(), nil
	case consolev1.RoomType_ROOM_TYPE_OFFICE:
		return model.RoomTypeOffice.String(), nil
	case consolev1.RoomType_ROOM_TYPE_WORKSHOP:
		return model.RoomTypeWorkshop.String(), nil
	case consolev1.RoomType_ROOM_TYPE_STORAGE:
		return model.RoomTypeStorage.String(), nil
	default:
		return "", errors.New("invalid room type")
	}
}

func (h *Handler) CreateRoom(
	ctx context.Context,
	req *connect.Request[consolev1.CreateRoomRequest],
) (*connect.Response[consolev1.CreateRoomResponse], error) {
	orgID, ok := domain.Value[model.OrganizationID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.WithMessage(domainerrors.ErrNotFound, "organization not found"))
	}

	roomTypeStr, err := convertRoomType(req.Msg.RoomType)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	input := dto.CreateRoomInput{
		OrganizationID: orgID,
		Name:           req.Msg.Name,
		BuildingName:   req.Msg.BuildingName,
		FloorNumber:    req.Msg.FloorNumber,
		RoomType:       roomTypeStr,
		Description:    req.Msg.Description,
	}

	roomID, err := h.useCase.CreateRoom(ctx, input)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&consolev1.CreateRoomResponse{
		Id: roomID,
	}), nil
}

func (h *Handler) AssignRoomToTenant(
	ctx context.Context,
	req *connect.Request[consolev1.AssignRoomToTenantRequest],
) (*connect.Response[consolev1.AssignRoomToTenantResponse], error) {
	tenantID, err := model.ParseTenantID(req.Msg.TenantId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err, "invalid tenant ID"))
	}

	roomID, err := model.ParseRoomID(req.Msg.RoomId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err, "invalid room ID"))
	}

	input := dto.AssignRoomToTenantInput{
		TenantID: tenantID,
		RoomID:   roomID,
	}

	if req.Msg.ExpiresAt != nil {
		expiryTime := req.Msg.ExpiresAt.AsTime()
		input.ExpiresAt = &expiryTime
	}

	assignmentID, err := h.useCase.AssignRoomToTenant(ctx, input)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&consolev1.AssignRoomToTenantResponse{
		AssignmentId: assignmentID,
	}), nil
}

func convertToProtoRoomType(roomType model.RoomType) consolev1.RoomType {
	switch roomType {
	case model.RoomTypeClassroom:
		return consolev1.RoomType_ROOM_TYPE_CLASSROOM
	case model.RoomTypeMeetingRoom:
		return consolev1.RoomType_ROOM_TYPE_MEETING_ROOM
	case model.RoomTypeLaboratory:
		return consolev1.RoomType_ROOM_TYPE_LABORATORY
	case model.RoomTypeOffice:
		return consolev1.RoomType_ROOM_TYPE_OFFICE
	case model.RoomTypeWorkshop:
		return consolev1.RoomType_ROOM_TYPE_WORKSHOP
	case model.RoomTypeStorage:
		return consolev1.RoomType_ROOM_TYPE_STORAGE
	default:
		return consolev1.RoomType_ROOM_TYPE_UNSPECIFIED
	}
}

func (h *Handler) GetAllRooms(
	ctx context.Context,
	req *connect.Request[consolev1.GetAllRoomsRequest],
) (*connect.Response[consolev1.GetAllRoomsResponse], error) {
	orgID, ok := domain.Value[model.OrganizationID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.WithMessage(domainerrors.ErrNotFound, "organization not found"))
	}

	rooms, err := h.useCase.GetAllRooms(ctx, orgID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoRooms := lo.Map(rooms, func(room model.Room, _ int) *consolev1.Room {
		return &consolev1.Room{
			Id:           room.ID.String(),
			Name:         room.Name.String(),
			BuildingName: room.BuildingName.String(),
			FloorNumber:  room.FloorNumber.String(),
			RoomType:     convertToProtoRoomType(room.Type),
			Description:  room.Description.String(),
		}
	})

	return connect.NewResponse(&consolev1.GetAllRoomsResponse{
		Rooms: protoRooms,
	}), nil
}
