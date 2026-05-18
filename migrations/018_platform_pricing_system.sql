-- Description: Platform margin configuration and Product price enhancements
-- Author: Gemini CLI

-- 1. Create platform_margin_settings table
CREATE TABLE IF NOT EXISTS platform_margin_settings (
    id              BIGSERIAL PRIMARY KEY,
    category_id     BIGINT, -- NULL means global
    margin_type     VARCHAR(20) NOT NULL DEFAULT 'FIXED', -- FIXED or PERCENT
    margin_value    DECIMAL(15, 2) NOT NULL DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Add original_price and platform_margin to products
ALTER TABLE products 
    ADD COLUMN IF NOT EXISTS original_price DECIMAL(18, 2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS platform_margin DECIMAL(18, 2) DEFAULT 0;

-- 3. Seed initial platform margin (Global Rp 200 per transaction)
INSERT INTO platform_margin_settings (margin_type, margin_value) VALUES ('FIXED', 200);
