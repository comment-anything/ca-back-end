// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: admins.sql

package generated

import (
	"context"
)

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