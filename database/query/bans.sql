-- name: GetDomainBans :many
SELECT * from "DomainBans"
WHERE user_id = $1;

-- name: BanUserFromDomain :exec
INSERT INTO "DomainBans"
(user_id, banned_from, banned_by)
VALUES ($1, $2, $3);

-- name: UnbanUserFromDomain :exec
DELETE FROM "DomainBans"
WHERE
user_id = $1 AND banned_from = $2;

-- name: AddBanRecord :exec
INSERT INTO "BanActions"
(taken_by, target_user, reason, domain, set_banned_to)
VALUES
($1, $2, $3, $4, $5);

-- name: GetUserBanStatus :one
SELECT banned from "Users"
WHERE id = $1;

-- name: BanUserGlobally :exec
UPDATE "Users"
SET banned = true
WHERE id = $1;

-- name: UnbanUserGlobally :exec
UPDATE "Users"
SET banned = false
WHERE id = $1;
;

