-- +goose Up
CREATE TABLE admin (
                       id         SERIAL PRIMARY KEY,
                       login      VARCHAR(255) NOT NULL UNIQUE,
                       password   VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO admin (login, password) VALUES ('admin', 'change_me');

-- +goose Down
DROP TABLE admin;