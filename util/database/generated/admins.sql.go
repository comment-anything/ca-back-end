// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: admins.sql

package generated

import (
	"context"
)

const assignAdmin = `-- name: AssignAdmin :exec
INSERT INTO "AdminAssignments" (assigned_to, assigned_by, is_deactivation)
VALUES ($1, $2, false)
`

type AssignAdminParams struct {
	AssignedTo int64 `json:"assigned_to"`
	AssignedBy int64 `json:"assigned_by"`
}

func (q *Queries) AssignAdmin(ctx context.Context, arg AssignAdminParams) error {
	_, err := q.db.ExecContext(ctx, assignAdmin, arg.AssignedTo, arg.AssignedBy)
	return err
}

const getAdminAssignment = `-- name: GetAdminAssignment :many
Select id, assigned_to, assigned_by, assigned_at, is_deactivation from "AdminAssignments" where assigned_to = $1 ORDER BY assigned_at DESC
`

func (q *Queries) GetAdminAssignment(ctx context.Context, assignedTo int64) ([]AdminAssignment, error) {
	rows, err := q.db.QueryContext(ctx, getAdminAssignment, assignedTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AdminAssignment
	for rows.Next() {
		var i AdminAssignment
		if err := rows.Scan(
			&i.ID,
			&i.AssignedTo,
			&i.AssignedBy,
			&i.AssignedAt,
			&i.IsDeactivation,
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