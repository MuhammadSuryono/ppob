-- +goose Up
-- ============================================
-- Trigger: Create wallet when role is assigned
-- ============================================

-- Function to create wallet for role
CREATE OR REPLACE FUNCTION create_wallet_for_role()
RETURNS TRIGGER AS $$
DECLARE
    parent_wallet_id UUID;
    role_name_val VARCHAR(50);
    mitra_user_id UUID;
BEGIN
    -- Get role name
    SELECT role_name INTO role_name_val 
    FROM roles 
    WHERE role_id = NEW.role_id;
    
    IF role_name_val = 'Mitra' THEN
        -- Create main wallet for Mitra
        INSERT INTO wallets (wallet_id, owner_id, balance_available, balance_held, is_main_wallet, parent_wallet_id, is_frozen)
        VALUES (gen_random_uuid(), NEW.user_id, 0, 0, TRUE, NULL, FALSE);
        
    ELSIF role_name_val = 'Staff' THEN
        -- Get assigned_by (Mitra) user_id
        mitra_user_id := NEW.assigned_by;
        
        -- Find Mitra's main wallet
        SELECT wallet_id INTO parent_wallet_id 
        FROM wallets 
        WHERE owner_id = mitra_user_id AND is_main_wallet = TRUE;
        
        IF parent_wallet_id IS NULL THEN
            RAISE EXCEPTION 'Mitra (%) has no main wallet', mitra_user_id;
        END IF;
        
        -- Create sub-wallet for Staff
        INSERT INTO wallets (wallet_id, owner_id, balance_available, balance_held, is_main_wallet, parent_wallet_id, is_frozen)
        VALUES (gen_random_uuid(), NEW.user_id, 0, 0, FALSE, parent_wallet_id, FALSE);
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
CREATE TRIGGER trg_create_wallet_on_role_assign
AFTER INSERT ON user_roles
FOR EACH ROW EXECUTE FUNCTION create_wallet_for_role();

-- +goose Down
DROP TRIGGER IF EXISTS trg_create_wallet_on_role_assign ON user_roles;
DROP FUNCTION IF EXISTS create_wallet_for_role();