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
SELECT * FROM tenants
WHERE id = $1;

-- name: GetTenantsByOrganization :many
SELECT * FROM tenants
WHERE organization_id = $1
ORDER BY created_at DESC;
