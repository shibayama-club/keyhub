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

func (h *Handler) CreateKey(
	ctx context.Context,
	req *connect.Request[consolev1.CreateKeyRequest],
) (*connect.Response[consolev1.CreateKeyResponse], error) {
	orgID, ok := domain.Value[model.OrganizationID](ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.WithMessage(domainerrors.ErrNotFound, "organization not found"))
	}

	roomID, err := model.ParseRoomID(req.Msg.RoomId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err, "invalid room ID"))
	}

	input := dto.CreateKeyInput{
		RoomID:         roomID,
		OrganizationID: orgID,
		KeyNumber:      req.Msg.KeyNumber,
	}

	keyID, err := h.useCase.CreateKey(ctx, input)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&consolev1.CreateKeyResponse{
		Id: keyID,
	}), nil
}

func convertToProtoKeyStatus(status model.KeyStatus) consolev1.KeyStatus {
	switch status {
	case model.KeyStatusAvailable:
		return consolev1.KeyStatus_KEY_STATUS_AVAILABLE
	case model.KeyStatusInUse:
		return consolev1.KeyStatus_KEY_STATUS_IN_USE
	case model.KeyStatusLost:
		return consolev1.KeyStatus_KEY_STATUS_LOST
	case model.KeyStatusDamaged:
		return consolev1.KeyStatus_KEY_STATUS_DAMAGED
	default:
		return consolev1.KeyStatus_KEY_STATUS_UNSPECIFIED
	}
}

func (h *Handler) GetKeysByRoom(
	ctx context.Context,
	req *connect.Request[consolev1.GetKeysByRoomRequest],
) (*connect.Response[consolev1.GetKeysByRoomResponse], error) {
	roomID, err := model.ParseRoomID(req.Msg.RoomId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err, "invalid room ID"))
	}

	keys, err := h.useCase.GetKeysByRoom(ctx, roomID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoKeys := lo.Map(keys, func(key model.Key, _ int) *consolev1.Key {
		return &consolev1.Key{
			Id:        key.ID.String(),
			KeyNumber: key.KeyNumber.String(),
			RoomId:    key.RoomID.String(),
			Status:    convertToProtoKeyStatus(key.Status),
		}
	})

	return connect.NewResponse(&consolev1.GetKeysByRoomResponse{
		Keys: protoKeys,
	}), nil
}
