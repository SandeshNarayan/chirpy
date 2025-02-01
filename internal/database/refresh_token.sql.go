// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: refresh_token.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, created_at, updated_at, expires_at, revoked_at)
VALUES (
    $1,
    $2,
    NOW(),
    NOW(),
    $3,
    NULL
)
RETURNING token, user_id, created_at, updated_at, expires_at, revoked_at
`

type CreateRefreshTokenParams struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, arg.Token, arg.UserID, arg.ExpiresAt)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getUserFromToken = `-- name: GetUserFromToken :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users 
WHERE id =(
    SELECT user_id 
    FROM refresh_tokens 
    WHERE token = $1 
    AND expires_at > NOW()
    AND revoked_at IS NULL
)
`

func (q *Queries) GetUserFromToken(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromToken, token)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}
