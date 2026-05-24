-- Migration: 019_make_email_nullable.sql
-- Description: Ensure email column is nullable and convert empty strings to NULL to avoid unique constraint issues
-- +goose Up
UPDATE users SET email = NULL WHERE email = '';
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- +goose Down
-- Note: Reverting might fail if there are multiple NULLs
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
