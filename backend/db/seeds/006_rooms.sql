-- +goose Up
-- +goose StatementBegin
SELECT 'Seed: Insert rooms';

-- organization_id: 550e8400-e29b-41d4-a716-446655440000 (default)
-- room_type: 'classroom', 'meeting_room', 'laboratory', 'office', 'workshop', 'storage'

INSERT INTO rooms (id, organization_id, name, building_name, floor_number, room_type, description, created_at, updated_at) VALUES
    ('40000000-0000-0000-0000-000000000001', '550e8400-e29b-41d4-a716-446655440000', '会議室A', '本館', '3F', 'meeting_room', '最大20名収容可能な大会議室です。プロジェクター完備。', NOW(), NOW()),
    ('40000000-0000-0000-0000-000000000002', '550e8400-e29b-41d4-a716-446655440000', '教室101', '東館', '1F', 'classroom', '40名収容可能な講義室です。', NOW(), NOW()),
    ('40000000-0000-0000-0000-000000000003', '550e8400-e29b-41d4-a716-446655440000', 'AI実験室', '研究棟', '2F', 'laboratory', 'GPUサーバーを備えたAI研究用実験室。', NOW(), NOW()),
    ('40000000-0000-0000-0000-000000000004', '550e8400-e29b-41d4-a716-446655440000', '総務課オフィス', '本館', '1F', 'office', '総務部の執務スペースです。', NOW(), NOW()),
    ('40000000-0000-0000-0000-000000000005', '550e8400-e29b-41d4-a716-446655440000', '工作室B', '西館', 'B1F', 'workshop', '電子工作・3Dプリンタを備えた工作室。', NOW(), NOW()),
    ('40000000-0000-0000-0000-000000000006', '550e8400-e29b-41d4-a716-446655440000', '資材倉庫', '西館', 'B2F', 'storage', '各種機材・消耗品の保管庫です。', NOW(), NOW());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'Seed Rollback: Delete rooms';

DELETE FROM rooms WHERE id IN (
    '40000000-0000-0000-0000-000000000001',
    '40000000-0000-0000-0000-000000000002',
    '40000000-0000-0000-0000-000000000003',
    '40000000-0000-0000-0000-000000000004',
    '40000000-0000-0000-0000-000000000005',
    '40000000-0000-0000-0000-000000000006'
);

-- +goose StatementEnd
