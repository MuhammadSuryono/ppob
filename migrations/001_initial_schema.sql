-- Migration: 001_initial_schema.sql
-- Description: Core tables: roles, users, user_roles, wallets, categories, products, transactions
-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- 1. roles
-- ============================================
CREATE TABLE IF NOT EXISTS roles (
    id          BIGSERIAL PRIMARY KEY,
    role_name   VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (role_name, description, is_active) VALUES
    ('Mitra', 'Mitra (Merchant/Agent)', TRUE),
    ('Staff', 'Staff member working under Mitra', TRUE),
    ('Admin', 'System administrator', TRUE)
ON CONFLICT (role_name) DO NOTHING;

-- ============================================
-- 2. users
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL PRIMARY KEY,
    email           VARCHAR(255) UNIQUE,
    phone           VARCHAR(20) UNIQUE,
    name            VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    pin_hash        VARCHAR(255),
    pin_salt        VARCHAR(32),
    role            VARCHAR(20) NOT NULL DEFAULT 'Mitra',
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    phone_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    avatar          VARCHAR(500),
    address         TEXT,
    date_of_birth   DATE,
    last_login_at   TIMESTAMP,
    deleted_at      TIMESTAMP,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- ============================================
-- 3. user_roles (many-to-many)
-- ============================================
CREATE TABLE IF NOT EXISTS user_roles (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id     BIGINT NOT NULL REFERENCES roles(id),
    assigned_by BIGINT REFERENCES users(id),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_user_roles_user ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role_id);

-- ============================================
-- 4. wallets
-- ============================================
CREATE TABLE IF NOT EXISTS wallets (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id         BIGINT,
    balance         DECIMAL(18, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    hold_amount     DECIMAL(18, 2) NOT NULL DEFAULT 0,
    currency        VARCHAR(10) NOT NULL DEFAULT 'IDR',
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    is_main_wallet  BOOLEAN NOT NULL DEFAULT FALSE,
    parent_wallet_id BIGINT REFERENCES wallets(id),
    is_frozen       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_wallets_user ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_wallets_parent ON wallets(parent_wallet_id) WHERE is_main_wallet = false;

-- ============================================
-- 5. categories
-- ============================================
CREATE TABLE IF NOT EXISTS categories (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(255) UNIQUE NOT NULL,
    icon_url       VARCHAR(500),
    display_order  INT DEFAULT 0,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 6. products
-- ============================================
CREATE TABLE IF NOT EXISTS products (
    id             BIGSERIAL PRIMARY KEY,
    code           VARCHAR(100) UNIQUE NOT NULL,
    name           VARCHAR(255) NOT NULL,
    category       VARCHAR(50),
    category_id    BIGINT REFERENCES categories(id),
    price          DECIMAL(18, 2) NOT NULL CHECK (price >= 0),
    price_api      DECIMAL(18, 2) NOT NULL DEFAULT 0,
    platform_price DECIMAL(18, 2) NOT NULL DEFAULT 0,
    provider       VARCHAR(50),
    product_type   VARCHAR(20),
    description    TEXT,
    stock          INT NOT NULL DEFAULT 0,
    status         VARCHAR(20) NOT NULL DEFAULT 'active',
    is_prepaid     BOOLEAN NOT NULL DEFAULT TRUE,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    last_sync_at   TIMESTAMP,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_code ON products(code);
CREATE INDEX IF NOT EXISTS idx_products_active ON products(is_active) WHERE is_active = true;

-- ============================================
-- 7. transactions
-- ============================================
CREATE TABLE IF NOT EXISTS transactions (
    id                    BIGSERIAL PRIMARY KEY,
    transaction_id        VARCHAR(100) UNIQUE NOT NULL,
    ref_id                VARCHAR(255) UNIQUE,
    user_id               BIGINT NOT NULL REFERENCES users(id),
    wallet_id             BIGINT REFERENCES wallets(id),
    product_id            BIGINT NOT NULL REFERENCES products(id),
    product_code          VARCHAR(100),
    customer_number       VARCHAR(50) NOT NULL,
    amount                DECIMAL(18, 2) NOT NULL CHECK (amount >= 0),
    price                 DECIMAL(18, 2) NOT NULL DEFAULT 0,
    selling_price         DECIMAL(18, 2) NOT NULL DEFAULT 0,
    margin                DECIMAL(15, 2) NOT NULL DEFAULT 0,
    status                VARCHAR(50) NOT NULL DEFAULT 'initiated',
    provider_ref          VARCHAR(100),
    provider_status      VARCHAR(20),
    message               TEXT,
    completed_at         TIMESTAMP,
    previous_status      VARCHAR(50),
    status_change_reason  VARCHAR(255),
    hold_released_at      TIMESTAMP,
    reconciled_at         TIMESTAMP,
    digiflazz_rc          VARCHAR(10),
    webhook_received_at   TIMESTAMP,
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transactions_user ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_ref ON transactions(ref_id);
CREATE INDEX IF NOT EXISTS idx_transactions_provider_ref ON transactions(provider_ref);
CREATE INDEX IF NOT EXISTS idx_transactions_created ON transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_pending_timeout
    ON transactions(status, created_at, updated_at) WHERE status = 'pending';

-- ============================================
-- Triggers for auto-updating timestamps
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_user_roles_updated_at BEFORE UPDATE ON user_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_wallets_updated_at BEFORE UPDATE ON wallets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS wallets CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
