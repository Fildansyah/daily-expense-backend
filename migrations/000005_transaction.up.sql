CREATE TABLE IF NOT EXISTS transactions (
    ID UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    wallet_id UUID NOT NULL,
    category_id UUID NOT NULL,
    transaction_type_id UUID NOT NULL,
    amount FLOAT NOT NULL,
    memo VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);