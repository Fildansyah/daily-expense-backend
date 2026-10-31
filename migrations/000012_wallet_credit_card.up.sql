ALTER TABLE wallets
    ADD COLUMN IF NOT EXISTS opening_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS credit_limit NUMERIC(18,2),
    ADD COLUMN IF NOT EXISTS statement_day SMALLINT,
    ADD COLUMN IF NOT EXISTS due_day SMALLINT,
    ADD COLUMN IF NOT EXISTS interest_rate NUMERIC(5,2),
    ADD COLUMN IF NOT EXISTS account_number VARCHAR(64),
    ADD COLUMN IF NOT EXISTS is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP;

ALTER TABLE wallets
    ALTER COLUMN balance TYPE NUMERIC(18,2) USING balance::NUMERIC(18,2);

-- Saldo yang sudah ada dianggap sebagai saldo awal agar tidak berubah setelah migrasi.
UPDATE wallets SET opening_balance = balance WHERE opening_balance = 0;

CREATE INDEX IF NOT EXISTS idx_wallets_user ON wallets (user_id, is_archived);
