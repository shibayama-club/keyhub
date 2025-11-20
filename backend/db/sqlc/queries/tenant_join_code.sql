-- name: CreateTenantJoinCode :one
INSERT INTO tenant_join_codes(
    id,
    tenant_id,
    code,
    expires_at,
    max_uses,
    used_count
)
VALUES(
    @id,
    @tenant_id,
    @code,
    @expires_at,
    @max_uses,
    @used_count
)
RETURNING sqlc.embed(tenant_join_codes);

-- name: GetTenantByJoinCode :one
SELECT
    sqlc.embed(t)
FROM tenant_join_codes tjc
INNER JOIN tenants t ON tjc.tenant_id = t.id
WHERE tjc.code = $1
    AND (tjc.expires_at IS NULL OR tjc.expires_at > CURRENT_TIMESTAMP)
    AND (tjc.max_uses = 0 OR tjc.used_count < tjc.max_uses);

-- name: UpdateTenantJoinCode :one
UPDATE tenant_join_codes
SET 
    code = @code,
    expires_at = @expires_at,
    max_uses = @max_uses
WHERE tenant_id = @tenant_id
RETURNING sqlc.embed(tenant_join_codes);