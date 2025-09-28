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

CREATE TRIGGER refresh_users_updated_at_keyhub1
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION refresh_updated_at_keyhub1();

CREATE TRIGGER refresh_users_updated_at_keyhub2
BEFORE UPDATE OF updated_at ON users
FOR EACH ROW EXECUTE FUNCTION refresh_updated_at_keyhub2();

CREATE TRIGGER refresh_users_updated_at_keyhub3
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION refresh_updated_at_keyhub3();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - users table rollback';

DROP TRIGGER IF EXISTS refresh_users_updated_at_keyhub3 ON users;
DROP TRIGGER IF EXISTS refresh_users_updated_at_keyhub2 ON users;
DROP TRIGGER IF EXISTS refresh_users_updated_at_keyhub1 ON users;

DROP TABLE IF EXISTS users;
-- +goose StatementEnd
