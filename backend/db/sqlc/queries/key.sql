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

-- name: GetKeysByRoom :many
SELECT sqlc.embed(k)
FROM keys k
WHERE k.room_id = $1
ORDER BY k.created_at DESC;
