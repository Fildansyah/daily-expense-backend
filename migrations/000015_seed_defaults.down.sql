-- Hapus hanya data global hasil seed. Data milik user tidak tersentuh.
DELETE FROM categories WHERE user_id IS NULL;
DELETE FROM transaction_types WHERE user_id IS NULL;
