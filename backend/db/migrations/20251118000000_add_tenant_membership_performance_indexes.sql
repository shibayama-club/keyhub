-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Add performance indexes for tenant memberships';

-- Composite index for efficient user tenant lookups with left_at filtering
-- This supports: WHERE user_id = ? AND left_at IS NULL
CREATE INDEX idx_memberships_user_left_tenant ON tenant_memberships(user_id, left_at, tenant_id)
WHERE left_at IS NULL;

-- Composite index for efficient member counting per tenant
-- This supports: WHERE tenant_id = ? AND left_at IS NULL for COUNT queries
CREATE INDEX idx_memberships_tenant_left ON tenant_memberships(tenant_id, left_at)
WHERE left_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - Remove performance indexes for tenant memberships';

DROP INDEX IF EXISTS idx_memberships_tenant_left;
DROP INDEX IF EXISTS idx_memberships_user_left_tenant;
-- +goose StatementEnd
