-- +goose Up
ALTER TABLE product
ADD COLUMN details TEXT;

-- +goose Down
ALTER TABLE product
DROP COLUMN details;
