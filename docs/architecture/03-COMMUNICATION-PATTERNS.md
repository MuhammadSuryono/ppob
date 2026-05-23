# 📡 Communication Patterns

## 1. Introduction
The PPOB microservices communicate through three primary channels: Synchronous APIs (gRPC/REST), Asynchronous Events (Redis Streams), and Orchestrated Sagas for cross-service transactions.

## 2. Synchronous Communication
Used for request-response patterns where immediate consistency or feedback is required.

### 2.1 External REST APIs
- **Gateway:** Kong API Gateway handles routing, SSL termination, and global rate limiting.
- **Client:** Mobile app communicates via JSON/REST over HTTPS.
- **Inter-service:** Some internal calls use REST with shared secret authentication.

### 2.2 Internal gRPC
Used for high-performance service-to-service calls:
- **Auth → User:** Fetching profile data during login.
- **Transaction → Wallet:** Placing/releasing holds and debits.
- **Transaction → Product:** Validating prices and SKU availability.

## 3. Asynchronous Events
Used for decoupling services and ensuring eventual consistency.

### 3.1 Event Bus: Redis Streams
Services publish events to shared streams (planned migration):
- `user.registered`: Triggers wallet creation and welcome notifications.
- `transaction.success`: Triggers commission calculation and push notifications.
- `product.updated`: Invalidates pricing caches.

## 4. Transaction Coordination: Saga Pattern
Since PPOB involves multi-step financial operations across services, we use an **Orchestration-based Saga** managed by the Transaction Service.

### 4.1 Happy Path Flow
1.  **Transaction Service:** Receives request, validates idempotency.
2.  **Product Service:** (gRPC) Validates product status and price.
3.  **Wallet Service:** (gRPC) Places a **Hold** on the user's available balance.
4.  **Integration Service:** (REST) Initiates transaction with Digiflazz.
5.  **Integration Service:** Receives Webhook/Poll result from provider.
6.  **Transaction Service:** Updates status to `Success`.
7.  **Wallet Service:** (Event/gRPC) Converts **Hold to Debit**.
8.  **Wallet Service:** (Event) Credits staff commission.

### 4.2 Failure & Compensation
If a step fails, the orchestrator triggers compensation actions:
- If provider rejects: Release the wallet hold.
- If debit fails after success response: Log to DLQ (Dead Letter Queue) and trigger a compensation job in the Integration Service.

## 5. Idempotency Mechanisms
To handle network retries safely:
- **Idempotency-Key:** Required header for all write operations.
- **Ref ID Uniqueness:** Digiflazz `ref_id` is used as a unique constraint in the `transactions` table.
- **State Machine Guards:** Status transitions are one-way (e.g., `Pending` → `Success` is allowed, but `Success` → `Pending` is rejected).
