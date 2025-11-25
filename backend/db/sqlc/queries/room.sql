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

-- name: GetAllRooms :many
SELECT sqlc.embed(r)
FROM rooms r
WHERE organization_id = $1
ORDER BY created_at DESC;

-- name: GetRoomsByTenant :many
SELECT sqlc.embed(r)
FROM rooms r
INNER JOIN room_assignments ra ON r.id = ra.room_id
WHERE ra.tenant_id = $1
  AND (ra.expires_at IS NULL OR ra.expires_at > NOW())
ORDER BY r.created_at DESC;
