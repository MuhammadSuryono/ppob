# ⚠️ Concurrency Control for PPOB Application

**Audience:** Backend developers, database administrators, system architects  
**Last Updated:** 2026-05-07  
**Status:** Draft — implementation requires DBA review

---

## 1. Overview

This document defines concurrency patterns, locking strategies, isolation levels, and compensation logic across microservices. The goal is to prevent race conditions, ensure data consistency, and maintain system integrity under concurrent load.

**Critical Focus Areas:**
- Wallet operations (debit, credit, hold, top-up)
- Product price synchronization
- Staff assignment and limit updates
- Transaction deduplication
- Distributed job coordination

---

## 2. Concurrency Challenges in PPOB

| Challenge | Risk if Uncontrolled | Service Affected |
|---|---|---|
| **Double-spend** — two simultaneous transactions from same wallet | Wallet overdrawn, financial loss | Wallet Service |
| **Lost update** — two staff top-ups updating same Mitra wallet balance | Balance calculated incorrectly | Wallet Service |
| **Race on product sync** — two sync jobs overwrite each other | Price inconsistencies, lost updates | Product Service |
| **Staff limit exceed** — concurrent transactions bypass daily limit check | Fraud, unmitigated loss | Transaction Service |
| **Duplicate webhook** — same webhook processed twice | Double debit/credit, inconsistent state | Integration Service |
| **Compensating transaction failure** — top-up: debit succeeds, credit fails | Money disappears (neither wallet has correct balance) | Wallet Service |

---

## 3. Event-Sourced Wallet Architecture

**Design Choice:** Event sourcing provides immutable audit trail, natural concurrency control, and perfect reconstruction of balance at any point in time.

### 3.1 Wallet Events Table

**Schema:**
```sql
CREATE TABLE wallet_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- 'WalletCreated', 'Credited', 'Debited', 'Held', 'HoldReleased', 'TopupAdded', 'Refunded', 'Compensated'
    amount DECIMAL(18,2) NOT NULL CHECK (amount > 0),
    balance_before DECIMAL(18,2) NOT NULL,
    balance_after DECIMAL(18,2) NOT NULL,
    reference_id VARCHAR(255), -- transaction_id, topup_id, hold_id
    reference_type VARCHAR(50), -- 'transaction', 'topup', 'hold'
    metadata JSONB, -- additional context (product_id, staff_id, rc_code)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID, -- user_id who triggered, or system
    FOREIGN KEY (wallet_id) REFERENCES wallets(wallet_id)
);

CREATE INDEX idx_wallet_events_wallet_created ON wallet_events(wallet_id, created_at DESC);
CREATE INDEX idx_wallet_events_reference ON wallet_events(reference_id) WHERE reference_id IS NOT NULL;
```

### 3.2 Balance Computation

Balance is **derived** from wallet_events — never stored directly in `wallets.balance` as source of truth. However, we maintain `wallets.balance_cached` for fast reads with periodic reconciliation.

```sql
-- Get current balance (authoritative)
SELECT COALESCE(SUM(CASE 
    WHEN event_type IN ('Credited', 'TopupAdded', 'Refunded', 'HoldReleased') THEN amount
    ELSE -amount 
END), 0) AS balance
FROM wallet_events
WHERE wallet_id = $1;
```

### 3.3 Materialized View / Cached Balance

**Table:** `wallets` (existing)

```sql
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS balance_cached DECIMAL(18,2) DEFAULT 0;
-- balance_cached is denormalized for fast reads; reconciled hourly
```

**Reconciliation Job (hourly):**
```sql
-- Rebuild cache from events
WITH computed AS (
    SELECT wallet_id,
           COALESCE(SUM(CASE WHEN event_type IN ('Credited','TopupAdded','Refunded','HoldReleased') THEN amount ELSE -amount END), 0) AS computed_balance
    FROM wallet_events
    GROUP BY wallet_id
)
UPDATE wallets w
SET balance_cached = c.computed_balance
FROM computed c
WHERE w.wallet_id = c.wallet_id
  AND w.balance_cached != c.computed_balance;
```

If drift > threshold (e.g., Rp 100), create audit alert.

---

## 4. Pessimistic Locking Pattern

**When to use:** When updating wallet balance during transaction success or top-up (high contention).

**Implementation (PostgreSQL + Go GORM):**

```go
func (s *WalletService) debitForTransaction(txID uuid.UUID, amount float64) error {
    // Get transaction to find wallet_id
    transaction := &models.Transaction{}
    if err := s.db.First(transaction, "id = ?", txID).Error; err != nil {
        return err
    }

    // Begin transaction with pessimistic lock on wallet row
    return s.db.Transaction(func(tx *gorm.DB) error {
        wallet := &models.Wallet{}
        // SELECT ... FOR UPDATE — locks row until transaction commits
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(wallet, "wallet_id = ?", transaction.WalletID).Error; err != nil {
            return err
        }

        // Compute new balance
        newBalance := wallet.BalanceCached - amount
        if newBalance < 0 {
            return errors.New("insufficient_balance")
        }

        // Update wallet
        wallet.BalanceCached = newBalance
        if err := tx.Save(wallet).Error; err != nil {
            return err
        }

        // Insert wallet_event for audit
        event := models.WalletEvent{
            WalletID:       wallet.WalletID,
            EventType:      "Debited",
            Amount:         amount,
            BalanceBefore:  wallet.BalanceCached + amount,
            BalanceAfter:   newBalance,
            ReferenceID:   txID.String(),
            ReferenceType: "transaction",
            CreatedBy:     transaction.UserID,
        }
        if err := tx.Create(&event).Error; err != nil {
            return err
        }

        return nil
    })
}
```

**Lock Timeout:** Set `statement_timeout = 5000` (5s) to avoid deadlock wait forever. On timeout, retry up to 3 times with 100ms backoff.

---

## 5. Optimistic Locking Pattern

**When to use:** Product sync updates where conflicts are rare but possible (e.g., two sync jobs overlapping).

**Implementation:**

Add `version` column to `products`:
```sql
ALTER TABLE products ADD COLUMN version INT NOT NULL DEFAULT 0;
```

Update with version check:
```go
func (s *ProductService) UpsertProduct(product *Product) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Try update with version check
        result := tx.Model(&Product{}).
            Where("product_id = ? AND version = ?", product.ID, product.Version).
            Updates(map[string]interface{}{
                "buyer_sku_code": product.SKU,
                "platform_price": product.Price,
                "version":        gorm.Expr("version + 1"),
            })

        if result.RowsAffected == 0 {
            // Conflict: another update changed row
            return ErrConcurrentUpdate
        }

        return nil
    })
}
```

On `ErrConcurrentUpdate`, refetch row, merge changes, retry (max 3 attempts).

---

## 6. Isolation Levels

**Default:** `READ COMMITTED` (PostgreSQL default) — prevents dirty reads, allows non-repeatable reads and phantom reads.

**Special Cases:**

### 6.1. Transaction Deduplication (Serializable needed)
When checking `idempotency_keys` within same transaction that creates transaction, use `SERIALIZABLE` isolation to prevent double-spend under race.

```sql
BEGIN ISOLATION LEVEL SERIALIZABLE;
-- Check idempotency key not exists
INSERT INTO idempotency_keys ...;
-- Create transaction
COMMIT;
-- If serialization_failure (SQLSTATE 40001), retry whole transaction
```

### 6.2. Daily Limit Check
```sql
-- Count staff's transactions today with lock to prevent race
SELECT COUNT(*) FROM transactions
WHERE user_id = $1
  AND DATE(created_at) = CURRENT_DATE
  AND status IN ('Success', 'Pending')
FOR SHARE;  -- lightweight lock, allows concurrent reads
```

Then check against limit. Still vulnerable if two transactions pass check simultaneously. **Fix:** Compute + insert in single atomic statement:

```sql
INSERT INTO transactions (transaction_id, user_id, amount, status, created_at)
SELECT gen_random_uuid(), $1, $2, 'Initiated', NOW()
WHERE (
    SELECT COALESCE(SUM(amount), 0)
    FROM transactions
    WHERE user_id = $1
      AND DATE(created_at) = CURRENT_DATE
      AND status IN ('Success', 'Pending')
) + $2 <= $3  -- $3 is daily limit
RETURNING transaction_id;
```

If `RETURNING` empty → limit exceeded, reject.

---

## 7. Distributed Locking (Redis Redlock)

**Use Case:** Product sync scheduler must not run concurrently across multiple pods.

**Implementation with Redlock:**

```go
import "github.com/red-go/redlock"

var redlock = redlock.NewRedlock([]redlock.RedisNode{{
    Address: "redis:6379",
    Pool:    redisPool,
}})

func syncProducts() error {
    ctx := context.Background()
    // Try to acquire lock (TTL 5 min — should finish within this)
    lock, err := redlock.Lock(ctx, "product-sync-lock", 5*time.Minute)
    if err != nil {
        return fmt.Errorf("could not acquire lock: %w", err)
    }
    defer redlock.Unlock(ctx, lock)

    // Critical section: do sync
    return performProductSync()
}
```

**Key:** `product-sync-lock` (for prepaid) and `product-sync-lock-postpaid` (for postpaid) — separate locks.

---

## 8. Deadlock Prevention & Detection

**Deadlock Scenarios:**
1. Transaction A: locks wallet X, then tries to lock wallet Y
2. Transaction B: locks wallet Y, then tries to lock wallet X → deadlock

**Prevention:** Always lock resources in **global order** (wallet_id sorted ascending).

```go
func transferBetweenWallets(fromID, toID uuid.UUID, amount float64) error {
    // Determine order
    first, second := fromID, toID
    if fromID.String() > toID.String() {
        first, second = toID, fromID
    }

    return db.Transaction(func(tx *gorm.DB) error {
        // Lock in canonical order
        lockWallet(tx, first)
        lockWallet(tx, second)

        // Perform debit/credit
        // ...
        return nil
    })
}

func lockWallet(tx *gorm.DB, walletID uuid.UUID) {
    tx.Clauses(clause.Locking{Strength: "UPDATE"}).
       First(&Wallet{}, "wallet_id = ?", walletID)
}
```

**Detection:** PostgreSQL automatically detects deadlocks and rolls back one transaction. Application should catch `err = sql.Error{Code: "40P01"}` (deadlock_detected) and **retry with backoff** (max 3 retries).

---

## 9. Compensation & Saga Pattern

**Scenario:** Mitra top-up staff wallet:
1. Debit Mitra main wallet (success)
2. Credit staff wallet (fails — DB error)

**Without compensation:** Money lost (debited but not credited).

**Solution:** Compensation transaction (reverse step 1).

### 9.1. Compensation Jobs Table

```sql
CREATE TABLE compensation_jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_type VARCHAR(50) NOT NULL, -- 'mitra_topup_staff'
    payload JSONB NOT NULL, -- { "mitra_wallet_id": "...", "staff_wallet_id": "...", "amount": 50000, "failed_step": "credit_staff" }
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'retrying', 'failed', 'completed'
    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 10,
    next_retry_at TIMESTAMP, -- for backoff scheduling
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (mitra_wallet_id) REFERENCES wallets(wallet_id)  -- optional, may be NULL if wallet deleted
);
CREATE INDEX idx_compensation_retry ON compensation_jobs(status, next_retry_at);
```

### 9.2. Compensation Worker

**Process:** Runs every 30s, fetches `pending` or `retrying` jobs where `next_retry_at <= NOW()` or `next_retry_at IS NULL`.

**Logic:**
- For `mitra_topup_staff` failure: attempt to **credit Mitra wallet back** (reverse of debit)
- On success: mark job `completed`, log audit
- On failure: increment `retry_count`, set `next_retry_at = NOW() + exponential_backoff`, if `retry_count >= max_retries` set `status='failed'` and alert admin

**Alert:** If any job reaches `failed` status, send PagerDuty alert with job details.

---

## 10. Concurrent Transaction Limits

**Requirement:** Staff limited to 50 transactions/day or Rp 5M daily turnover.

**Naive Approach (Race Condition):**
```go
count, _ := db.CountTodaySuccessTransactions(staffID)
if count >= 50 {
    return errors.New("daily limit reached")
}
// create transaction...
```
Two concurrent requests both see count=49, both succeed → 51 total. **Wrong.**

**Correct Approach — Atomic Check-and-Increment:**

**Method A — Advisory Lock (PostgreSQL):**
```sql
BEGIN;
SELECT pg_advisory_xact_lock(hashtext($1));  -- lock key based on staff_id
-- Now safe to count and insert
INSERT INTO transactions ...;
COMMIT;
```

**Method B — Serialized Update with Unique Partial Index:**
```sql
-- Create table tracking daily usage
CREATE TABLE daily_limits (
    user_id UUID NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    transaction_count INT NOT NULL DEFAULT 0,
    total_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, date)
);

-- Atomic increment via INSERT ... ON CONFLICT ... DO UPDATE
INSERT INTO daily_limits (user_id, date, transaction_count, total_amount)
VALUES ($1, CURRENT_DATE, 1, $2)
ON CONFLICT (user_id, date)
DO UPDATE SET 
    transaction_count = daily_limits.transaction_count + 1,
    total_amount = daily_limits.total_amount + EXCLUDED.total_amount
WHERE daily_limits.transaction_count < 50
  AND daily_limits.total_amount + EXCLUDED.total_amount <= 5000000
RETURNING transaction_count;
```
If `RETURNING` yields no rows → limit exceeded. This single statement is atomic.

**Adopt Method B** — no application-level race; DB enforces limit.

---

## 11. Isolation Levels by Use Case

| Use Case | Recommended Isolation | Reason |
|---|---|---|
| Wallet debit/credit operations | `READ COMMITTED` with `SELECT FOR UPDATE` | Prevents dirty reads, explicit row locks prevent lost updates |
| Product sync upserts | `READ COMMITTED` with optimistic lock (`version` column) | Conflicts rare; retry cheap |
| Daily limit check | `SERIALIZABLE` for the check-and-insert statement only | Guarantees atomicity; short transaction |
| Transaction history reads | `READ COMMITTED` (no lock needed) | Reads can be stale; eventual consistency acceptable |
| Report queries (analytics) | `REPEATABLE READ` or use read replica | Consistent snapshot for report duration |

---

## 12. Lock Timeouts & Retry

**Pessimistic Lock Timeout:**
```sql
SET statement_timeout = '5s';  -- fail after 5s waiting for lock
```

**Application Retry Logic:**
```go
func withRetry(operation func() error) error {
    backoff := 100 * time.Millisecond
    for i := 0; i < 3; i++ {
        err := operation()
        if err == nil {
            return nil
        }
        if isDeadlock(err) || isLockTimeout(err) {
            time.Sleep(backoff)
            backoff *= 2
            continue
        }
        return err  // non-retryable error
    }
    return errors.New("max retries exceeded due to deadlocks")
}
```

---

## 13. Observability for Concurrency

**Metrics to Monitor:**
- `db_lock_wait_seconds{lock_type="pessimistic_wallet"}` — histogram; p99 should be <100ms
- `db_deadlocks_total` — counter; should be 0 or near-zero; spike indicates design flaw
- `wallet_event_insert_rate` — rate of wallet_event writes; should match transaction rate
- `hold_release_duration_seconds` — time between hold and release; target < 2min for successful txn
- `compensation_job_retries_total` — if high, indicates systemic failure in synchronous path

**Logging:**
- Log all deadlock occurrences with full SQL stack (using `pg_stat_activity`)
- Log lock acquisition timeouts with trace ID for root cause

---

## 14. Code Review Checklist for Concurrency

Reviewers must verify:
- [ ] Wallet debit/credit uses `SELECT FOR UPDATE` on wallet row within DB transaction
- [ ] Event is inserted in same DB transaction as balance update (atomicity)
- [ ] All wallet operations go through service layer (no direct DB access from controllers)
- [ ] Idempotency key check is done before creating transaction (within same transaction)
- [ ] Daily limit check uses single atomic `INSERT ... ON CONFLICT` statement, not read-then-write
- [ ] Product sync uses distributed lock (Redlock) to prevent concurrent runs
- [ webhook processing uses advisory lock or `SELECT FOR UPDATE` on transaction row to prevent duplicate processing
- [ ] No N+1 queries inside critical section (keep transaction short)
- [ ] Compensation logic exists for all multi-step financial operations
- [ ] Retry logic has exponential backoff and max retries (no infinite loop)

---

## 15. Performance Implications

**Lock Contention:**
- Wallet locks are shortest lived (<10ms) if business logic efficient
- Expect 100 TPS → 100 wallet lock acquisitions/sec — easily handled by PostgreSQL (10k+ TPS lock capability)
- Monitor `pg_locks` view; high `wait_event_type = 'Lock'` indicates bottleneck

**Scalability:**
- If wallet contention becomes bottleneck, consider **sharding wallets** across multiple databases (shard key: wallet_id hash)
- Alternative: Event sourcing with CQRS — separate write model (event store) from read model (balance cache)

**Event Store Size:**
- Estimate: 1M transactions/month → 1M wallet_events (plus additional events for holds, refunds, etc.)
- At ~200 bytes per event → 200MB/year; manageable in primary DB; archive after 5 years to cold storage

---

## Appendix A — Migration to Event-Sourced Wallet

**Step 1:** Create `wallet_events` table.
**Step 2:** Backfill historical events from existing `wallets` and `transactions`:
```sql
INSERT INTO wallet_events (wallet_id, event_type, amount, balance_before, balance_after, reference_id, reference_type, created_at)
SELECT 
    w.wallet_id,
    'InitialBalance' as event_type,
    0 as amount,  -- no amount change
    0 as balance_before,
    w.balance as balance_after,
    NULL,
    'initial',
    NOW()
FROM wallets w;
```
**Step 3:** Change application to **append-only** to `wallet_events`. Compute balance from events for all reads.
**Step 4:** Keep `wallets.balance_cached` for fast reads; set up hourly reconciliation job.
**Step 5:** After 30 days of stable operation, drop `wallets.balance` direct column if desired.

---

## Appendix B — Advisory Lock Example (PostgreSQL)

```sql
-- Acquire lock (application-side)
SELECT pg_advisory_xact_lock(hashtext('wallet:' || $1));  -- $1 = wallet_id

-- This transaction now has exclusive lock on this wallet_id
-- Other transactions trying to lock same wallet_id will wait
-- Different wallet_ids can proceed concurrently (no global lock)
```

**Key Function:** `hashtext()` converts string to 32-bit int; collision possible but negligible for our key space (<10k wallets).

---

## Appendix C — Redlock Algorithm Parameters

**Redis Nodes:** 3 master nodes (no replicas for lock algorithm)  
**Quorum:** `N/2 + 1 = 2` (majority)  
**Lock TTL:** 5 minutes (longer than expected critical section)  
**Clock drift:** 1% of TTL (3s)  
**Retry:** Up to 3 times with random delay 50–200ms

---

## Open Questions

1. **Hold Balance Column?** Should `wallets` have `balance_held DECIMAL` column separate from `balance_available`?  
   **Yes** — simplifies UI display and queries. Update both within same `SELECT FOR UPDATE` transaction.

2. **Event Sourcing Complexity:** Is full event sourcing worth it for just 3 event types?  
   **Yes** — provides immutable audit trail for financial compliance; simplifies debugging; enables point-in-time balance reconstruction.

3. **Sharding Plan?** Not needed initially. If >10k TPS, consider sharding by `wallet_id` hash across 4 DB instances.

---

**Implementation Order:**
1. Add `wallet_events` table + backfill
2. Implement `pessimistic lock` pattern in WalletService
3. Add daily_limits table + atomic check logic
4. Add product sync Redlock
5. Set up reconciliation jobs (cron)
6. Add monitoring metrics

---

**References:**
- PostgreSQL Advisory Locks: https://www.postgresql.org/docs/current/explicit-locking.html#ADVISORY-LOCKS
- Redlock Algorithm: https://redis.io/docs/reference/patterns/distributed-locks/
- Event Sourcing Pattern: Martin Kleppmann, "Designing Data-Intensive Applications"
