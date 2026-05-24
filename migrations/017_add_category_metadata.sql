-- +goose Up
-- Description: Add metadata fields to categories for dynamic UI generation
-- Author: Gemini CLI

ALTER TABLE categories
    ADD COLUMN IF NOT EXISTS input_type VARCHAR(50) DEFAULT 'TEXT',
    ADD COLUMN IF NOT EXISTS input_label VARCHAR(100),
    ADD COLUMN IF NOT EXISTS placeholder VARCHAR(255),
    ADD COLUMN IF NOT EXISTS validation_regex VARCHAR(255);

-- Update existing categories with some defaults
UPDATE categories SET input_type = 'NUMBER', input_label = 'Nomor HP', placeholder = '08xxxxxxxxxx', validation_regex = '^[0-9]{10,14}$' WHERE code IN ('pulsa', 'data', 'paket-sms-telpon');
UPDATE categories SET input_type = 'NUMBER', input_label = 'ID Pelanggan', placeholder = '12 digit nomor meter', validation_regex = '^[0-9]{11,12}$' WHERE code = 'pln';
UPDATE categories SET input_type = 'NUMBER', input_label = 'ID Pelanggan', placeholder = 'Nomor pelanggan E-Money', validation_regex = '^[0-9]{10,14}$' WHERE code = 'e-money';
UPDATE categories SET input_type = 'TEXT', input_label = 'ID Game', placeholder = 'User ID (Zone ID)', validation_regex = '^[a-zA-Z0-9]{5,20}$' WHERE code = 'games';

-- +goose Down
-- ALTER TABLE categories DROP COLUMN IF EXISTS input_type;
-- ALTER TABLE categories DROP COLUMN IF EXISTS input_label;
-- ALTER TABLE categories DROP COLUMN IF EXISTS placeholder;
-- ALTER TABLE categories DROP COLUMN IF EXISTS validation_regex;
