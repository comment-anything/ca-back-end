-- name: CreateLog :one
INSERT INTO "Logs" ("ip", "url") VALUES ($1, $2) RETURNING 1;

-- name: UpdateLogUser :exec
UPDATE "Logs" SET 
 "user_id" = $2 
WHERE "id" = $1;

-- name: GetLogsInRange :many
SELECT * FROM "Logs"
WHERE
    at_time > $1
AND 
    at_time < $2;