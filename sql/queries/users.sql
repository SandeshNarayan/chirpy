-- name: CreateUser :one
INSERT INTO "user" (id, created_at, updated_at, email)
Values (
    gen_random_uuid(),
    NOW(),
    NOW(),
     $1
     )
RETURNING *;
