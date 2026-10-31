DROP INDEX IF EXISTS idx_transactions_user_transfer;
DROP INDEX IF EXISTS idx_transactions_user_wallet;
DROP INDEX IF EXISTS idx_transactions_user_category;
DROP INDEX IF EXISTS idx_transactions_user_date;

ALTER TABLE transactions
    DROP COLUMN IF EXISTS created_by_rule_id,
    DROP COLUMN IF EXISTS note,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS transfer_wallet_id;

ALTER TABLE transactions ALTER COLUMN amount TYPE FLOAT USING amount::FLOAT;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'transactions' AND column_name = 'transaction_date'
    ) THEN
        ALTER TABLE transactions RENAME COLUMN transaction_date TO date;
    END IF;
END $$;
