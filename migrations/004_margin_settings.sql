-- +goose Up
-- Staff margin settings

CREATE TABLE IF NOT EXISTS staff_global_margin_settings (
    id                   BIGSERIAL PRIMARY KEY,
    mitra_id              BIGINT,
    staff_id              BIGINT NOT NULL,
    commission_type       VARCHAR(20) NOT NULL DEFAULT 'MarginShare',
    global_margin_percent DECIMAL(5, 2) NOT NULL DEFAULT 0,
    fixed_allowance       DECIMAL(15, 2) NOT NULL DEFAULT 0,
    is_active             BOOLEAN NOT NULL DEFAULT TRUE,
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (mitra_id, staff_id)
);

CREATE INDEX IF NOT EXISTS idx_staff_margin_staff ON staff_global_margin_settings(staff_id);

CREATE TABLE IF NOT EXISTS staff_product_margin_overrides (
    id              BIGSERIAL PRIMARY KEY,
    mitra_id        BIGINT,
    staff_id        BIGINT NOT NULL,
    product_code    VARCHAR(100) NOT NULL,
    margin_percent  DECIMAL(5, 2) NOT NULL DEFAULT 0,
    fixed_margin     DECIMAL(15, 2) NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (staff_id, product_code)
);

CREATE INDEX IF NOT EXISTS idx_product_override_staff ON staff_product_margin_overrides(staff_id);
CREATE INDEX IF NOT EXISTS idx_product_override_code ON staff_product_margin_overrides(product_code);

-- +goose Down
DROP TABLE IF EXISTS staff_product_margin_overrides CASCADE;
DROP TABLE IF EXISTS staff_global_margin_settings CASCADE;
