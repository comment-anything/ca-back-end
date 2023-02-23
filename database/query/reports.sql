
-- name: GetAllFeedbacks :many
SELECT * FROM "Feedbacks" WHERE "Feedbacks"."hidden" = $1;

-- name: GetAllFeedbacksInRange :many
SELECT * FROM "Feedbacks"
WHERE submitted_at > $1 and submitted_at < $2 and "Feedbacks"."hidden" = $3;

-- name: GetFeedbacksInRangeOfType :many
SELECT * FROM "Feedbacks"
WHERE "Feedbacks"."type" = $1 and submitted_at > $2 and submitted_at < $3 and "Feedbacks"."hidden" = $4;

-- name: GetAllFeedbacksOfType :many
SELECT * FROM "Feedbacks"
WHERE "Feedbacks"."type" = $1 and "Feedbacks"."hidden" = $2;

-- name: GetUsersCount :one
SELECT count(id) from "Users";

-- name: GetNewestUser :one
SELECT id, username, created_at from "Users" order by created_at desc LIMIT 1;