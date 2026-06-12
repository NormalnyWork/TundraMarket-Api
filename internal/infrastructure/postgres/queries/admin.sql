-- internal/infrastructure/postgres/queries/admin.sql

-- name: GetAdminByLogin :one
SELECT * FROM admin WHERE login = $1;

-- name: UpdateAdminPassword :exec
UPDATE admin SET password = $2 WHERE login = $1;