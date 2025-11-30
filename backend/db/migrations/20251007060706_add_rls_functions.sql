-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - RLS Functions (Basic)';

-- These functions only read session variables, no table dependencies
CREATE OR REPLACE FUNCTION current_membership_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT current_setting('keyhub.membership_id', true)::uuid
$$;

CREATE OR REPLACE FUNCTION current_organization_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT current_setting('keyhub.organization_id', true)::uuid
$$;

GRANT EXECUTE ON FUNCTION current_membership_id() TO keyhub;
GRANT EXECUTE ON FUNCTION current_organization_id() TO keyhub;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - RLS Functions rollback';

DROP FUNCTION IF EXISTS current_organization_id();
DROP FUNCTION IF EXISTS current_membership_id();
-- +goose StatementEnd
