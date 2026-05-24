-- +goose Up
-- Migration: 011_integration_enhancements.sql
-- Description: Add compensation jobs and webhook tracking tables

-- Compensation Jobs Table
CREATE TABLE IF NOT EXISTS compensation_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(100) NOT NULL UNIQUE,
    transaction_id VARCHAR(100),
    action VARCHAR(50) NOT NULL,
    payload TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,
    error_message TEXT,
    next_retry_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_compensation_job_id ON compensation_jobs(job_id);
CREATE INDEX IF NOT EXISTS idx_compensation_transaction ON compensation_jobs(transaction_id);
CREATE INDEX IF NOT EXISTS idx_compensation_status ON compensation_jobs(status);
CREATE INDEX IF NOT EXISTS idx_compensation_next_retry ON compensation_jobs(next_retry_at);

-- Webhook Processing Log Table
CREATE TABLE IF NOT EXISTS webhook_logs (
    id BIGSERIAL PRIMARY KEY,
    webhook_id VARCHAR(100) NOT NULL,
    ref_id VARCHAR(100) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    provider_status VARCHAR(50),
    request_body TEXT,
    response_body TEXT,
    error_message TEXT,
    processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhook_ref_id ON webhook_logs(ref_id);
CREATE INDEX IF NOT EXISTS idx_webhook_provider ON webhook_logs(provider);
CREATE INDEX IF NOT EXISTS idx_webhook_processed ON webhook_logs(processed_at);

-- Add provider_ref to transactions table if not exists
-- ALTER TABLE transactions ADD COLUMN IF NOT EXISTS provider_ref VARCHAR(100);

-- Create digiflazz error codes table for reference
CREATE TABLE IF NOT EXISTS digiflazz_error_codes (
    id BIGSERIAL PRIMARY KEY,
    rc_code VARCHAR(10) NOT NULL UNIQUE,
    error_type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    user_message_id VARCHAR(50) NOT NULL,
    suggested_action TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert standard error codes
INSERT INTO digiflazz_error_codes (rc_code, error_type, message, user_message_id, suggested_action) VALUES
('00', 'success', 'Transaction successful', 'success', 'Transaction completed successfully'),
('01', 'error', 'Invalid product code', 'invalid_product', 'Please check the product code and try again'),
('02', 'error', 'Invalid phone number or customer ID', 'invalid_phone', 'Please verify the phone number or customer ID'),
('03', 'pending', 'Transaction is being processed', 'pending', 'Please wait while we process your transaction'),
('06', 'error', 'Product is currently inactive', 'product_inactive', 'Please choose a different product'),
('39', 'error', 'Insufficient balance in provider', 'insufficient_balance', 'Please try again later or use a different product'),
('69', 'error', 'System error occurred', 'system_error', 'Please try again in a few minutes'),
('99', 'error', 'Transaction timeout', 'timeout', 'Please try again or contact support')
ON CONFLICT (rc_code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS digiflazz_error_codes CASCADE;
DROP TABLE IF EXISTS webhook_logs CASCADE;