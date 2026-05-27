-- +goose Up
-- Description: Add inquiry support flags to categories and products
-- Author: Gemini CLI

-- 1. Update categories table
ALTER TABLE categories ADD COLUMN IF NOT EXISTS needs_inquiry BOOLEAN DEFAULT FALSE;

-- 2. Update products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS is_inquiry BOOLEAN DEFAULT FALSE;

-- 3. Flag categories that need inquiry
UPDATE categories SET needs_inquiry = TRUE WHERE code IN ('pln', 'e-money', 'tv', 'gas');

-- 4. Flag specific inquiry products (SKU for checking names/bills)
UPDATE products SET is_inquiry = TRUE WHERE code IN ('danacek', 'gopaycek', 'ovocek', 'plncek');

-- 5. Data refinement for PLN brands (Sub-categories)
-- We'll assume the migration 016 already inserted some PLN products.
-- Let's ensure they have descriptive brands.
UPDATE products SET brand = 'PLN Token' WHERE category = 'PLN' AND code LIKE 'pln%';

-- +goose Down
ALTER TABLE categories DROP COLUMN IF EXISTS needs_inquiry;
ALTER TABLE products DROP COLUMN IF EXISTS is_inquiry;
