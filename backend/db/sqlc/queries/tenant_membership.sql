-- name: CreateTenantMembership :exec
INSERT INTO tenant_memberships(
    id,
    tenant_id,
    user_id,
    role
)
VALUES(
    @id,
    @tenant_id,
    @user_id,
    @role
);

-- name: IncrementJoinCodeUsedCount :exec
UPDATE tenant_join_codes
SET used_count = used_count + 1
WHERE code = $1;

-- name: GetTenantMembershipByTenantAndUser :one
SELECT sqlc.embed(tenant_memberships)
FROM tenant_memberships
WHERE tenant_id = $1 AND user_id = $2;

-- name: ClearActiveMembershipByTenantID :exec
UPDATE sessions 
SET active_membership_id = NULL 
WHERE active_membership_id IN (
    SELECT id FROM tenant_memberships WHERE tenant_id = $1
);