-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Users Table';

CREATE TABLE users (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    email TEXT NOT NULL,
    name TEXT NOT NULL,
    icon TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

GRANT SELECT,INSERT,UPDATE ON TABLE users TO keyhub;

CREATE TRIGGER refresh_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - users table rollback';

DROP TRIGGER IF EXISTS refresh_users_updated_at ON users;

DROP TABLE IF EXISTS users;
-- +goose StatementEnd
