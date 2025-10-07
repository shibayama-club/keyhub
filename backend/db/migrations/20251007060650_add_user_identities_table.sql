-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - User Identities Table';

CREATE TABLE user_identities (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    user_id UUID NOT NULL,
    provider TEXT NOT NULL,
    provider_sub TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (provider, provider_sub)
);

GRANT SELECT,INSERT,UPDATE ON TABLE user_identities TO keyhub;

-- Create trigger for updating the updated_at column
CREATE TRIGGER refresh_user_identities_updated_at
BEFORE UPDATE ON user_identities
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - user_identities table rollback';

DROP TRIGGER IF EXISTS refresh_user_identities_updated_at ON user_identities;

DROP TABLE IF EXISTS user_identities;
-- +goose StatementEnd
