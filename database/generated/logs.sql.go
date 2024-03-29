// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: logs.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const createLog = `-- name: CreateLog :one
INSERT INTO "Logs" ("ip", "url") VALUES ($1, $2) RETURNING id, user_id, ip, url, at_time
`

type CreateLogParams struct {
	Ip  sql.NullString `json:"ip"`
	Url sql.NullString `json:"url"`
}

func (q *Queries) CreateLog(ctx context.Context, arg CreateLogParams) (Log, error) {
	row := q.db.QueryRowContext(ctx, createLog, arg.Ip, arg.Url)
	var i Log
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Ip,
		&i.Url,
		&i.AtTime,
	)
	return i, err
}

const getLogsInRange = `-- name: GetLogsInRange :many
SELECT id, user_id, ip, url, at_time FROM "Logs"
WHERE
    at_time > $1
AND 
    at_time < $2
`

type GetLogsInRangeParams struct {
	AtTime   time.Time `json:"at_time"`
	AtTime_2 time.Time `json:"at_time_2"`
}

func (q *Queries) GetLogsInRange(ctx context.Context, arg GetLogsInRangeParams) ([]Log, error) {
	rows, err := q.db.QueryContext(ctx, getLogsInRange, arg.AtTime, arg.AtTime_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Log
	for rows.Next() {
		var i Log
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Ip,
			&i.Url,
			&i.AtTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateLogUser = `-- name: UpdateLogUser :exec
UPDATE "Logs" SET 
 "user_id" = $2 
WHERE "id" = $1
`

type UpdateLogUserParams struct {
	ID     int64         `json:"id"`
	UserID sql.NullInt64 `json:"user_id"`
}

func (q *Queries) UpdateLogUser(ctx context.Context, arg UpdateLogUserParams) error {
	_, err := q.db.ExecContext(ctx, updateLogUser, arg.ID, arg.UserID)
	return err
}
