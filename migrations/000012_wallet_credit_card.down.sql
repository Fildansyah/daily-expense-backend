DROP INDEX IF EXISTS idx_wallets_user;

ALTER TABLE wallets
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS is_archived,
    DROP COLUMN IF EXISTS account_number,
    DROP COLUMN IF EXISTS interest_rate,
    DROP COLUMN IF EXISTS due_day,
    DROP COLUMN IF EXISTS statement_day,
    DROP COLUMN IF EXISTS credit_limit,
    DROP COLUMN IF EXISTS opening_balance;

ALTER TABLE wallets ALTER COLUMN balance TYPE FLOAT USING balance::FLOAT;
