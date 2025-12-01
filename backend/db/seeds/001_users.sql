-- +goose Up
-- +goose StatementBegin
SELECT 'Seed: Insert users';

INSERT INTO users (id, email, name, icon, created_at, updated_at) VALUES
    ('11111111-1111-1111-1111-111111111111', 'yamada.taro@example.com', '山田太郎', 'https://example.com/icons/yamada.png', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'suzuki.hanako@example.com', '鈴木花子', 'https://example.com/icons/suzuki.png', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', 'sato.ichiro@example.com', '佐藤一郎', 'https://example.com/icons/sato.png', NOW(), NOW()),
    ('44444444-4444-4444-4444-444444444444', 'tanaka.yuki@example.com', '田中雪', 'https://example.com/icons/tanaka.png', NOW(), NOW()),
    ('55555555-5555-5555-5555-555555555555', 'watanabe.ken@example.com', '渡辺健', 'https://example.com/icons/watanabe.png', NOW(), NOW());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'Seed Rollback: Delete users';

DELETE FROM users WHERE id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    '44444444-4444-4444-4444-444444444444',
    '55555555-5555-5555-5555-555555555555'
);

-- +goose StatementEnd
