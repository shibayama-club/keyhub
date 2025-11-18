-- name: CreateTenantMembership :one
INSERT INTO tenant_memberships(
    id,
    tenant_id,
    user_id,
    role,
    created_at
)
VALUES(
    @id,
    @tenant_id,
    @user_id,
    @role,
    CURRENT_TIMESTAMP
)
RETURNING sqlc.embed(tenant_memberships);

-- name: IncrementJoinCodeUsedCount :exec
UPDATE tenant_join_codes
SET used_count = used_count + 1
WHERE code = $1;

-- name: GetTenantMembershipByTenantAndUser :one
SELECT sqlc.embed(tenant_memberships)
FROM tenant_memberships
WHERE tenant_id = $1 AND user_id = $2;
