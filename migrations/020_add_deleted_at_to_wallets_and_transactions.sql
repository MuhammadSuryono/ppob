-- Migration: 020_add_deleted_at_to_wallets_and_transactions.sql
-- Description: Add deleted_at column for soft delete support in several tables
-- +goose Up
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE roles ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE products ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE daily_limits ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_wallets_deleted_at ON wallets(deleted_at);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at ON transactions(deleted_at);
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_roles_deleted_at ON user_roles(deleted_at);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);
CREATE INDEX IF NOT EXISTS idx_daily_limits_deleted_at ON daily_limits(deleted_at);

-- +goose Down
DROP INDEX IF EXISTS idx_daily_limits_deleted_at;
DROP INDEX IF EXISTS idx_products_deleted_at;
DROP INDEX IF EXISTS idx_categories_deleted_at;
DROP INDEX IF EXISTS idx_user_roles_deleted_at;
DROP INDEX IF EXISTS idx_roles_deleted_at;
DROP INDEX IF EXISTS idx_transactions_deleted_at;
DROP INDEX IF EXISTS idx_wallets_deleted_at;

ALTER TABLE daily_limits DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE products DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE categories DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE user_roles DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE roles DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE transactions DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE wallets DROP COLUMN IF EXISTS deleted_at;
