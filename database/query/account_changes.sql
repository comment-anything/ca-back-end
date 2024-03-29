-- name: UpdateUserEmail :exec
UPDATE "Users"
SET email = $2
WHERE id = $1;

-- name: UpdateUserProfileBlurb :exec
UPDATE "Users"
SET profile_blurb = $2
WHERE id = $1;

-- name: CreatePWResetCode :one
INSERT INTO "PasswordResetCodes" (
    user_id, id
) VALUES (
    $1, $2
) RETURNING *;

-- name: CreateVerificationCode :one
INSERT INTO "VerificationCodes" (
    user_id, id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetPWResetCodeEntry :one
SELECT * FROM "PasswordResetCodes"
WHERE "id" = $1 LIMIT 1;

-- name: GetVerificationCodeEntry :one
SELECT * FROM "VerificationCodes"
WHERE "id" = $1 LIMIT 1;

-- name: DeletePreviousPWRestCodesForUser :exec
DELETE FROM "PasswordResetCodes"
WHERE user_id = $1;

-- name: DeletePreviousVerificationCodesForUser :exec
DELETE FROM "VerificationCodes"
WHERE user_id = $1;

-- name: UpdateVerification :exec
UPDATE "Users"
SET "is_verified" = $2
WHERE "id" = $1;

-- name: SetNewUserPassword :exec
UPDATE "Users"
SET password = $1
WHERE email = $2;