-- name: CreateKey :exec
INSERT INTO keys(
    id,
    room_id,
    organization_id,
    key_number,
    status
)
VALUES(
    @id,
    @room_id,
    @organization_id,
    @key_number,
    @status
);
