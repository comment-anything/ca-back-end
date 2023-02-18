-- name: GetCommentsForPath :many
SELECT * FROM "Comments"
WHERE id = $1 ORDER BY id;

-- name: GetUpVotesForComment :many
SELECT category, SUM(value) FROM "VoteRecords"
WHERE comment_id = $1 and value > 0 GROUP BY category;

-- name: GetDownVotesForComment :many
SELECT category, SUM(value) FROM "VoteRecords"
WHERE comment_id = $1 and value < 0 GROUP BY category;

-- name: GetVotesForCommentByUser :many
SELECT category FROM "VoteRecords"
WHERE comment_id = $1 and user_id = $2;



