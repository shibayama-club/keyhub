-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Console Sessions Table';

CREATE TABLE console_sessions (
    session_id TEXT NOT NULL,
    tenant_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (session_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE console_sessions TO keyhub;

CREATE INDEX idx_console_sessions_tenant ON console_sessions(tenant_id);
CREATE INDEX idx_console_sessions_expires ON console_sessions(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - console_sessions table rollback';

DROP INDEX IF EXISTS idx_console_sessions_expires;
DROP INDEX IF EXISTS idx_console_sessions_tenant;

DROP TABLE IF EXISTS console_sessions;
-- +goose StatementEnd
