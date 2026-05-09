-- +goose Up
-- OTP, Device Fingerprints, Audit Logs
-- ============================================

CREATE TABLE IF NOT EXISTS otp_codes (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES users(id) ON DELETE CASCADE,
    phone       VARCHAR(20) NOT NULL,
    code        VARCHAR(10) NOT NULL,
    type        VARCHAR(20) NOT NULL DEFAULT 'verification',
    attempts    INT NOT NULL DEFAULT 0 CHECK (attempts <= 5),
    salt        VARCHAR(32),
    expires_at  TIMESTAMP NOT NULL,
    used_at     TIMESTAMP,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_otp_phone_expires ON otp_codes(phone, expires_at);

CREATE TABLE IF NOT EXISTS device_fingerprints (
    id                 BIGSERIAL PRIMARY KEY,
    user_id            BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fingerprint_hash   VARCHAR(64) NOT NULL,
    user_agent         TEXT,
    ip_address         INET,
    trust_score        INT NOT NULL DEFAULT 0 CHECK (trust_score >= 0 AND trust_score <= 100),
    is_trusted         BOOLEAN NOT NULL DEFAULT FALSE,
    first_seen         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_device_user ON device_fingerprints(user_id);
CREATE INDEX IF NOT EXISTS idx_device_last_seen ON device_fingerprints(last_seen DESC);

CREATE TABLE IF NOT EXISTS audit_logs (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT REFERENCES users(id),
    action        VARCHAR(255) NOT NULL,
    resource_type VARCHAR(100),
    resource_id   BIGINT,
    old_value    JSONB,
    new_value    JSONB,
    details      JSONB,
    severity     VARCHAR(20) DEFAULT 'INFO' CHECK (severity IN ('INFO','WARN','ERROR','CRITICAL')),
    trace_id     VARCHAR(64),
    ip_address   INET,
    user_agent   TEXT,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_audit_user_action ON audit_logs(user_id, action, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_trace ON audit_logs(trace_id);
CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_resource ON audit_logs(resource_type, resource_id);

-- +goose Down
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS device_fingerprints CASCADE;
DROP TABLE IF EXISTS otp_codes CASCADE;
