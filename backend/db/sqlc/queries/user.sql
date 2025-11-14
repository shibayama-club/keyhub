-- name: GetUser :one
SELECT sqlc.embed(u)
FROM users u
WHERE u.id = $1;

-- name: GetUserByProviderIdentity :one
SELECT sqlc.embed(u)
FROM users u
INNER JOIN user_identities ui ON u.id = ui.user_id
WHERE ui.provider = $1 AND ui.provider_sub = $2;

-- name: UpsertUser :one
INSERT INTO users (
    email,
    name,
    icon
) VALUES (
    $1, $2, $3
)
ON CONFLICT (email)
DO UPDATE SET
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    updated_at = NOW()
RETURNING sqlc.embed(users);

-- name: UpsertUserIdentity :exec
INSERT INTO user_identities (
    user_id,
    provider,
    provider_sub
) VALUES (
    $1, $2, $3
)
ON CONFLICT (provider, provider_sub)
DO UPDATE SET
    user_id = EXCLUDED.user_id,
    updated_at = NOW();
