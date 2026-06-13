-- internal/infrastructure/postgres/queries/admin.sql

-- name: GetAdminByLogin :one
SELECT * FROM admin WHERE login = $1;

-- name: UpdateAdminPasswordHash :exec
UPDATE admin SET password_hash = $2 WHERE login = $1;
