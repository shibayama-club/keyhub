-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Tenants Table';

CREATE TABLE tenants (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL  DEFAULT '',   
    password_hash TEXT NOT NULL,
    tenant_type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

GRANT SELECT,INSERT,UPDATE ON TABLE tenants TO keyhub;

-- Create unique index for non-empty description
CREATE UNIQUE INDEX idx_tenants_description_nonempty
  ON tenants (description)
  WHERE description <> '';

-- Create trigger for updating the updated_at column
CREATE TRIGGER refresh_tenants_updated_at
BEFORE UPDATE ON tenants
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - tenants table rollback';

DROP TRIGGER IF EXISTS refresh_tenants_updated_at ON tenants;

DROP INDEX IF EXISTS idx_tenants_description_nonempty;

DROP TABLE IF EXISTS tenants;
-- +goose StatementEnd
