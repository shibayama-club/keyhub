-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Tenant Domains Table';

CREATE TABLE tenant_domains (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    tenant_id UUID NOT NULL,
    domain TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE tenant_domains TO keyhub;

CREATE INDEX idx_tenant_domains_tenant ON tenant_domains(tenant_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - tenant_domains table rollback';

DROP INDEX IF EXISTS idx_tenant_domains_tenant;

DROP TABLE IF EXISTS tenant_domains;
-- +goose StatementEnd
