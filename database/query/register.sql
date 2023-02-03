-- name: CreateUser :one
INSERT INTO "Users" (
    username,
    password,
    email
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserByUserName :one
SELECT * FROM "Users"
WHERE "username" = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM "Users"
WHERE "email" = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM "Users"
WHERE "id" = $1 LIMIT 1;

-- name: Tst_DeleteUser :exec
DELETE FROM "Users"
WHERE id = $1;