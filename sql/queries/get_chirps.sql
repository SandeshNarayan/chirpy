-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;


-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;


-- name: DeleteChirpByID :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2;