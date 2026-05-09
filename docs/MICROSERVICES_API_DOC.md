# PPOB Microservices API Documentation

This document provides a detailed reference for all API endpoints across the PPOB (Payment Point Online Bank) microservices ecosystem.

---

## 1. Auth Service
**Port:** 8081  
**Base Path:** `/api/v1/auth`

Handles user registration, authentication, OTP verification, and session management.

### Endpoints

| Method | Endpoint | Auth | Description |
|:---|:---|:---|:---|
| `POST` | `/register` | None | Register a new user with email, phone, name, password, and PIN. |
| `POST` | `/login` | None | Login using Email or Phone + Password. |
| `POST` | `/verify-otp` | None | Verify OTP code for login or registration. |
| `POST` | `/refresh` | None | Refresh access token using a valid refresh token. |
| `POST` | `/logout` | Bearer | Logout and invalidate session. |
| `POST` | `/change-password` | Bearer | Update account password. |
| `POST` | `/change-pin` | Bearer | Update 6-digit transaction PIN. |

---

## 2. User Service
**Port:** 8082  
**Base Path:** `/api/v1`

Manages user profiles, roles, staff management, and administrative lists.

### Endpoints

| Method | Endpoint | Role | Description |
|:---|:---|:---|
| `GET` | `/users/:id` | Owner/Admin | Get detailed profile information for a user. |
| `PUT` | `/users/:id` | Owner/Admin | Update profile fields (name, phone, address, DOB). |
| `GET` | `/users` | Admin/Staff | Paginated list of all users. |
| `GET` | `/users/:id/roles` | Owner/Admin | Get roles assigned to a user. |
| `POST` | `/users/:id/roles` | Admin | Assign a role (e.g., Mitra, Staff) to a user. |
| `GET` | `/roles` | Admin | List available system roles. |
| `POST` | `/roles` | Admin | Create a new role definition. |
| `GET` | `/staff` | Mitra | List staff users with stats (transactions, wallet, limits). |
| `POST` | `/staff` | Mitra | Create a new staff user. |
| `GET` | `/staff/:id` | Mitra/Admin | Get detailed staff info including margin settings. |
| `PUT` | `/staff/:id` | Mitra | Update staff details and limits. |
| `GET` | `/staff/:id/stats` | Mitra/Admin | Get staff performance stats. |
| `GET` | `/staff/pending-count` | Mitra | Count of pending staff invitations. |
| `GET` | `/notifications` | Bearer | List user notifications (optional query: unread_only). |
| `GET` | `/notifications/uncount` | Bearer | Get unread notifications count. |
| `PATCH` | `/notifications/:id/read` | Bearer | Mark a notification as read. |
| `POST` | `/notifications/mark-all-read` | Bearer | Mark all notifications as read. |

---

## 3. Product Service
**Port:** 8085  
**Base Path:** `/api/v1`

Manages product catalog and synchronization with Digiflazz.

### Endpoints

| Method | Endpoint | Role | Description |
|:---|:---|:---|:---|
| `GET` | `/products` | None | List products with category/status filters. |
| `GET` | `/products/:id` | None | Get product by ID. |
| `GET` | `/products/search` | None | Search products by name/code. |
| `GET` | `/products/by-code/:code` | None | Get product by SKU code. |
| `GET` | `/products/validate-price` | None | Check if a price meets margin requirements. |
| `POST` | `/products` | Admin | Create product manually. |
| `PUT` | `/products/:id` | Admin | Update product information. |
| `DELETE` | `/products/:id` | Admin | Soft-delete a product. |
| `POST` | `/sync/prepaid` | Bearer | Trigger prepaid product sync. |
| `POST` | `/sync/postpaid` | Bearer | Trigger postpaid product sync. |
| `GET` | `/sync/status` | None | Get last sync timestamps. |
| `GET` | `/categories` | None | List all categories. |
| `POST` | `/categories` | Admin | Create a category. |

---

## 4. Wallet Service
**Port:** 8083  
**Base Path:** `/api/v1/wallets`

Financial ledger handling balance, holds, and transfers.

### Endpoints

| Method | Endpoint | Auth | Description |
|:---|:---|:---|:---|
| `GET` | `/:id/balance` | Bearer | Get current available and held balance. |
| `GET` | `/:id/balance-events` | Bearer | Get balance reconstructed from events. |
| `POST` | `/:id/hold` | Bearer | Place a hold on funds for a transaction. |
| `POST` | `/:id/release-hold` | Bearer | Release a previously held amount. |
| `POST` | `/:id/debit` | Bearer | Direct debit from wallet. |
| `POST` | `/:id/credit` | Bearer | Direct credit to wallet. |
| `POST` | `/transfer` | Bearer | Internal transfer between users. |
| `POST` | `/staff/topup` | Bearer | Mitra tops up linked Staff wallet. |
| `POST` | `/:id/topup` | Bearer | Mitra tops up own wallet. |
| `GET` | `/:id/events` | Bearer | Paginated transaction history of the wallet. |
| `GET` | `/:id/reconcile` | Bearer | Check balance vs events drift. |

---

## 5. Transaction Service
**Port:** 8084  
**Base Path:** `/api/v1/transactions`

Orchestrates the lifecycle of a PPOB transaction.

### Endpoints

| Method | Endpoint | Auth | Description |
|:---|:---|:---|:---|
| `POST` | `/initiate` | Bearer | Start new transaction (requires `Idempotency-Key`). |
| `GET` | `/:id` | Bearer | Get transaction status by internal ID. |
| `GET` | `/by-id/:id` | Bearer | Get transaction by Transaction UUID. |
| `GET` | `/` | Bearer | List transactions with filters. |
| `GET` | `/history` | Bearer | Paginated transaction history (cursor-based). |
| `POST` | `/:id/status` | Bearer | Manually update transaction status (Admin). |
| `POST` | `/:id/cancel` | Bearer | Cancel a pending transaction. |
| `POST` | `/webhook/digiflazz` | None | Webhook endpoint for Digiflazz updates. |
| `GET` | `/reports` | Bearer (Mitra/Admin) | Get aggregated reports: KPIs, sales trend, staff performance. |

---

## 6. Integration Service
**Port:** 8086  
**Base Path:** `/api/v1/integrations`

Gateway for external provider (Digiflazz) communication.

### Endpoints

| Method | Endpoint | Auth | Description |
|:---|:---|:---|:---|
| `POST` | `/digiflazz/transaction` | Bearer | Forward transaction to Digiflazz API. |
| `GET` | `/providers` | Bearer | List provider configurations. |
| `GET` | `/errors` | Bearer | Retrieve mapped error code catalog. |
| `GET` | `/compensation/jobs` | Bearer | List background compensation/retry jobs. |
| `GET` | `/compensation/dead-letter` | Bearer | List failed jobs in dead letter queue. |
| `POST` | `/webhook/digiflazz` | None | Webhook verification and processing. |

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
- `GET /health/ready`: Liveness/Readiness probe for Kubernetes.
- `GET /metrics`: Prometheus metrics exporter.

---

### Global Headers
- `Authorization: Bearer <JWT_ACCESS_TOKEN>`
- `Idempotency-Key: <UUID>` (Required for `/transactions/initiate`)
- `X-Trace-ID: <UUID>` (Optional, for request tracing)

<br/>

**Generated on:** Friday, May 8, 2026
