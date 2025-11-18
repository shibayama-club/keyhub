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
RETURNING sqlc.embed(tenants);

-- name: GetTenant :one
SELECT sqlc.embed(t) 
FROM tenants t
WHERE id = $1;

-- name: GetAllTenants :many
SELECT sqlc.embed(t)
FROM tenants t
WHERE organization_id = $1
ORDER BY created_at DESC;

-- name: GetTenantsByUserID :many
SELECT
    sqlc.embed(t),
    COUNT(tm_all.id)::INT AS member_count
FROM tenants t
INNER JOIN tenant_memberships tm ON t.id = tm.tenant_id
LEFT JOIN tenant_memberships tm_all ON t.id = tm_all.tenant_id AND tm_all.left_at IS NULL
WHERE tm.user_id = $1
  AND tm.left_at IS NULL
GROUP BY t.id
ORDER BY t.created_at DESC;
