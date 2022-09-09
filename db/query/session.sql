-- name: CreateSession :one
INSERT INTO sessions(
    id,
    username,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    expires_at
) VALUES(
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: ToggleBlockSession :exec
UPDATE sessions
SET is_blocked = $2
WHERE id = $1;

-- name: ExpireSession :exec
UPDATE sessions
SET expires_at = $2
WHERE id = $1;

-- name: UpdateSession :one
UPDATE sessions
SET
    is_blocked = CASE WHEN @is_blocked_to_update::boolean
    THEN @is_blocked::bool ELSE is_blocked END,

    expires_at = CASE WHEN @expires_at_to_update::boolean
    THEN @expires_at::timestamp ELSE expires_at END
WHERE 
    id = @id
RETURNING *;

