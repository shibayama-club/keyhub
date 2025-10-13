-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - RLS Functions and Policies';

-- Create app schema if not exists
CREATE SCHEMA IF NOT EXISTS app;

-- Create helper functions for RLS
CREATE OR REPLACE FUNCTION app.current_membership_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT current_setting('app.membership_id', true)::uuid
$$;

CREATE OR REPLACE FUNCTION app.current_tenant_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT tm.tenant_id
  FROM tenant_memberships tm
  WHERE tm.id = app.current_membership_id()
$$;

-- Enable RLS on multitenant tables
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_join_codes ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_memberships ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
CREATE POLICY tenant_is_current ON tenants
  USING (id = app.current_tenant_id());

CREATE POLICY tenant_is_current ON tenant_join_codes
  USING (tenant_id = app.current_tenant_id());

CREATE POLICY tenant_is_current ON tenant_memberships
  USING (tenant_id = app.current_tenant_id());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - RLS Functions and Policies rollback';

-- Drop policies
DROP POLICY IF EXISTS tenant_is_current ON tenant_memberships;
DROP POLICY IF EXISTS tenant_is_current ON tenant_join_codes;
DROP POLICY IF EXISTS tenant_is_current ON tenants;

-- Disable RLS
ALTER TABLE tenant_memberships DISABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_join_codes DISABLE ROW LEVEL SECURITY;
ALTER TABLE tenants DISABLE ROW LEVEL SECURITY;

-- Drop functions
DROP FUNCTION IF EXISTS app.current_tenant_id();
DROP FUNCTION IF EXISTS app.current_membership_id();

-- Drop schema (only if empty)
DROP SCHEMA IF EXISTS app CASCADE;
-- +goose StatementEnd
