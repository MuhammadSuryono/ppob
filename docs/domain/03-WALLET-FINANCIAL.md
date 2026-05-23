# 💰 Domain: Wallet & Financials

## 1. Core Mission
The Financial domain acts as the source of truth for all monetary value within the platform. It ensures every Rupiah is tracked via an immutable ledger.

## 2. Key Concepts

### 2.1 Wallet Types
- **Main Wallet:** Owned by a Mitra. Can be used for transactions and funding staff wallets.
- **Sub-Wallet:** Owned by a Staff member. Funded solely by the Mitra.

### 2.2 Balance States
- **Available Balance:** Funds ready for immediate use.
- **Held Balance:** Funds reserved for pending transactions (not spendable).

### 2.3 Financial Events
Instead of simple increments/decrements, every change is an event.
- **Event Types:** `Credited`, `Debited`, `Held`, `Released`, `TopUp`.

## 3. Business Rules

### 3.1 Balance Integrity
- **No Overdrafts:** `Available Balance` must always be ≥ 0.
- **Top-Ups:** Staff wallets can only be topped up by their respective Mitra. Mitra wallets are topped up via bank transfer or platform admin.

### 3.2 Commission & Profit
- **Revenue Distribution:** Every successful transaction must distribute profit to the Mitra and commission to the Staff immediately.
- **Schemes:** Supports `FixedAllowance` (Rp/txn) and `MarginShare` (% of margin).

## 4. Domain Logic

### 4.1 Event Sourcing
The system reconstructs the current balance by summing the history of financial events. Denormalized "cached" balances are maintained for performance but reconciled daily against the event log.

### 4.2 Atomicity
Financial operations (like a Staff top-up) must be atomic. A debit from the Mitra and a credit to the Staff must happen together or not at all.
