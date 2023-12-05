-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
                           id VARCHAR(256) PRIMARY KEY,
                           userid VARCHAR(256),
                           iscredit INT,
                           balance INT,
                           currency ENUM(0,1,2),
                           isblocked bool
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounts;
-- +goose StatementEnd
