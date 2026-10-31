ALTER TABLE transactions ADD COLUMN IF NOT EXISTS transfer_wallet_id UUID;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS note TEXT;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS created_by_rule_id UUID;

-- Model Go memakai transaction_date; kolom lama bernama date.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'transactions' AND column_name = 'date'
    ) THEN
        ALTER TABLE transactions RENAME COLUMN date TO transaction_date;
    END IF;
END $$;

ALTER TABLE transactions
    ALTER COLUMN amount TYPE NUMERIC(18,2) USING amount::NUMERIC(18,2);

CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions (user_id, transaction_date DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_user_category ON transactions (user_id, category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_wallet ON transactions (user_id, wallet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_transfer ON transactions (user_id, transfer_wallet_id);
