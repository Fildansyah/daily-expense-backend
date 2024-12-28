CREATE TABLE IF NOT EXISTS categories (
    ID UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    transaction_type_id UUID NOT NULL,
    bg_color VARCHAR(255) NOT NULL,
    icon VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);