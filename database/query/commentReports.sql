-- name: GetAllCommentReports :many
select "CR".id, "P".domain, "CR".reporting_user, "CR".comment, "CR".reason, "CR".action_taken, "CR".time_created
from "CommentReports" as "CR"
inner join
(select id, path_id from "Comments") as "C"
on "CR".comment = "C".id
inner join
(select id, domain from "Paths") as "P"
on "C".path_id = "P".id
;

-- name: GetCommentReportsForDomain :many
select "CR".id, "P".domain, "CR".reporting_user, "CR".comment, "CR".reason, "CR".action_taken, "CR".time_created
from "CommentReports" as "CR"
inner join
(select id, path_id from "Comments") as "C"
on "CR".comment = "C".id
inner join
(select id, domain from "Paths" where domain = $1) as "P"
on "C".path_id = "P".id
;