-- name: CreateTenant :one
INSERT INTO tenants(
    id,
    organization_id,
    name,
    description,
    tenant_type,
    created_at,
    updated_at
)
VALUES(
    @id,
    @organization_id,
    @name,
    @description,
    @tenant_type,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
RETURNING *;

-- name: GetTenant :one
SELECT sqlc.embed(t) 
FROM tenants t
WHERE id = $1;

-- name: GetAllTenants :many
SELECT sqlc.embed(t) 
FROM tenants t
WHERE organization_id = $1
ORDER BY created_at DESC;
