ALTER TABLE categories DROP CONSTRAINT IF EXISTS chk_categories_no_self_parent;
DROP INDEX IF EXISTS idx_categories_transaction_type;
DROP INDEX IF EXISTS idx_categories_user_parent;

ALTER TABLE categories
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS is_archived,
    DROP COLUMN IF EXISTS sort_order,
    DROP COLUMN IF EXISTS parent_id;
