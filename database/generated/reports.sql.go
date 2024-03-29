// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: reports.sql

package generated

import (
	"context"
	"time"
)

const getAllFeedbacks = `-- name: GetAllFeedbacks :many
SELECT id, user_id, type, submitted_at, content, hidden FROM "Feedbacks" WHERE "Feedbacks"."hidden" = $1
`

func (q *Queries) GetAllFeedbacks(ctx context.Context, hidden bool) ([]Feedback, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeedbacks, hidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feedback
	for rows.Next() {
		var i Feedback
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.SubmittedAt,
			&i.Content,
			&i.Hidden,
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

const getAllFeedbacksInRange = `-- name: GetAllFeedbacksInRange :many
SELECT id, user_id, type, submitted_at, content, hidden FROM "Feedbacks"
WHERE submitted_at > $1 and submitted_at < $2 and "Feedbacks"."hidden" = $3
`

type GetAllFeedbacksInRangeParams struct {
	SubmittedAt   time.Time `json:"submitted_at"`
	SubmittedAt_2 time.Time `json:"submitted_at_2"`
	Hidden        bool      `json:"hidden"`
}

func (q *Queries) GetAllFeedbacksInRange(ctx context.Context, arg GetAllFeedbacksInRangeParams) ([]Feedback, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeedbacksInRange, arg.SubmittedAt, arg.SubmittedAt_2, arg.Hidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feedback
	for rows.Next() {
		var i Feedback
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.SubmittedAt,
			&i.Content,
			&i.Hidden,
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

const getAllFeedbacksOfType = `-- name: GetAllFeedbacksOfType :many
SELECT id, user_id, type, submitted_at, content, hidden FROM "Feedbacks"
WHERE "Feedbacks"."type" = $1 and "Feedbacks"."hidden" = $2
`

type GetAllFeedbacksOfTypeParams struct {
	Type   string `json:"type"`
	Hidden bool   `json:"hidden"`
}

func (q *Queries) GetAllFeedbacksOfType(ctx context.Context, arg GetAllFeedbacksOfTypeParams) ([]Feedback, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeedbacksOfType, arg.Type, arg.Hidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feedback
	for rows.Next() {
		var i Feedback
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.SubmittedAt,
			&i.Content,
			&i.Hidden,
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

const getFeedbacksInRangeOfType = `-- name: GetFeedbacksInRangeOfType :many
SELECT id, user_id, type, submitted_at, content, hidden FROM "Feedbacks"
WHERE "Feedbacks"."type" = $1 and submitted_at > $2 and submitted_at < $3 and "Feedbacks"."hidden" = $4
`

type GetFeedbacksInRangeOfTypeParams struct {
	Type          string    `json:"type"`
	SubmittedAt   time.Time `json:"submitted_at"`
	SubmittedAt_2 time.Time `json:"submitted_at_2"`
	Hidden        bool      `json:"hidden"`
}

func (q *Queries) GetFeedbacksInRangeOfType(ctx context.Context, arg GetFeedbacksInRangeOfTypeParams) ([]Feedback, error) {
	rows, err := q.db.QueryContext(ctx, getFeedbacksInRangeOfType,
		arg.Type,
		arg.SubmittedAt,
		arg.SubmittedAt_2,
		arg.Hidden,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feedback
	for rows.Next() {
		var i Feedback
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.SubmittedAt,
			&i.Content,
			&i.Hidden,
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

const getNewestUser = `-- name: GetNewestUser :one
SELECT id, username, created_at from "Users" order by created_at desc LIMIT 1
`

type GetNewestUserRow struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) GetNewestUser(ctx context.Context) (GetNewestUserRow, error) {
	row := q.db.QueryRowContext(ctx, getNewestUser)
	var i GetNewestUserRow
	err := row.Scan(&i.ID, &i.Username, &i.CreatedAt)
	return i, err
}

const getUsersCount = `-- name: GetUsersCount :one
SELECT count(id) from "Users"
`

func (q *Queries) GetUsersCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getUsersCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const setFeedbackHiddenTo = `-- name: SetFeedbackHiddenTo :exec
UPDATE "Feedbacks"
SET "hidden" = $2
WHERE id = $1
`

type SetFeedbackHiddenToParams struct {
	ID     int64 `json:"id"`
	Hidden bool  `json:"hidden"`
}

func (q *Queries) SetFeedbackHiddenTo(ctx context.Context, arg SetFeedbackHiddenToParams) error {
	_, err := q.db.ExecContext(ctx, setFeedbackHiddenTo, arg.ID, arg.Hidden)
	return err
}
