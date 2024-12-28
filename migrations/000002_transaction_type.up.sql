CREATE TABLE IF NOT EXISTS transaction_type (
    ID UUID PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);