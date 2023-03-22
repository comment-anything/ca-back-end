
-- name: CreateModActionRecord :exec
INSERT INTO "CommentModerationActions" (
    taken_by,
    comment_id,
    reason,
    set_hidden_to,
    set_removed_to,
    associated_report
) VALUES (
    $1,$2,$3,$4,$5,$6
);

-- name: ModerateComment :exec
UPDATE "Comments"
SET "hidden" = $2, "removed" = $3
WHERE "id" = $1;


-- name: GetModerationRecords :many
SELECT * FROM "CommentModerationActions"
WHERE taken_on > $1 AND taken_on < $2;