-- name: GetNomadByPhone :one
SELECT * FROM nomad WHERE phone = $1;

-- name: CreateNomad :one
INSERT INTO nomad (phone)
VALUES ($1)
RETURNING *;

-- name: GetNomadByID :one
SELECT * FROM nomad WHERE id = $1;