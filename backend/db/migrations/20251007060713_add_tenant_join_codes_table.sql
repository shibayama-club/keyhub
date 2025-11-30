-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Tenant Join Codes Table';

CREATE TABLE tenant_join_codes (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ,
    max_uses INTEGER NOT NULL DEFAULT 0,
    used_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE tenant_join_codes TO keyhub;

CREATE INDEX idx_join_codes_tenant ON tenant_join_codes(tenant_id);
CREATE INDEX idx_join_codes_exp ON tenant_join_codes(expires_at);

-- Enable RLS
ALTER TABLE tenant_join_codes ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_join_codes FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_join_codes_org_isolation ON tenant_join_codes
    FOR ALL
    TO keyhub
    USING (
        current_organization_id() IS NULL
        OR tenant_id IN (
            SELECT id FROM tenants WHERE organization_id = current_organization_id()
        )
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - tenant_join_codes table rollback';

DROP INDEX IF EXISTS idx_join_codes_exp;
DROP INDEX IF EXISTS idx_join_codes_tenant;

DROP POLICY IF EXISTS tenant_join_codes_org_isolation ON tenant_join_codes;

DROP TABLE IF EXISTS tenant_join_codes;
-- +goose StatementEnd
