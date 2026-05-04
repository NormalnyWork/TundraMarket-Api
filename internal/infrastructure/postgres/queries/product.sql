-- name: GetAllProducts :many
SELECT id, name, details, price, weight, volume, created_at FROM product ORDER BY id;

-- name: GetProductByID :one
SELECT id, name, details, price, weight, volume, created_at FROM product WHERE id = $1;

-- name: GetProductsByIDs :many
SELECT id, name, details, price, weight, volume, created_at FROM product WHERE id = ANY($1::int[]);
