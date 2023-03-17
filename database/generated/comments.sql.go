// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: comments.sql

package generated

import (
	"context"
	"database/sql"
)

const createComment = `-- name: CreateComment :one
INSERT INTO "Comments"
("path_id", "author", "content", "parent")
VALUES
( $1, $2, $3, $4) RETURNING id, path_id, author, content, created_at, parent, hidden, removed
`

type CreateCommentParams struct {
	PathID  int64         `json:"path_id"`
	Author  int64         `json:"author"`
	Content string        `json:"content"`
	Parent  sql.NullInt64 `json:"parent"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment,
		arg.PathID,
		arg.Author,
		arg.Content,
		arg.Parent,
	)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.PathID,
		&i.Author,
		&i.Content,
		&i.CreatedAt,
		&i.Parent,
		&i.Hidden,
		&i.Removed,
	)
	return i, err
}

const createCommentReport = `-- name: CreateCommentReport :exec
INSERT INTO "CommentReports"
("reporting_user", "comment", "reason", "action_taken")
VALUES
($1, $2, $3, false)
`

type CreateCommentReportParams struct {
	ReportingUser int64          `json:"reporting_user"`
	Comment       int64          `json:"comment"`
	Reason        sql.NullString `json:"reason"`
}

func (q *Queries) CreateCommentReport(ctx context.Context, arg CreateCommentReportParams) error {
	_, err := q.db.ExecContext(ctx, createCommentReport, arg.ReportingUser, arg.Comment, arg.Reason)
	return err
}

const createCommentVote = `-- name: CreateCommentVote :exec
INSERT INTO "VoteRecords"
("comment_id", "user_id", "category", "value")
VALUES
($1, $2, $3, $4)
`

type CreateCommentVoteParams struct {
	CommentID int64  `json:"comment_id"`
	UserID    int64  `json:"user_id"`
	Category  string `json:"category"`
	Value     int16  `json:"value"`
}

func (q *Queries) CreateCommentVote(ctx context.Context, arg CreateCommentVoteParams) error {
	_, err := q.db.ExecContext(ctx, createCommentVote,
		arg.CommentID,
		arg.UserID,
		arg.Category,
		arg.Value,
	)
	return err
}

const deleteCommentVote = `-- name: DeleteCommentVote :exec
DELETE FROM "VoteRecords"
WHERE comment_id = $1
AND category = $2
AND user_id = $3
`

type DeleteCommentVoteParams struct {
	CommentID int64  `json:"comment_id"`
	Category  string `json:"category"`
	UserID    int64  `json:"user_id"`
}

func (q *Queries) DeleteCommentVote(ctx context.Context, arg DeleteCommentVoteParams) error {
	_, err := q.db.ExecContext(ctx, deleteCommentVote, arg.CommentID, arg.Category, arg.UserID)
	return err
}

const getCommentByID = `-- name: GetCommentByID :one
SELECT id, path_id, author, content, created_at, parent, hidden, removed from "Comments"
WHERE id = $1
`

func (q *Queries) GetCommentByID(ctx context.Context, id int64) (Comment, error) {
	row := q.db.QueryRowContext(ctx, getCommentByID, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.PathID,
		&i.Author,
		&i.Content,
		&i.CreatedAt,
		&i.Parent,
		&i.Hidden,
		&i.Removed,
	)
	return i, err
}

const getCommentsForPath = `-- name: GetCommentsForPath :many
SELECT id, path_id, author, content, created_at, parent, hidden, removed FROM "Comments"
WHERE path_id = $1 ORDER BY id
`

func (q *Queries) GetCommentsForPath(ctx context.Context, pathID int64) ([]Comment, error) {
	rows, err := q.db.QueryContext(ctx, getCommentsForPath, pathID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.ID,
			&i.PathID,
			&i.Author,
			&i.Content,
			&i.CreatedAt,
			&i.Parent,
			&i.Hidden,
			&i.Removed,
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

const getDownVotesForComment = `-- name: GetDownVotesForComment :many
SELECT category, SUM(value) FROM "VoteRecords"
WHERE comment_id = $1 and value < 0 GROUP BY category
`

type GetDownVotesForCommentRow struct {
	Category string `json:"category"`
	Sum      int64  `json:"sum"`
}

func (q *Queries) GetDownVotesForComment(ctx context.Context, commentID int64) ([]GetDownVotesForCommentRow, error) {
	rows, err := q.db.QueryContext(ctx, getDownVotesForComment, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDownVotesForCommentRow
	for rows.Next() {
		var i GetDownVotesForCommentRow
		if err := rows.Scan(&i.Category, &i.Sum); err != nil {
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

const getUpVotesForComment = `-- name: GetUpVotesForComment :many
SELECT category, SUM(value) FROM "VoteRecords"
WHERE comment_id = $1 and value > 0 GROUP BY category
`

type GetUpVotesForCommentRow struct {
	Category string `json:"category"`
	Sum      int64  `json:"sum"`
}

func (q *Queries) GetUpVotesForComment(ctx context.Context, commentID int64) ([]GetUpVotesForCommentRow, error) {
	rows, err := q.db.QueryContext(ctx, getUpVotesForComment, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUpVotesForCommentRow
	for rows.Next() {
		var i GetUpVotesForCommentRow
		if err := rows.Scan(&i.Category, &i.Sum); err != nil {
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

const getVotesForCommentAndCategoryByUser = `-- name: GetVotesForCommentAndCategoryByUser :one
SELECT comment_id, category, user_id, value FROM "VoteRecords"
WHERE comment_id = $1 and user_id = $2 and category = $3
`

type GetVotesForCommentAndCategoryByUserParams struct {
	CommentID int64  `json:"comment_id"`
	UserID    int64  `json:"user_id"`
	Category  string `json:"category"`
}

func (q *Queries) GetVotesForCommentAndCategoryByUser(ctx context.Context, arg GetVotesForCommentAndCategoryByUserParams) (VoteRecord, error) {
	row := q.db.QueryRowContext(ctx, getVotesForCommentAndCategoryByUser, arg.CommentID, arg.UserID, arg.Category)
	var i VoteRecord
	err := row.Scan(
		&i.CommentID,
		&i.Category,
		&i.UserID,
		&i.Value,
	)
	return i, err
}

const updateCommentVote = `-- name: UpdateCommentVote :exec
UPDATE "VoteRecords"
SET value = $1
WHERE comment_id = $1
AND category = $2
AND user_id = $3
`

type UpdateCommentVoteParams struct {
	Value    int16  `json:"value"`
	Category string `json:"category"`
	UserID   int64  `json:"user_id"`
}

func (q *Queries) UpdateCommentVote(ctx context.Context, arg UpdateCommentVoteParams) error {
	_, err := q.db.ExecContext(ctx, updateCommentVote, arg.Value, arg.Category, arg.UserID)
	return err
}
