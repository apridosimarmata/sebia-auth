-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ms_wallet (
    id VARCHAR(36) PRIMARY KEY,
    owned_by VARCHAR(36) NOT NULL,
    enabled_at VARCHAR(30),
    balance INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ms_wallet;
-- +goose StatementEnd
