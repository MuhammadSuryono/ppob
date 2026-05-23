# 🚨 Error Code Reference

## 1. Error Taxonomy
The PPOB system uses a prefixed error code system to categorize failures across microservices.

| Prefix | Domain | Purpose |
|---|---|---|
| `AUTH_` | Identity | Authentication and Authorization failures. |
| `VALIDATION_` | Schema | Input data format and business constraint violations. |
| `TRANSACTION_` | Business | Operation-level failures (e.g., balance, limits). |
| `DIGIFLAZZ_` | Upstream | Direct mapping from the primary provider. |
| `SYSTEM_` | Internal | Infrastructure, DB, or unhandled failures. |

## 2. Global Error Catalog

### 2.1 Identity Errors (`AUTH_`)
- `AUTH_INVALID_CREDENTIALS`: Wrong phone/password/PIN.
- `AUTH_TOKEN_EXPIRED`: JWT session has timed out.
- `AUTH_DEVICE_NOT_TRUSTED`: Action requires OTP (Trust Score too low).
- `AUTH_ACCOUNT_LOCKED`: 5+ failed attempts; locked for 1 hour.

### 2.2 Business Errors (`TRANSACTION_`)
- `TRANSACTION_INSUFFICIENT_BALANCE`: Wallet available funds < Selling Price.
- `TRANSACTION_DAILY_LIMIT_EXCEEDED`: Staff turnover or count limit reached.
- `TRANSACTION_EXPIRED`: Order Pending for > 15 minutes.
- `TRANSACTION_ALREADY_PROCESSING`: Idempotency key conflict.

### 2.3 Upstream Errors (`DIGIFLAZZ_`)
- `DIGIFLAZZ_OUT_OF_STOCK` (RC 68): Product currently unavailable.
- `DIGIFLAZZ_TIMEOUT` (RC 01/70): Provider took too long (Auto-retried).
- `DIGIFLAZZ_INSUFFICIENT_DEPOSIT` (RC 44): Platform-level funding error.

## 3. Implementation Policy
- **Tracing:** Every error returned to a client MUST include a `trace_id` for support lookup.
- **Messages:** Internal technical errors (e.g., SQL syntax) MUST NOT be returned in the `message` field.
- **Formatting:** Codes are always `UPPERCASE_SNAKE_CASE`.
