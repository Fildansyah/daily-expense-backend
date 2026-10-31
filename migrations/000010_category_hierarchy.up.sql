ALTER TABLE categories
    ADD COLUMN IF NOT EXISTS parent_id UUID,
    ADD COLUMN IF NOT EXISTS sort_order INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_categories_user_parent ON categories (user_id, parent_id);
CREATE INDEX IF NOT EXISTS idx_categories_transaction_type ON categories (transaction_type_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_categories_no_self_parent'
    ) THEN
        ALTER TABLE categories
            ADD CONSTRAINT chk_categories_no_self_parent
            CHECK (parent_id IS NULL OR parent_id <> id);
    END IF;
END $$;
