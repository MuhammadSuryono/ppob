# 💼 Business Logic Specification for PPOB Application

**Audience:** Product managers, backend developers, QA  
**Last Updated:** 2026-05-07  
**Status:** Draft — subject to legal/finance review

---

## 1. Overview

This document defines all business rules, formulas, validation logic, and decision matrices for the PPOB application. It serves as the single source of truth for financial calculations, role behaviors, and edge-case handling.

---

## 2. Pricing & Margin Calculation

### 2.1 Price Hierarchy

```
Digiflazz Base Price (from pricelist)
    ↓ + Platform Markup (configurable %)
Platform Price (price charged to Mitra)
    ↓ + Mitra Markup (custom per Mitra per product, optional)
Mitra Selling Price (price charged to end-user by Mitra/Staff)
```

**Formulas:**

1. **Platform Price:** `platform_price = base_price × (1 + platform_markup_percent)`
   - Example: base_price = Rp 25,000, platform_markup = 5% → platform_price = Rp 26,250
   - `platform_markup_percent` stored in config table (global, e.g., 0.05)

2. **Mitra Selling Price:**
   - If Mitra set custom price for this product → use `mitra_product_prices.selling_price`
   - Else if Mitra has default markup percent → `platform_price × (1 + mitra_default_markup_percent)`
   - Else → use `platform_price` (no additional margin)

3. **Profit & Commission:**
   ```
   Margin = Mitra Selling Price − Platform Price
   
   If scheme = FixedAllowance:
       Staff Commission = fixed_amount (per config)
       Mitra Profit = Mitra Selling Price − Platform Price − Staff Commission
   If scheme = MarginShare:
       Staff Commission = Margin × (staff_percentage / 100)
       Mitra Profit = Margin × (100 − staff_percentage) / 100
   ```

**Validation:** `selling_price ≥ platform_price` (enforced at API layer). Reject transaction if selling_price < platform_price (configurable override for promotions? Not initially).

### 2.2 Margin Scheme Configuration

**Tables:** Two tables replace original `staff_margin_settings` split:

1. `staff_global_margin_settings` (default for all products):
   ```sql
   CREATE TABLE staff_global_margin_settings (
       setting_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       mitra_id UUID NOT NULL,
       staff_id UUID NOT NULL,
       scheme_type VARCHAR(20) NOT NULL CHECK (scheme_type IN ('MarginShare', 'FixedAllowance')),
       value DECIMAL(10,2) NOT NULL, -- For MarginShare: 0–100 (%); For FixedAllowance: 0–1,000,000 (IDR)
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       UNIQUE(mitra_id, staff_id),
       FOREIGN KEY (mitra_id) REFERENCES users(user_id),
       FOREIGN KEY (staff_id) REFERENCES users(user_id)
   );
   ```

2. `staff_product_margin_overrides` (per-product overrides):
   ```sql
   CREATE TABLE staff_product_margin_overrides (
       override_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       mitra_id UUID NOT NULL,
       staff_id UUID NOT NULL,
       product_id UUID NOT NULL,
       scheme_type VARCHAR(20) NOT NULL,
       value DECIMAL(10,2) NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       UNIQUE(mitra_id, staff_id, product_id),
       FOREIGN KEY (mitra_id) REFERENCES users(user_id),
       FOREIGN KEY (staff_id) REFERENCES users(user_id),
       FOREIGN KEY (product_id) REFERENCES products(product_id)
   );
   ```

**Lookup Logic (per transaction):**
```go
func GetStaffCommission(staffID, productID uuid.UUID) (scheme string, value float64, err error) {
    // 1. Check product-specific override
    override := &StaffProductMarginOverride{}
    err = db.Where("staff_id = ? AND product_id = ?", staffID, productID).First(&override).Error
    if err == nil {
        return override.SchemeType, override.Value, nil
    }
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return "", 0, err
    }

    // 2. Fall back to global setting
    global := &StaffGlobalMarginSetting{}
    err = db.Where("staff_id = ?", staffID).First(&global).Error
    if err != nil {
        return "", 0, ErrNoMarginSetting
    }
    return global.SchemeType, global.Value, nil
}
```

**Default on Staff Creation:** `FixedAllowance` with `value = 0` (Mitra must set value later). Or pre-configure default per Mitra in `mitra_settings` table.

---

## 3. Wallet & Balance Rules

### 3.1 Wallet Types

| Wallet Type | Owner | Can Debit? | Top-up Source | Parent Wallet |
|---|---|---|---|---|
| Mitra Main Wallet | Mitra | Yes (transactions) | Direct bank transfer? Not in initial scope; manual admin credit | NULL |
| Staff Sub-Wallet | Staff | Yes (transactions) | From Mitra main wallet only (`POST /wallets/topup-staff`) | Mitra main wallet ID |

**Important:** Staff wallet cannot be topped up directly via bank transfer — only from Mitra's main wallet. This enforces Mitra as funding source.

### 3.2 Balance Fields

**`wallets` Table:**
```sql
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS balance_available DECIMAL(18,2) NOT NULL DEFAULT 0;
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS balance_held DECIMAL(18,2) NOT NULL DEFAULT 0;
-- balance_available = spendable now
-- balance_held = reserved for pending transactions (will be debited or released)
-- Total visible to user: balance_available + balance_held (or show separately)
-- Event sourcing: events compute these; these columns are cache only
```

**Update Path:** All balance changes go through `wallet_events` table. Triggers or application logic updates `balance_available`/`balance_held` synchronously within same DB transaction as event insert.

### 3.3 Hold Lifecycle

1. **Place Hold** (transaction initiated, RC 03 received or pre-debit):
   - `UPDATE wallets SET balance_held = balance_held + amount WHERE wallet_id = ? AND balance_available >= amount`
   - Check affected rows > 0 (sufficient available balance)
   - Insert `wallet_event` type `Held` with `balance_before`, `balance_after` (available decreases, held increases)
   - Return hold ID (use `ref_id`)

2. **Convert Hold to Debit** (transaction succeeds):
   - `UPDATE wallets SET balance_held = balance_held - amount, balance_available = balance_available - amount WHERE wallet_id = ?`
   - Insert `wallet_event` type `Debited` (amount same as held, now permanently deducted)
   - Link to transaction `ref_id`

3. **Release Hold** (transaction failed/expired/cancelled):
   - `UPDATE wallets SET balance_held = balance_held - amount WHERE wallet_id = ?`
   - Insert `wallet_event` type `HoldReleased` (amount returns to available)
   - Link to transaction `ref_id`

**Atomicity:** Hold placement and transaction creation must be in same DB transaction to avoid orphan holds.

### 3.4 Overdraft Prohibited

`balance_available` cannot go negative. Any debit or hold that would make it <0 is rejected with error `INSUFFICIENT_BALANCE`.

---

## 4. Transaction Flow Business Rules

### 4.1 Pre-Transaction Validation

**Before calling Digiflazz:**
1. ✅ User is active (not suspended)
2. ✅ Product is active (`products.is_active = true`)
3. ✅ Product `buyer_sku_code` matches allowed pattern for category
4. ✅ Wallet has sufficient `balance_available` for `selling_price` (not `platform_price` — user pays selling price)
5. ✅ Staff daily limit not exceeded (if staff role)
6. ✅ No duplicate `idempotency_key` already used
7. ✅ `selling_price >= platform_price` (no negative margin)
8. ✅ Customer number format valid per product (regex in `product_validation_patterns` table)

### 4.2 During Transaction

- Call Digiflazz with timeout 30s
- If network error or retryable RC → retry up to 3 times
- If all retries fail → transaction status = `Failed`
- If RC 03 "Pending" → status = `Pending`, hold remains

### 4.3 Post-Transaction

**Success (RC 00):**
- Convert hold to debit
- Calculate commission:
  ```sql
  INSERT INTO commissions (user_id, transaction_id, amount, scheme, earned_at)
  VALUES (?, ?, ?, ?, NOW())
  ```
- Update `staff_daily_usage` counters
- Push notification to user

**Failure (RC != 00/03):**
- Release hold immediately
- No commission
- Notify user with reason

**Pending (RC 03):**
- Hold stays until webhook update
- If webhook not received within 15min → reconciliation job marks `Expired`, releases hold

---

## 5. Staff Daily Limits (Fraud Detection)

### 5.1 Limit Structure

Two independent limits (whichever hits first stops staff):

| Limit Type | Default Value | Unit | Purpose |
|---|---|---|---|
| Transaction Count | 50 | count/day | Prevent spam/fat-finger |
| Daily Turnover | 5,000,000 | IDR/day | Cap financial exposure |

**Per-staff override:** Mitra can set custom limits per staff in `staff_limit_overrides` table (or add columns to `staff_global_margin_settings`).

### 5.2 Enforcement

**Atomic Check-and-Increment:**

```sql
WITH new_usage AS (
    INSERT INTO daily_limits (user_id, date, transaction_count, total_amount)
    VALUES ($1, CURRENT_DATE, 1, $2)
    ON CONFLICT (user_id, date)
    DO UPDATE SET
        transaction_count = daily_limits.transaction_count + 1,
        total_amount = daily_limits.total_amount + EXCLUDED.total_amount
    WHERE 
        daily_limits.transaction_count < 50
        AND daily_limits.total_amount + EXCLUDED.total_amount <= 5000000
    RETURNING transaction_count, total_amount
)
SELECT * FROM new_usage;
```

If query returns 0 rows → limit exceeded, reject transaction with error code `DAILY_LIMIT_EXCEEDED`.

**Reset:** Daily limits reset at midnight Indonesia timezone (UTC+7). Cron job at 00:00 could clear counters, but better to use `DATE` column with `CURRENT_DATE` automatically groups by day.

**View Current Usage:**
```sql
SELECT transaction_count, total_amount 
FROM daily_limits 
WHERE user_id = $1 AND date = CURRENT_DATE;
```

---

## 6. Multi-Role & Multi-Tenant Behavior

### 6.1 Active Role Determines Context

After login, user selects an active role (from their assigned roles). This active role dictates:

- Which wallet is used for transactions (`wallets` joined via `users.user_id` + active `role_id`)
- Which staff they can manage (Mitra role → can manage staff assigned by them)
- Which reports they see (Mitra sees all staff; Staff sees only self)

**Role Switch Request:**
```http
POST /users/switch-role
{
  "role_id": "uuid_of_target_role"
}
```

**Server validates:**
- User has this role in `user_roles` table
- If switching to Staff role, verify `user_roles.assigned_by` exists (which Mitra they belong to)

**Response:** New JWT with updated `role` claim and `wallet_id` for that role.

### 6.2 Wallet Resolution by Role

```go
func GetActiveWalletForRole(userID, roleID uuid.UUID) (*Wallet, error) {
    wallet := &Wallet{}
    err := db.Joins("JOIN user_roles ur ON wallets.owner_id = ur.user_id").
        Where("wallets.owner_id = ? AND ur.role_id = ?", userID, roleID).
        First(wallet).Error
    return wallet, err
}
```

**Note:** Each user-role combination should have exactly one wallet (enforced at role assignment time).

### 6.3 Data Isolation

**Mitra sees:**
- Their own wallet
- All staff assigned to them (`user_roles.assigned_by = mitra_user_id`)
- Transactions from: (their own wallet) + (all staff wallets under them)

**Staff sees:**
- Their own wallet only
- Their own transactions only

**Query filter pattern:**
```go
// Mitra query
txns := db.Where("wallet_id IN (SELECT wallet_id FROM wallets WHERE owner_id IN (SELECT user_id FROM user_roles WHERE assigned_by = ?))", mitraID)

// Staff query
txns := db.Where("wallet_id = ?", staffWalletID)
```

---

## 7. Postpaid Inquiry Validity

**Rule:** Postpaid inquiry result is valid for **24 hours** only. After that, must re-inquire before payment.

**Implementation:**
- Store inquiry in `postpaid_inquiries` table with `expires_at = created_at + INTERVAL '24 hours'`
- Payment endpoint `POST /transactions/pay-postpaid` checks `expires_at > NOW()` before allowing
- If expired → return error `INQUIRY_EXPIRED`, require user to re-inquire (fresh `ref_id`)

**Auto-cleanup:** Cron job every hour deletes expired inquiries:
```sql
DELETE FROM postpaid_inquiries WHERE expires_at < NOW();
```

---

## 8. Refund Handling

**Source:** Digiflazz sends refund webhook if transaction later reversed (rare).

**Refund Flow:**
1. Webhook received with `status="Refund"` or RC 74
2. System credits original wallet with refunded amount
3. If commission already paid to staff, issue **negative commission** to reverse
4. Status → `Refunded`
5. Notify user

**Negative Commission:**
```sql
INSERT INTO commissions (user_id, transaction_id, amount, type, description)
VALUES (staff_id, tx_id, -commission_amount, 'refund', 'Refund for transaction')
```

**Mitra Profit Impact:** Refund reduces Mitra profit (original selling_price - platform_price − commission). If refund is full, Mitra gets 0 profit from that transaction.

---

## 9. Commission Calculation & Payout

### 9.1 Calculation Timing

**Real-time:** Commission calculated at transaction success and recorded immediately in `commissions` table.

**Payout:** Nightly batch job at 02:00 AM transfers accumulated commissions to staff wallet balance (or separate `commission_balance`?). 

**Simpler Approach:** Commission directly credited to staff wallet at success time (no separate accumulation). This is immediate and transparent.

**Decision:** **Credit staff wallet immediately** upon transaction success. No separate accrual → payout table needed.

### 9.2 Commission Record

```sql
CREATE TABLE commissions (
    commission_id UUID PRIMARY KEY,
    staff_id UUID NOT NULL,
    transaction_id UUID NOT NULL,
    amount DECIMAL(18,2) NOT NULL, -- positive for credit, negative for refund
    scheme_used VARCHAR(20) NOT NULL, -- 'MarginShare' or 'FixedAllowance'
    scheme_value DECIMAL(10,2) NOT NULL, -- % or fixed amount
    margin_amount DECIMAL(18,2), -- total margin (selling - platform) for reference
    earned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMP, -- NULL until transferred (if using accrual model)
    FOREIGN KEY (staff_id) REFERENCES users(user_id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
);
```

**If immediate wallet credit:**
```go
func (s *Service) creditCommission(txID uuid.UUID) error {
    // Get transaction + staff + commission amount
    // ...
    // Credit staff wallet (in same DB transaction as commission record insert)
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. Insert commission record
        // 2. Insert wallet_event type 'Credited' with reference=commission_id
        // 3. Update wallets.balance_available
        return nil
    })
}
```

---

## 10. Platform Markup Policy

`platform_markup_percent` is **global** (config table `system_settings`).

**Table:**
```sql
CREATE TABLE system_settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Insert: ('platform_markup_percent', '0.05')
```

**How Mitra sees price:**
- Admin dashboard: Mitra can view but NOT change platform markup
- Mitra's selling price calculation uses current `platform_markup_percent` at time of transaction (snapshot stored in transaction record)

**Change Management:**
- Changing platform_markup affects **future** transactions only
- Existing product `platform_price` cached; updated on next sync (prepaid hourly, postpaid 5min)
- Consider archiving old prices? Not needed — transactions store `platform_price` as snapshot

---

## 11. Product Category & SKU Management

### 11.1 Category Hierarchy

Flat list (no subcategories in scope). Categories:
- Pulsa
- Paket Data
- PLN
- PDAM
- E-Wallet
- etc.

**Added via seed data; admin can add via admin panel.**

### 11.2 Product Active/Inactive

Digiflazz may disable products (`buyer_product_status = false` or `seller_product_status = false`). During sync, mark local product `is_active = false` if Digiflazz indicates inactive.

**Soft delete:** Do not hard-delete; keep for historical transaction integrity.

### 11.3 Price Update Strategy

On sync:
- If `price` changed by >1% → log warning (possible market shift)
- If `price` changed by >10% → alert product manager (verify with Digiflazz)
- If `unlimited_stock = false` and `stock = 0` → mark `is_active = false`

---

## 12. Validation Rules Cheat Sheet

| Field | Rule | Error Code | Message (ID) |
|---|---|---|---|
| `phone_number` | Regex `^\+62[0-9]{8,12}$` | `INVALID_PHONE` | "Nomor HP tidak valid" |
| `password` | Min 8 chars, 1 upper, 1 lower, 1 digit | `WEAK_PASSWORD` | "Password minimal 8 karakter, huruf besar & angka" |
| `pin` | 6 digits, not sequential, not all-same | `INVALID_PIN` | "PIN tidak valid" |
| `customer_no` (prepaid) | Per-operator regex (stored in validation_rules) | `INVALID_CUSTOMER_NO` | "Nomor tujuan salah" |
| `selling_price` | >= `platform_price` | `PRICE_BELOW_COST` | "Harga jual minimal Rp X" |
| `otp_code` | 6 digits numeric | `INVALID_OTP` | "OTP tidak valid" |
| `staff_daily_count` | < 50 | `DAILY_TXN_LIMIT` | "Limit harian transaksi tercapai" |
| `staff_daily_amount` | < 5,000,000 | `DAILY_AMOUNT_LIMIT` | "Limit harian nilai transaksi tercapai" |
| `wallet_balance` | available >= amount | `INSUFFICIENT_BALANCE` | "Saldo tidak mencukupi" |
| `product_active` | is_active = true | `PRODUCT_INACTIVE` | "Produk tidak aktif" |
| `postpaid_inquiry` | not expired | `INQUIRY_EXPIRED` | "Inquiry sudah expired, silakan cek ulang" |
| `ref_id_uniqueness` | unique per transaction | `DUPLICATE_REF_ID` | "ID transaksi sudah digunakan" |

---

## 13. Cron Jobs & Scheduled Tasks

| Job Name | Schedule | Description |
|---|---|---|
| `product-sync-prepaid` | `5 * * * *` (hourly) | Sync prepaid products from Digiflazz |
| `product-sync-postpaid` | `*/5 * * * *` (every 5min) | Sync postpaid products |
| `reconcile-pending-transactions` | `* * * * *` (every minute) | Find Pending >15min, mark Expired, release holds |
| `daily-limit-reset` | `0 0 * * *` (midnight) | Insert fresh daily_limit rows for all active staff? Not needed — daily_limits uses DATE column auto-group |
| `balance-reconciliation` | `0 2 * * *` (2 AM daily) | Verify sum(wallet_events) = wallet.balance_cached for all wallets |
| `digiflazz-deposit-check` | `0 * * * *` (hourly) | Call Cek Saldo, compare to sum(Mitra main balances), alert drift |
| `postpaid-inquiry-cleanup` | `0 * * * *` (hourly) | Delete `postpaid_inquiries` where `expires_at < NOW()` |
| `compensation-retry` | `*/2 * * * *` (every 2min) | Retry failed compensation jobs with exponential backoff |
| `idempotency-cleanup` | `0 1 * * *` (1 AM daily) | Delete expired `idempotency_keys` (older than 24h) |

**Kubernetes CronJob Example:**
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: reconcile-pending-txns
spec:
  schedule: "* * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: worker
            image: ppob/transaction-worker:latest
            args: ["--reconcile-pending"]
          restartPolicy: OnFailure
```

---

## 14. Audit Trail Requirements

**Every state-changing operation must record:**
- Actor (user_id or 'system')
- Action (verb: create, update, delete, credit, debit)
- Resource type and ID
- Old value (if update)
- New value
- Timestamp
- IP address, User-Agent
- Trace ID (for cross-service correlation)

**Example:** Wallet debit:
```json
{
  "user_id": "staff-123",
  "action": "wallet_debit",
  "resource": {"type": "wallet_event", "id": "evt-456"},
  "details": {
    "amount": 25000,
    "balance_before": 150000,
    "balance_after": 125000,
    "transaction_id": "txn-789",
    "wallet_id": "wal-001"
  }
}
```

**Immutable:** Use database triggers or application guarantees (append-only table).

---

## 15. Financial Integrity Invariants

These **must always hold true** (checked by daily reconciliation jobs):

1. **Σ(wallet_events per wallet) = wallets.balance_cached**
   - For each wallet, sum of credited amounts minus debited amounts equals cached balance.

2. **Σ(platform_price across success transactions) ≤ Total Digiflazz deposit**
   - Platform cannot pay Digiflazz more than available deposit.

3. **Mitra Profit + Staff Commission = Total Margin**
   - Margin (SellingPrice − PlatformPrice) fully distributed (no leakage).

4. **Daily Limits Enforced:**
   - `SELECT COUNT(*) FROM daily_limits WHERE user_id=? AND date=?` should equal count of Success transactions that day.

5. **No Negative Balances:**
   - `wallets.balance_available >= 0` for all wallets.

**Reconciliation Job:** 
- Runs daily at 2 AM
- If any invariant violated, create `audit_logs` entry with severity `CRITICAL` and send alert to finance@company.com

---

## 16. Edge Cases & Decision Matrix

| Scenario | Decision | Rationale |
|---|---|---|
| **Mitra sets selling_price < platform_price** | Reject transaction with `PRICE_BELOW_COST` | Prevents financial loss |
| **Staff leaves company** | Deactivate user (soft delete), keep transactions | Historical integrity; cannot reuse phone number |
| **Mitra goes bankrupt** | Freeze all wallets, admin review | Legal requirement to hold funds |
| **Digiflazz partial bill payment (PLN multiple installments)** | Treat as single transaction; if some installments fail, overall status = PartialSuccess? Currently: all-or-nothing per Digiflazz response | Digiflazz handles partiality; we get single status |
| **Two staff share same phone number** | Disallowed — phone_number UNIQUE constraint | Account uniqueness |
| **User has roles from two Mitras** | Allowed — user can be Mitra of A and Staff of B; wallets separate; active role determines which wallet is used | Multi-tenant design |
| **Mitra wants to top-up staff wallet with amount > own balance** | Reject with `INSUFFICIENT_BALANCE` | Cannot overdraw Mitra |
| **Refund larger than original transaction** | Reject — refunds only up to original amount | Financial control |
| **Product removed from Digiflazz but already sold** | Keep product record with `is_active=false` for history; new transactions cannot use it | Historical integrity |
| **User requests data deletion (PDP law)** | Anonymize PII (phone, name) but keep transactions (legal requirement to retain 5 years) | Compliance |

---

## 17. Business Metrics Definitions

| Metric | Formula | Dashboard | Owner |
|---|---|---|---|
| Daily Transaction Volume | COUNT(transactions WHERE status='Success' AND DATE(created_at)=today) | Transaction Dashboard | Ops |
| Gross Revenue | SUM(transactions.selling_price WHERE status='Success' AND DATE=...) | Finance | CFO |
| Platform Profit | SUM(transactions.platform_price − base_cost) — sum of staff commissions | Finance | CFO |
| Staff Commission Payout | SUM(commissions.amount WHERE type='credit') | Finance | Payroll |
| Active Mitra Count | COUNT(DISTINCT users WHERE has_role('Mitra') AND last_login > 30d) | Growth | PM |
| Average Transaction Value | AVG(selling_price) where success | Ops | Analyst |
| Top Selling Products | TOP 10 products by transaction count | PM | — |
| Failed Transaction Rate | COUNT(failed)/COUNT(total) | SRE | — |
| Pending-to-Success Latency | AVG(webhook_received_at − created_at) | SRE | — |

---

## 18. Configurable Parameters

**System-wide (in `system_settings` table):**

| Key | Type | Default | Description |
|---|---|---|---|
| `platform_markup_percent` | DECIMAL(5,4) | `0.05` | 5% platform markup on base price |
| `staff_default_daily_txn_limit` | INT | 50 | Default daily transaction limit per staff |
| `staff_default_daily_amount_limit` | DECIMAL(18,2) | 5000000 | Default daily turnover limit (IDR) |
| `pending_timeout_minutes` | INT | 15 | Auto-expire pending after N minutes |
| `max_idempotency_ttl_hours` | INT | 24 | Idempotency key validity window |
| `inquiry_expiry_hours` | INT | 24 | Postpaid inquiry valid for 24h |
| `otp_expiry_minutes` | INT | 5 | OTP valid for 5 minutes |
| `password_min_length` | INT | 8 | Min password length |
| `pin_max_attempts` | INT | 5 | Max wrong PIN tries before lock |
| `pin_lock_duration_hours` | INT | 1 | Account lock duration after PIN failures |

**Change Management:** These settings editable only via admin panel (super-admin role). Changes take effect immediately (except markup — affects only new transactions; existing product prices updated on next sync).

---

## 19. Compliance Rules

### 19.1 Data Retention

- **Transactions:** Keep for 5 years (tax law)
- **Audit Logs:** Keep for 5 years
- **Wallet Events:** Keep for 5 years (source of truth for financial reconciliation)
- **Login History (device_fingerprints):** Keep for 90 days, then archive to cold storage
- **OTP Codes:** Auto-delete after 24h (TTL)

**Archival Strategy:**
- Partition `transactions` table by month (PostgreSQL table partitioning)
- After 12 months, move partitions to `transactions_archive_2025` table with same structure but on cheaper storage
- After 5 years, export to S3 Glacier and delete from DB (retain only summary for reporting)

### 19.2 Right to be Forgotten (PDP)

If user requests data deletion:
1. Anonymize PII: set `phone_number` = `' anonymized-' || user_id`, `name` = `'Anonymized'`
2. Keep all financial records (transactions, wallet_events) intact (legal requirement)
3. Mark `users.deleted_at = NOW()` and `is_active = false`
4. Cannot reuse phone number (unique constraint still on masked phone)

---

## 20. Testing Matrix for Business Logic

| Test Case | Preconditions | Steps | Expected Result |
|---|---|---|---|
| **Margin non-negative** | Mitra sets selling_price < platform_price | Attempt transaction | Rejected with PRICE_BELOW_COST |
| **Daily limit count-based** | Staff has 49 success today | Initiate 2nd txn same day | 1st succeeds, 2nd rejected DAILY_LIMIT |
| **Daily limit amount-based** | Staff turnover Rp 4,999,999 today | Initiate txn of Rp 1,001 | Rejected |
| **Hold release on failure** | Pending txn with held amount | Webhook RC 68 | Held amount released to available |
| **Hold release on expiry** | Pending txn >15min old | Reconciliation job runs | Hold released, status Expired |
| **Commission rate** | Global MarginShare 60% | Success txn margin Rp 10k | Staff gets Rp 6,000; Mitra Rp 4,000 |
| **Product override** | Staff product override 70% | Success txn | Staff gets 70% of margin |
| **Postpaid inquiry expiry** | Inquiry created 25h ago | Attempt pay with same ref_id | Rejected INQUIRY_EXPIRED |
| **Refund reversal** | Success txn with commission paid | Refund webhook received | Staff wallet debited commission, wallet credited full amount |
| **Multi-role wallet selection** | User is Mitra A and Staff B | Switch role → initiate txn | Txn uses wallet corresponding to active role |

---

## Appendix A — SQL: Compute Real-Time Staff Commission

```sql
SELECT 
    s.user_id AS staff_id,
    COALESCE(SUM(
        CASE 
            WHEN c.amount > 0 THEN c.amount 
            ELSE 0 
        END
    ), 0) AS total_commission
FROM users u
JOIN user_roles ur ON u.user_id = ur.user_id AND ur.role_name = 'Staff'
LEFT JOIN commissions c ON u.user_id = c.staff_id AND c.earned_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY s.user_id;
```

---

## Appendix B — Financial Closing Checklist (Month-End)

- [ ] Reconcile all wallet balances to sum of wallet_events (audit verified)
- [ ] Verify total platform profit = Σ(margin) across all success transactions
- [ ] Generate staff commission statements for payroll
- [ ] Backup all tables with `pg_dump --format=custom`
- [ ] Archive transactions older than 2 years to data warehouse (Redshift/BigQuery) for reporting
- [ ] Run anomaly detection: look for negative margin transactions
- [ ] Tax report generation (PPh 23/26 as applicable)

---

## Open Questions

1. **Should platform_markup_percent be per-category?** Initially global; later could vary by category (PLN lower margin, E-wallet higher). Design for future extension: `category_markups` table.

2. **Should commission be paid immediately or nightly batch?** Immediate provides transparency. Nightly batch allows reversal of failed transactions before payout. **Decision:** Immediate credit to staff wallet visible in app; payroll runs on upcoming Tuesday.

3. **What about failed transaction cost?** Digiflazz usually does not charge for failed RC (except RC 44 insufficient deposit). We don't pass cost to user; no deduction from wallet on failure.

4. **Can Mitra set selling price per-customer (negotiated)?** Not in initial scope. Only per-product global Mitra price OR per-staff per-product override.

---

**Implementation Priority:**
1. Implement pricing formula in Product Service (core)
2. Add margin setting CRUD in User Service
3. Update Transaction Service to calculate commission at success
4. Add daily limits atomic check
5. Build admin UI for Mitra to set staff margins & limits
6. Financial reports (export CSV)

---

**Related Documents:**
- `ERROR_HANDLING.md` (error codes for business rule violations)
- `TRANSACTION_STATE_MACHINE.md` (when to apply charges)
- `API_CONTRACTS.md` (endpoint request/response)

---

**Owner:** Product Team + Backend Tech Lead  
**Next Review:** After 1 month of pilot with 10 Mitra testers
