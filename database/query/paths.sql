-- name: GetPath :one
SELECT * FROM "Paths"
WHERE domain = $1 AND path = $2 LIMIT 1;

-- name: CreatePath :one
INSERT INTO "Paths"
( "domain", "path" )
VALUES
( $1 , $2)
RETURNING *;

-- name: GetPathById :one
SELECT * FROM "Paths"
WHERE id = $1 LIMIT 1;

