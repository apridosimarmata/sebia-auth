-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tr_wallet_transaction (
    id VARCHAR(36) PRIMARY KEY,
    wallet_id VARCHAR(36) NOT NULL,
    amount INTEGER NOT NULL,
    created_at VARCHAR(30) NOT NULL,
    created_by VARCHAR(36) NOT NULL,
    type VARCHAR(15) NOT NULL,
    status VARCHAR(15) NOT NULL,
    reference_id VARCHAR(36) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallet_transactions;
-- +goose StatementEnd
