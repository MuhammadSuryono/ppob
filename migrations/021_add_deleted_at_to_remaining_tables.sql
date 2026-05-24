-- Migration: 021_add_deleted_at_to_remaining_tables.sql
-- Description: Add deleted_at column to remaining tables for consistent soft delete support
-- +goose Up
ALTER TABLE device_fingerprints ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE otp_codes ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE staff_global_margin_settings ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE staff_product_margin_overrides ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_device_fingerprints_deleted_at ON device_fingerprints(deleted_at);
CREATE INDEX IF NOT EXISTS idx_otp_codes_deleted_at ON otp_codes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_staff_global_margin_deleted_at ON staff_global_margin_settings(deleted_at);
CREATE INDEX IF NOT EXISTS idx_staff_product_margin_deleted_at ON staff_product_margin_overrides(deleted_at);
CREATE INDEX IF NOT EXISTS idx_notifications_deleted_at ON notifications(deleted_at);

-- +goose Down
DROP INDEX IF EXISTS idx_notifications_deleted_at;
DROP INDEX IF EXISTS idx_staff_product_margin_deleted_at;
DROP INDEX IF EXISTS idx_staff_global_margin_deleted_at;
DROP INDEX IF EXISTS idx_otp_codes_deleted_at;
DROP INDEX IF EXISTS idx_device_fingerprints_deleted_at;

ALTER TABLE notifications DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE staff_product_margin_overrides DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE staff_global_margin_settings DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE otp_codes DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE device_fingerprints DROP COLUMN IF EXISTS deleted_at;
