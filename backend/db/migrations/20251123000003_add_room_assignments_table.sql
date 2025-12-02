-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Room Assignments Table';

CREATE TABLE room_assignments (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    tenant_id UUID NOT NULL,
    room_id UUID NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT room_assignments_date_check CHECK (expires_at IS NULL OR expires_at > assigned_at)
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE room_assignments TO keyhub;

-- Enable RLS
ALTER TABLE room_assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE room_assignments FORCE ROW LEVEL SECURITY;

CREATE POLICY room_assignments_org_isolation ON room_assignments
    FOR ALL
    TO keyhub
    USING (
        current_organization_id() IS NULL
        OR tenant_id IN (
            SELECT id FROM tenants WHERE organization_id = current_organization_id()
        )
    );

CREATE TRIGGER refresh_room_assignments_updated_at
BEFORE UPDATE ON room_assignments
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - room_assignments table rollback';

DROP TRIGGER IF EXISTS refresh_room_assignments_updated_at ON room_assignments;

DROP POLICY IF EXISTS room_assignments_org_isolation ON room_assignments;

DROP TABLE IF EXISTS room_assignments;
-- +goose StatementEnd
