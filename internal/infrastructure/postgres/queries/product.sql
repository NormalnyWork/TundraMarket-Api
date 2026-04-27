-- name: GetAllProducts :many
SELECT * FROM product ORDER BY id;

-- name: GetProductByID :one
SELECT * FROM product WHERE id = $1;

-- name: GetProductsByIDs :many
SELECT * FROM product WHERE id = ANY($1::int[]);