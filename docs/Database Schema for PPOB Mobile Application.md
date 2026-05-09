# Database Schema for PPOB Mobile Application (Updated)

This document outlines the database schema for the PPOB mobile application, designed to support multi-tenancy, user authentication, wallet management, product catalog, and transaction processing. The chosen database is PostgreSQL.

**Latest Updates:**  
- Added `pin_salt`, `balance_available`/`balance_held` for event-sourced wallets  
- Replaced `staff_margin_settings` with two-table design  
- Added new tables: `wallet_events`, `idempotency_keys`, `device_fingerprints`, `daily_limits`, `compensation_jobs`, `postpaid_inquiries`, `system_settings`  
- Added missing indexes for performance  
- Added triggers for automatic wallet creation on role assignment

---

## 1. Table: `users`

Stores user information, including Mitra and Staff.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `user_id` | UUID | PRIMARY KEY | Unique identifier for the user |
| `phone_number` | VARCHAR(20) | UNIQUE, NOT NULL | User's phone number, used for login |
| `name` | VARCHAR(255) | NOT NULL | User's full name |
| `password_hash` | VARCHAR(255) | NOT NULL | Hashed password (bcrypt) |
| `pin_hash` | VARCHAR(255) | NOT NULL | Hashed transaction PIN (argon2id) |
| `pin_salt` | VARCHAR(32) | NOT NULL | Salt for PIN hashing |
| `is_active` | BOOLEAN | NOT NULL DEFAULT TRUE | Soft delete / active flag |
| `deleted_at` | TIMESTAMP | | Soft delete timestamp |
| `active_role_id` | UUID | FOREIGN KEY (`roles.role_id`) | Currently selected role |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of user creation |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of last update |

**Additional Constraints:**
```sql
ALTER TABLE users ADD CONSTRAINT phone_format CHECK (phone_number ~ '^\+62[0-9]{8,12}$');
ALTER TABLE users ADD CONSTRAINT name_nonempty CHECK (char_length(trim(name)) >= 2);
```

**Indexes:**
```sql
CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_users_active_role ON users(active_role_id);
```

---

## 2. Table: `roles`

Defines the different roles within the system (e.g., Mitra, Staff).

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `role_id` | UUID | PRIMARY KEY | Unique identifier for the role |
| `role_name` | VARCHAR(50) | UNIQUE, NOT NULL | Name of the role (e.g., 'Mitra', 'Staff') |
| `description` | TEXT | | Description of the role |
| `is_active` | BOOLEAN | NOT NULL DEFAULT TRUE | Whether role can be assigned |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of role creation |

**Seed Data:** Insert `Mitra`, `Staff`, `Admin` rows on migration.

---

## 3. Table: `user_roles`

Links users to their assigned roles, supporting multi-role capability.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `user_id` | UUID | FOREIGN KEY (`users.user_id`), NOT NULL | Reference to the user |
| `role_id` | UUID | FOREIGN KEY (`roles.role_id`), NOT NULL | Reference to the role |
| `assigned_by` | UUID | FOREIGN KEY (`users.user_id`) | User who assigned this role (e.g., Mitra assigning staff) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of role assignment |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of last update |
| `PRIMARY KEY (user_id, role_id)` | | | Composite primary key |

**Triggers:** After INSERT, triggers `create_wallet_for_role` automatically create corresponding wallet.

**Indexes:**
```sql
CREATE INDEX idx_user_roles_assigned_by ON user_roles(assigned_by);
```

---

## 4. Table: `wallets`

Stores wallet information for Mitra and Staff.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `wallet_id` | UUID | PRIMARY KEY | Unique identifier for the wallet |
| `owner_id` | UUID | UNIQUE, NOT NULL | Reference to the user who owns the wallet |
| `balance_available` | DECIMAL(18, 2) | NOT NULL DEFAULT 0 CHECK (balance_available >= 0) | Available balance (not held) |
| `balance_held` | DECIMAL(18, 2) | NOT NULL DEFAULT 0 CHECK (balance_held >= 0) | Amount held for pending transactions |
| `is_main_wallet` | BOOLEAN | NOT NULL DEFAULT FALSE | True if Mitra's main wallet, False for staff sub-wallets |
| `parent_wallet_id` | UUID | FOREIGN KEY (`wallets.wallet_id`) | For staff wallets, Mitra's main wallet reference |
| `is_frozen` | BOOLEAN | NOT NULL DEFAULT FALSE | Freeze flag for suspicious activity |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of wallet creation |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of last update |

**Constraints:**
```sql
-- Staff must have parent; Mitra must not
ALTER TABLE wallets ADD CONSTRAINT staff_parent_check CHECK (
    (is_main_wallet = false AND parent_wallet_id IS NOT NULL) OR
    (is_main_wallet = true AND parent_wallet_id IS NULL)
);
```

**Indexes:**
```sql
CREATE INDEX idx_wallets_owner ON wallets(owner_id);
CREATE INDEX idx_wallets_parent ON wallets(parent_wallet_id) WHERE is_main_wallet = false;
```

---

## 5. Table: `wallet_events`

Immutable event log for all wallet balance changes (event sourcing). Balance computed on-read from this table.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `event_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Unique event identifier |
| `wallet_id` | UUID | NOT NULL FOREIGN KEY REFERENCES wallets(wallet_id) | Wallet affected |
| `event_type` | VARCHAR(50) | NOT NULL CHECK (event_type IN ('WalletCreated','Credited','Debited','Held','HoldReleased','TopupAdded','Refunded','Compensated')) | Type of balance change |
| `amount` | DECIMAL(18,2) | NOT NULL CHECK (amount > 0) | Amount involved (always positive; sign from event_type) |
| `balance_before` | DECIMAL(18,2) | NOT NULL | Balance before this event |
| `balance_after` | DECIMAL(18,2) | NOT NULL | Balance after this event |
| `reference_id` | VARCHAR(255) | | Linked transaction_id / topup_id / hold_id |
| `reference_type` | VARCHAR(50) | | 'transaction', 'topup', 'hold', etc. |
| `metadata` | JSONB | | Additional context (product_id, rc_code, etc.) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Event timestamp |
| `created_by` | UUID | FOREIGN KEY (users.user_id) | Actor who triggered (or system) |

**Indexes:**
```sql
CREATE INDEX idx_wallet_events_wallet_created ON wallet_events(wallet_id, created_at DESC);
CREATE INDEX idx_wallet_events_reference ON wallet_events(reference_id) WHERE reference_id IS NOT NULL;
```

---

## 6. Table: `transactions`

Records all PPOB transactions.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `transaction_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Unique identifier for the transaction |
| `ref_id` | VARCHAR(255) | UNIQUE, NOT NULL | Unique reference ID from Digiflazz |
| `wallet_id` | UUID | FOREIGN KEY (`wallets.wallet_id`), NOT NULL | Wallet from which funds were deducted |
| `user_id` | UUID | FOREIGN KEY (`users.user_id`), NOT NULL | User who initiated the transaction |
| `product_id` | UUID | FOREIGN KEY (`products.product_id`), NOT NULL | Product purchased |
| `customer_no` | VARCHAR(255) | NOT NULL | Target customer number/ID |
| `amount` | DECIMAL(18, 2) | NOT NULL CHECK (amount >= 0) | Amount (platform price) |
| `selling_price` | DECIMAL(18, 2) | NOT NULL CHECK (selling_price >= 0) | Price charged to end-user |
| `platform_price` | DECIMAL(18, 2) | NOT NULL CHECK (platform_price >= 0) | Snapshot of product.platform_price at transaction time |
| `status` | VARCHAR(50) | NOT NULL CHECK (status IN ('Initiated','Pending','Success','Failed','Expired','Cancelled','Refunded')) | Current status |
| `previous_status` | VARCHAR(50) | | Prior status before transition |
| `status_change_reason` | VARCHAR(255) | | Why status changed (e.g., 'webhook_rc_00') |
| `digiflazz_response` | JSONB | | Raw response from Digiflazz |
| `hold_released_at` | TIMESTAMP | | When hold (if any) was released |
| `reconciled_at` | TIMESTAMP | | When balance reconciliation verified this txn |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of initiation |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Timestamp of last update |

**Additional CHECK:**
```sql
ALTER TABLE transactions ADD CONSTRAINT selling_price_minimum CHECK (selling_price >= platform_price);
ALTER TABLE transactions ADD CONSTRAINT amount_matches_platform 
    CHECK (abs(amount - platform_price) < 1); -- allow rounding diff < 1
```

**Indexes:**
```sql
CREATE INDEX idx_transactions_user_created ON transactions(user_id, created_at DESC);
CREATE INDEX idx_transactions_wallet ON transactions(wallet_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_ref_id ON transactions(ref_id);
CREATE INDEX idx_transactions_created ON transactions(created_at DESC);
CREATE INDEX idx_transactions_recent ON transactions(created_at DESC) 
    WHERE status IN ('Success','Pending');
```

---

## 7. Table: `products`

Stores information about digital products.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `product_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Unique identifier |
| `category_id` | UUID | FOREIGN KEY (`categories.category_id`), NOT NULL | Category reference |
| `buyer_sku_code` | VARCHAR(255) | UNIQUE, NOT NULL | Digiflazz SKU code |
| `product_name` | VARCHAR(255) | NOT NULL | Display name |
| `description` | TEXT | | Description |
| `base_price` | DECIMAL(18, 2) | NOT NULL CHECK (base_price >= 0) | Price from Digiflazz |
| `platform_markup` | DECIMAL(5, 2) | NOT NULL DEFAULT 0.05 CHECK (platform_markup >= 0 AND platform_markup <= 1) | Markup decimal (5% = 0.05) |
| `platform_price` | DECIMAL(18, 2) | NOT NULL CHECK (platform_price >= base_price) | base_price × (1 + platform_markup) |
| `is_prepaid` | BOOLEAN | NOT NULL | True = prepaid, False = postpaid |
| `is_active` | BOOLEAN | NOT NULL DEFAULT TRUE | Soft delete flag |
| `last_sync_at` | TIMESTAMP | | When last fetched from Digiflazz |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation time |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update |

**Constraints:**
```sql
ALTER TABLE products ADD CONSTRAINT price_consistency_check 
    CHECK (abs(platform_price - (base_price * (1 + platform_markup))) < 1);
```

**Indexes:**
```sql
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_sku ON products(buyer_sku_code);
CREATE INDEX idx_products_active ON products(is_active) WHERE is_active = true;
```

---

## 8. Table: `categories`

Organizes products into categories.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `category_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Unique identifier |
| `category_name` | VARCHAR(255) | UNIQUE, NOT NULL | Category name |
| `icon_url` | VARCHAR(500) | | Icon asset URL |
| `display_order` | INT DEFAULT 0 | | Sort order |
| `is_active` | BOOLEAN | NOT NULL DEFAULT TRUE | Active flag |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

---

## 9. Table: `staff_global_margin_settings`

Default margin/revenue sharing scheme for staff (applies to all products unless overridden).

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `setting_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Setting ID |
| `mitra_id` | UUID | FOREIGN KEY (`users.user_id`), NOT NULL | Mitra owner |
| `staff_id` | UUID | FOREIGN KEY (`users.user_id`), NOT NULL | Staff member |
| `scheme_type` | VARCHAR(20) | NOT NULL CHECK (scheme_type IN ('MarginShare','FixedAllowance')) | Scheme |
| `value` | DECIMAL(10,2) | NOT NULL | MarginShare: 0–100 (%), FixedAllowance: 0–1,000,000 (IDR) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Unique:** `UNIQUE(mitra_id, staff_id)`

**Check:**
```sql
ALTER TABLE staff_global_margin_settings ADD CONSTRAINT scheme_value_check 
    CHECK (
        (scheme_type = 'MarginShare' AND value >= 0 AND value <= 100) OR
        (scheme_type = 'FixedAllowance' AND value >= 0 AND value <= 1000000)
    );
```

---

## 10. Table: `staff_product_margin_overrides`

Per-product overrides for staff margin (overrides global setting).

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `override_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Override ID |
| `mitra_id` | UUID | FOREIGN KEY (`users.user_id`), NOT NULL | Mitra owner |
| `staff_id` | UUID | FOREIGN KEY (`users.user_id`), NOT NULL | Staff member |
| `product_id` | UUID | FOREIGN KEY (`products.product_id`), NOT NULL | Specific product |
| `scheme_type` | VARCHAR(20) | NOT NULL | 'MarginShare' or 'FixedAllowance' |
| `value` | DECIMAL(10,2) | NOT NULL | Same range as global |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Unique:** `UNIQUE(mitra_id, staff_id, product_id)`

**Check:** Same as global.

**Trigger:** Before insert/update, validate `selling_price >= platform_price` (if needed for FixedAllowance? No, value is allowance, not price).

---

## 11. Table: `otp_codes`

Stores OTP codes for verification.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `otp_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | OTP record ID |
| `phone_number` | VARCHAR(20) | NOT NULL | Phone number |
| `code` | VARCHAR(10) | NOT NULL | Hashed OTP (bcrypt) |
| `attempts` | INT | NOT NULL DEFAULT 0 CHECK (attempts <= 3) | Failed verification attempts |
| `expires_at` | TIMESTAMP | NOT NULL | Expiration time |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Constraints:**
```sql
ALTER TABLE otp_codes ADD CONSTRAINT otp_format CHECK (code ~ '^[0-9]{6}$');
CREATE INDEX idx_otp_phone_expires ON otp_codes(phone_number, expires_at);
```

---

## 12. Table: `audit_logs`

Records all significant actions for auditing.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `log_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Log entry ID |
| `user_id` | UUID | FOREIGN KEY (`users.user_id`) | Actor (NULL for system) |
| `action` | VARCHAR(255) | NOT NULL | Action performed |
| `resource_type` | VARCHAR(100) | | Type of resource affected |
| `resource_id` | UUID | | ID of affected resource |
| `old_value` | JSONB | | State before change |
| `new_value` | JSONB | | State after change |
| `details` | JSONB | | Additional context |
| `trace_id` | VARCHAR(32) | | OpenTelemetry trace ID |
| `ip_address` | INET | | Actor IP |
| `user_agent` | TEXT | | Browser/app UA |
| `severity` | VARCHAR(20) | DEFAULT 'INFO' CHECK (severity IN ('INFO','WARN','ERROR','CRITICAL')) | Severity |
| `schema_version` | INT DEFAULT 1 | | For future log format changes |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Indexes:**
```sql
CREATE INDEX idx_audit_user_action ON audit_logs(user_id, action, created_at DESC);
CREATE INDEX idx_audit_trace ON audit_logs(trace_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at DESC);
```

---

## 13. Table: `idempotency_keys`

Prevents duplicate transaction initiation.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `key_hash` | VARCHAR(64) | PRIMARY KEY | SHA256(Idempotency-Key header) |
| `user_id` | UUID | NOT NULL FOREIGN KEY REFERENCES users(user_id) | User who initiated |
| `action` | VARCHAR(50) | NOT NULL | e.g., 'transaction_initiate' |
| `resource_id` | UUID | | Created resource ID (transaction_id) |
| `expires_at` | TIMESTAMP | NOT NULL | Expiry (typically 24h) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Indexes:**
```sql
CREATE INDEX idx_idem_user_action ON idempotency_keys(user_id, action, expires_at);
CREATE INDEX idx_idem_expires ON idempotency_keys(expires_at);
```

---

## 14. Table: `device_fingerprints`

Stores device trust data for anomaly detection.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `device_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Persistent device identifier |
| `user_id` | UUID | NOT NULL FOREIGN KEY REFERENCES users(user_id) ON DELETE CASCADE | Owner |
| `fingerprint_hash` | VARCHAR(64) | NOT NULL | SHA256 of combined signals |
| `user_agent` | TEXT | | Full user-agent string |
| `ip_address` | INET | | Last known IP |
| `trust_score` | INT | NOT NULL DEFAULT 0 CHECK (trust_score >= 0 AND trust_score <= 100) | 0–100 score |
| `is_trusted` | BOOLEAN | NOT NULL DEFAULT FALSE | ≥70? |
| `first_seen` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |
| `last_seen` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Indexes:**
```sql
CREATE INDEX idx_device_user ON device_fingerprints(user_id);
CREATE INDEX idx_device_last_seen ON device_fingerprints(last_seen DESC);
```

---

## 15. Table: `daily_limits`

Atomic enforcement of per-staff daily transaction count/amount limits.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `user_id` | UUID | NOT NULL FOREIGN KEY REFERENCES users(user_id) ON DELETE CASCADE | Staff user |
| `date` | DATE | NOT NULL DEFAULT CURRENT_DATE | Business date (UTC+7) |
| `transaction_count` | INT | NOT NULL DEFAULT 0 CHECK (transaction_count >= 0) | Success count today |
| `total_amount` | DECIMAL(18,2) | NOT NULL DEFAULT 0 CHECK (total_amount >= 0) | Sum of selling_price today |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Primary Key:** `(user_id, date)`

**Usage:** Atomic increment via `INSERT ... ON CONFLICT DO UPDATE` prevents race conditions.

---

## 16. Table: `compensation_jobs`

Tracks failed multi-step operations needing compensation.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `job_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Job ID |
| `job_type` | VARCHAR(50) | NOT NULL | e.g., 'mitra_topup_staff' |
| `payload` | JSONB | NOT NULL | Full context {mitra_wallet_id, staff_wallet_id, amount, failed_step} |
| `status` | VARCHAR(20) | NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','retrying','failed','completed')) | Current status |
| `retry_count` | INT | NOT NULL DEFAULT 0 | Attempts made |
| `max_retries` | INT | NOT NULL DEFAULT 10 | Max attempts |
| `next_retry_at` | TIMESTAMP | | Scheduled retry time (exponential backoff) |
| `error_message` | TEXT | | Last error |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Indexes:**
```sql
CREATE INDEX idx_compensation_retry ON compensation_jobs(status, next_retry_at);
```

---

## 17. Table: `postpaid_inquiries`

Stores postpaid bill inquiry results (valid 24h).

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `inquiry_id` | UUID | PRIMARY KEY DEFAULT gen_random_uuid() | Inquiry ID |
| `ref_id` | VARCHAR(255) | UNIQUE, NOT NULL | Same ref_id used for payment |
| `product_id` | UUID | NOT NULL FOREIGN KEY REFERENCES products(product_id) | Product |
| `customer_no` | VARCHAR(255) | NOT NULL | Bill customer number |
| `customer_name` | TEXT | | Customer name from provider |
| `bill_details` | JSONB | NOT NULL | Full Digiflazz `desc` object |
| `admin_amount` | DECIMAL(18,2) | NOT NULL | Admin fee |
| `total_amount` | DECIMAL(18,2) | NOT NULL | Total bill (excluding admin?) |
| `selling_price` | DECIMAL(18,2) | NOT NULL | What we charge user |
| `expires_at` | TIMESTAMP | NOT NULL | 24h from creation |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Indexes:**
```sql
CREATE INDEX idx_postpaid_inquiry_ref ON postpaid_inquiries(ref_id);
CREATE INDEX idx_postpaid_inquiry_expires ON postpaid_inquiries(expires_at);
```

---

## 18. Table: `system_settings`

Global configuration key-value store.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `key` | VARCHAR(100) | PRIMARY KEY | Setting key |
| `value` | TEXT | NOT NULL | Setting value (JSON string for complex) |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last modified |

**Seed Rows:**
```sql
INSERT INTO system_settings (key, value) VALUES 
('platform_markup_percent', '0.05'),
('staff_default_daily_txn_limit', '50'),
('staff_default_daily_amount_limit', '5000000'),
('pending_timeout_minutes', '15'),
('max_idempotency_ttl_hours', '24');
```

---

## 19. Deprecated Table: `staff_margin_settings`

**Deprecated.** This table has been replaced by `staff_global_margin_settings` and `staff_product_margin_overrides`.

| Column Name | Data Type | Status | Replacement |
|---|---|---|---|
| `product_id` (nullable) | UUID | Deprecated | Split into global (NULL) and override (NOT NULL) |
| `scheme_type` | VARCHAR(50) | Moved | Same |
| `value` | DECIMAL(5,2) | Moved | Same |

**Migration Path:**
```sql
-- Migrate global settings (product_id IS NULL)
INSERT INTO staff_global_margin_settings (mitra_id, staff_id, scheme_type, value, created_at, updated_at)
SELECT mitra_id, staff_id, scheme_type, value, created_at, updated_at
FROM staff_margin_settings
WHERE product_id IS NULL;

-- Migrate product-specific overrides
INSERT INTO staff_product_margin_overrides (mitra_id, staff_id, product_id, scheme_type, value, created_at, updated_at)
SELECT mitra_id, staff_id, product_id, scheme_type, value, created_at, updated_at
FROM staff_margin_settings
WHERE product_id IS NOT NULL;

-- After verification, drop old table:
-- DROP TABLE staff_margin_settings;
```

---

## 20. Triggers

### 20.1 Auto-Create Wallet on Role Assignment

See `CONCURRENCY_CONTROL.md` for full trigger function `create_wallet_for_role()`.

**Logic:**
- Mitra role → create main wallet (`is_main_wallet=true`)
- Staff role → create sub-wallet with `parent_wallet_id` from assigned_by Mitra

### 20.2 Validate Mitra Product Price

```sql
CREATE OR REPLACE FUNCTION validate_mitra_price()
RETURNS TRIGGER AS $$
DECLARE
    platform_price DECIMAL(18,2);
BEGIN
    SELECT p.platform_price INTO platform_price 
    FROM products p WHERE p.product_id = NEW.product_id;
    
    IF NEW.selling_price < platform_price THEN
        RAISE EXCEPTION 'selling_price cannot be less than platform_price (Rp %, got Rp %)', 
            platform_price, NEW.selling_price;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_mitra_price
BEFORE INSERT OR UPDATE ON mitra_product_prices
FOR EACH ROW EXECUTE FUNCTION validate_mitra_price();
```

### 20.3 Update `updated_at` Timestamp ( reusable)

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to tables: users, roles, user_roles, wallets, transactions, products, categories, 
-- mitra_product_prices, staff_global_margin_settings, staff_product_margin_overrides
-- (Use GENERATED ALWAYS AS ROW START/END in PG 10+ for system-time period tables if needed)
```

---

## 21. Partitioning Strategy (For Large Tables)

**When to partition:** `transactions` exceeds 5M rows.

**Method:** Range partitioning by `created_at` (monthly).

```sql
-- Parent table (already exists). Create partitions:
CREATE TABLE transactions_2026_05 PARTITION OF transactions
FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');

CREATE TABLE transactions_2026_06 PARTITION OF transactions
FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
```

**Cron job:** At month start, create next month's partition automatically.

**Indexes on partitions:** Same as parent; automatically created if `CREATE INDEX` on parent before partitioning.

---

## 22. Data Retention & Archival

| Table | Online Retention | Archive Strategy |
|---|---|---|
| `transactions` | 2 years hot DB | Monthly pg_dump to S3; after 2y delete from hot, keep in data warehouse |
| `wallet_events` | 5 years (partitioned) | After 5y → S3 Glacier |
| `audit_logs` | 5 years (partitioned) | After 5y → S3 Glacier (WORM) |
| `otp_codes` | 24h (TTL index) | Auto-delete via cron |
| `idempotency_keys` | 7 days | Auto-delete via cron |
| `postpaid_inquiries` | 30 days after expiry | Auto-delete via cron |
| `device_fingerprints` | 90 days | Archive last_seen >90d to `device_fingerprints_archive` table |

**Archival Job (daily):**
```sql
-- Move old transactions to archive table
INSERT INTO transactions_archive SELECT * FROM transactions 
WHERE created_at < NOW() - INTERVAL '2 years';
DELETE FROM transactions WHERE created_at < NOW() - INTERVAL '2 years';
-- Then dump archive table to S3 and drop partition
```

---

## 23. Migrations

**Tool:** `golang-migrate` or ` goose`. Migration files numbered: `001_initial.sql`, `002_add_pin_salt.sql`, etc.

**All migrations source-controlled.** Never modify existing migrations; create new for changes.

**Production migration:**
```bash
kubectl exec -it ppob-migration-job -- migrate -path /migrations -database $DATABASE_URL up
```

Or use CI/CD pipeline migration step that runs before deploying new code.

---

## 24. Summary of New Tables vs Original

| Original Table | Status | Replacement / Extension |
|---|---|---|
| `staff_margin_settings` | ⚠️ Deprecated | Split into `staff_global_margin_settings` + `staff_product_margin_overrides` |
| `wallets` | ✅ Extended | Added `balance_available`, `balance_held`, `is_frozen` |
| `users` | ✅ Extended | Added `pin_salt`, `is_active`, `deleted_at`, `active_role_id` |
| `transactions` | ✅ Extended | Added `previous_status`, `status_change_reason`, `hold_released_at`, `reconciled_at` |
| `audit_logs` | ✅ Extended | Added `trace_id`, `ip_address`, `user_agent`, `severity`, `schema_version` |

**New Tables (10):**
- `wallet_events`
- `idempotency_keys`
- `device_fingerprints`
- `daily_limits`
- `compensation_jobs`
- `postpaid_inquiries`
- `staff_global_margin_settings`
- `staff_product_margin_overrides`
- `system_settings`

---

## Appendix A — Full Index List

```sql
-- users
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone_number);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active) WHERE is_active;
CREATE INDEX IF NOT EXISTS idx_users_active_role ON users(active_role_id);

-- user_roles
CREATE INDEX IF NOT EXISTS idx_user_roles_assigned_by ON user_roles(assigned_by);

-- wallets
CREATE INDEX IF NOT EXISTS idx_wallets_owner ON wallets(owner_id);
CREATE INDEX IF NOT EXISTS idx_wallets_parent ON wallets(parent_wallet_id) WHERE is_main_wallet = false;

-- transactions
CREATE INDEX IF NOT EXISTS idx_transactions_user_created ON transactions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_wallet ON transactions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_ref_id ON transactions(ref_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created ON transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_recent ON transactions(created_at DESC) 
    WHERE status IN ('Success','Pending');

-- products
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(buyer_sku_code);
CREATE INDEX IF NOT EXISTS idx_products_active ON products(is_active) WHERE is_active = true;

-- otp_codes
CREATE INDEX IF NOT EXISTS idx_otp_phone_expires ON otp_codes(phone_number, expires_at);

-- wallet_events
CREATE INDEX IF NOT EXISTS idx_wallet_events_wallet_created ON wallet_events(wallet_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_wallet_events_reference ON wallet_events(reference_id) WHERE reference_id IS NOT NULL;

-- idempotency_keys
CREATE INDEX IF NOT EXISTS idx_idem_user_action ON idempotency_keys(user_id, action, expires_at);
CREATE INDEX IF NOT EXISTS idx_idem_expires ON idempotency_keys(expires_at);

-- device_fingerprints
CREATE INDEX IF NOT EXISTS idx_device_user ON device_fingerprints(user_id);
CREATE INDEX IF NOT EXISTS idx_device_last_seen ON device_fingerprints(last_seen DESC);

-- daily_limits
CREATE INDEX IF NOT EXISTS idx_daily_limits_user_date ON daily_limits(user_id, date DESC);
```

---

## Appendix B — Migration Checklist

- [ ] `002_add_pin_salt.sql` — add column to users
- [ ] `003_split_wallet_balance.sql` — add available/held, backfill
- [ ] `004_create_wallet_events.sql`
- [ ] `005_create_idempotency_keys.sql`
- [ ] `006_create_device_fingerprints.sql`
- [ ] `007_create_daily_limits.sql`
- [ ] `008_create_compensation_jobs.sql`
- [ ] `009_split_staff_margin_settings.sql` — create two new tables, migrate data
- [ ] `010_create_postpaid_inquiries.sql`
- [ ] `011_create_system_settings.sql` — seed defaults
- [ ] `012_add_indexes.sql` — all indexes listed above
- [ ] `013_add_triggers.sql` — wallet creation trigger, price validation trigger
- [ ] `014_add_audit_columns.sql` — extend audit_logs with trace_id etc.

**Each migration must be:**
- Backwards compatible (old code still works)
- Reversible (DOWN script provided)
- Tested on staging before prod

---

**Owner:** Database Administrator  
**Next Review:** After 3 months in production for index performance analysis (pg_stat_user_indexes)

---

**Related:** `CONCURRENCY_CONTROL.md` (event sourcing), `DATA_MODEL_VALIDATION.md` (constraints), `BUSINESS_LOGIC_SPEC.md` (field meanings)
