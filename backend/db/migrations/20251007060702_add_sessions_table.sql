-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Sessions Table';

CREATE TABLE sessions (
    session_id TEXT NOT NULL,
    user_id UUID NOT NULL,
    active_membership_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    csrf_token TEXT,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (session_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE sessions TO keyhub;

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - sessions table rollback';

DROP INDEX IF EXISTS idx_sessions_expires;
DROP INDEX IF EXISTS idx_sessions_user;

DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
