-- name: CreateRoom :exec
INSERT INTO rooms(
    id,
    organization_id,
    name,
    building_name,
    floor_number,
    room_type,
    description
)
VALUES(
    @id,
    @organization_id,
    @name,
    @building_name,
    @floor_number,
    @room_type,
    @description
);

-- name: GetRoomById :one
SELECT sqlc.embed(r)
FROM rooms r
WHERE r.id = $1;
