-- +goose Up
-- Integration tables: integration_logs, provider_configs, postpaid_inquiries, compensation_jobs

CREATE TABLE IF NOT EXISTS integration_logs (
    id            BIGSERIAL PRIMARY KEY,
    log_id        VARCHAR(36) UNIQUE NOT NULL,
    provider      VARCHAR(50) NOT NULL,
    action        VARCHAR(50) NOT NULL,
    reference_id  VARCHAR(100),
    request       TEXT,
    response      TEXT,
    status        VARCHAR(20) NOT NULL,
    error_code    VARCHAR(50),
    error_message VARCHAR(255),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_integration_provider ON integration_logs(provider);
CREATE INDEX IF NOT EXISTS idx_integration_reference ON integration_logs(reference_id);
CREATE INDEX IF NOT EXISTS idx_integration_created ON integration_logs(created_at DESC);

CREATE TABLE IF NOT EXISTS provider_configs (
    id          BIGSERIAL PRIMARY KEY,
    provider    VARCHAR(50) UNIQUE NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    config      TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS postpaid_inquiries (
    id            BIGSERIAL PRIMARY KEY,
    inquiry_id    VARCHAR(100) UNIQUE NOT NULL,
    ref_id        VARCHAR(255),
    product_id    BIGINT REFERENCES products(id),
    customer_no   VARCHAR(255) NOT NULL,
    customer_name TEXT,
    bill_details  TEXT,
    admin_amount  DECIMAL(18, 2) NOT NULL DEFAULT 0,
    total_amount  DECIMAL(18, 2) NOT NULL DEFAULT 0,
    selling_price DECIMAL(18, 2) NOT NULL DEFAULT 0,
    expires_at    TIMESTAMP NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_postpaid_ref ON postpaid_inquiries(ref_id);
CREATE INDEX IF NOT EXISTS idx_postpaid_expires ON postpaid_inquiries(expires_at);

CREATE TABLE IF NOT EXISTS compensation_jobs (
    id             BIGSERIAL PRIMARY KEY,
    job_id         VARCHAR(100) NOT NULL UNIQUE,
    transaction_id VARCHAR(100),
    action         VARCHAR(50) NOT NULL,
    payload       TEXT,
    status         VARCHAR(20) NOT NULL DEFAULT 'pending',
    retry_count    INT NOT NULL DEFAULT 0,
    max_retries    INT NOT NULL DEFAULT 3,
    error_message  TEXT,
    next_retry_at  TIMESTAMP,
    completed_at   TIMESTAMP,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_compensation_status_retry ON compensation_jobs(status, next_retry_at);

-- +goose Down
DROP TABLE IF EXISTS compensation_jobs CASCADE;
DROP TABLE IF EXISTS postpaid_inquiries CASCADE;
DROP TABLE IF EXISTS provider_configs CASCADE;
DROP TABLE IF EXISTS integration_logs CASCADE;
