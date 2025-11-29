package v1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/samber/lo"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	appv1 "github.com/shibayama-club/keyhub/internal/interface/gen/keyhub/app/v1"
)

func (h *Handler) GetRoomsByTenant(
	ctx context.Context,
	req *connect.Request[appv1.GetRoomsByTenantRequest],
) (*connect.Response[appv1.GetRoomsByTenantResponse], error) {
	tenantID, err := model.ParseTenantID(req.Msg.TenantId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err, "invalid tenant ID"))
	}

	rooms, err := h.useCase.GetRoomsByTenant(ctx, tenantID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoRooms := make([]*appv1.Room, 0, len(rooms))
	for _, room := range rooms {
		keys, err := h.useCase.GetKeysByRoom(ctx, room.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err, "failed to get keys for room"))
		}

		protoKeys := lo.Map(keys, func(key model.Key, _ int) *appv1.Key {
			return &appv1.Key{
				Id:        key.ID.String(),
				KeyNumber: key.KeyNumber.String(),
				RoomId:    key.RoomID.String(),
				Status:    convertToProtoKeyStatus(key.Status),
			}
		})

		protoRooms = append(protoRooms, &appv1.Room{
			Id:           room.ID.String(),
			Name:         room.Name.String(),
			BuildingName: room.BuildingName.String(),
			FloorNumber:  room.FloorNumber.String(),
			RoomType:     convertToProtoRoomType(room.Type),
			Description:  room.Description.String(),
			Keys:         protoKeys,
		})
	}

	return connect.NewResponse(&appv1.GetRoomsByTenantResponse{
		Rooms: protoRooms,
	}), nil
}

func convertToProtoRoomType(roomType model.RoomType) appv1.RoomType {
	switch roomType {
	case model.RoomTypeClassroom:
		return appv1.RoomType_ROOM_TYPE_CLASSROOM
	case model.RoomTypeMeetingRoom:
		return appv1.RoomType_ROOM_TYPE_MEETING_ROOM
	case model.RoomTypeLaboratory:
		return appv1.RoomType_ROOM_TYPE_LABORATORY
	case model.RoomTypeOffice:
		return appv1.RoomType_ROOM_TYPE_OFFICE
	case model.RoomTypeWorkshop:
		return appv1.RoomType_ROOM_TYPE_WORKSHOP
	case model.RoomTypeStorage:
		return appv1.RoomType_ROOM_TYPE_STORAGE
	default:
		return appv1.RoomType_ROOM_TYPE_UNSPECIFIED
	}
}

func convertToProtoKeyStatus(status model.KeyStatus) appv1.KeyStatus {
	switch status {
	case model.KeyStatusAvailable:
		return appv1.KeyStatus_KEY_STATUS_AVAILABLE
	case model.KeyStatusInUse:
		return appv1.KeyStatus_KEY_STATUS_IN_USE
	case model.KeyStatusLost:
		return appv1.KeyStatus_KEY_STATUS_LOST
	case model.KeyStatusDamaged:
		return appv1.KeyStatus_KEY_STATUS_DAMAGED
	default:
		return appv1.KeyStatus_KEY_STATUS_UNSPECIFIED
	}
}
