-- +goose Up
ALTER TABLE status_history
    ADD COLUMN comment VARCHAR(255);

-- +goose Down
ALTER TABLE status_history
    DROP COLUMN comment;
