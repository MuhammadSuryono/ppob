# Flow Reconciliation & Auth Scoring Plan

## 1. Background & Motivation
An analysis of the current implementation against the documentation in `docs/flow/` and `docs/architecture/` revealed a few functional gaps. The goal of this plan is to reconcile these discrepancies, specifically focusing on the Device Trust Scoring mechanism in the Auth Service, and consolidating the Transaction Expiry and Commission distribution flows.

## 2. Identified Gaps

### 2.1 Device Trust Scoring (Auth Service)
**Blueprint:** `01-AUTH-REGISTRATION.md` and `04-SECURITY-ARCHITECTURE.md` require an Adaptive Authentication flow based on a Trust Score (0-100).
- Trusted (Score >= 70): PIN only.
- Semi-trusted (30-69): Password + OTP.
- Untrusted (< 30): Password + OTP + Challenge.
**Current State:** The `DeviceFingerprint` model has a `TrustScore` field, but the `AuthService` (`InitiateAuth` and `upsertDeviceTrust`) ignores it. It simply sets `IsTrusted = true` immediately upon the first successful login, bypassing the graduated scoring logic.

### 2.2 Expiry Compensation (Reconciliation Service)
**Blueprint:** `06-RECONCILIATION.md` states that when a transaction is marked as `Expired` by the background job, it must trigger a compensation action to **Release the Wallet Hold**.
**Current State:** The `expireTransaction` function in `reconciliation_service.go` only updates the database status to `Expired`. It does not call the `WalletClient` via gRPC to release the hold, leaving funds locked indefinitely.

### 2.3 Commission Distribution Event (Transaction / Wallet Services)
**Blueprint:** `04-WALLET-OPERATIONS.md` dictates that the Wallet Service should listen to the `transaction.success` event (via Redis Streams) and credit the staff commission.
**Current State:** `TransactionService` publishes the event, but also pushes a task to a legacy `commission_queue` (Redis List) which is processed by a `CommissionWorker`. The `WalletEventHandler` in `wallet-service` receives the `transaction.success` event but has a `TODO` for the commission logic.

## 3. Implementation Phases

### Phase 1: Device Trust Scoring Implementation
- **Refactor `InitiateAuth`:** 
  - Read the `TrustScore` from the `DeviceFingerprint`.
  - Apply the rules: `>= 70` is Trusted (PIN), `< 70` Requires OTP.
- **Refactor `upsertDeviceTrust`:**
  - Implement a scoring algorithm.
  - New devices start with a base score (e.g., 20).
  - Successful Password+OTP logins add points (e.g., +30 points).
  - If the score crosses the 70 threshold, set `IsTrusted = true`.

### Phase 2: Expiry Compensation
- **Update `ReconciliationService`:**
  - Inject the `WalletClient` (gRPC) into the `ReconciliationService`.
  - In `expireTransaction`, after successfully updating the DB status to `Expired`, call `walletClient.ReleaseHoldForTransaction(ctx, tx.UserID, tx.TransactionID)`.

### Phase 3: Event-Driven Commission Consolidation
- **Update `wallet-service/internal/services/event_handler.go`:**
  - Implement the `handleTransactionSuccess` logic to query the margin/commission rules and credit the user's wallet.
- **Update `transaction-service`:**
  - Remove the legacy `CommissionWorker` and the `publishCommissionTask` logic (Redis List).

## 4. Verification
- After Phase 1: Test auth flow with a new device to ensure it requires OTP until the score hits 70.
- After Phase 2: Simulate a pending transaction timeout and verify the wallet hold is released.
- After Phase 3: Perform a successful transaction and verify the commission is distributed entirely via the Redis Stream event.
- Copy this file to `docs/plan/flow-reconciliation.md` and rename it to `flow-reconciliation_COMPLETED.md` upon full implementation.