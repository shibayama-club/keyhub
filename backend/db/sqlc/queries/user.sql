-- name: InsertUser :one
INSERT INTO users (
    id,
    email,
    name,
    icon,
    created_at,
    updated_at
)
VALUES (
    @id,
    @email,
    @name,
    @icon,
    @created_at,
    @updated_at
)
RETURNING *;
