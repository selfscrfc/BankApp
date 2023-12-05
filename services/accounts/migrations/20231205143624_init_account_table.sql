-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
       id VARCHAR(256) PRIMARY KEY,
       userid VARCHAR(256),
       iscredit BOOL,
       balance INT,
       currency INT,
       isblocked BOOL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounts;
-- +goose StatementEnd
