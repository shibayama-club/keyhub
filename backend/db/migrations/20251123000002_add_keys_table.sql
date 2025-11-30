-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Keys Table';

CREATE TABLE keys (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    room_id UUID NOT NULL,
    organization_id UUID NOT NULL DEFAULT '550e8400-e29b-41d4-a716-446655440000',
    key_number TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'available',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT keys_status_check CHECK (status IN ('available', 'in_use', 'lost', 'damaged'))
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE keys TO keyhub;

-- Enable RLS
ALTER TABLE keys ENABLE ROW LEVEL SECURITY;
ALTER TABLE keys FORCE ROW LEVEL SECURITY;

CREATE POLICY keys_org_isolation ON keys
    FOR ALL
    TO keyhub
    USING (organization_id = current_organization_id());

CREATE TRIGGER refresh_keys_updated_at
BEFORE UPDATE ON keys
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - keys table rollback';

DROP TRIGGER IF EXISTS refresh_keys_updated_at ON keys;

DROP POLICY IF EXISTS keys_org_isolation ON keys;

DROP TABLE IF EXISTS keys;
-- +goose StatementEnd
