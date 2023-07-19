-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart (
    user_id bigint,
    sku bigint,
    "count" bigint,
    PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart;
-- +goose StatementEnd
