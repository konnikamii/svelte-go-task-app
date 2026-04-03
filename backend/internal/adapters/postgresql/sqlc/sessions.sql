-- name: CreateSession :one
INSERT INTO user_sessions (user_id, device_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id,
  user_id,
  device_id,
  token_hash,
  expires_at,
  revoked_at,
  created_at,
  updated_at;

-- name: GetActiveSessionByTokenHash :one
SELECT id,
  user_id,
  device_id,
  token_hash,
  expires_at,
  revoked_at,
  created_at,
  updated_at
FROM user_sessions
WHERE token_hash = $1
  AND revoked_at IS NULL
  AND expires_at > NOW();

-- name: RevokeSessionByTokenHash :execrows
UPDATE user_sessions
SET revoked_at = NOW(),
  updated_at = CURRENT_TIMESTAMP
WHERE token_hash = $1
  AND revoked_at IS NULL;

-- name: RevokeSessionsByUserAndDevice :execrows
UPDATE user_sessions
SET revoked_at = NOW(),
  updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
  AND device_id = $2
  AND revoked_at IS NULL;

-- name: DeleteStaleSessions :execrows
DELETE FROM user_sessions
WHERE expires_at < NOW()
  OR (
    revoked_at IS NOT NULL
    AND revoked_at < NOW() - INTERVAL '30 days'
  );