-- +goose Up
INSERT INTO trading_station (name, latitude, longitude)
VALUES
    ('Паюта', 67.9914965, 68.5914509),
    ('Степина', 67.6477483, 67.9260914);

-- +goose Down
DELETE FROM trading_station
WHERE name IN ('Паюта', 'Степина');
