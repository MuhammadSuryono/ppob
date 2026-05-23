# 🗄️ Data Architecture

## 1. Overview
The PPOB system utilizes a primary PostgreSQL database for persistent storage and a Redis instance for caching, rate limiting, and event streaming. The data architecture focuses on consistency, auditability, and performance.

## 2. Persistence Layer: PostgreSQL 15
While services are logically decoupled, they currently share a PostgreSQL instance with independent schemas.

### 2.1 Key Design Patterns
- **Soft Deletes:** Critical tables (`users`, `products`) use `deleted_at` for historical integrity.
- **Strict Constraints:** Heavy use of `CHECK` constraints, `UNIQUE` indexes, and `FOREIGN KEY`s to enforce business rules at the database level.
- **Triggers:** Automatic wallet creation on role assignment and `selling_price` validation via PL/pgSQL triggers.

### 2.2 Event Sourcing (Wallet Ledger)
Wallet balances are managed using an append-only event sourcing pattern in the `wallet_events` table.
- **Source of Truth:** The sum of all events for a `wallet_id`.
- **Cached Balance:** Denormalized `balance_available` and `balance_held` in the `wallets` table for fast reads, reconciled hourly.

## 3. Caching & Messaging: Redis 7.2
Redis serves multiple roles in the architecture:
- **Distributed Locking:** Uses Redlock to coordinate product synchronization and prevent concurrent webhook processing.
- **Rate Limiting:** Sliding window counters for auth attempts and transaction initiation.
- **Session Store:** JWT refresh token allowlist and blocklist.
- **Event Bus:** Future migration to Redis Streams for inter-service communication.

## 4. Concurrency Control
- **Pessimistic Locking:** `SELECT FOR UPDATE` is used during wallet debit/hold operations to prevent double-spending.
- **Optimistic Locking:** `version` column in `products` table for catalog updates.
- **Canonical Ordering:** Resources are always locked in a global sorted order (e.g., by `wallet_id`) to prevent deadlocks.

## 5. Data Integrity Invariants
The system enforces several critical invariants checked daily by reconciliation jobs:
1.  **Balance Integrity:** `Σ(wallet_events) == wallets.balance_cached`.
2.  **Margin Integrity:** `Mitra Profit + Staff Commission == Total Margin`.
3.  **Deposit Integrity:** `Σ(Platform Prices) ≤ Total Digiflazz Deposit`.

## 6. Retention & Archival
- **Hot Storage:** 2 years of transactions and audit logs kept in the primary DB.
- **Cold Storage:** Data older than 2 years moved to partitioned archive tables or S3 Glacier for compliance (5-year tax law requirement).
