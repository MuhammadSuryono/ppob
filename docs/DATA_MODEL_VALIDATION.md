# Data Model Validation & Integrity for PPOB

**Audience:** Database administrators, backend developers, QA  
**Last Updated:** 2026-05-07  
**Status:** Draft — requires DBA review and performance testing

---

## 1. Overview

This document specifies field-level validation rules, database-level constraints (CHECK, UNIQUE, FOREIGN KEY), application validation flow, and data integrity invariants for the PPOB PostgreSQL database.

**Goal:** Ensure data correctness at every layer — application → database — with defensive constraints that prevent invalid states even if application code has bugs.

---

## 2. Validation Layers Strategy

| Layer | Responsibility | Tools |
|---|---|---|
| **Client (Flutter)** | Format hints (keyboard type, required fields) | Form widgets, regex on input |
| **API Controller** | Validate JSON structure, required fields | JSON Schema validation middleware |
| **Service Layer** | Business rule validation, cross-field checks | Go `validator` struct tags, custom Validator interface |
| **Database** | Final integrity enforcement (never trust app) | CHECK constraints, NOT NULL, UNIQUE, FOREIGN KEY, triggers |

---

## 3. Table-by-Table Validation Rules

### 3.1 `users`

**Columns & Constraints:**
```sql
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL CHECK (char_length(name) >= 2),
    password_hash VARCHAR(255) NOT NULL,
    pin_hash VARCHAR(255) NOT NULL,
    pin_salt VARCHAR(32),  -- new column for argon2id salt per PIN change
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    deleted_at TIMESTAMP,  -- soft delete
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT phone_format CHECK (phone_number ~ '^\+62[0-9]{8,12}$'),
    CONSTRAINT name_no_whitespace CHECK (name !~ '^\s*$')
);
```

**Indexes:**
```sql
CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = TRUE;
```

**Application Validation (Go):**
```go
type UserCreate struct {
    PhoneNumber string `json:"phone_number" validate:"required,phone_id"`
    Name        string `json:"name" validate:"required,min=2,max=255"`
    Password    string `json:"password" validate:"required,min=8,complex"`
    PIN         string `json:"pin" validate:"required,len=6,numeric,pinformat"`
}

// Custom validators
var _ = validator.New().RegisterValidation("phone_id", func(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    return regexp.MustCompile(`^\+62[0-9]{8,12}$`).MatchString(phone)
})

var _ = validator.New().RegisterValidation("pinformat", func(fl validator.FieldLevel) bool {
    pin := fl.Field().String()
    // exactly 6 digits, not sequential, not all same
    if !regexp.MustCompile(`^[0-9]{6}$`).MatchString(pin) {
        return false
    }
    return pin != "123456" && pin != "654321" && pin != "111111" && pin != "000000"
})
```

---

### 3.2 `roles`

```sql
CREATE TABLE roles (
    role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT role_name_check CHECK (role_name IN ('Mitra', 'Staff', 'Admin'))  -- enum
);
```

**Seed Data:** Insert `Mitra`, `Staff`, `Admin` at migration time.

---

### 3.3 `user_roles`

```sql
CREATE TABLE user_roles (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    assigned_by UUID, -- NULL if self-assigned (e.g., user becomes Mitra via registration)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE RESTRICT,
    FOREIGN KEY (assigned_by) REFERENCES users(user_id) ON DELETE SET NULL,
    
    -- Business rule: cannot assign Staff role to self (assigned_by must be non-null for Staff)
    CONSTRAINT staff_must_be_assigned CHECK (
        (role_id = (SELECT role_id FROM roles WHERE role_name = 'Staff') => assigned_by IS NOT NULL)
    )
);
```

**Indexes:**
```sql
CREATE INDEX idx_user_roles_assigned_by ON user_roles(assigned_by);
```

**Trigger:** After insert, ensure user's wallet created for that role (if Mitra → main wallet; if Staff → sub-wallet linked to assigned_by Mitra's wallet).

```sql
CREATE OR REPLACE FUNCTION create_wallet_for_role()
RETURNS TRIGGER AS $$
DECLARE
    parent_wallet_id UUID;
BEGIN
    IF (SELECT role_name FROM roles WHERE role_id = NEW.role_id) = 'Mitra' THEN
        INSERT INTO wallets (wallet_id, owner_id, balance, is_main_wallet) 
        VALUES (gen_random_uuid(), NEW.user_id, 0, true);
    ELSIF (SELECT role_name FROM roles WHERE role_id = NEW.role_id) = 'Staff' THEN
        -- Find Mitra's main wallet (assigned_by)
        SELECT wallet_id INTO parent_wallet_id 
        FROM wallets 
        WHERE owner_id = NEW.assigned_by AND is_main_wallet = true;
        
        IF parent_wallet_id IS NULL THEN
            RAISE EXCEPTION 'Mitra (assigned_by) has no main wallet';
        END IF;
        
        INSERT INTO wallets (wallet_id, owner_id, balance, is_main_wallet, parent_wallet_id)
        VALUES (gen_random_uuid(), NEW.user_id, 0, false, parent_wallet_id);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_wallet_on_role_assign
AFTER INSERT ON user_roles
FOR EACH ROW EXECUTE FUNCTION create_wallet_for_role();
```

**Note:** Trigger may fail if Mitra wallet does not exist; ensure Mitra always created with wallet at registration.

---

### 3.4 `wallets`

```sql
CREATE TABLE wallets (
    wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    balance_available DECIMAL(18,2) NOT NULL DEFAULT 0 CHECK (balance_available >= 0),
    balance_held DECIMAL(18,2) NOT NULL DEFAULT 0 CHECK (balance_held >= 0),
    is_main_wallet BOOLEAN NOT NULL DEFAULT FALSE,
    parent_wallet_id UUID, -- for staff sub-wallets, references Mitra's main wallet
    is_frozen BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (owner_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_wallet_id) REFERENCES wallets(wallet_id) ON DELETE RESTRICT,
    
    -- A user has exactly one main wallet (if Mitra) or one sub-wallet (if Staff)
    UNIQUE(owner_id),
    
    -- Staff wallet must have parent; Mitra wallet must NOT have parent
    CONSTRAINT staff_parent_check CHECK (
        (is_main_wallet = false AND parent_wallet_id IS NOT NULL) OR
        (is_main_wallet = true AND parent_wallet_id IS NULL)
    )
);
```

**Indexes:**
```sql
CREATE INDEX idx_wallets_owner ON wallets(owner_id);
CREATE INDEX idx_wallets_parent ON wallets(parent_wallet_id) WHERE is_main_wallet = false;
```

**Check Constraint:** `balance_available + balance_held <= total_reconciled_balance` — enforced by application logic and reconciliation job.

---

### 3.5 `transactions`

```sql
CREATE TABLE transactions (
    transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ref_id VARCHAR(255) UNIQUE NOT NULL,
    wallet_id UUID NOT NULL,
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    customer_no VARCHAR(255) NOT NULL,
    amount DECIMAL(18,2) NOT NULL CHECK (amount >= 0), -- platform cost from Digiflazz
    selling_price DECIMAL(18,2) NOT NULL CHECK (selling_price >= 0),
    platform_price DECIMAL(18,2) NOT NULL CHECK (platform_price >= 0), -- snapshot at transaction time
    status VARCHAR(50) NOT NULL CHECK (status IN ('Initiated', 'Pending', 'Success', 'Failed', 'Expired', 'Cancelled', 'Refunded')),
    previous_status VARCHAR(50),
    status_change_reason VARCHAR(255),
    digiflazz_response JSONB,
    hold_released_at TIMESTAMP,
    reconciled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (wallet_id) REFERENCES wallets(wallet_id) ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE RESTRICT,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE RESTRICT,
    
    -- selling_price must be >= platform_price (business rule)
    CONSTRAINT selling_price_minimum CHECK (selling_price >= platform_price),
    
    -- amount should equal platform_price (or close, due to rounding)
    CONSTRAINT amount_matches_platform CHECK (abs(amount - platform_price) < 1)
);
```

**Indexes:**
```sql
CREATE INDEX idx_transactions_user_created ON transactions(user_id, created_at DESC);
CREATE INDEX idx_transactions_wallet ON transactions(wallet_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_ref_id ON transactions(ref_id);
CREATE INDEX idx_transactions_created ON transactions(created_at DESC);
```

**Partial Index for Active Queries:**
```sql
CREATE INDEX idx_transactions_recent ON transactions(created_at DESC) 
WHERE status IN ('Success', 'Pending');
```

---

### 3.6 `products`

```sql
CREATE TABLE products (
    product_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL,
    buyer_sku_code VARCHAR(255) UNIQUE NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    description TEXT,
    base_price DECIMAL(18,2) NOT NULL CHECK (base_price >= 0),
    platform_markup DECIMAL(5,2) NOT NULL DEFAULT 0.05 CHECK (platform_markup >= 0 AND platform_markup <= 1),
    platform_price DECIMAL(18,2) NOT NULL CHECK (platform_price >= base_price),
    is_prepaid BOOLEAN NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE RESTRICT
);

-- Add check: platform_price == base_price * (1 + platform_markup) within tolerance
ALTER TABLE products ADD CONSTRAINT price_consistency_check 
CHECK (abs(platform_price - (base_price * (1 + platform_markup))) < 1);

CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_sku ON products(buyer_sku_code);
CREATE INDEX idx_products_active ON products(is_active) WHERE is_active = true;
```

---

### 3.7 `categories`

```sql
CREATE TABLE categories (
    category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_name VARCHAR(255) UNIQUE NOT NULL,
    icon_url VARCHAR(500),
    display_order INT DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT category_name_not_empty CHECK (char_length(category_name) >= 1)
);
```

---

### 3.8 `mitra_product_prices` (Mitra Custom Pricing)

```sql
CREATE TABLE mitra_product_prices (
    mitra_price_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mitra_id UUID NOT NULL,
    product_id UUID NOT NULL,
    selling_price DECIMAL(18,2) NOT NULL CHECK (selling_price >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(mitra_id, product_id),
    FOREIGN KEY (mitra_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    
    CONSTRAINT selling_price_minimum CHECK (selling_price >= 
        (SELECT platform_price FROM products WHERE product_id = mitra_product_prices.product_id)
    )
);
```

**Note:** Check constraint references another table — PostgreSQL allows this via subquery in CHECK?

**No** — CHECK cannot reference other tables. Use trigger instead:

```sql
CREATE OR REPLACE FUNCTION validate_mitra_price()
RETURNS TRIGGER AS $$
DECLARE
    platform_price DECIMAL(18,2);
BEGIN
    SELECT platform_price INTO platform_price FROM products WHERE product_id = NEW.product_id;
    IF NEW.selling_price < platform_price THEN
        RAISE EXCEPTION 'selling_price cannot be less than platform_price';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_mitra_price
BEFORE INSERT OR UPDATE ON mitra_product_prices
FOR EACH ROW EXECUTE FUNCTION validate_mitra_price();
```

---

### 3.9 `staff_margin_settings` (Replaced by Two Tables)

**Original table no longer used.** Instead use:

- `staff_global_margin_settings` (see Section 5.1)
- `staff_product_margin_overrides` (see Section 5.2)

Both have appropriate CHECK constraints:

```sql
-- staff_global_margin_settings
ALTER TABLE staff_global_margin_settings ADD CONSTRAINT scheme_value_check 
CHECK (
    (scheme_type = 'MarginShare' AND value >= 0 AND value <= 100) OR
    (scheme_type = 'FixedAllowance' AND value >= 0 AND value <= 1000000)
);

-- staff_product_margin_overrides
ALTER TABLE staff_product_margin_overrides ADD CONSTRAINT override_scheme_value_check
CHECK (
    (scheme_type = 'MarginShare' AND value >= 0 AND value <= 100) OR
    (scheme_type = 'FixedAllowance' AND value >= 0 AND value <= 1000000)
);
```

---

### 3.10 `otp_codes`

```sql
CREATE TABLE otp_codes (
    otp_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL,
    code VARCHAR(10) NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT otp_format CHECK (code ~ '^[0-9]{6}$'),
    CONSTRAINT max_attempts CHECK (attempts <= 3)
);

CREATE INDEX idx_otp_phone_expires ON otp_codes(phone_number, expires_at);
```

**Trigger:** After 3 failed verification attempts, mark as `used` (or delete).

---

### 3.11 `audit_logs`

```sql
CREATE TABLE audit_logs (
    log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID, -- NULL for system actions
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(100),
    resource_id UUID,
    old_value JSONB,
    new_value JSONB,
    details JSONB,
    trace_id VARCHAR(32),
    ip_address INET,
    user_agent TEXT,
    severity VARCHAR(20) NOT NULL DEFAULT 'INFO' CHECK (severity IN ('INFO','WARN','ERROR','CRITICAL')),
    schema_version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL
);

CREATE INDEX idx_audit_user_action ON audit_logs(user_id, action, created_at DESC);
CREATE INDEX idx_audit_trace ON audit_logs(trace_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at DESC);
```

---

### 3.12 New Tables Required

Already referenced in other docs, but include here for completeness:

#### `idempotency_keys`

```sql
CREATE TABLE idempotency_keys (
    key_hash VARCHAR(64) PRIMARY KEY, -- SHA256(Idempotency-Key header)
    user_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL, -- 'transaction_initiate'
    resource_id UUID, -- e.g., transaction_id after creation
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    INDEX idx_idem_user_action (user_id, action, expires_at)
);

-- Prevent duplicate keys via PK; TTL cleanup via daily cron
CREATE INDEX idx_idem_expires ON idempotency_keys(expires_at);
```

#### `device_fingerprints`

```sql
CREATE TABLE device_fingerprints (
    device_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    fingerprint_hash VARCHAR(64) NOT NULL,
    user_agent TEXT,
    ip_address INET,
    trust_score INT NOT NULL DEFAULT 0 CHECK (trust_score >= 0 AND trust_score <= 100),
    is_trusted BOOLEAN NOT NULL DEFAULT FALSE,
    first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    INDEX idx_device_user (user_id),
    INDEX idx_device_last_seen (last_seen DESC)
);
```

#### `wallet_events` (See CONCURRENCY_CONTROL.md for full schema)

```sql
CREATE TABLE wallet_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN (
        'WalletCreated', 'Credited', 'Debited', 'Held', 'HoldReleased', 
        'TopupAdded', 'Refunded', 'Compensated'
    )),
    amount DECIMAL(18,2) NOT NULL CHECK (amount > 0),
    balance_before DECIMAL(18,2) NOT NULL,
    balance_after DECIMAL(18,2) NOT NULL,
    reference_id VARCHAR(255),
    reference_type VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    
    FOREIGN KEY (wallet_id) REFERENCES wallets(wallet_id),
    INDEX idx_wallet_events_wallet_created (wallet_id, created_at DESC)
);
```

#### `transaction_limits` (Daily tracking)

```sql
CREATE TABLE daily_limits (
    user_id UUID NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    transaction_count INT NOT NULL DEFAULT 0 CHECK (transaction_count >= 0),
    total_amount DECIMAL(18,2) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (user_id, date),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- Daily reconciliation can check this table matches transactions
```

#### `compensation_jobs`

```sql
CREATE TABLE compensation_jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','retrying','failed','completed')),
    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 10,
    next_retry_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_compensation_retry (status, next_retry_at)
);
```

#### `postpaid_inquiries`

```sql
CREATE TABLE postpaid_inquiries (
    inquiry_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ref_id VARCHAR(255) UNIQUE NOT NULL,
    product_id UUID NOT NULL,
    customer_no VARCHAR(255) NOT NULL,
    customer_name TEXT,
    bill_details JSONB NOT NULL,
    admin_amount DECIMAL(18,2) NOT NULL,
    total_amount DECIMAL(18,2) NOT NULL,
    selling_price DECIMAL(18,2) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (product_id) REFERENCES products(product_id),
    INDEX idx_postpaid_inquiry_ref (ref_id),
    INDEX idx_postpaid_inquiry_expires (expires_at)
);
```

---

## 4. Data Integrity Invariants (Must Always Hold)

### Invariant 1: Wallet Balance Consistency
```
Δ balance_available + Δ balance_held = Σ(wallet_events WHERE event_type IN ('Debited','Credited','Held','HoldReleased'))
```
**Check:** Daily reconciliation job computes sum(wallet_events) for each wallet and compares to `balance_available + balance_held`. Mismatch → CRITICAL alert.

### Invariant 2: No Negative Available Balance
```
balance_available ≥ 0
```
**Check:** CHECK constraint + application validation before every debit/hold.

### Invariant 3: Staff Commission ≤ Margin
```
commission_amount ≤ (selling_price − platform_price)
```
For all success transactions. Enforced in commission calculation service.

### Invariant 4: Daily Limit Exceedance
```
COUNT(txn WHERE user_id=? AND DATE(created_at)=CURRENT_DATE AND status='Success') 
  ≤ staff_daily_txn_limit
```
Enforced via atomic `INSERT ... ON CONFLICT` with WHERE clause in `daily_limits` table.

### Invariant 5: Transaction Amount Matches Snapshot
```
transactions.amount = products.platform_price at time of txn
```
`platform_price` stored in transaction row; future product price changes do not affect historical records.

### Invariant 6: Reference ID Uniqueness
```
SELECT COUNT(*) FROM transactions WHERE ref_id = ? ≤ 1
```
UNIQUE constraint at DB level prevents duplicates.

### Invariant 7: Wallet Ownership
```
wallets.owner_id ∈ (SELECT user_id FROM user_roles WHERE role_id matches wallet's role type)
```
Trigger ensures wallet created only when role assigned.

### Invariant 8: All Transaction Statuses are Valid
```
status ∈ {'Initiated','Pending','Success','Failed','Expired','Cancelled','Refunded'}
```
CHECK constraint.

---

## 5. Application Validation Flow

Every incoming API request passes through these stages:

```
HTTP Request
    ↓
Routing ( Gin / Echo )
    ↓
JSON Schema Validation (required fields, types)
    ↓
Struct Validation (go-playground/validator)
    ↓
Service Layer Business Validation
    ↓
Database Constraints (final safety net)
    ↓
Handler Logic
    ↓
DB Transaction Commit
```

**Example: InitiateTransaction**
```go
type InitiateTxRequest struct {
    ProductID     uuid.UUID `json:"product_id" validate:"required"`
    CustomerNo    string    `json:"customer_no" validate:"required,customer_no"`
    SellingPrice  float64   `json:"selling_price" validate:"required,gt=0"`
    PIN           string    `json:"pin" validate:"required,len=6,numeric,pinformat"`
    IdempotencyKey string  `json:"idempotency_key" validate:"required,uuid4"`
}

func (h *Handler) InitiateTransaction(w http.ResponseWriter, r *http.Request) {
    // 1. Parse & validate struct
    var req InitiateTxRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "INVALID_JSON", err.Error())
        return
    }
    if err := validator.New().Struct(req); err != nil {
        respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
        return
    }
    
    // 2. Business validation
    ctx := r.Context()
    userID := auth.UserID(ctx)
    
    // Check user active
    if !userService.IsActive(ctx, userID) {
        respondError(w, http.StatusForbidden, "USER_INACTIVE", "Account suspended")
        return
    }
    
    // Verify product exists & active
    product, err := productService.Get(ctx, req.ProductID)
    if err != nil || !product.IsActive {
        respondError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product inactive")
        return
    }
    
    // Validate selling_price >= product.platform_price (with Mitra override check)
    actualSellingPrice := product.GetSellingPriceForMitra(mitraID)
    if req.SellingPrice < actualSellingPrice {
        respondError(w, http.StatusBadRequest, "PRICE_BELOW_COST", 
            fmt.Sprintf("Minimum allowed: Rp %.0f", actualSellingPrice))
        return
    }
    
    // Check wallet balance (available)
    wallet, err := walletService.GetActive(ctx, userID, activeRoleID)
    if err != nil {
        respondError(w, http.StatusInternalServerError, "WALLET_ERROR", err.Error())
        return
    }
    if wallet.BalanceAvailable < actualSellingPrice {
        respondError(w, http.StatusBadRequest, "INSUFFICIENT_BALANCE", 
            fmt.Sprintf("Saldo tidak cukup. Saldo: Rp %.0f, dibutuhkan: Rp %.0f", 
                wallet.BalanceAvailable, actualSellingPrice))
        return
    }
    
    // Check daily limits if user is Staff
    if isStaff {
        if exceeded, limitType := dailyLimitService.Exceeded(ctx, userID, actualSellingPrice); exceeded {
            respondError(w, http.StatusBadRequest, "DAILY_LIMIT_EXCEEDED", 
                fmt.Sprintf("Limit harian %s tercapai", limitType))
            return
        }
    }
    
    // 3. Validation passed → proceed to transaction service
    // ...
}
```

---

## 6. Database Migration Scripts

**Tool:** Use `golang-migrate` or `migrate` CLI. SQL files numbered.

**Migration 001: Initial Schema** (already exists in original docs)  
**Migration 002: Add pin_salt to users**
```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS pin_salt VARCHAR(32);

-- Backfill: for existing users without PIN yet, set random salt
UPDATE users SET pin_salt = encode(gen_random_bytes(16), 'hex') WHERE pin_salt IS NULL;
```

**Migration 003: Add balance_available/held to wallets**
```sql
ALTER TABLE wallets 
    ADD COLUMN IF NOT EXISTS balance_available DECIMAL(18,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS balance_held DECIMAL(18,2) NOT NULL DEFAULT 0;

-- Copy existing balance to available
UPDATE wallets SET balance_available = balance, balance_held = 0 WHERE balance_available = 0;

-- Drop old balance column after verification
-- ALTER TABLE wallets DROP COLUMN balance;
```

**Migration 004: Create wallet_events table** (full schema seen earlier)

**Migration 005: Create idempotency_keys table**  
**Migration 006: Create device_fingerprints table**  
**Migration 007: Create daily_limits table**

```sql
-- Populate daily_limits for existing active users for today (one-time backfill)
INSERT INTO daily_limits (user_id, date, transaction_count, total_amount)
SELECT 
    u.user_id,
    CURRENT_DATE,
    COUNT(t.transaction_id) FILTER (WHERE t.status = 'Success'),
    COALESCE(SUM(t.selling_price) FILTER (WHERE t.status = 'Success'), 0)
FROM users u
JOIN user_roles ur ON u.user_id = ur.user_id AND ur.role_name = 'Staff'
LEFT JOIN transactions t ON u.user_id = t.user_id AND DATE(t.created_at) = CURRENT_DATE
GROUP BY u.user_id
ON CONFLICT (user_id, date) DO UPDATE SET 
    transaction_count = EXCLUDED.transaction_count,
    total_amount = EXCLUDED.total_amount;
```

**Migration 008: Alter staff_margin_settings → split into two tables** (data migration required)

```sql
-- 1. Create new tables
-- (run SQL from sections 5.1 and 5.2 above)

-- 2. Copy data: global settings (product_id is NULL) → staff_global_margin_settings
INSERT INTO staff_global_margin_settings (mitra_id, staff_id, scheme_type, value, created_at, updated_at)
SELECT mitra_id, staff_id, scheme_type, value, created_at, updated_at
FROM staff_margin_settings
WHERE product_id IS NULL;

-- 3. Copy data: product-specific overrides
INSERT INTO staff_product_margin_overrides (mitra_id, staff_id, product_id, scheme_type, value, created_at, updated_at)
SELECT mitra_id, staff_id, product_id, scheme_type, value, created_at, updated_at
FROM staff_margin_settings
WHERE product_id IS NOT NULL;

-- 4. Verify counts, then drop old table
-- DROP TABLE staff_margin_settings;
```

**Migration 009: Add trigger for wallet creation on role assignment**  
**Migration 010: Add trigger to validate mitra_product_prices selling_price vs platform_price**  
**Migration 011: Create postpaid_inquiries table**  
**Migration 012: Create compensation_jobs table**

**Future migrations:** Partition `transactions` by month (if volume > 1M rows).

---

## 7. Validation in Code (Go Example)

**Standard validation pipeline:**

```go
package validation

import (
    "github.com/go-playground/validator/v10"
    "regexp"
)

var validate = validator.New()

func init() {
    // Register custom validators
    _ = validate.RegisterValidation("phone_id", validatePhone)
    _ = validate.RegisterValidation("pinformat", validatePIN)
    _ = validate.RegisterValidation("customer_no", validateCustomerNo)
}

func ValidateUserCreate(req CreateUserRequest) error {
    return validate.Struct(req)
}

func validatePIN(fl validator.FieldLevel) bool {
    pin := fl.Field().String()
    if len(pin) != 6 {
        return false
    }
    if !regexp.MustCompile(`^[0-9]{6}$`).MatchString(pin) {
        return false
    }
    // Sequential check
    seq := []string{"123456", "654321", "111111", "000000"}
    for _, s := range seq {
        if pin == s {
            return false
        }
    }
    return true
}
```

**Service-level guard:**
```go
func (s *TransactionService) validateAndCreate(ctx context.Context, req InitiateRequest) error {
    // Load product with lock to prevent concurrent price change
    product := &Product{}
    err := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(product, "product_id = ?", req.ProductID).Error
    if err != nil {
        return err
    }
    
    // Re-check selling price against current platform_price (may have changed since UI fetch)
    sellingPrice := product.GetSellingPriceForMitra(req.MitraID)
    if req.SellingPrice < sellingPrice {
        return ErrSellingPriceTooLow
    }
    
    // Atomic daily limit check (see BUSINESS_LOGIC_SPEC)
    // ...
}
```

---

## 8. Test Data & Fixtures

**Seed Scripts** for dev environment:

- **Roles:** Insert `Mitra`, `Staff`, `Admin`
- **Categories:** Insert `Pulsa`, `Paket Data`, `PLN`, `PDAM`, `E-Wallet`
- **Sample Products:** 10 prepaid products (XL, Telkomsel, Axis)
- **Sample Mitra:** `+6281234567890` (phone), password `Test1234!`, PIN `123456`
- **Sample Staff:** `+6281234567891` (belongs to Mitra), password `Staff123!`, PIN `654321`

**Factory for Tests:** Use `testify` or `factory-go` to generate random valid data:
```go
func RandomUser(role string) *User {
    return &User{
        PhoneNumber: "+62" + randomDigits(8, 12),
        Name:       randomName(),
        // ...
    }
}
```

---

## 9. Data Migration & Backfill Strategy

When adding new columns or tables, provide zero-downtime migrations.

**Example: Add `balance_available`/`balance_held` to existing `wallets`:**

1. **Step 1:** Add columns nullable, default 0 (instant, no lock):
   ```sql
   ALTER TABLE wallets ADD COLUMN balance_available DECIMAL(18,2);
   ALTER TABLE wallets ADD COLUMN balance_held DECIMAL(18,2);
   ```

2. **Step 2:** Backfill in batches of 1000 to avoid long transaction:
   ```sql
   UPDATE wallets 
   SET balance_available = balance, balance_held = 0
   WHERE balance_available IS NULL
   LIMIT 1000;
   ```
   Run in loop until all rows updated.

3. **Step 3:** Add NOT NULL constraints:
   ```sql
   ALTER TABLE wallets ALTER COLUMN balance_available SET NOT NULL;
   ALTER TABLE wallets ALTER COLUMN balance_held SET NOT NULL;
   ```

4. **Step 4:** Add CHECK constraints (requires table rebuild? In PG 12+, `ADD CONSTRAINT` is instant if no violations):
   ```sql
   ALTER TABLE wallets ADD CONSTRAINT positive_available CHECK (balance_available >= 0);
   ```

5. **Step 5:** Deploy application code that writes to new columns; keep legacy column read-only until confirmed safe.

6. **Step 6:** (Optional) Drop old `balance` column after 1 month.

---

## 10. Monitoring Data Quality

**Daily Health Check SQL:**
```sql
-- 1. Orphaned transactions
SELECT COUNT(*) FROM transactions t 
LEFT JOIN wallets w ON t.wallet_id = w.wallet_id 
WHERE w.wallet_id IS NULL;

-- 2. Negative balances
SELECT wallet_id, balance_available FROM wallets WHERE balance_available < 0;

-- 3. Transactions without product
SELECT COUNT(*) FROM transactions WHERE product_id NOT IN (SELECT product_id FROM products);

-- 4. Wallet events summing mismatch
SELECT w.wallet_id, 
       COALESCE(SUM(CASE WHEN event_type IN ('Credited','TopupAdded','Refunded') THEN amount ELSE -amount END),0) AS computed,
       w.balance_available + w.balance_held AS cached
FROM wallets w
LEFT JOIN wallet_events we ON w.wallet_id = we.wallet_id
GROUP BY w.wallet_id, w.balance_available, w.balance_held
HAVING ABS(computed - (w.balance_available + w.balance_held)) > 1;  -- tolerance 1 IDR

-- 5. Duplicate ref_id (should be 0)
SELECT ref_id, COUNT(*) FROM transactions GROUP BY ref_id HAVING COUNT(*) > 1;

-- 6. Users with no wallet (should be 0)
SELECT u.user_id FROM users u LEFT JOIN wallets w ON u.user_id = w.owner_id WHERE w.wallet_id IS NULL;
```

**Automate:** Run nightly; if any query returns rows → alert finance+eng.

---

## 11. Archival & Purging Policy

| Table | Retention | Purge Strategy |
|---|---|---|
| `otp_codes` | 24 hours | TTL index + daily DELETE WHERE expires_at < NOW() |
| `idempotency_keys` | 7 days | TTL index + daily DELETE WHERE expires_at < NOW() |
| `postpaid_inquiries` | 30 days (after expiry) | Daily DELETE where expires_at + 30d < NOW() |
| `device_fingerprints` | 90 days | Archive to S3, delete from DB (last_seen < 90d ago) |
| `transactions` | 5 years | Partition by month; after 2 years, move to `transactions_archive` schema; after 5 years → S3 Glacier |
| `wallet_events` | 5 years | Same partition strategy as transactions |
| `audit_logs` | 5 years | Same partition strategy |

**Partitioning Example (`transactions`):**
```sql
-- Create parent table
CREATE TABLE transactions (
    -- columns
) PARTITION BY RANGE (created_at);

-- Create monthly partitions automatically (cron job)
CREATE TABLE transactions_2026_05 PARTITION OF transactions
FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
```

---

## 12. PII Masking Functions

For display and logs:

```sql
CREATE OR REPLACE FUNCTION mask_phone(phone TEXT) RETURNS TEXT AS $$
BEGIN
    -- +62 812-345-***-****
    IF length(phone) >= 10 THEN
        RETURN overlay(phone placing '***' from length(phone)-3 for 3);
    END IF;
    RETURN '***';
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION mask_customer_no(no TEXT) RETURNS TEXT AS $$
BEGIN
    -- Show last 4
    IF length(no) >= 4 THEN
        RETURN '****' || right(no, 4);
    END IF;
    RETURN '****';
END;
$$ LANGUAGE plpgsql IMMUTABLE;
```

---

## Appendix A — SQL: Full Invariant Check Script

```sql
-- Run as daily job
\echo '=== Wallet Balance Consistency ==='
SELECT w.wallet_id, 
       COALESCE(SUM(CASE WHEN we.event_type IN ('Credited','TopupAdded','Refunded','HoldReleased') THEN we.amount ELSE -we.amount END),0) AS computed,
       w.balance_available + w.balance_held AS cached,
       (COALESCE(SUM(CASE WHEN we.event_type IN ('Credited','TopupAdded','Refunded','HoldReleased') THEN we.amount ELSE -we.amount END),0) - (w.balance_available + w.balance_held)) AS diff
FROM wallets w
LEFT JOIN wallet_events we ON w.wallet_id = we.wallet_id
GROUP BY w.wallet_id, w.balance_available, w.balance_held
HAVING ABS(diff) > 1;

\echo '=== Negative Balances ==='
SELECT wallet_id, balance_available FROM wallets WHERE balance_available < 0;

\echo '=== Duplicate Ref IDs ==='
SELECT ref_id, COUNT(*) FROM transactions GROUP BY ref_id HAVING COUNT(*) > 1;

\echo '=== Orphaned Transactions ==='
SELECT COUNT(*) FROM transactions t LEFT JOIN wallets w ON t.wallet_id = w.wallet_id WHERE w.wallet_id IS NULL;

\echo '=== Orphaned Wallet Events ==='
SELECT COUNT(*) FROM wallet_events we LEFT JOIN wallets w ON we.wallet_id = w.wallet_id WHERE w.wallet_id IS NULL;
```

---

## Appendix B — Test Cases Matrix

| Table | Valid Data | Invalid Data | Expected |
|---|---|---|---|
| users.phone_number | "+6281234567890" | "081234567890" (no +62) | ❌ Reject CHECK violation |
| users.phone_number | "+628123456789012345" (too long) | >15 digits after +62 | ❌ Reject |
| wallets.balance | 0, positive | -1000 | ❌ CHECK violation |
| transactions.status | "Success" | "success" (lowercase) | ❌ CHECK enum violation |
| otp_codes.code | "123456" | "12a456" | ❌ CHECK regex |
| staff_margin_settings.value | 60 (MarginShare) | 150 (>100%) | ❌ CHECK violation |
| mitra_product_prices.selling_price | ≥ platform_price | < platform_price | ❌ Trigger error |

---

## Open Questions

1. **Should platform_price be rounded to nearest 100?** Digiflazz prices are integers. Our platform_price may have decimals due to percentage markup. Round to nearest integer? **Yes** — use `ROUND(base * (1+markup), 0)`.

2. **What precision for DECIMAL?** `DECIMAL(18,2)` handles up to Rp 999,999,999,999,999.99 — sufficient for high-value PLN bills.

3. **Currency Support?** Only IDR (Rupiah) in v1. Future multi-currency? Not planned.

4. **Soft delete vs hard delete?** Soft delete (`deleted_at`) for users, products. Hard delete for logs? No, logs are immutable.

---

**Owner:** Database Administrator + Backend Tech Lead  
**Reviewed By:** Security Architect, Finance Controller (for financial invariants)  
**Effective Date:** Upon deployment of Migration 012

---

**Related Documents:**
- `CONCURRENCY_CONTROL.md` (wallet_events usage)
- `BUSINESS_LOGIC_SPEC.md` (business rules that these constraints enforce)
- `ERROR_HANDLING.md` (validation error codes mapping)
