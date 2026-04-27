-- name: AddProductToOrder :exec
INSERT INTO orders_to_product (orders_id, product_id, quantity)
VALUES ($1, $2, $3);

-- name: GetProductsByOrderID :many
SELECT p.id, p.name, p.weight, p.volume, otp.quantity
FROM orders_to_product otp
         JOIN product p ON p.id = otp.product_id
WHERE otp.orders_id = $1;