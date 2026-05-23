# 🔌 Integration Architecture

## 1. Introduction
The Integration Service acts as the gateway between internal microservices and external PPOB providers, primarily **Digiflazz**. It encapsulates all logic related to provider-specific protocols, signatures, and error mapping.

## 2. Digiflazz Gateway
### 2.1 Authentication
Requests are signed using an MD5-based HMAC strategy:
- `sign = md5(username + apiKey + cmd)`
- `cmd` varies per endpoint (e.g., `deposit`, `pricelist`, or the transaction `ref_id`).

### 2.2 Synchronization Logic
- **Prepaid Catalog:** Hourly sync (at minute :05) using a distributed Redlock to prevent parallel jobs.
- **Postpaid Catalog:** Sync every 5 minutes due to higher price/status volatility.
- **Price Calculation:** `Platform Price = Digiflazz Price + Platform Markup (5%)`.

### 2.3 Transaction Handling
- **Prepaid:** Synchronous request-response.
- **Postpaid:** Two-step flow:
    1.  **Inquiry:** Fetch bill details; stored in `postpaid_inquiries` for 24h.
    2.  **Payment:** Same `ref_id` used for final checkout.

## 3. Webhook Architecture
External status updates are received via an idempotent webhook endpoint.
- **Verification:** HMAC-SHA1 signature verification in the `X-Hub-Signature` header.
- **Processing:** Asynchronous hand-off to a task queue to ensure a <5s response time back to the provider.
- **Deduplication:** Redis-based locking ensures a webhook is never processed more than once for a single `ref_id`.

## 4. Resilience & Error Mapping
- **Circuit Breaker:** Trips after a 50% failure rate to Digiflazz, failing fast to internal services.
- **Error Catalog:** Standardizes provider Response Codes (RC) into internal error types (e.g., RC 68 → `DIGIFLAZZ_OUT_OF_STOCK`).
- **Retry Policy:** Automatic 3-retry attempt with exponential backoff for transient provider timeouts (RC 01, 70, 99).
- **Dead Letter Queue:** Webhooks that fail processing after verification are logged to a DLQ for manual admin intervention.

## 5. Security Protocols
- **IP Whitelisting:** Digiflazz webhooks are only accepted from known provider IP ranges.
- **Secret Rotation:** API keys and Webhook secrets are rotated quarterly via AWS Secrets Manager.
- **Redaction:** Logs redact all `customer_no` and `sn` (serial numbers) to protect user privacy.
