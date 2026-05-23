# 🧹 Reconciliation & Recovery Flow

## 1. Overview
Automated background processes ensure that the system remains in a consistent state, handling network timeouts, provider delays, and financial drifts.

## 2. Stale Transaction Reconciliation
**Job:** `reconcile_pending_transactions`
**Frequency:** Every 1 minute

1.  **Identify:** Find transactions in `Pending` state for more than 15 minutes.
2.  **Poll Provider:** Call Digiflazz `Cek Status` API for each identified transaction.
3.  **Update State:**
    - If status is now `Success` or `Failed` → Advance state normally.
    - If still `Pending` at provider → Keep as `Pending` but log warning.
    - If status is `Not Found` or definitive failure → Mark as `Expired`.
4.  **Recover:** For `Expired` transactions, trigger a compensation action to **Release the Wallet Hold**.

## 3. Financial Balance Reconciliation
**Job:** `daily_balance_check`
**Frequency:** Daily at 02:00 WIB

1.  **Sum Events:** Calculate the sum of all `wallet_events` for each wallet.
2.  **Verify Cache:** Compare the calculated sum to the `balance_available + balance_held` columns in the `wallets` table.
3.  **Log Drift:** Any mismatch > Rp 1 is logged as an anomaly.
4.  **Alert:** Mismatch > Rp 1,000 triggers a `CRITICAL` alert to the finance team.

## 4. Provider Deposit Reconciliation
**Job:** `digiflazz_deposit_sync`
**Frequency:** Hourly

1.  **Internal Sum:** Calculate the sum of all `Platform Price` for successful transactions in the last hour.
2.  **Provider Check:** Fetch current deposit balance from Digiflazz.
3.  **Verification:** Ensure `Internal Sum <= (Initial Provider Balance - Final Provider Balance)`.
4.  **Action:** Discrepancies indicate potential leakage or incorrect pricing snapshots.

## 5. Idempotency Cleanup
**Job:** `idempotency_ttl_cleanup`
**Frequency:** Daily at 01:00 WIB

- Deletes records from the `idempotency_keys` table older than 24 hours to keep the table size manageable.
