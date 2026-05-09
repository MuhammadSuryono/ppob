-- +goose Up
-- Add salt column to otp_codes for secure OTP hashing
ALTER TABLE otp_codes ADD COLUMN salt VARCHAR(32) NOT NULL DEFAULT '';

-- Update existing records with random salt (they will be invalid, but that's okay as they expire quickly)
UPDATE otp_codes SET salt = substr(md5(random()::text), 1, 32) WHERE salt = '';

-- Make salt NOT NULL after populating (already NOT NULL with default)
-- +goose Down
ALTER TABLE otp_codes DROP COLUMN IF EXISTS salt;
