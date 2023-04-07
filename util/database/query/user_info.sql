
-- name: GetUsername :one
SELECT username from "Users" where id=$1;