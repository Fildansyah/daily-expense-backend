CREATE TABLE IF NOT EXISTS bill_payments (
    id UUID PRIMARY KEY,
    bill_id UUID NOT NULL,
    user_id UUID NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    paid_at DATE NOT NULL,
    wallet_id UUID,
    note VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bill_payments_bill ON bill_payments (bill_id, paid_at);
CREATE INDEX IF NOT EXISTS idx_bill_payments_user ON bill_payments (user_id, paid_at);
CREATE INDEX IF NOT EXISTS idx_bill_payments_wallet ON bill_payments (wallet_id);

ALTER TABLE bills ALTER COLUMN amount TYPE NUMERIC(18,2) USING amount::NUMERIC(18,2);
ALTER TABLE bills ALTER COLUMN paid_amount TYPE NUMERIC(18,2) USING paid_amount::NUMERIC(18,2);
