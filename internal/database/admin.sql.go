// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: admin.sql

package database

import (
	"context"
)

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM "user"
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}
