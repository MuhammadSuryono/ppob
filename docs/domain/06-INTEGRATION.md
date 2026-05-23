# 🔌 Domain: Integration Boundaries

## 1. Core Mission
This domain abstracts the complexity of external digital product providers, providing a unified internal interface for the system.

## 2. Key Concepts

### 2.1 Provider (Aggregator)
The external business supplying the digital products (e.g., Digiflazz).

### 2.2 Provider Response Codes (RC)
The unique set of codes used by the provider to indicate success, pending status, or failure.

### 2.3 Webhooks
An asynchronous communication channel used by the provider to push status updates.

## 3. Business Rules

### 3.1 Security
- **HMAC Signatures:** All outbound requests and inbound webhooks must be cryptographically signed.
- **Provider IP Whitelist:** Webhooks are only accepted from known provider IP ranges.

### 3.2 Resilience
- **Fail Fast:** If the provider is unreachable, the system must trigger a circuit breaker to protect the database and other microservices.
- **Deduplication:** Webhooks must be handled idempotently using the provider's `ref_id`.

## 4. Domain Logic

### 4.1 Error Normalization
The domain is responsible for mapping provider-specific RC codes to internal system statuses. For example, Digiflazz RC `68` is mapped to internal status `Failed` with reason `OUT_OF_STOCK`.

### 4.2 Compensation Jobs
If an operation succeeds at the provider but fails to update locally (or vice versa), the domain manages background compensation jobs to ensure eventual consistency.
