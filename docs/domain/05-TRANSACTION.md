# 🛒 Domain: Transaction Operations

## 1. Core Mission
The Transaction domain manages the full lifecycle of a purchase, ensuring that money is exchanged for digital products accurately and reliably.

## 2. Key Concepts

### 2.1 Transaction (Order)
The core business record capturing who bought what, for how much, and when.
- **Attributes:** Transaction ID, Ref ID, Product, Customer No, Amount, Selling Price, Status.

### 2.2 Transaction States
The status of a purchase at any given time.
- **Terminal States:** `Success`, `Failed`, `Expired`, `Cancelled`, `Refunded`.
- **Active States:** `Initiated`, `Pending`.

### 2.3 Idempotency
Ensuring that a single business intent (e.g., buying 10k credit) is only executed once, even if the network fails and the user retries.

## 3. Business Rules

### 3.1 Order Validation
Before execution, the domain validates:
- Product is active and in stock.
- Customer Number matches the category regex.
- User has enough `Available Balance`.
- Staff has not exceeded daily transaction or amount limits.

### 3.2 Timeouts
- **Pending Expiry:** If an order remains in `Pending` for > 15 minutes without a provider update, it is automatically marked as `Expired`.

## 4. Domain Logic

### 4.1 State Machine
The domain enforces a rigid state machine. For example, a `Success` transaction cannot be transitioned to `Failed`. A `Refunded` status is only reachable from `Success`.

### 4.2 Financial Coupling
Transactions are tightly coupled with the Wallet domain.
- `Initiated` → No wallet impact.
- `Pending` → Move from `Available` to `Held`.
- `Success` → Deduct from `Held` permanently.
- `Failed/Expired` → Return from `Held` to `Available`.
