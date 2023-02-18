-- name: EnsureDomainRecordExits :exec
INSERT INTO "Domains" ("id") VALUES ($1);