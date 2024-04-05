-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tx_disbursements (
    id VARCHAR(36) PRIMARY KEY,
    group_id VARCHAR(36),
    amount BIGINT NOT NULL,
    account_number VARCHAR(20) NOT NULL,
    account_name VARCHAR(50) NOT NULL,
    status VARCHAR(15) NOT NULL,
    message VARCHAR(100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tx_disbursements;
-- +goose StatementEnd
