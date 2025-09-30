-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Users Table';

-- Create function for updating the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

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

-- Create a single trigger for updating the updated_at column
CREATE TRIGGER refresh_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - users table rollback';

DROP TRIGGER IF EXISTS refresh_users_updated_at ON users;

DROP TABLE IF EXISTS users;

DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd
