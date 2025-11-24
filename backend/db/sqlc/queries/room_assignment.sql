-- name: CreateRoomAssignment :exec
INSERT INTO room_assignments(
    id,
    tenant_id,
    room_id,
    assigned_at,
    expires_at
)
VALUES(
    @id,
    @tenant_id,
    @room_id,
    @assigned_at,
    @expires_at
);
