-- name: AddStatusHistory :exec
INSERT INTO status_history (orders_id, status, comment)
VALUES ($1, $2, $3);

-- name: GetStatusHistoryByOrderID :many
SELECT id, orders_id, status, comment, created_at FROM status_history
WHERE orders_id = $1
ORDER BY created_at ASC;

-- name: GetStatusHistoryAfter :many
SELECT id, orders_id, status, comment, created_at FROM status_history
WHERE orders_id = $1
  AND created_at > to_timestamp($2)
ORDER BY created_at ASC;
