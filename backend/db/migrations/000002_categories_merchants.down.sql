ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_category;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_merchant;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS merchants;
