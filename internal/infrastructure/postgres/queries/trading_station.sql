-- name: GetTradingStationByPhone :one
SELECT * FROM trading_station WHERE phone = $1;

-- name: GetTradingStationByID :one
SELECT * FROM trading_station WHERE id = $1;

-- name: GetAllTradingStations :many
SELECT * FROM trading_station ORDER BY id;

-- name: SetTradingStationPhone :one
UPDATE trading_station
SET phone = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
