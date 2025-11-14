-- name: SaveOAuthState :exec
INSERT INTO oauth_states (
    state,
    code_verifier,
    nonce,
    created_at
) VALUES (
    $1, $2, $3, NOW()
);

-- name: GetOAuthState :one
SELECT sqlc.embed(os)
FROM oauth_states os
WHERE os.state = $1
AND os.consumed_at IS NULL;

-- name: ConsumeOAuthState :exec
UPDATE oauth_states
SET consumed_at = NOW()
WHERE state = $1;

-- name: CleanupExpiredOAuthStates :exec
DELETE FROM oauth_states
WHERE created_at < NOW() - INTERVAL '10 minutes';
