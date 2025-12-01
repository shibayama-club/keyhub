-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Tenants Table';

CREATE TABLE tenants (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    organization_id UUID NOT NULL DEFAULT '550e8400-e29b-41d4-a716-446655440000',
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL  DEFAULT '',
    tenant_type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

GRANT SELECT,INSERT,UPDATE ON TABLE tenants TO keyhub;

CREATE INDEX idx_tenants_organization_id ON tenants(organization_id);

-- Enable RLS
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenants FORCE ROW LEVEL SECURITY;

CREATE POLICY tenants_org_isolation ON tenants
    FOR ALL
    TO keyhub
    USING (
        current_organization_id() IS NULL
        OR organization_id = current_organization_id()
    );

CREATE TRIGGER refresh_tenants_updated_at
BEFORE UPDATE ON tenants
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - tenants table rollback';

DROP TRIGGER IF EXISTS refresh_tenants_updated_at ON tenants;

DROP POLICY IF EXISTS tenants_org_isolation ON tenants;

DROP INDEX IF EXISTS idx_tenants_organization_id;

DROP TABLE IF EXISTS tenants;
-- +goose StatementEnd
