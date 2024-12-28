CREATE TABLE IF NOT EXISTS transaction (
    ID UUID PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    wallet_id VARCHAR(255) NOT NULL,
    category_id VARCHAR(255) NOT NULL,
    transaction_type_id VARCHAR(255) NOT NULL,
    amount FLOAT NOT NULL,
    memo VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);