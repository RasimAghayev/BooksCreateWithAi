-- +goose Up
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT,
    account_id BIGINT,
    amount BIGINT,
    status VARCHAR(20) NOT NULL,
    type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transaction_entries (
    id SERIAL PRIMARY KEY,
    transaction_id BIGINT REFERENCES transactions(id),
    account_id BIGINT,
    amount BIGINT,
    type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transaction_entries_tx_id ON transaction_entries(transaction_id);

-- +goose Down
DROP TABLE IF EXISTS transaction_entries;
DROP TABLE IF EXISTS transactions;
