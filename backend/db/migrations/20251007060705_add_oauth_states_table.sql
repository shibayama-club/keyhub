-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - OAuth States Table';

CREATE TABLE oauth_states (
    state TEXT NOT NULL,
    code_verifier TEXT NOT NULL,
    nonce TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    consumed_at TIMESTAMPTZ,
    PRIMARY KEY (state)
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE oauth_states TO keyhub;

CREATE INDEX idx_oauth_states_created ON oauth_states(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - oauth_states table rollback';

DROP INDEX IF EXISTS idx_oauth_states_created;

DROP TABLE IF EXISTS oauth_states;
-- +goose StatementEnd
