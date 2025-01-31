-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, created_at, updated_at, expires_at, revoked_at)
VALUES (
    $1,
    $2,
    NOW(),
    NOW(),
    $3,
    NULL
)
RETURNING *;



-- name: GetUserFromRefreshToken :one
SELECT * FROM users 
WHERE id =(
    SELECT user_id 
    FROM refresh_tokens 
    WHERE token = $1 
    AND expires_at >NOW()
    AND revoked_at IS NULL
);

