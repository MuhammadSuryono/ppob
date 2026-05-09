-- Migration: 010_product_sync_tracking.sql
-- Description: Add product sync tracking and enhance product tables

-- Add provider_ref column to transactions if not exists
-- ALTER TABLE transactions ADD COLUMN IF NOT EXISTS provider_ref VARCHAR(100);

-- Add UpdatedAt to transactions for tracking
-- ALTER TABLE transactions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Add last_sync_at to products if not exists
-- ALTER TABLE products ADD COLUMN IF NOT EXISTS last_sync_at TIMESTAMP;

-- Create product sync log table
CREATE TABLE IF NOT EXISTS product_sync_logs (
    id BIGSERIAL PRIMARY KEY,
    sync_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    total_products INTEGER NOT NULL DEFAULT 0,
    created_products INTEGER NOT NULL DEFAULT 0,
    updated_products INTEGER NOT NULL DEFAULT 0,
    failed_products INTEGER NOT NULL DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    duration_seconds INTEGER
);

CREATE INDEX idx_sync_logs_type ON product_sync_logs(sync_type);
CREATE INDEX idx_sync_logs_status ON product_sync_logs(status);
CREATE INDEX idx_sync_logs_started ON product_sync_logs(started_at);

-- Create idempotency keys table
CREATE TABLE IF NOT EXISTS idempotency_keys (
    id BIGSERIAL PRIMARY KEY,
    key_value VARCHAR(255) NOT NULL UNIQUE,
    transaction_id VARCHAR(100) NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_idempotency_key ON idempotency_keys(key_value);
CREATE INDEX idx_idempotency_expires ON idempotency_keys(expires_at);

-- Create index on provider_ref in transactions
CREATE INDEX IF NOT EXISTS idx_transactions_provider_ref ON transactions(provider_ref);

-- Create index on customer_number in transactions for history lookup
CREATE INDEX IF NOT EXISTS idx_transactions_customer ON transactions(customer_number);