-- +goose Up
-- Migration: 008_wallet_tables.sql
-- Holds and mitra_product_prices (duplicates removed — consolidated into 003/004)

CREATE TABLE IF NOT EXISTS holds (
    id            BIGSERIAL PRIMARY KEY,
    wallet_id     BIGINT NOT NULL REFERENCES wallets(id),
    amount        DECIMAL(15, 2) NOT NULL DEFAULT 0,
    reference_id  VARCHAR(100) NOT NULL,
    reference_type VARCHAR(50) NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'active',
    expires_at    TIMESTAMP,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    released_at   TIMESTAMP,
    UNIQUE (reference_id, reference_type, status)
);

CREATE INDEX IF NOT EXISTS idx_holds_wallet ON holds(wallet_id);
CREATE INDEX IF NOT EXISTS idx_holds_reference ON holds(reference_id, reference_type);
CREATE INDEX IF NOT EXISTS idx_holds_status ON holds(status);
CREATE INDEX IF NOT EXISTS idx_holds_expires ON holds(expires_at) WHERE status = 'active';

CREATE TABLE IF NOT EXISTS mitra_product_prices (
    id            BIGSERIAL PRIMARY KEY,
    mitra_id      BIGINT NOT NULL,
    product_id    BIGINT NOT NULL,
    selling_price DECIMAL(18, 2) NOT NULL DEFAULT 0,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (mitra_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_mitra_product_mitra ON mitra_product_prices(mitra_id);
CREATE INDEX IF NOT EXISTS idx_mitra_product ON mitra_product_prices(product_id);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION validate_mitra_price()
RETURNS TRIGGER AS $$
DECLARE
    platform_price_val DECIMAL(18, 2);
BEGIN
    SELECT platform_price INTO platform_price_val
    FROM products
    WHERE id = NEW.product_id;

    IF platform_price_val IS NULL THEN
        RAISE EXCEPTION 'Product not found: %', NEW.product_id;
    END IF;

    IF NEW.selling_price < platform_price_val THEN
        RAISE EXCEPTION 'selling_price (%) cannot be less than platform_price (%)',
            NEW.selling_price, platform_price_val;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE TRIGGER trg_validate_mitra_price
    BEFORE INSERT OR UPDATE ON mitra_product_prices
    FOR EACH ROW EXECUTE FUNCTION validate_mitra_price();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS trg_validate_mitra_price ON mitra_product_prices;
DROP FUNCTION IF EXISTS validate_mitra_price();
DROP TABLE IF EXISTS mitra_product_prices CASCADE;
DROP TABLE IF EXISTS holds CASCADE;
