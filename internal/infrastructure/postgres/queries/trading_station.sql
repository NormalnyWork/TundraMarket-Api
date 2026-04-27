-- name: GetTradingStationByPhone :one
SELECT * FROM trading_station WHERE phone = $1;

-- name: GetTradingStationByID :one
SELECT * FROM trading_station WHERE id = $1;

-- name: GetAllTradingStations :many
SELECT * FROM trading_station ORDER BY id;