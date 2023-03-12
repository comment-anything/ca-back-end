-- name: GetAdminAssignment :many
Select * from "AdminAssignments" where assigned_to = $1 ORDER BY assigned_at DESC;

-- name: AssignAdmin :exec
INSERT INTO "AdminAssignments" (assigned_to, assigned_by, is_deactivation)
VALUES ($1, $2, false);