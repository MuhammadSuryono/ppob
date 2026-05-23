# 💰 Wallet & Financial Operations Flow

## 1. Overview
The Wallet Service is the source of truth for all financial balances. It uses a combination of an event-sourced ledger and cached balances for performance.

## 2. Event Sourcing Pattern
Every balance change is recorded in the `wallet_events` table before updating the cached balance in the `wallets` table.

- **Credit Event:** Increases `balance_available`.
- **Debit Event:** Decreases `balance_available`.
- **Hold Event:** Decreases `balance_available` and increases `balance_held`.
- **Release Event:** Decreases `balance_held` and increases `balance_available`.

## 3. Atomic Debit Flow
To prevent race conditions, the Wallet Service uses pessimistic locking during updates.

1.  **Begin Transaction.**
2.  **Lock Row:** `SELECT * FROM wallets WHERE wallet_id = ? FOR UPDATE`.
3.  **Check Balance:** Verify `balance_available >= amount`.
4.  **Insert Event:** Append record to `wallet_events`.
5.  **Update Cache:** Update `balance_available` in the `wallets` table.
6.  **Commit Transaction.**

## 4. Commission Distribution Flow
Commissions are distributed asynchronously after a successful transaction.

1.  **Event:** Transaction Service publishes `transaction.success`.
2.  **Worker:** A background worker picks up the event.
3.  **Calculation:** Worker calculates Staff Commission and Mitra Profit based on the configured scheme.
4.  **Execution:**
    - Worker calls Wallet Service to **Credit Staff** wallet.
    - Worker calls Wallet Service to **Credit Mitra** wallet (if profit > 0).
5.  **History:** Commissions are logged in the `commissions` table for reporting.

## 5. Reconciliation Flow
A background job runs daily to ensure the event-sourced ledger matches the cached balance.

1.  **Query:** `SELECT SUM(amount) FROM wallet_events WHERE wallet_id = ?`.
2.  **Compare:** Compare sum to `wallets.balance_available + wallets.balance_held`.
3.  **Action:** If a drift > Rp 1,000 is detected, the system creates a CRITICAL audit alert and notifies the finance team.
