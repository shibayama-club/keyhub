-- name: CreateTenantJoinCode :one
INSERT INTO tenant_join_codes(
    id,
    tenant_id,
    code,
    expires_at,
    max_uses,
    used_count,
    created_at
)
VALUES(
    @id,
    @tenant_id,
    @code,
    @expires_at,
    @max_uses,
    @used_count,
    CURRENT_TIMESTAMP
)
RETURNING sqlc.embed(tenant_join_codes);
