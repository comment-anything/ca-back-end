-- name: CreateFeedback :exec
INSERT INTO "Feedbacks" ("user_id", "type", "content") VALUES ($1, $2, $3);