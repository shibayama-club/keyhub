-- name: CreateConsoleSession :exec
INSERT INTO console_sessions (
    session_id,
    organization_id,
    created_at,
    expires_at
) VALUES (
    $1, $2, NOW(), NOW() + INTERVAL '24 hours'
);

-- name: GetConsoleSession :one
SELECT sqlc.embed(cs)
FROM console_sessions cs
WHERE cs.session_id = $1
AND cs.expires_at > NOW();

-- name: DeleteConsoleSession :exec
DELETE FROM console_sessions
WHERE session_id = $1;

-- name: CleanupExpiredConsoleSessions :exec
DELETE FROM console_sessions
WHERE expires_at < NOW();