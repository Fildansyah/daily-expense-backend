CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    category_id UUID NOT NULL,
    period_month DATE NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_budgets_user_category_month
    ON budgets (user_id, category_id, period_month)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_budgets_user_month ON budgets (user_id, period_month);
