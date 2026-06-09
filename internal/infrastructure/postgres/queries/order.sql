-- name: CreateOrder :one
INSERT INTO orders (nomad_id, trading_station_id, longitude, latitude, comment)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: GetOrderById :one
SELECT * FROM orders WHERE id = $1;


-- name: GetCurrentOrderByNomadID :one
SELECT * FROM orders
WHERE nomad_id = $1
  AND status NOT IN ('CANCELLED')
ORDER BY created_at DESC
    LIMIT 1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2
WHERE id = $1
    RETURNING *;

-- name: GetOrdersByNomadIDAndCategory :many
SELECT * FROM orders
WHERE nomad_id = $1
  AND CASE sqlc.arg(order_category)::text
    WHEN 'NEW'        THEN status = 'CREATED'
    WHEN 'PROCESSING' THEN status IN ('PROCESSING', 'SENT')
    WHEN 'HISTORY'    THEN status IN ('COMPLETED', 'CANCELLED', 'DENIED')
END
ORDER BY created_at DESC
LIMIT  $2
OFFSET COALESCE(sqlc.arg(anchor)::int, 0);


-- name: GetOrdersByNomadIDUpdatedAfter :many
SELECT DISTINCT o.* FROM orders o
LEFT JOIN status_history sh ON sh.orders_id = o.id
WHERE o.nomad_id = $1
  AND (o.created_at > to_timestamp($2) OR sh.created_at > to_timestamp($2))
ORDER BY o.created_at DESC;

-- name: GetOrdersByStationAndCategory :many
SELECT * FROM orders
WHERE trading_station_id = $1
  AND CASE sqlc.arg(order_category)::text
    WHEN 'NEW'        THEN status = 'CREATED'
    WHEN 'PROCESSING' THEN status IN ('PROCESSING', 'SENT')
    WHEN 'HISTORY'    THEN status IN ('COMPLETED', 'CANCELLED', 'DENIED')
END
ORDER BY created_at DESC
LIMIT  $2
OFFSET COALESCE(sqlc.arg(anchor)::int, 0);

-- name: GetOrdersByStationUpdatedAfter :many
SELECT DISTINCT o.* FROM orders o
LEFT JOIN status_history sh ON sh.orders_id = o.id
WHERE o.trading_station_id = $1
  AND (o.created_at > to_timestamp($2) OR sh.created_at > to_timestamp($2))
ORDER BY o.created_at DESC;
