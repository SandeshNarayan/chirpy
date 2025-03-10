// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: create_chirrp.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirps(id, body, user_id, created_at, updated_at) 
VALUES (
    gen_random_uuid(),
    $1,
    $2, 
    NOW(), 
    NOW()
    )
RETURNING id, created_at, updated_at, body, user_id
`

type CreateChirpParams struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.Body, arg.UserID)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}
