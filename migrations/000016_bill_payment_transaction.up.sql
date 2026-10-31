ALTER TABLE bill_payments ADD COLUMN IF NOT EXISTS transaction_id UUID;

CREATE INDEX IF NOT EXISTS idx_bill_payments_transaction ON bill_payments (transaction_id);
