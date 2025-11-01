-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Baseline database setup';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

GRANT EXECUTE ON FUNCTION update_updated_at_column() TO keyhub;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - Baseline rollback';

DROP FUNCTION IF EXISTS update_updated_at_column();
DROP EXTENSION IF EXISTS "uuid-ossp";

-- +goose StatementEnd
