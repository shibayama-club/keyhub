-- +goose Up
-- +goose StatementBegin
SELECT 'Seed: Insert tenants';

-- organization_id: 550e8400-e29b-41d4-a716-446655440000 (default)
-- tenant_type: TENANT_TYPE_TEAM, TENANT_TYPE_DEPARTMENT, TENANT_TYPE_PROJECT, TENANT_TYPE_LABORATORY

INSERT INTO tenants (id, organization_id, name, description, tenant_type, created_at, updated_at) VALUES
    ('10000000-0000-0000-0000-000000000001', '550e8400-e29b-41d4-a716-446655440000', '開発チームAlpha', 'フロントエンド・バックエンド開発を担当するチームです。', 'TENANT_TYPE_TEAM', NOW(), NOW()),
    ('10000000-0000-0000-0000-000000000002', '550e8400-e29b-41d4-a716-446655440000', '総務部', '社内の総務業務を担当する部署です。', 'TENANT_TYPE_DEPARTMENT', NOW(), NOW()),
    ('10000000-0000-0000-0000-000000000003', '550e8400-e29b-41d4-a716-446655440000', 'KeyHubプロジェクト', '鍵管理システムの開発プロジェクトです。', 'TENANT_TYPE_PROJECT', NOW(), NOW()),
    ('10000000-0000-0000-0000-000000000004', '550e8400-e29b-41d4-a716-446655440000', 'AI研究室', 'AI・機械学習の研究を行う研究室です。', 'TENANT_TYPE_LABORATORY', NOW(), NOW()),
    ('10000000-0000-0000-0000-000000000005', '550e8400-e29b-41d4-a716-446655440000', 'インフラチームBeta', 'サーバー・インフラ管理を担当するチームです。', 'TENANT_TYPE_TEAM', NOW(), NOW());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'Seed Rollback: Delete tenants';

DELETE FROM tenants WHERE id IN (
    '10000000-0000-0000-0000-000000000001',
    '10000000-0000-0000-0000-000000000002',
    '10000000-0000-0000-0000-000000000003',
    '10000000-0000-0000-0000-000000000004',
    '10000000-0000-0000-0000-000000000005'
);

-- +goose StatementEnd
