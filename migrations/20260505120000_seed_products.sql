-- +goose Up
INSERT INTO product (name, details, price, weight, volume)
VALUES
    ('Хлеб', 'буханка', 49, 800, 2000),
    ('Молоко', 'пакет', 109, 1000, 1000),
    ('Крупы', 'герметичная упаковка', 189, 2000, 3000),
    ('Средства гигиены', NULL, 299, 1000, 2000);

-- +goose Down
DELETE FROM product
WHERE name IN ('Хлеб', 'Молоко', 'Крупы', 'Средства гигиены');
