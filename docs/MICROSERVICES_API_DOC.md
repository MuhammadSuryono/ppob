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
| `POST` | `/initiate` | `.../auth/initiate` | None | Check phone registration and device trust status. |
| `POST` | `/register` | `.../auth/register` | None | Register a new user (OTP must be verified first). |
| `POST` | `/login` | `.../auth/login` | None | Classic login using Email/Phone + Password. |
| `POST` | `/send-otp` | `.../auth/send-otp` | None | Request OTP for registration or login (returns `request_id`). |
| `POST` | `/verify-otp` | `.../auth/verify-otp` | None | Verify OTP code and obtain verified `request_id`. |
| `POST` | `/verify-credential` | `.../auth/verify-credential` | None | Verify PIN or Password after OTP (unified post-OTP auth for existing users on untrusted devices). |
| `POST` | `/verify-password` | `.../auth/verify-password` | None | Validate password for existing users on untrusted devices (requires verified `request_id`). **Consider using `/verify-credential` instead.** |
| `POST` | `/verify-pin` | `.../auth/verify-pin` | None | Final auth step for trusted devices using PIN. |
| `POST` | `/refresh` | `.../auth/refresh` | None | Refresh access token using a valid refresh token. |
| `POST` | `/logout` | `.../auth/logout` | Bearer | Logout and invalidate session. |
| `POST` | `/change-password` | `.../auth/change-password` | Bearer | Update account password. |
| `POST` | `/change-pin` | `.../auth/change-pin` | Bearer | Update 6-digit transaction PIN. |

---

### Detailed Endpoint Specifications

---

#### `POST /auth/initiate`

**Description:** Checks whether the phone number belongs to a registered user, and whether the device is trusted. This is the **entry point** for the adaptive auth flow.

**Request Body:**
```json
{
  "phone": "+6281234567890",
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "fingerprint": "(optional) additional device fingerprint data"
}
```

**Response Body:**
```json
{
  "user_id": 1,
  "is_registered": true,
  "is_trusted": false,
  "requires_otp": true
}
```

**Field Explanation:**
| Field | Description |
|---|---|
| `is_registered` | Whether the phone number is already registered in the system |
| `is_trusted` | Whether the device fingerprint is already marked as trusted for this user |
| `requires_otp` | Whether OTP verification is required before proceeding (true for new users or untrusted devices) |

**Flow Decision:**
- `is_trusted=true` → Navigate directly to **PIN Login** (Screen 4)
- `requires_otp=true` → Navigate to **OTP Input** (Screen 2)
- `is_registered=false` → After OTP, navigate to **Registration** (Screen 3: Create PIN + Password)
- `is_registered=true` + `requires_otp=true` → After OTP, navigate to **Credential Input** (Screen 3: PIN or Password)

**Errors:**
- `500`: System error

---

#### `POST /auth/send-otp`

**Description:** Triggers OTP delivery via SMS/WhatsApp for either registration or login flow. Returns a unique `request_id` that must be used in subsequent OTP verification and final auth steps.

**Security:** Rate limited to 5 requests per minute per IP address to prevent abuse.

**Request Body:**
```json
{
  "phone": "+6281234567890",
  "type": "login"
}
```

**Request Body:**
```json
{
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "expires_at": 1746975300
}
```

**Errors:**
- `400`: Invalid phone format or missing type
- `429`: Rate limit exceeded
- `500`: OTP delivery failed or system error

**Note:** In development, OTP codes are printed to server logs. Production deployments should integrate with SMS/WhatsApp gateway provider.

---

#### `POST /auth/verify-otp`

**Description:** Validates the 6-digit OTP code. Upon success, marks the `request_id` as verified for 10 minutes. **No token is returned** — client must proceed to `/register` (new user) or `/verify-credential` (existing user) with the verified `request_id`.

**Request Body:**
```json
{
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "phone": "+6281234567890",
  "code": "123456",
  "type": "login"
}
```

**Response Body:**
```json
{
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "is_verified": true,
  "is_new_user": false
}
```

**Field Explanation:**
| Field | Description |
|---|---|
| `is_verified` | Whether OTP verification succeeded |
| `is_new_user` | `true` if the phone number is **not** registered yet (→ navigate to Registration); `false` if the phone number exists (→ navigate to Credential Input) |

**Usage:** Mobile app uses `is_new_user` to decide the next screen:
- `is_new_user=true` → Show **Create PIN + Create Password** screen, then call `/register`
- `is_new_user=false` → Show **Input PIN or Password** screen, then call `/verify-credential`

**Errors:**
- `400`: Invalid or expired OTP
- `500`: System error

---

#### `POST /auth/register`

**Description:** Completes registration for a **new user**. **Requires a verified `request_id`** from a prior OTP verification.

**Request Body:**
```json
{
  "email": "user@example.com",
  "phone": "+6281234567890",
  "full_name": "John Doe",
  "password": "SecurePass123!",
  "pin": "123456",
  "device_id": "device-fingerprint-hash",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response Body:**
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "phone": "+6281234567890",
  "full_name": "John Doe",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": 1746975300,
  "refresh_expires_at": 1747575300
}
```

**Notes:**
- The `request_id` must exist under Redis key `verified:{request_id}` and match the phone number.
- OTP verification must have been completed within the last 10 minutes.
- On success, creates the user, wallet, and marks the device as trusted.
- **Both access and refresh tokens are returned** — unlike `/verify-password` and `/login`, this endpoint previously only returned an access token.

**Errors:**
- `400`: `request_id` not verified, phone mismatch, or user already exists (`AUTH_OTP_NOT_VERIFIED`, `AUTH_USER_EXISTS`)
- `500`: System error

---

#### `POST /auth/verify-credential`

**Description:** **NEW** — Verifies credentials (PIN or Password) for an **existing user** on an untrusted device after OTP verification. This is the unified endpoint for the post-OTP credential step, replacing the need to call `/verify-password` or `/verify-pin` separately in the OTP flow.

**Request Body:**
```json
{
  "phone": "+6281234567890",
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "auth_method": "pin",
  "value": "123456"
}
```

**Field Explanation:**
| Field | Required | Description |
|---|---|---|
| `phone` | Yes | User phone number (must start with `+62`) |
| `device_id` | Yes | Device fingerprint hash |
| `request_id` | Yes | Verified request ID from OTP step |
| `auth_method` | Yes | `"pin"` or `"password"` — determines which credential is being verified |
| `value` | Yes | The PIN (6-digit) or password (8+ char) value |

**Response Body:** Same as `/auth/login` and `/auth/verify-password`:
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "phone": "+6281234567890",
  "full_name": "John Doe",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": 1746975300
}
```

**Behavior:**
- Validates the `request_id` verification flag (must exist and match phone)
- Verifies the credential against the stored hash (bcrypt)
- **Automatically marks the device as trusted** for future PIN-only logins
- Consumes (deletes) the `verified:{request_id}` flag — single use

**Errors:**
- `400`: `request_id` not verified (`AUTH_OTP_NOT_VERIFIED`)
- `401`: Invalid credential (`AUTH_INVALID_CREDENTIALS`)
- `500`: System error

---

#### `POST /auth/verify-password`

**Description:** Validates password for an existing user on an untrusted device. **Requires a verified `request_id`**. On success, returns tokens and marks the device as trusted. **NOTE:** Consider using `/verify-credential` (with `auth_method: "password"`) for a unified flow instead.

**Request Body:**
```json
{
  "phone": "+6281234567890",
  "password": "UserPassword123!",
  "device_id": "device-fingerprint-hash",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response Body:**
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "phone": "+6281234567890",
  "full_name": "John Doe",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": 1746975300
}
```

**Notes:**
- The `request_id` verification flag is consumed (deleted) after successful verification.
- The device is marked as trusted upon successful password validation.

**Errors:**
- `400`: `request_id` not verified
- `401`: Invalid password
- `500`: System error

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

### Detailed Endpoint Specifications

---

#### `GET /product/categories`

**Description:** Lists all product categories with metadata for dynamic UI rendering.

**Response Body:**
```json
{
  "categories": [
    {
      "id": 1,
      "name": "Pulsa",
      "code": "pulsa",
      "icon": "https://...",
      "input_type": "NUMBER",
      "input_label": "Nomor HP",
      "placeholder": "08xxxxxxxxxx",
      "validation_regex": "^[0-9]{10,14}$"
    },
    {
      "id": 2,
      "name": "PLN",
      "code": "pln",
      "icon": "https://...",
      "input_type": "NUMBER",
      "input_label": "ID Pelanggan",
      "placeholder": "12 digit nomor meter",
      "validation_regex": "^[0-9]{11,12}$"
    }
  ]
}
```

**Field Explanation (Metadata):**
| Field | Description |
|---|---|
| `input_type` | UI input type: `TEXT`, `NUMBER`, `PHONE`. |
| `input_label` | Label to display above the input field. |
| `placeholder` | Placeholder text for the input field. |
| `validation_regex` | Regex pattern for frontend validation. |

---

---

## 4. Wallet Service
**Base URL:** `https://fedora.sinauplatform.id/api/v1/wallet`

Financial ledger handling balance, holds, and transfers.

### Endpoints

| Method | Endpoint | Full URL | Auth | Description |
|:---|:---|:---|:---|:---|
| `GET` | `/me/balance` | `.../wallet/me/balance` | Bearer | Get current available and held balance for logged-in user. |
| `GET` | `/me/balance-events` | `.../wallet/me/balance-events` | Bearer | Get balance reconstructed from events for logged-in user. |
| `POST` | `/:id/hold` | `.../wallet/:id/hold` | Bearer | Place a hold on funds for a transaction. |
| `POST` | `/:id/release-hold` | `.../wallet/:id/release-hold` | Bearer | Release a previously held amount. |
| `POST` | `/:id/debit` | `.../wallet/:id/debit` | Bearer | Direct debit from wallet. |
| `POST` | `/:id/credit` | `.../wallet/:id/credit` | Bearer | Direct credit to wallet. |
| `POST` | `/transfer` | `.../wallet/transfer` | Bearer | Internal transfer between users. |
| `POST` | `/staff/topup` | `.../wallet/staff/topup` | Bearer | Mitra tops up linked Staff wallet. |
| `POST` | `/me/topup` | `.../wallet/me/topup` | Bearer | Mitra tops up own wallet. |
| `GET` | `/me/events` | `.../wallet/me/events` | Bearer | Paginated transaction history of the wallet for logged-in user. |
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

**Last Updated:** 2026-05-17
ated:** 2026-05-09
BALANCE",
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
