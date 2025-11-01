-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - RLS Functions and Policies (Placeholder for future use)';

CREATE OR REPLACE FUNCTION current_membership_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT current_setting('keyhub.membership_id', true)::uuid
$$;

CREATE OR REPLACE FUNCTION current_tenant_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT tm.tenant_id
  FROM tenant_memberships tm
  WHERE tm.id = current_membership_id()
$$;

CREATE OR REPLACE FUNCTION current_organization_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT current_setting('keyhub.organization_id', true)::uuid
$$;

GRANT EXECUTE ON FUNCTION current_membership_id() TO keyhub;
GRANT EXECUTE ON FUNCTION current_tenant_id() TO keyhub;
GRANT EXECUTE ON FUNCTION current_organization_id() TO keyhub;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - RLS Functions and Policies rollback';

-- Drop functions (all in public schema now)
DROP FUNCTION IF EXISTS current_organization_id();
DROP FUNCTION IF EXISTS current_tenant_id();
DROP FUNCTION IF EXISTS current_membership_id();
-- +goose StatementEnd
