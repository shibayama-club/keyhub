-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Tenant Memberships Table';

CREATE TABLE tenant_memberships (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'member')) DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMPTZ,
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (tenant_id, user_id)
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE tenant_memberships TO keyhub;

CREATE INDEX idx_memberships_user ON tenant_memberships(user_id);

-- Add foreign key constraint and index to sessions.active_membership_id
ALTER TABLE sessions
    ADD CONSTRAINT fk_sessions_active_membership
    FOREIGN KEY (active_membership_id) REFERENCES tenant_memberships(id);

CREATE INDEX idx_sessions_active_membership ON sessions(active_membership_id);

-- Create current_tenant_id function (depends on tenant_memberships table)
CREATE OR REPLACE FUNCTION current_tenant_id()
RETURNS uuid LANGUAGE sql STABLE AS $$
  SELECT tm.tenant_id
  FROM tenant_memberships tm
  WHERE tm.id = current_membership_id()
$$;

GRANT EXECUTE ON FUNCTION current_tenant_id() TO keyhub;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - tenant_memberships table rollback';

-- Remove sessions foreign key and index first
DROP INDEX IF EXISTS idx_sessions_active_membership;

ALTER TABLE sessions
    DROP CONSTRAINT IF EXISTS fk_sessions_active_membership;

DROP INDEX IF EXISTS idx_memberships_user;

DROP FUNCTION IF EXISTS current_tenant_id();

DROP TABLE IF EXISTS tenant_memberships;
-- +goose StatementEnd
