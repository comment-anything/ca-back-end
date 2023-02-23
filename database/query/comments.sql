-- name: GetCommentsForPath :many
SELECT * FROM "Comments"
WHERE path_id = $1 ORDER BY id;

-- name: GetUpVotesForComment :many
SELECT category, SUM(value) FROM "VoteRecords"
WHERE comment_id = $1 and value > 0 GROUP BY category;

-- name: GetDownVotesForComment :many
SELECT category, SUM(value) FROM "VoteRecords"
WHERE comment_id = $1 and value < 0 GROUP BY category;

-- name: GetVotesForCommentAndCategoryByUser :one
SELECT * FROM "VoteRecords"
WHERE comment_id = $1 and user_id = $2 and category = $3;

-- name: CreateCommentVote :exec
INSERT INTO "VoteRecords"
("comment_id", "user_id", "category", "value")
VALUES
($1, $2, $3, $4);

-- name: DeleteCommentVote :exec
DELETE FROM "VoteRecords"
WHERE comment_id = $1
AND category = $2
AND user_id = $3;

-- name: UpdateCommentVote :exec
UPDATE "VoteRecords"
SET value = $1
WHERE comment_id = $1
AND category = $2
AND user_id = $3;

-- name: CreateComment :one
INSERT INTO "Comments"
("path_id", "author", "content", "parent")
VALUES
( $1, $2, $3, $4) RETURNING *;

-- name: GetCommentByID :one
SELECT * from "Comments"
WHERE id = $1;

