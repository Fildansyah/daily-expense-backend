ALTER TABLE categories ADD COLUMN IF NOT EXISTS user_id UUID;
ALTER TABLE transaction_types ADD COLUMN IF NOT EXISTS user_id UUID;
ALTER TABLE lenders ADD COLUMN IF NOT EXISTS user_id UUID;

CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);
CREATE INDEX IF NOT EXISTS idx_transaction_types_user_id ON transaction_types(user_id);
CREATE INDEX IF NOT EXISTS idx_lenders_user_id ON lenders(user_id);
