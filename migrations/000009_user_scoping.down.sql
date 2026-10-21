DROP INDEX IF EXISTS idx_lenders_user_id;
DROP INDEX IF EXISTS idx_transaction_types_user_id;
DROP INDEX IF EXISTS idx_categories_user_id;

ALTER TABLE lenders DROP COLUMN IF EXISTS user_id;
ALTER TABLE transaction_types DROP COLUMN IF EXISTS user_id;
ALTER TABLE categories DROP COLUMN IF EXISTS user_id;
