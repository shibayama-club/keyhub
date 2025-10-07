-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Tenant Memberships Table';

CREATE TABLE tenant_memberships (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'member')) DEFAULT 'member',
    status TEXT NOT NULL CHECK (status IN ('active', 'invited', 'left')) DEFAULT 'active',
    joined_via TEXT CHECK (joined_via IN ('domain', 'code', 'manual')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMPTZ,
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (tenant_id, user_id)
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE tenant_memberships TO keyhub;

CREATE INDEX idx_memberships_user ON tenant_memberships(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - tenant_memberships table rollback';

DROP INDEX IF EXISTS idx_memberships_user;

DROP TABLE IF EXISTS tenant_memberships;
-- +goose StatementEnd
