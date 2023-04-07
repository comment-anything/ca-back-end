// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: bans.sql

package generated

import (
	"context"
	"database/sql"
)

const addBanRecord = `-- name: AddBanRecord :exec
INSERT INTO "BanActions"
(taken_by, target_user, reason, domain, set_banned_to)
VALUES
($1, $2, $3, $4, $5)
`

type AddBanRecordParams struct {
	TakenBy     int64          `json:"taken_by"`
	TargetUser  int64          `json:"target_user"`
	Reason      sql.NullString `json:"reason"`
	Domain      sql.NullString `json:"domain"`
	SetBannedTo sql.NullBool   `json:"set_banned_to"`
}

func (q *Queries) AddBanRecord(ctx context.Context, arg AddBanRecordParams) error {
	_, err := q.db.ExecContext(ctx, addBanRecord,
		arg.TakenBy,
		arg.TargetUser,
		arg.Reason,
		arg.Domain,
		arg.SetBannedTo,
	)
	return err
}

const banUserFromDomain = `-- name: BanUserFromDomain :exec
INSERT INTO "DomainBans"
(user_id, banned_from, banned_by)
VALUES ($1, $2, $3)
`

type BanUserFromDomainParams struct {
	UserID     int64  `json:"user_id"`
	BannedFrom string `json:"banned_from"`
	BannedBy   int64  `json:"banned_by"`
}

func (q *Queries) BanUserFromDomain(ctx context.Context, arg BanUserFromDomainParams) error {
	_, err := q.db.ExecContext(ctx, banUserFromDomain, arg.UserID, arg.BannedFrom, arg.BannedBy)
	return err
}

const banUserGlobally = `-- name: BanUserGlobally :exec
UPDATE "Users"
SET banned = true
WHERE id = $1
`

func (q *Queries) BanUserGlobally(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, banUserGlobally, id)
	return err
}

const getDomainBans = `-- name: GetDomainBans :many
SELECT user_id, banned_from, banned_by, banned_at from "DomainBans"
WHERE user_id = $1
`

func (q *Queries) GetDomainBans(ctx context.Context, userID int64) ([]DomainBan, error) {
	rows, err := q.db.QueryContext(ctx, getDomainBans, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DomainBan
	for rows.Next() {
		var i DomainBan
		if err := rows.Scan(
			&i.UserID,
			&i.BannedFrom,
			&i.BannedBy,
			&i.BannedAt,
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

const getUserBanStatus = `-- name: GetUserBanStatus :one
SELECT banned from "Users"
WHERE id = $1
`

func (q *Queries) GetUserBanStatus(ctx context.Context, id int64) (bool, error) {
	row := q.db.QueryRowContext(ctx, getUserBanStatus, id)
	var banned bool
	err := row.Scan(&banned)
	return banned, err
}

const unbanUserFromDomain = `-- name: UnbanUserFromDomain :exec
DELETE FROM "DomainBans"
WHERE
user_id = $1 AND banned_from = $2
`

type UnbanUserFromDomainParams struct {
	UserID     int64  `json:"user_id"`
	BannedFrom string `json:"banned_from"`
}

func (q *Queries) UnbanUserFromDomain(ctx context.Context, arg UnbanUserFromDomainParams) error {
	_, err := q.db.ExecContext(ctx, unbanUserFromDomain, arg.UserID, arg.BannedFrom)
	return err
}

const unbanUserGlobally = `-- name: UnbanUserGlobally :exec
UPDATE "Users"
SET banned = false
WHERE id = $1
`

func (q *Queries) UnbanUserGlobally(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, unbanUserGlobally, id)
	return err
}
