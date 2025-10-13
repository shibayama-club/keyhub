-- name: InsertTenant :one
INSERT INTO tenants(
    id,
    name,
    description,
    tenant_type,
    created_at,
    updated_at
)
VALUES(
    @id,
    @name,
    @slug,
    @password_hash,
    @created_at,
    @updated_at
)
RETURNING *;