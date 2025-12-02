-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query - Rooms Table';

CREATE TABLE rooms (
    id UUID NOT NULL DEFAULT UUID_GENERATE_V4(),
    organization_id UUID NOT NULL DEFAULT '550e8400-e29b-41d4-a716-446655440000',
    name TEXT NOT NULL,
    building_name TEXT NOT NULL,
    floor_number TEXT NOT NULL,
    room_type TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT rooms_room_type_check CHECK (room_type IN ('classroom', 'meeting_room', 'laboratory', 'office', 'workshop', 'storage'))
);

GRANT SELECT,INSERT,UPDATE,DELETE ON TABLE rooms TO keyhub;

-- Enable RLS
ALTER TABLE rooms ENABLE ROW LEVEL SECURITY;
ALTER TABLE rooms FORCE ROW LEVEL SECURITY;

CREATE POLICY rooms_org_isolation ON rooms
    FOR ALL
    TO keyhub
    USING (
        current_organization_id() IS NULL
        OR organization_id = current_organization_id()
    );

CREATE TRIGGER refresh_rooms_updated_at
BEFORE UPDATE ON rooms
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - rooms table rollback';

DROP TRIGGER IF EXISTS refresh_rooms_updated_at ON rooms;

DROP POLICY IF EXISTS rooms_org_isolation ON rooms;

DROP TABLE IF EXISTS rooms;
-- +goose StatementEnd
