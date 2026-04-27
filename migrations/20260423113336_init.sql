-- +goose Up
CREATE TABLE nomad (
    id          SERIAL      PRIMARY KEY,
    phone       VARCHAR(25) NOT NULL UNIQUE,
    longitude   DECIMAL(9, 6),
    latitude    DECIMAL(9, 6),
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   DEFAULT NOW()
);

CREATE TABLE trading_station (
    id          SERIAL      PRIMARY KEY,
    phone       VARCHAR(25) UNIQUE,
    name        VARCHAR(255),
    longitude   DECIMAL(9, 6),
    latitude    DECIMAL(9, 6),
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   DEFAULT NOW()
);

CREATE TYPE status AS ENUM(
    'CREATED',
    'PROCESSING',
    'SENT',
    'COMPLETED',
    'CANCELLED',
    'DENIED'
);

CREATE TABLE orders (
    id                  SERIAL        PRIMARY KEY,
    nomad_id            INT           REFERENCES nomad(id) ON DELETE CASCADE,
    trading_station_id  INT           REFERENCES trading_station(id) ON DELETE CASCADE,
    status              status        NOT NULL DEFAULT 'CREATED',
    longitude           DECIMAL(9, 6) NOT NULL,
    latitude            DECIMAL(9, 6) NOT NULL,
    comment             VARCHAR(255),
    created_at          TIMESTAMP     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_nomad_id ON orders(nomad_id);
CREATE INDEX idx_orders_station_id ON orders(trading_station_id);
CREATE INDEX idx_orders_status ON orders(status);

CREATE TABLE product (
    id          SERIAL      PRIMARY KEY,
    name        VARCHAR     NOT NULL,
    price       INT,
    weight      INT,
    volume      INT,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE TABLE orders_to_product (
    id          SERIAL      PRIMARY KEY,
    orders_id   INT         NOT NULL REFERENCES orders(id) ON DELETE CASCADE ,
    product_id  INT         NOT NULL REFERENCES product(id) ON DELETE CASCADE,
    UNIQUE (orders_id, product_id),
    quantity    INT         NOT NULL DEFAULT 1,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_otp_orders_id ON orders_to_product(orders_id);
CREATE INDEX idx_otp_product_id ON orders_to_product(product_id);

CREATE TABLE status_history (
    id          SERIAL      PRIMARY KEY,
    orders_id   INT         NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status      status,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_status_history_order_id ON status_history(orders_id);

-- +goose Down
DROP INDEX idx_status_history_order_id;
DROP INDEX idx_otp_product_id;
DROP INDEX idx_otp_orders_id;
DROP INDEX idx_orders_status;
DROP INDEX idx_orders_station_id;
DROP INDEX idx_orders_nomad_id;


DROP TABLE status_history;
DROP TABLE orders_to_product;
DROP TABLE product;
DROP TABLE orders;
DROP TYPE status;
DROP TABLE trading_station;
DROP TABLE nomad;

