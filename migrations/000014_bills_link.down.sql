DROP INDEX IF EXISTS idx_bills_user_wallet;
DROP INDEX IF EXISTS idx_bills_user_due;

ALTER TABLE bills
    DROP COLUMN IF EXISTS currency,
    DROP COLUMN IF EXISTS remind_due,
    DROP COLUMN IF EXISTS remind_h1,
    DROP COLUMN IF EXISTS remind_h3,
    DROP COLUMN IF EXISTS interest_rate,
    DROP COLUMN IF EXISTS installment_index,
    DROP COLUMN IF EXISTS installment_total,
    DROP COLUMN IF EXISTS category_id,
    DROP COLUMN IF EXISTS wallet_id;
