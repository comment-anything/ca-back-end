-- name: GetAdminAssignment :many
Select * from "AdminAssignments" where assigned_to = $1;