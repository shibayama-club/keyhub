package v1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
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
