-- name: GetModActionsInRange :many
SELECT 
"MA".id, "MA".taken_by, "MA".comment_id, "MA".reason as "ActionReason", "MA".taken_on, 
"MA".set_hidden_to, "MA".set_removed_to, 

"MA".associated_report,
"CR".reporting_user,
"CR".reason as "ReportReason",
"CR".time_created "ReportCreated"
FROM 
(SELECT * FROM "CommentModerationActions" 
WHERE taken_on > $1 AND taken_on <$2)
as "MA"
LEFT JOIN
"CommentReports" as "CR"
ON
"MA".associated_report = "CR".id
;