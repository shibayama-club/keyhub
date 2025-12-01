-- +goose Up
-- +goose StatementBegin
SELECT 'Seed: Insert user_identities';

INSERT INTO user_identities (id, user_id, provider, provider_sub, created_at, updated_at) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'google', 'google-sub-111111', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222', 'google', 'google-sub-222222', NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '33333333-3333-3333-3333-333333333333', 'google', 'google-sub-333333', NOW(), NOW()),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', '44444444-4444-4444-4444-444444444444', 'google', 'google-sub-444444', NOW(), NOW()),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '55555555-5555-5555-5555-555555555555', 'google', 'google-sub-555555', NOW(), NOW());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'Seed Rollback: Delete user_identities';

DELETE FROM user_identities WHERE id IN (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee'
);

-- +goose StatementEnd
