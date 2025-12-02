-- +goose Up
-- +goose StatementBegin
SELECT 'Seed: Insert tenant_join_codes';

-- code: 6-20 alphanumeric characters
-- max_uses: 0 means unlimited
-- expires_at: NULL means never expires

INSERT INTO tenant_join_codes (id, tenant_id, code, expires_at, max_uses, used_count, created_at) VALUES
    ('30000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'ALPHA001', NOW() + INTERVAL '30 days', 10, 2, NOW()),
    ('30000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', 'SOUMU002', NOW() + INTERVAL '60 days', 20, 1, NOW()),
    ('30000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000003', 'KEYHUB03', NULL, 0, 3, NOW()),
    ('30000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000004', 'AILAB004', NOW() + INTERVAL '90 days', 5, 0, NOW()),
    ('30000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000005', 'INFRA005', NOW() + INTERVAL '45 days', 15, 5, NOW());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'Seed Rollback: Delete tenant_join_codes';

DELETE FROM tenant_join_codes WHERE id IN (
    '30000000-0000-0000-0000-000000000001',
    '30000000-0000-0000-0000-000000000002',
    '30000000-0000-0000-0000-000000000003',
    '30000000-0000-0000-0000-000000000004',
    '30000000-0000-0000-0000-000000000005'
);

-- +goose StatementEnd
