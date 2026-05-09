-- +goose Up
-- Transaction support tables: wallet_events, daily_limits, idempotency_keys, transaction_status_events

CREATE TABLE IF NOT EXISTS wallet_events (
    id              BIGSERIAL PRIMARY KEY,
    wallet_id       BIGINT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    event_type      VARCHAR(50) NOT NULL,
    amount          DECIMAL(18, 2) NOT NULL DEFAULT 0,
    balance_before  DECIMAL(18, 2) NOT NULL DEFAULT 0,
    balance_after   DECIMAL(18, 2) NOT NULL DEFAULT 0,
    reference_id    VARCHAR(100),
    reference_type  VARCHAR(50),
    metadata        TEXT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_wallet_events_wallet ON wallet_events(wallet_id);
CREATE INDEX IF NOT EXISTS idx_wallet_events_created ON wallet_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_wallet_events_reference ON wallet_events(reference_id, reference_type);

CREATE TABLE IF NOT EXISTS daily_limits (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL,
    date          VARCHAR(10) NOT NULL,
    count         INT NOT NULL DEFAULT 0,
    total_amount  DECIMAL(18, 2) NOT NULL DEFAULT 0,
    max_count     INT NOT NULL DEFAULT 100,
    max_amount    DECIMAL(18, 2) NOT NULL DEFAULT 10000000,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, date)
);

CREATE INDEX IF NOT EXISTS idx_daily_limits_user_date ON daily_limits(user_id, date);

CREATE TABLE IF NOT EXISTS idempotency_keys (
    id              BIGSERIAL PRIMARY KEY,
    key_value       VARCHAR(255) NOT NULL UNIQUE,
    transaction_id  VARCHAR(100),
    user_id         BIGINT NOT NULL,
    action          VARCHAR(50) NOT NULL,
    expires_at     TIMESTAMP NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_idem_user_action ON idempotency_keys(user_id, action, expires_at);

CREATE TABLE IF NOT EXISTS transaction_status_events (
    id              BIGSERIAL PRIMARY KEY,
    transaction_id  BIGINT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    event_type      VARCHAR(50) NOT NULL,
    status          VARCHAR(50) NOT NULL,
    metadata        TEXT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_txn_events_txn ON transaction_status_events(transaction_id);
CREATE INDEX IF NOT EXISTS idx_txn_events_created ON transaction_status_events(created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS transaction_status_events CASCADE;
DROP TABLE IF EXISTS idempotency_keys CASCADE;
DROP TABLE IF EXISTS daily_limits CASCADE;
DROP TABLE IF EXISTS wallet_events CASCADE;
