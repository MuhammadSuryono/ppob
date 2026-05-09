# API Contracts for PPOB Mobile Application Backend (Updated)

This document outlines the RESTful API contracts for the PPOB mobile application backend, implemented using Golang microservices. It covers authentication, user and role management, wallet operations, product information, and transaction processing.

**Latest Updates:**
- Unified login endpoint (single `/auth/login` with adaptive trust logic)
- Added pagination & filtering to list endpoints
- Standardized error response format (see `ERROR_HANDLING.md`)
- Added idempotency-key requirement for transaction initiation
- Fixed top-up endpoint to be staff-specific with ownership verification
- Added device fingerprint field to login flow
- Added transaction state machine details

---

## 1. Authentication Service API

**Base URL:** `/auth/v1` (versioned)

### 1.1. Register User

**Endpoint:** `POST /auth/v1/register`

**Description:** Registers a new user with a phone number. Sends OTP via SMS.

**Request Body:**
```json
{
  "phone_number": "+6281234567890",
  "name": "John Doe"
}
```

**Response (200 OK):**
```json
{
  "message": "Kode OTP telah dikirim ke nomor Anda"
}
```

**Rate Limit:** 3 per hour per phone number.

---

### 1.2. Verify OTP and Set Password/PIN

**Endpoint:** `POST /auth/v1/verify-otp`

**Description:** Verifies OTP and creates user account with password and PIN.

**Request Body:**
```json
{
  "phone_number": "+6281234567890",
  "otp_code": "123456",
  "password": "StrongPassword123!",
  "pin": "123456"
}
```

**Validation Rules:**
- `password`: min 8 chars, 1 uppercase, 1 lowercase, 1 digit
- `pin`: exactly 6 digits, not sequential, not all same

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "random256bitstring",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "user_id": "uuid_user_id",
    "phone_number": "+6281234567890",
    "name": "John Doe",
    "roles": [
      {"role_id": "uuid_mitra_role", "role_name": "Mitra"}
    ],
    "active_role": {"role_id": "uuid_mitra_role", "role_name": "Mitra"},
    "wallet_id": "uuid_wallet_id"
  }
}
```

---

### 1.3. Login User (Unified Endpoint)

**Endpoint:** `POST /auth/v1/login`

**Description:** Single login endpoint that adapts based on device trust score. Always requires password. May require OTP based on trust level. PIN always required.

**Request Body:**
```json
{
  "phone_number": "+6281234567890",
  "password": "StrongPassword123!",
  "pin": "123456",
  "device_fingerprint": {
    "device_id": "550e8400-e29b-41d4-a716-446655440000",
    "user_agent": "Android/14 (Pixel 7) PPOB/2.1.0",
    "app_version": "2.1.0",
    "install_ts": 1700000000000,
    "last_login_ts": 1704800000000
  },
  "otp_code": "654321"  // Optional: only required if device trust level demands it
}
```

**Trust Logic (Server-Side):**
- Trusted (score ≥70): password + PIN only
- Semi-trusted (30–69): password + OTP + PIN
- Untrusted (<30): password + OTP + PIN (may add CAPTCHA future)

**Response (200 OK):** Same as register response.

**Error Responses:**
- `400 AUTH_INVALID_CREDENTIALS` — wrong phone or password
- `400 AUTH_OTP_INVALID` — OTP wrong (if required)
- `403 AUTH_ACCOUNT_LOCKED` — too many PIN attempts
- `400 AUTH_DEVICE_NOT_TRUSTED` — response includes `requires_otp: true`

---

### 1.4. Refresh Token

**Endpoint:** `POST /auth/v1/refresh`

**Description:** Obtain new access token using refresh token. Refresh token is single-use; new one returned.

**Request Body:**
```json
{
  "refresh_token": "old_refresh_token_here"
}
```

**Response (200 OK):**
```json
{
  "access_token": "new_access_token",
  "refresh_token": "new_refresh_token",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**Error:** `401 AUTH_TOKEN_EXPIRED` if refresh token expired or revoked.

---

### 1.5. Logout

**Endpoint:** `POST /auth/v1/logout`

**Description:** Invalidate current refresh token.

**Headers:** `Authorization: Bearer <refresh_token>`

**Response (204 No Content)**

---

### 1.6. Logout All Devices

**Endpoint:** `POST /auth/v1/logout-all`

**Description:** Invalidate all refresh tokens for user.

**Headers:** `Authorization: Bearer <access_token>`

**Response (204)**

---

### 1.7. Change PIN

**Endpoint:** `POST /auth/v1/change-pin`

**Description:** Update transaction PIN. Requires re-authentication (current PIN OR password+OTP).

**Request Body (Option A — current PIN):**
```json
{
  "current_pin": "123456",
  "new_pin": "654321",
  "confirm_pin": "654321"
}
```

**Request Body (Option B — password+OTP for forgotten PIN):**
```json
{
  "password": "StrongPassword123!",
  "otp_code": "654321",
  "new_pin": "654321",
  "confirm_pin": "654321"
}
```

**Response (200 OK):**
```json
{
  "message": "PIN berhasil diubah"
}
```

**Side Effect:** All existing sessions invalidated (security precaution).

---

### 1.8. Forgot Password

**Endpoint:** `POST /auth/v1/forgot-password`

**Description:** Initiate password reset via OTP.

**Request Body:**
```json
{
  "phone_number": "+6281234567890"
}
```

**Response (200):** `{"message":"OTP untuk reset password dikirim"}`

Then call `/auth/v1/reset-password` with OTP + new password.

---

## 2. User & Role Service API

**Base URL:** `/users/v1`

**Headers:** All endpoints require `Authorization: Bearer <access_token>`

### 2.1. Get User Profile

**Endpoint:** `GET /users/v1/profile`

**Response (200 OK):**
```json
{
  "user_id": "uuid_user_id",
  "phone_number": "+6281234567890",
  "name": "John Doe",
  "roles": [
    {"role_id": "uuid_mitra", "role_name": "Mitra"},
    {"role_id": "uuid_staff_b", "role_name": "Staff"}
  ],
  "active_role": {"role_id": "uuid_mitra", "role_name": "Mitra"},
  "wallet": {
    "wallet_id": "uuid_wallet",
    "balance_available": 1500000.50,
    "balance_held": 25000.00
  }
}
```

---

### 2.2. Switch Active Role

**Endpoint:** `POST /users/v1/switch-role`

**Description:** Changes the active role for the session. Returns new JWT with updated claims.

**Request Body:**
```json
{
  "role_id": "uuid_target_role_id"
}
```

**Validation:** User must have this role in `user_roles`.

**Response (200 OK):**
```json
{
  "access_token": "new_jwt_with_new_active_role",
  "refresh_token": "new_refresh_token",  // optional: rotate refresh too
  "expires_in": 900,
  "active_role": {"role_id": "uuid_staff_role", "role_name": "Staff"},
  "wallet_id": "uuid_staff_wallet"
}
```

---

### 2.3. Add Staff (Mitra Only)

**Endpoint:** `POST /users/v1/staff`

**Authorization:** Requires Mitra role with permission `staff:create`.

**Request Body:**
```json
{
  "phone_number": "+6281234567891",
  "name": "Budi Santoso",
  "password": "StaffPass123!",
  "pin": "654321",
  "margin_scheme": "FixedAllowance",  // or "MarginShare"
  "margin_value": 10000,  // Fixed: Rp 10,000 per txn; Share: 60 (60%)
  "daily_txn_limit": 50,  // optional; defaults to system setting
  "daily_amount_limit": 5000000
}
```

**Process:**
1. Create user (phone_number, hashed password, hashed PIN)
2. Assign Staff role with `assigned_by = current_mitra_user_id`
3. Create staff wallet (trigger does it)
4. Create `staff_global_margin_settings` with provided scheme

**Response (201 Created):**
```json
{
  "message": "Staff berhasil ditambahkan",
  "staff": {
    "user_id": "uuid_staff_id",
    "phone_number": "+6281234567891",
    "name": "Budi Santoso"
  }
}
```

**Error:** `403 AUTH_INSUFFICIENT_PERMISSION` if caller not Mitra.

---

### 2.4. List Staff (Mitra Only)

**Endpoint:** `GET /users/v1/staff`

**Query Parameters:** `?limit=20&offset=0&sort=name_asc`

**Response (200 OK):**
```json
{
  "staff": [
    {
      "user_id": "uuid",
      "name": "Budi Santoso",
      "phone_number": "+6281234567891",
      "wallet_balance": 125000.00,
      "daily_txn_count": 12,
      "daily_txn_amount": 340000,
      "margin_scheme": "FixedAllowance",
      "margin_value": 10000,
      "is_active": true
    }
  ],
  "pagination": {
    "total": 15,
    "limit": 20,
    "offset": 0,
    "has_more": false
  }
}
```

---

### 2.5. Update Staff Settings (Mitra Only)

**Endpoint:** `PUT /users/v1/staff/{staff_id}`

**Request Body:**
```json
{
  "name": "Budi Santoso Updated",  // optional
  "margin_scheme": "MarginShare",
  "margin_value": 60.00,
  "daily_txn_limit": 100,
  "daily_amount_limit": 10000000,
  "is_active": true
}
```

**Validation:** Mitra must be `assigned_by` of this staff.

**Response (200 OK):** `{"message":"Staff updated"}`

---

### 2.6. Get Trusted Devices (User)

**Endpoint:** `GET /users/v1/devices`

**Response:**
```json
{
  "devices": [
    {
      "device_id": "uuid",
      "user_agent": "Android/14...",
      "last_seen": "2026-05-07T00:30:15Z",
      "is_trusted": true,
      "trust_score": 85
    }
  ]
}
```

---

### 2.7. Revoke Device Trust

**Endpoint:** `DELETE /users/v1/devices/{device_id}`

**Response (204)** — removes device from trusted list; next login from that device requires full OTP.

---

## 3. Wallet Service API

**Base URL:** `/wallets/v1`

### 3.1. Get Wallet Balance

**Endpoint:** `GET /wallets/v1/balance`

**Response (200 OK):**
```json
{
  "wallet_id": "uuid_wallet_id",
  "balance_available": 1500000.50,
  "balance_held": 25000.00,
  "balance_total": 1525000.50,
  "currency": "IDR",
  "updated_at": "2026-05-07T00:25:10Z"
}
```

---

### 3.2. Top-Up Staff Wallet (Mitra Only)

**Endpoint:** `POST /wallets/v1/staff/{staff_id}/topup`

**Authorization:** Mitra role;必须是该staff的assigned_by Mitra.

**Request Body:**
```json
{
  "amount": 500000.00
}
```

**Validation:**
- `amount > 0`
- Mitra main wallet `balance_available >= amount`
- Staff exists and belongs to this Mitra

**Process (atomic within DB transaction):**
1. Debit Mitra main wallet (`balance_available -= amount`)
2. Credit staff wallet (`balance_available += amount`)
3. Create wallet_events for both wallets
4. Audit log `wallet_topup`

**Response (200 OK):**
```json
{
  "message": "Top up berhasil",
  "staff_wallet": {
    "wallet_id": "uuid_staff_wallet",
    "balance_available": 500000.00
  },
  "mitra_wallet": {
    "wallet_id": "uuid_mitra_wallet",
    "balance_available": 1000000.50
  }
}
```

**Error:** `403 TRANSACTION_INSUFFICIENT_BALANCE` if Mitra lacks funds.

---

## 4. Product Service API

**Base URL:** `/products/v1`

### 4.1. Get Product Categories

**Endpoint:** `GET /products/v1/categories`

**Response (200 OK):**
```json
{
  "categories": [
    {
      "category_id": "uuid_1",
      "category_name": "Pulsa",
      "icon_url": "https://cdn.ppob.co.id/icons/pulsa.png",
      "display_order": 1
    },
    {
      "category_id": "uuid_2",
      "category_name": "PLN",
      "icon_url": "https://cdn.ppob.co.id/icons/pln.png",
      "display_order": 2
    }
  ]
}
```

---

### 4.2. Get Products by Category (with Pagination)

**Endpoint:** `GET /products/v1?category_id={category_id}&limit=20&cursor={cursor}`

**Query Parameters:**
- `category_id` (required)
- `limit` (default 20, max 100)
- `cursor` (optional, for pagination — base64 encoded last product_id)
- `sort` (default `name_asc`)

**Response (200 OK):**
```json
{
  "products": [
    {
      "product_id": "uuid_prod1",
      "buyer_sku_code": "xld25",
      "product_name": "XL Data 25GB",
      "base_price": 25000.00,
      "platform_price": 26250.00,
      "is_prepaid": true,
      "mitra_selling_price": 27000.00  // only if Mitra set custom; else null
    }
  ],
  "pagination": {
    "next_cursor": "base64_of_last_id",
    "has_more": true
  }
}
```

**Note:** `mitra_selling_price` returned only for Mitra role (staff see platform or Mitra's set price automatically applied at transaction time, not listed).

---

## 5. Transaction Service API

**Base URL:** `/transactions/v1`

### 5.1. Initiate Transaction

**Endpoint:** `POST /transactions/v1/initiate`

**Headers:** `Idempotency-Key: uuidv4` (required, unique per client attempt)

**Request Body:**
```json
{
  "product_id": "uuid_product_id",
  "customer_no": "081234567890",
  "selling_price": 27000.00,
  "pin": "123456",
  "device_fingerprint": "sha256:..."  // optional: for fraud scoring
}
```

**Validations (server):**
- Product exists and active
- `selling_price >= product.platform_price` (enforced)
- Wallet balance (available) sufficient for `selling_price`
- Daily limit not exceeded (if staff)
- PIN correct
- Idempotency key not used before
- Customer number format valid per product regex

**Process:**
1. Check idempotency → if exists, return existing transaction
2. Place hold on wallet (`balance_held += selling_price`, `balance_available -= selling_price`)
3. Call Digiflazz transaction endpoint
4. On RC 00 → debit hold immediately, mark Success
5. On RC 03 → mark Pending, hold remains
6. On other RC → release hold, mark Failed
7. Store transaction record

**Response:**
- **202 Accepted** if Pending (async webhook expected)
- **200 OK** if Success (immediate)
- **4xx** if validation/business rule fails
- **502/503** if Digiflazz unavailable

```json
{
  "transaction_id": "uuid_txn",
  "ref_id": "digiflazz_ref_123",
  "status": "Pending",  // or "Success"
  "message": "Transaksi sedang diproses",
  "selling_price": 27000,
  "platform_price": 26250,
  "commission_amount": 150  // if margin share: (27000-26250)*0.6; null if FixedAllowance
}
```

**Error Codes:**
- `400 VALIDATION_*` — input errors
- `400 TRANSACTION_INSUFFICIENT_BALANCE`
- `403 TRANSACTION_DAILY_LIMIT_EXCEEDED`
- `409 TRANSACTION_ALREADY_PROCESSING` (idempotency conflict)
- `502 DIGIFLAZZ_TIMEOUT` (retryable)
- `503 DIGIFLAZZ_OUT_OF_STOCK`

---

### 5.2. Get Transaction Status

**Endpoint:** `GET /transactions/v1/{transaction_id}/status`

**Response (200 OK):**
```json
{
  "transaction_id": "uuid",
  "ref_id": "digiflazz_ref_123",
  "status": "Success",
  "message": "Transaksi berhasil",
  "serial_number": "SN123456789",  // if applicable
  "details": {
    "product_name": "XL Data 25GB",
    "customer_no": "+6281234567890",
    "amount": 26250.00,  // platform price (cost)
    "selling_price": 27000.00,
    "commission_amount": 150.00,
    "margin_amount": 750.00,  // selling - platform
    "created_at": "2026-05-07T00:30:15Z",
    "completed_at": "2026-05-07T00:30:20Z"
  }
}
```

---

### 5.3. Get Transaction History (Paginated)

**Endpoint:** `GET /transactions/v1/history`

**Query Parameters:**
- `limit` (default 20, max 100)
- `cursor` (pagination token)
- `start_date` (ISO 8601, optional)
- `end_date` (ISO 8601)
- `status` (optional: Success, Pending, Failed, etc.)
- `category_id` (optional filter by product category)

**Response (200 OK):**
```json
{
  "transactions": [ /* array of transaction summaries */ ],
  "pagination": {
    "next_cursor": "base64...",
    "has_more": true
  }
}
```

---

### 5.4. Cancel Transaction (Pending Only)

**Endpoint:** `POST /transactions/v1/{transaction_id}/cancel`

**Description:** Cancel a pending transaction (before webhook success).

**Request Body:**
```json
{
  "reason": "User requested cancellation"
}
```

**Response (200 OK):**
```json
{
  "message": "Transaksi dibatalkan",
  "status": "Cancelled"
}
```

**Errors:**
- `400 TRANSACTION_NOT_CANCELLABLE` if already Success/Failed
- `403 AUTH_INSUFFICIENT_PERMISSION` if not owner or Mitra

---

## 6. Wallet & Commission Reporting

### 6.1. Get Commission History (Staff)

**Endpoint:** `GET /wallets/v1/commissions`

**Query:** `?start_date=&end_date=`

**Response:**
```json
{
  "commissions": [
    {
      "commission_id": "uuid",
      "transaction_id": "uuid_txn",
      "amount": 150.00,
      "scheme_used": "MarginShare",
      "scheme_value": 60,
      "margin_amount": 750.00,
      "earned_at": "2026-05-07T00:30:20Z",
      "paid_at": null  // if not yet transferred (but credited to wallet immediately)
    }
  ]
}
```

---

## 7. Error Response Standard

**All errors follow format defined in `ERROR_HANDLING.md`:**

```json
{
  "error": {
    "code": "TRANSACTION_DAILY_LIMIT_EXCEEDED",
    "message": "Limit harian transaksi tercapai. Batas: 50 transaksi / Rp 5.000.000 per hari",
    "details": {
      "limit_type": "count",
      "used_today": 50,
      "reset_at": "2026-05-08T00:00:00+07:00"
    },
    "trace_id": "abc123",
    "timestamp": "2026-05-07T00:35:00Z"
  }
}
```

**HTTP Status Mapping:**
- `2xx` — Success
- `4xx` — Client error (validation, business rule)
- `401/403` — Authz failure
- `429` — Rate limit (include `Retry-After` header)
- `5xx` — Server / upstream (Digiflazz, DB, etc.)

---

## 8. Pagination

**Cursor-based pagination recommended** for performance (avoid OFFSET).

**Request:**
```
GET /transactions/v1/history?limit=20&cursor=base64EncodedLastID
```

**Response includes:**
```json
{
  "transactions": [...],
  "pagination": {
    "next_cursor": "base64EncodedLastID_of_next_page",
    "has_more": true
  }
}
```

If `has_more=false`, client knows end of list.

---

## 9. Filtering

List endpoints support common filters:

- `status` — comma-separated list: `?status=Success,Pending`
- `start_date`, `end_date` — ISO 8601 timestamps
- `category_id` — UUID (for products, transactions filtered by product category)

Server constructs SQL: `WHERE (status IN (...)) AND (created_at BETWEEN ...)`.

---

## 10. Rate Limiting Headers

Responses include rate limit headers where applicable:

```
X-RateLimit-Limit: 30
X-RateLimit-Remaining: 29
X-RateLimit-Reset: 1704558000  // epoch seconds when window resets
```

On limit exceeded:
```
HTTP/1.1 429 Too Many Requests
Retry-After: 60
```

---

## 11. OpenAPI / Swagger

Full OpenAPI 3.0 specification to be generated separately in `openapi.yaml`. This document provides human-readable summary.

**Spec features:**
- All endpoints with request/response schemas
- Auth: `Bearer` JWT
- Error schemas
- Example calls

---

## Appendix A — Common Response Codes

| Code | Meaning | Retry? |
|---|---|---|
| 200 | Success | No |
| 201 | Created | No |
| 202 | Accepted (pending) | No (poll) |
| 204 | No Content (deleted) | No |
| 400 | Bad Request (validation) | No (fix input) |
| 401 | Unauthorized | No (re-login) |
| 403 | Forbidden (authz) | No |
| 404 | Not Found | No |
| 409 | Conflict (duplicate) | No (new id) |
| 429 | Too Many Requests | Yes (after Retry-After) |
| 500 | Internal Server Error | Yes (circuit breaker) |
| 502 | Bad Gateway (Digiflazz error) | Yes (retry) |
| 503 | Service Unavailable | Yes (after delay) |

---

**Owner:** API Team  
**Maintained by:** Backend Tech Lead  
**Related:** `ERROR_HANDLING.md` for full error code list, `API_AUTHENTICATION_FLOW.md` for auth details
