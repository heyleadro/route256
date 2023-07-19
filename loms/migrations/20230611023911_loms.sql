-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_id serial PRIMARY KEY,
    user_id bigint,
    "status" text
);

CREATE TABLE IF NOT EXISTS user_items (
    sku bigint,
    order_id bigint,
    warehouse_id bigint,
    "count" bigint,
    PRIMARY KEY (order_id)
);

CREATE TABLE IF NOT EXISTS stock_items (
    sku bigint,
    warehouse_id bigint,
    "count" bigint,
    PRIMARY KEY (warehouse_id, sku)
);

CREATE TABLE IF NOT EXISTS reserved_items (
    sku bigint,
    warehouse_id bigint,
    "count" bigint,
    PRIMARY KEY (warehouse_id, sku)
);

INSERT INTO stock_items (sku, warehouse_id, count)
    VALUES (773297411, 10, 100),
    (1076963, 10, 100),
    (1148162, 10, 100),
    (1625903, 10, 100),
    (2618151, 10, 100),
    (2956315, 10, 100),
    (2958025, 10, 100),
    (3596599, 10, 100),
    (3618852, 10, 100),
    (4288068, 10, 100),
    (4465995, 10, 100),
    (4487693, 10, 100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reserved_items;
DROP TABLE IF EXISTS stock_items;
DROP TABLE IF EXISTS user_items;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
