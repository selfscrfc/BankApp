-- +goose Up
-- +goose StatementBegin
CREATE TABLE customers (
    id VARCHAR(256) PRIMARY KEY,
    fullname VARCHAR(128),
    time BIGINT,
    login VARCHAR(64) UNIQUE,
    password VARCHAR(256),
    isblocked bool
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE customers;
-- +goose StatementEnd
