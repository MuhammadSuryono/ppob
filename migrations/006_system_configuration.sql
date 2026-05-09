-- +goose Up
-- ============================================
-- 21. system_settings (global key-value configuration)
-- ============================================
CREATE TABLE system_settings (
    key       VARCHAR(100) PRIMARY KEY,
    value     TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed data
INSERT INTO system_settings (key, value) VALUES 
('platform_markup_percent', '0.05'),
('staff_default_daily_txn_limit', '50'),
('staff_default_daily_amount_limit', '5000000'),
('pending_timeout_minutes', '15'),
('max_idempotency_ttl_hours', '24');

-- +goose Down
DELETE FROM system_settings;
DROP TABLE IF EXISTS system_settings;
