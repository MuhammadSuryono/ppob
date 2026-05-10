# PPOB Microservices API Documentation

This document provides a detailed reference for all API endpoints across the PPOB (Payment Point Online Bank) microservices ecosystem.

---

## Base URL

All services are exposed behind the API Gateway at:

```
https://fedora.sinauplatform.id/api/v1/{api_name}/{endpoint}
```

**Service Mapping:**

| `api_name` | Service | Port |
|---|---|---|
| `auth` | Auth Service | 8081 |
| `user` | User Service | 8082 |
| `wallet` | Wallet Service | 8083 |
| `transaction` | Transaction Service | 8084 |
| `product` | Product Service | 8085 |
| `integration` | Integration Service | 8086 |

**Example:** `GET https://fedora.sinauplatform.id/api/v1/auth/login`

---

## Global Headers

All requests must include the following headers:

| Header | Required | Description |
|---|---|---|
| `Authorization` | Yes | `Bearer <JWT_ACCESS_TOKEN>` |
| `Idempotency-Key` | For writes | `<UUID>` — Required for `/transactions/initiate` and any state-changing operation |
| `X-Trace-ID` | Optional | `<UUID>` — For distributed request tracing |
| `Content-Type` | Yes | `application/json` |
| `Accept` | Yes | `application/json` |

---

## 1. Auth Service
**Base URL:** `https://fedora.sinauplatform.id/api/v1/auth`

Handles user registration, authentication, OTP verification, and session management.

### Endpoints

| Method | Endpoint | Full URL | Auth | Description |
|:---|:---|:---|:---|:---|
| `POST` | `/register` | `.../auth/register` | None | Register a new user with email, phone, name, password, and PIN. |
| `POST` | `/login` | `.../auth/login` | None | Login using Email or Phone + Password. |
| `POST` | `/verify-otp` | `.../auth/verify-otp` | None | Verify OTP code for login or registration. |
| `POST` | `/refresh` | `.../auth/refresh` | None | Refresh access token using a valid refresh token. |
| `POST` | `/logout` | `.../auth/logout` | Bearer | Logout and invalidate session. |
| `POST` | `/change-password` | `.../auth/change-password` | Bearer | Update account password. |
| `POST` | `/change-pin` | `.../auth/change-pin` | Bearer | Update 6-digit transaction PIN. |

---

## 2. User Service
**Base URL:** `https://fedora.sinauplatform.id/api/v1/user`

Manages user profiles, roles, staff management, and administrative lists.

### Endpoints

| Method | Endpoint | Full URL | Role | Description |
|:---|:---|:---|:---|:---|
| `GET` | `/users/:id` | `.../user/users/:id` | Owner/Admin | Get detailed profile information for a user. |
| `PUT` | `/users/:id` | `.../user/users/:id` | Owner/Admin | Update profile fields (name, phone, address, DOB). |
| `GET` | `/users` | `.../user/users` | Admin/Staff | Paginated list of all users. |
| `GET` | `/users/:id/roles` | `.../user/users/:id/roles` | Owner/Admin | Get roles assigned to a user. |
| `POST` | `/users/:id/roles` | `.../user/users/:id/roles` | Admin | Assign a role (e.g., Mitra, Staff) to a user. |
| `GET` | `/roles` | `.../user/roles` | Admin | List available system roles. |
| `POST` | `/roles` | `.../user/roles` | Admin | Create a new role definition. |
| `GET` | `/staff` | `.../user/staff` | Mitra | List staff users with stats (transactions, wallet, limits). |
| `POST` | `/staff` | `.../user/staff` | Mitra | Create a new staff user. |
| `GET` | `/staff/:id` | `.../user/staff/:id` | Mitra/Admin | Get detailed staff info including margin settings. |
| `PUT` | `/staff/:id` | `.../user/staff/:id` | Mitra | Update staff details and limits. |
| `GET` | `/staff/:id/stats` | `.../user/staff/:id/stats` | Mitra/Admin | Get staff performance stats. |
| `GET` | `/staff/pending-count` | `.../user/staff/pending-count` | Mitra | Count of pending staff invitations. |
| `GET` | `/notifications` | `.../user/notifications` | Bearer | List user notifications (optional query: `unread_only=true`). |
| `GET` | `/notifications/uncount` | `.../user/notifications/uncount` | Bearer | Get unread notifications count. |
| `PATCH` | `/notifications/:id/read` | `.../user/notifications/:id/read` | Bearer | Mark a notification as read. |
| `POST` | `/notifications/mark-all-read` | `.../user/notifications/mark-all-read` | Bearer | Mark all notifications as read. |

---

## 3. Product Service
**Base URL:** `https://fedora.sinauplatform.id/api/v1/product`

Manages product catalog and synchronization with Digiflazz.

### Endpoints

| Method | Endpoint | Full URL | Role | Description |
|:---|:---|:---|:---|:---|
| `GET` | `/products` | `.../product/products` | None | List products with category/status filters. |
| `GET` | `/products/:id` | `.../product/products/:id` | None | Get product by ID. |
| `GET` | `/products/search` | `.../product/products/search` | None | Search products by name/code. |
| `GET` | `/products/by-code/:code` | `.../product/products/by-code/:code` | None | Get product by SKU code. |
| `GET` | `/products/validate-price` | `.../product/products/validate-price` | None | Check if a price meets margin requirements. |
| `POST` | `/products` | `.../product/products` | Admin | Create product manually. |
| `PUT` | `/products/:id` | `.../product/products/:id` | Admin | Update product information. |
| `DELETE` | `/products/:id` | `.../product/products/:id` | Admin | Soft-delete a product. |
| `POST` | `/sync/prepaid` | `.../product/sync/prepaid` | Bearer | Trigger prepaid product sync. |
| `POST` | `/sync/postpaid` | `.../product/sync/postpaid` | Bearer | Trigger postpaid product sync. |
| `GET` | `/sync/status` | `.../product/sync/status` | None | Get last sync timestamps. |
| `GET` | `/categories` | `.../product/categories` | None | List all categories. |
| `POST` | `/categories` | `.../product/categories` | Admin | Create a category. |

---

## 4. Wallet Service
**Base URL:** `https://fedora.sinauplatform.id/api/v1/wallet`

Financial ledger handling balance, holds, and transfers.

### Endpoints

| Method | Endpoint | Full URL | Auth | Description |
|:---|:---|:---|:---|:---|
| `GET` | `/:id/balance` | `.../wallet/:id/balance` | Bearer | Get current available and held balance. |
| `GET` | `/:id/balance-events` | `.../wallet/:id/balance-events` | Bearer | Get balance reconstructed from events. |
| `POST` | `/:id/hold` | `.../wallet/:id/hold` | Bearer | Place a hold on funds for a transaction. |
| `POST` | `/:id/release-hold` | `.../wallet/:id/release-hold` | Bearer | Release a previously held amount. |
| `POST` | `/:id/debit` | `.../wallet/:id/debit` | Bearer | Direct debit from wallet. |
| `POST` | `/:id/credit` | `.../wallet/:id/credit` | Bearer | Direct credit to wallet. |
| `POST` | `/transfer` | `.../wallet/transfer` | Bearer | Internal transfer between users. |
| `POST` | `/staff/topup` | `.../wallet/staff/topup` | Bearer | Mitra tops up linked Staff wallet. |
| `POST` | `/:id/topup` | `.../wallet/:id/topup` | Bearer | Mitra tops up own wallet. |
| `GET` | `/:id/events` | `.../wallet/:id/events` | Bearer | Paginated transaction history of the wallet. |
| `GET` | `/:id/reconcile` | `.../wallet/:id/reconcile` | Bearer | Check balance vs events drift. |

---

## 5. Transaction Service
**Base URL:** `https://fedora.sinauplatform.id/api/v1/transaction`

Orchestrates the lifecycle of a PPOB transaction.

### Endpoints

| Method | Endpoint | Full URL | Auth | Description |
|:---|:---|:---|:---|:---|
| `POST` | `/initiate` | `.../transaction/initiate` | Bearer | Start new transaction (requires `Idempotency-Key`). |
| `GET` | `/:id` | `.../transaction/:id` | Bearer | Get transaction status by internal ID. |
| `GET` | `/by-id/:id` | `.../transaction/by-id/:id` | Bearer | Get transaction by Transaction UUID. |
| `GET` | `/` | `.../transaction` | Bearer | List transactions with filters. |
| `GET` | `/history` | `.../transaction/history` | Bearer | Paginated transaction history (cursor-based). |
| `POST` | `/:id/status` | `.../transaction/:id/status` | Bearer | Manually update transaction status (Admin). |
| `POST` | `/:id/cancel` | `.../transaction/:id/cancel` | Bearer | Cancel a pending transaction. |
| `POST` | `/webhook/digiflazz` | `.../transaction/webhook/digiflazz` | None | Webhook endpoint for Digiflazz updates. |
| `GET` | `/reports` | `.../transaction/reports` | Bearer (Mitra/Admin) | Get aggregated reports: KPIs, sales trend, staff performance. |

---

## 6. Integration Service
**Base URL:** `https://fedora.sinauplatform.id/api/v1/integration`

Gateway for external provider (Digiflazz) communication.

### Endpoints

| Method | Endpoint | Full URL | Auth | Description |
|:---|:---|:---|:---|:---|
| `POST` | `/digiflazz/transaction` | `.../integration/digiflazz/transaction` | Bearer | Forward transaction to Digiflazz API. |
| `GET` | `/providers` | `.../integration/providers` | Bearer | List provider configurations. |
| `GET` | `/errors` | `.../integration/errors` | Bearer | Retrieve mapped error code catalog. |
| `GET` | `/compensation/jobs` | `.../integration/compensation/jobs` | Bearer | List background compensation/retry jobs. |
| `GET` | `/compensation/dead-letter` | `.../integration/compensation/dead-letter` | Bearer | List failed jobs in dead letter queue. |
| `POST` | `/webhook/digiflazz` | `.../integration/webhook/digiflazz` | None | Webhook verification and processing. |

---

## Error Handling

All services return errors in a standard JSON format:

```json
{
  "error": {
    "code": "TRANSACTION_INSUFFICIENT_BALANCE",
    "message": "Saldo tidak mencukupi",
    "details": { "required": 50000, "current": 25000 },
    "trace_id": "uuid-string",
    "timestamp": "2026-05-08T..."
  }
}
```

### Common Status Codes
- `400 Bad Request`: Validation failure or business rule violation.
- `401 Unauthorized`: Token expired or invalid.
- `403 Forbidden`: Insufficient permissions or account locked.
- `429 Too Many Requests`: Rate limit exceeded.
- `502 Bad Gateway`: Digiflazz API error or timeout.
- `503 Service Unavailable`: Circuit breaker open or system maintenance.

---

## Monitoring
- `GET https://fedora.sinauplatform.id/api/v1/health/ready`: Liveness/Readiness probe for Kubernetes.
- `GET https://fedora.sinauplatform.id/api/v1/metrics`: Prometheus metrics exporter.

<br/>

**Last Updated:** 2026-05-09
