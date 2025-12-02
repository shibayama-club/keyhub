package model

import (
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type RoomID uuid.UUID

func (id RoomID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id RoomID) String() string {
	return uuid.UUID(id).String()
}

func ParseRoomID(value string) (RoomID, error) {
	u, err := uuid.Parse(value)
	if err != nil {
		return RoomID{}, errors.WithHint(
			errors.Wrap(err, "failed to parse room ID"),
			"部屋IDの形式が正しくありません。",
		)
	}
	return RoomID(u), nil
}

type RoomName string

func (n RoomName) String() string {
	return string(n)
}

func (n RoomName) Validate() error {
	if n == "" {
		return errors.WithHint(
			errors.New("room name is required"),
			"部屋名は必須です。",
		)
	}

	if utf8.RuneCountInString(string(n)) > 20 {
		return errors.WithHint(
			errors.New("room name must be within 20 characters"),
			"部屋名は20文字以内で入力してください。",
		)
	}
	return nil
}

func NewRoomName(value string) (RoomName, error) {
	n := RoomName(value)
	if err := n.Validate(); err != nil {
		return "", err
	}
	return n, nil
}

type BuildingName string

func (n BuildingName) String() string {
	return string(n)
}

func (n BuildingName) Validate() error {
	if n == "" {
		return errors.WithHint(
			errors.New("building name is required"),
			"建物名は必須です。",
		)
	}

	if utf8.RuneCountInString(string(n)) > 20 {
		return errors.WithHint(
			errors.New("building name must be within 20 characters"),
			"建物名は20文字以内で入力してください。",
		)
	}
	return nil
}

func NewBuildingName(value string) (BuildingName, error) {
	n := BuildingName(value)
	if err := n.Validate(); err != nil {
		return "", err
	}
	return n, nil
}

type FloorNumber string

func (n FloorNumber) String() string {
	return string(n)
}

func (n FloorNumber) Validate() error {
	if n == "" {
		return errors.WithHint(
			errors.New("floor number is required"),
			"階数は必須です。",
		)
	}

	if utf8.RuneCountInString(string(n)) > 10 {
		return errors.WithHint(
			errors.New("floor number must be within 10 characters"),
			"階数は10文字以内で入力してください。",
		)
	}
	return nil
}

func NewFloorNumber(value string) (FloorNumber, error) {
	n := FloorNumber(value)
	if err := n.Validate(); err != nil {
		return "", err
	}
	return n, nil
}

type RoomType string

const (
	RoomTypeUnspecified RoomType = "ROOM_TYPE_UNSPECIFIED"
	RoomTypeClassroom   RoomType = "classroom"
	RoomTypeMeetingRoom RoomType = "meeting_room"
	RoomTypeLaboratory  RoomType = "laboratory"
	RoomTypeOffice      RoomType = "office"
	RoomTypeWorkshop    RoomType = "workshop"
	RoomTypeStorage     RoomType = "storage"
)

func (t RoomType) String() string {
	return string(t)
}

func (t RoomType) Validate() error {
	switch t {
	case RoomTypeClassroom, RoomTypeMeetingRoom, RoomTypeLaboratory, RoomTypeOffice, RoomTypeWorkshop, RoomTypeStorage:
		return nil
	case RoomTypeUnspecified:
		return errors.WithHint(
			errors.New("room type must be specified"),
			"部屋タイプを指定してください。",
		)
	default:
		return errors.WithHintf(
			errors.New("invalid room type"),
			"無効な部屋タイプです: %s", t,
		)
	}
}

func NewRoomType(value string) (RoomType, error) {
	t := RoomType(value)
	if err := t.Validate(); err != nil {
		return "", err
	}
	return t, nil
}

type RoomDescription string

func (d RoomDescription) String() string {
	return string(d)
}

func (d RoomDescription) Validate() error {
	if utf8.RuneCountInString(string(d)) > 200 {
		return errors.WithHint(
			errors.New("room description must be within 200 characters"),
			"部屋の説明は200文字以内で入力してください。",
		)
	}
	return nil
}

func NewRoomDescription(value string) (RoomDescription, error) {
	d := RoomDescription(value)
	if err := d.Validate(); err != nil {
		return "", err
	}
	return d, nil
}

type Room struct {
	ID             RoomID
	OrganizationID OrganizationID
	Name           RoomName
	BuildingName   BuildingName
	FloorNumber    FloorNumber
	Type           RoomType
	Description    RoomDescription
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (r Room) Validate() error {
	if err := r.OrganizationID.Validate(); err != nil {
		return err
	}

	if err := r.Name.Validate(); err != nil {
		return err
	}

	if err := r.BuildingName.Validate(); err != nil {
		return err
	}

	if err := r.FloorNumber.Validate(); err != nil {
		return err
	}

	if err := r.Type.Validate(); err != nil {
		return err
	}

	if err := r.Description.Validate(); err != nil {
		return err
	}

	if r.CreatedAt.IsZero() {
		return errors.WithHint(
			errors.New("created_at is required"),
			"作成日時は必須です。",
		)
	}

	if r.UpdatedAt.IsZero() {
		return errors.WithHint(
			errors.New("updated_at is required"),
			"更新日時は必須です。",
		)
	}

	return nil
}

func NewRoom(
	organizationID OrganizationID,
	name RoomName,
	buildingName BuildingName,
	floorNumber FloorNumber,
	roomType RoomType,
	description RoomDescription,
) (Room, error) {
	now := time.Now()
	room := Room{
		ID:             RoomID(uuid.New()),
		OrganizationID: organizationID,
		Name:           name,
		BuildingName:   buildingName,
		FloorNumber:    floorNumber,
		Type:           roomType,
		Description:    description,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := room.Validate(); err != nil {
		return Room{}, err
	}

	return room, nil
}
