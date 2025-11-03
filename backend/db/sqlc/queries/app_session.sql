-- name: CreateAppSession :one
INSERT INTO sessions (
    session_id,
    user_id,
    active_membership_id,
    created_at,
    expires_at,
    csrf_token,
    revoked
) VALUES (
    $1, $2, $3, $4, $5, $6, FALSE
) RETURNING sqlc.embed(sessions);

-- name: GetAppSession :one
SELECT sqlc.embed(s)
FROM sessions s
WHERE s.session_id = $1
AND s.revoked = FALSE;

-- name: RevokeAppSession :exec
UPDATE sessions
SET revoked = TRUE
WHERE session_id = $1;

-- name: CleanupExpiredAppSessions :exec
-- 期限切れまたは無効化されたセッションを物理削除する（バッチ処理用）
DELETE FROM sessions
WHERE expires_at < NOW() OR revoked = TRUE;
