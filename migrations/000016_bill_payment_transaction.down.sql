DROP INDEX IF EXISTS idx_bill_payments_transaction;

ALTER TABLE bill_payments DROP COLUMN IF EXISTS transaction_id;
