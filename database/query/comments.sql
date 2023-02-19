-- name: GetCommentsForPath :many
SELECT * FROM "Comments"
WHERE path_id = $1 ORDER BY id;

-- name: GetUpVotesForComment :many
SELECT category, SUM(value) FROM "VoteRecords"
WHERE comment_id = $1 and value > 0 GROUP BY category;

-- name: GetDownVotesForComment :many
SELECT category, SUM(value) FROM "VoteRecords"
WHERE comment_id = $1 and value < 0 GROUP BY category;

-- name: GetVotesForCommentByUser :many
SELECT category FROM "VoteRecords"
WHERE comment_id = $1 and user_id = $2;

-- name: CreateComment :one
INSERT INTO "Comments"
("path_id", "author", "content", "parent")
VALUES
( $1, $2, $3, $4) RETURNING *;

-- name: GetCommentByID :one
SELECT * from "Comments"
WHERE id = $1;

