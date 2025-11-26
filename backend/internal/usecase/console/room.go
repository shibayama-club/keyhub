package console

import (
	"context"

	"github.com/cockroachdb/errors"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"github.com/shibayama-club/keyhub/internal/domain/model"
	"github.com/shibayama-club/keyhub/internal/domain/repository"
	"github.com/shibayama-club/keyhub/internal/usecase/console/dto"
)

func (u *UseCase) CreateRoom(ctx context.Context, input dto.CreateRoomInput) (string, error) {
	roomName, err := model.NewRoomName(input.Name)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid room name")
	}

	buildingName, err := model.NewBuildingName(input.BuildingName)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid building name")
	}

	floorNumber, err := model.NewFloorNumber(input.FloorNumber)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid floor number")
	}

	roomType, err := model.NewRoomType(input.RoomType)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid room type")
	}

	roomDescription, err := model.NewRoomDescription(input.Description)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "invalid room description")
	}

	room, err := model.NewRoom(
		input.OrganizationID,
		roomName,
		buildingName,
		floorNumber,
		roomType,
		roomDescription,
	)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create room")
	}

	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		err = tx.CreateRoom(ctx, repository.CreateRoomArg{
			ID:             room.ID,
			OrganizationID: room.OrganizationID,
			Name:           room.Name,
			BuildingName:   room.BuildingName,
			FloorNumber:    room.FloorNumber,
			Type:           room.Type,
			Description:    room.Description,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create room in repository")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return room.ID.String(), nil
}

func (u *UseCase) AssignRoomToTenant(ctx context.Context, input dto.AssignRoomToTenantInput) (string, error) {
	// Verify room exists
	_, err := u.repo.GetRoomByID(ctx, input.RoomID)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrNotFound), "room not found")
	}

	// Verify tenant exists
	_, err = u.repo.GetTenantByID(ctx, input.TenantID)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrNotFound), "tenant not found")
	}

	assignment, err := model.NewRoomAssignment(
		input.TenantID,
		input.RoomID,
		input.ExpiresAt,
	)
	if err != nil {
		return "", errors.Wrap(errors.Mark(err, domainerrors.ErrValidation), "failed to create room assignment")
	}

	err = u.repo.WithTransaction(ctx, func(ctx context.Context, tx repository.Transaction) error {
		err = tx.CreateRoomAssignment(ctx, repository.CreateRoomAssignmentArg{
			ID:         assignment.ID,
			TenantID:   assignment.TenantID,
			RoomID:     assignment.RoomID,
			AssignedAt: assignment.AssignedAt,
			ExpiresAt:  assignment.ExpiresAt,
		})
		if err != nil {
			return errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to create room assignment in repository")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return assignment.ID.String(), nil
}

func (u *UseCase) GetAllRooms(ctx context.Context, organizationID model.OrganizationID) ([]model.Room, error) {
	rooms, err := u.repo.GetAllRooms(ctx, organizationID)
	if err != nil {
		return nil, errors.Wrap(errors.Mark(err, domainerrors.ErrInternal), "failed to get all rooms")
	}
	return rooms, nil
}
