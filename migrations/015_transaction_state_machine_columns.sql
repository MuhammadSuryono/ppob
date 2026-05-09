-- Migration: Add transaction state machine columns
-- Created: 2026-05-07
-- Purpose: Add columns required for state machine tracking per TRANSACTION_STATE_MACHINE.md

ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS wallet_id UUID REFERENCES wallets(wallet_id),
    ADD COLUMN IF NOT EXISTS hold_released_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS previous_status VARCHAR(50),
    ADD COLUMN IF NOT EXISTS status_change_reason VARCHAR(255),
    ADD COLUMN IF NOT EXISTS reconciled_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS digiflazz_rc VARCHAR(10),
    ADD COLUMN IF NOT EXISTS webhook_received_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_transactions_status_created
    ON transactions(status, created_at)
    WHERE status IN ('pending', 'Initiated');

CREATE INDEX IF NOT EXISTS idx_transactions_reconciled_at
    ON transactions(reconciled_at);

CREATE INDEX IF NOT EXISTS idx_transactions_pending_timeout
    ON transactions(status, created_at, updated_at)
    WHERE status = 'pending';

COMMENT ON COLUMN transactions.wallet_id IS 'Wallet associated with the transaction (for hold tracking)';
COMMENT ON COLUMN transactions.hold_released_at IS 'Timestamp when wallet hold was released (due to expiry/failure/cancel)';
COMMENT ON COLUMN transactions.previous_status IS 'Previous status before current transition (for audit trail)';
COMMENT ON COLUMN transactions.status_change_reason IS 'Reason for status change: webhook, timeout, user_cancel, admin_action';
COMMENT ON COLUMN transactions.reconciled_at IS 'Timestamp when transaction was last reconciled by background job';
COMMENT ON COLUMN transactions.digiflazz_rc IS 'Raw response code from Digiflazz';
COMMENT ON COLUMN transactions.webhook_received_at IS 'Timestamp when webhook was received from Digiflazz';