-- name: GetGlobalModeratorAssignments :many
SELECT * from "GlobalModeratorAssignments" where assigned_to = $1 ORDER BY assigned_at DESC;

-- name: GetDomainModeratorAssignments :many
select distinct on (domain) domain, is_deactivation, id
 from "DomainModeratorAssignments"
 WHERE assigned_to=$1
 order by domain, assigned_at DESC, id;
