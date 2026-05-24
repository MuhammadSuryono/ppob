-- +goose Up
-- ============================================
-- PII Masking Functions
-- ============================================

-- Mask phone number: +62 812-345-***-****
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION mask_phone(phone TEXT) RETURNS TEXT AS $$
BEGIN
    IF phone IS NULL THEN
        RETURN NULL;
    END IF;
    
    -- Remove all non-digit characters for processing
    phone := regexp_replace(phone, '[^0-9]', '', 'g');
    
    -- If length >= 10, mask middle digits
    IF length(phone) >= 10 THEN
        RETURN '+62 ' || 
               substr(phone, 3, 3) || '-' || 
               substr(phone, 6, 3) || '-***-****';
    END IF;
    
    RETURN '***';
END;
$$ LANGUAGE plpgsql IMMUTABLE;
-- +goose StatementEnd

-- Mask customer number: Show last 4 digits
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION mask_customer_no(no TEXT) RETURNS TEXT AS $$
BEGIN
    IF no IS NULL THEN
        RETURN NULL;
    END IF;
    
    -- Remove spaces
    no := regexp_replace(no, '\s', '', 'g');
    
    IF length(no) >= 4 THEN
        RETURN '****' || substr(no, length(no) - 3, 4);
    END IF;
    
    RETURN '****';
END;
$$ LANGUAGE plpgsql IMMUTABLE;
-- +goose StatementEnd

-- Mask name: Show only first letter
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION mask_name(name TEXT) RETURNS TEXT AS $$
BEGIN
    IF name IS NULL THEN
        RETURN NULL;
    END IF;
    
    name := trim(name);
    
    IF length(name) <= 1 THEN
        RETURN '*';
    END IF;
    
    RETURN substr(name, 1, 1) || repeat('*', length(name) - 1);
END;
$$ LANGUAGE plpgsql IMMUTABLE;
-- +goose StatementEnd

-- Mask email
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION mask_email(email TEXT) RETURNS TEXT AS $$
DECLARE
    at_pos INTEGER;
    domain_part TEXT;
BEGIN
    IF email IS NULL THEN
        RETURN NULL;
    END IF;
    
    at_pos := position('@' IN email);
    
    IF at_pos = 0 THEN
        RETURN mask_customer_no(email);
    END IF;
    
    -- Mask username part, keep domain
    domain_part := substr(email, at_pos);
    
    IF at_pos <= 3 THEN
        RETURN '***' || domain_part;
    END IF;
    
    RETURN substr(email, 1, 2) || repeat('*', at_pos - 3) || domain_part;
END;
$$ LANGUAGE plpgsql IMMUTABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS mask_phone(TEXT);
DROP FUNCTION IF EXISTS mask_customer_no(TEXT);
DROP FUNCTION IF EXISTS mask_name(TEXT);
DROP FUNCTION IF EXISTS mask_email(TEXT);