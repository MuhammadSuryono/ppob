# 🛠️ Internal API Reference (v1)

## 1. Global Specifications

### 1.1 Base URL
All requests are routed through the API Gateway:
`https://fedora.sinauplatform.id/api/v1/{service}/{endpoint}`

### 1.2 Mandatory Headers
| Header | Required | Description |
|---|---|---|
| `Authorization` | Yes | `Bearer <JWT_ACCESS_TOKEN>` |
| `Idempotency-Key` | Writes only | Unique UUIDv4 for duplicate prevention. |
| `X-Trace-ID` | No | Distributed tracing ID. |
| `Content-Type` | Yes | `application/json` |

## 2. Service Catalog

### 2.1 Auth Service (`/auth`)
- `POST /register`: New user registration.
- `POST /login`: Primary login (adaptive flow).
- `POST /verify-otp`: Validates 6-digit SMS/WA code.
- `POST /verify-credential`: Validates PIN or Password for existing users.
- `POST /verify-pin`: PIN login for trusted devices.
- `POST /refresh`: Issues a new JWT using a refresh token.

### 2.2 User Service (`/user`)
- `GET /users/:id`: Detailed profile fetch.
- `GET /staff`: (Mitra) List assigned staff with stats.
- `POST /staff`: (Mitra) Create a new staff organization entry.
- `GET /notifications`: User notification history.

### 2.3 Wallet Service (`/wallet`)
- `GET /me/balance`: Current spendable and held balance.
- `POST /staff/topup`: (Mitra) Fund a staff member's wallet.
- `GET /me/events`: Immutable transaction log for the wallet.

### 2.4 Transaction Service (`/transaction`)
- `POST /initiate`: Start a purchase (Prepaid or Postpaid Inquiry).
- `POST /confirm`: (Postpaid) Execute final payment.
- `GET /history`: Paginated transaction records.
- `GET /reports`: (Mitra) Aggregated sales and profit analytics.

### 2.5 Product Service (`/product`)
- `GET /products`: List products with category/SKU filters.
- `GET /categories`: UI-optimized category list.

## 3. Response Standardization

### 3.1 Success Envelope
```json
{
  "data": { ... },
  "metadata": { "page": 1, "has_more": true }
}
```

### 3.2 Error Envelope
```json
{
  "error": {
    "code": "TRANSACTION_INSUFFICIENT_BALANCE",
    "message": "Saldo tidak mencukupi",
    "trace_id": "abc-123-xyz",
    "timestamp": "2026-05-23T00:00:00Z"
  }
}
```
