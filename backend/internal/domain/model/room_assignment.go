package model

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type RoomAssignmentID uuid.UUID

func (id RoomAssignmentID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id RoomAssignmentID) String() string {
	return uuid.UUID(id).String()
}

func ParseRoomAssignmentID(value string) (RoomAssignmentID, error) {
	u, err := uuid.Parse(value)
	if err != nil {
		return RoomAssignmentID{}, errors.WithHint(
			errors.Wrap(err, "failed to parse room assignment ID"),
			"部屋割り当てIDの形式が正しくありません。",
		)
	}
	return RoomAssignmentID(u), nil
}

type RoomAssignment struct {
	ID         RoomAssignmentID
	TenantID   TenantID
	RoomID     RoomID
	AssignedAt time.Time
	ExpiresAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (ra RoomAssignment) Validate() error {
	if ra.AssignedAt.IsZero() {
		return errors.WithHint(
			errors.New("assigned_at is required"),
			"割り当て日時は必須です。",
		)
	}

	if ra.ExpiresAt != nil && !ra.ExpiresAt.After(ra.AssignedAt) {
		return errors.WithHint(
			errors.New("expires_at must be after assigned_at"),
			"有効期限は割り当て日時より後である必要があります。",
		)
	}

	if ra.CreatedAt.IsZero() {
		return errors.WithHint(
			errors.New("created_at is required"),
			"作成日時は必須です。",
		)
	}

	if ra.UpdatedAt.IsZero() {
		return errors.WithHint(
			errors.New("updated_at is required"),
			"更新日時は必須です。",
		)
	}

	return nil
}

func NewRoomAssignment(
	tenantID TenantID,
	roomID RoomID,
	expiresAt *time.Time,
) (RoomAssignment, error) {
	now := time.Now()
	assignment := RoomAssignment{
		ID:         RoomAssignmentID(uuid.New()),
		TenantID:   tenantID,
		RoomID:     roomID,
		AssignedAt: now,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := assignment.Validate(); err != nil {
		return RoomAssignment{}, err
	}

	return assignment, nil
}
