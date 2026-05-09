# Error Handling Specification for PPOB

**Audience:** Backend developers, mobile app developers, QA  
**Last Updated:** 2026-05-07  
**Status:** Draft — Indonesian translations to be refined by native speaker

---

## 1. Overview

This document defines a unified error handling strategy across all microservices, including error taxonomy, standardized error response format, error code mapping (especially Digiflazz RC codes), retry policies, and user-facing messages in Indonesian.

**Goals:**
- Consistent error experience across all endpoints
- Actionable error messages for end-users
- Machine-parsable error codes for mobile app logic
- Supportability: every error includes `trace_id` for debugging

---

## 2. Error Taxonomy

Errors are classified into five top-level categories:

| Category | Prefix | Description | Examples |
|---|---|---|---|
| `AUTH_` | Authentication/Authorization | Login, token, permission failures | `AUTH_INVALID_CREDENTIALS`, `AUTH_TOKEN_EXPIRED` |
| `TRANSACTION_` | Transaction Processing | Wallet, transaction, payment errors | `TRANSACTION_INSUFFICIENT_BALANCE`, `TRANSACTION_DAILY_LIMIT_EXCEEDED` |
| `DIGIFLAZZ_` | Upstream Provider | Errors from Digiflazz API (mapped by RC) | `DIGIFLAZZ_TIMEOUT`, `DIGIFLAZZ_OUT_OF_STOCK` |
| `VALIDATION_` | Input Validation | Bad request data, schema violations | `VALIDATION_PHONE_FORMAT`, `VALIDATION_PIN_FORMAT` |
| `SYSTEM_` | Internal System | Database, network, internal errors | `SYSTEM_DB_UNAVAILABLE`, `SYSTEM_REDIS_DOWN` |

---

## 3. Standardized Error Response Format

All API error responses follow this JSON structure:

```json
{
  "error": {
    "code": "TRANSACTION_INSUFFICIENT_BALANCE",
    "message": "Saldo tidak mencukupi",
    "details": {
      "current_balance": 25000,
      "required": 50000,
      "wallet_id": "550e8400-e29b-41d4-a716-446655440000"
    },
    "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
    "timestamp": "2026-05-07T00:30:15.123Z"
  }
}
```

**Fields:**
- `code` — Machine-parsable error code (uppercase_snake_case)
- `message` — User-friendly message in Indonesian (for mobile display)
- `details` — Optional; additional context (key-value pairs, never include secrets)
- `trace_id` — OpenTelemetry trace ID for support/debugging
- `timestamp` — ISO 8601 UTC

---

## 4. Error Code Catalog

### 4.1 AUTH Errors

| Code | HTTP Status | Message (Indonesian) | When to Use | Retryable? |
|---|---|---|---|---|
| `AUTH_INVALID_CREDENTIALS` | 401 | "Nomor HP atau PIN/Password salah" | Login failed after OTP verification | No |
| `AUTH_TOKEN_EXPIRED` | 401 | "Sesi berakhir, silakan login ulang" | Access token expired | No (must re-login) |
| `AUTH_TOKEN_REVOKED` | 401 | "Token tidak valid" | Refresh token used after logout | No |
| `AUTH_OTP_EXPIRED` | 400 | "Kode OTP sudah kadaluarsa" | OTP verification after 5min | No (resend OTP) |
| `AUTH_OTP_INVALID` | 400 | "Kode OTP tidak valid" | Wrong OTP | No |
| `AUTH_OTP_RATE_LIMIT` | 429 | "Terlalu banyak percobaan OTP, coba lagi dalam 1 menit" | OTP brute-force protection | Yes (after wait) |
| `AUTH_ACCOUNT_LOCKED` | 403 | "Akun diblokir karena terlalu banyak percobaan PIN" | PIN locked after 5 Wrong tries | No (require password+OTP) |
| `AUTH_DEVICE_NOT_TRUSTED` | 403 | "Perangkat tidak terpercaya, silakan Otentikasi ulang" | Device trust score too low | No |
| `AUTH_INSUFFICIENT_PERMISSION` | 403 | "Anda tidak memiliki izin untuk mengakses sumber daya ini" | RBAC failure | No |

### 4.2 TRANSACTION Errors

| Code | HTTP Status | Message (Indonesian) | When to Use | Retryable? |
|---|---|---|---|---|
| `TRANSACTION_INSUFFICIENT_BALANCE` | 400 | "Saldo tidak mencukupi" | Wallet balance too low | No (add funds) |
| `TRANSACTION_DAILY_LIMIT_EXCEEDED` | 403 | "Limit harian transaksi tercapai" | Staff daily count/amount limit hit | No (reset midnight) |
| `TRANSACTION_PRODUCT_INACTIVE` | 400 | "Produk tidak aktif" | Product is_active=false | No |
| `TRANSACTION_PRICE_BELOW_COST` | 400 | "Harga jual minimal Rp {platform_price}" | selling_price < platform_price | No |
| `TRANSACTION_HOLD_FAILED` | 500 | "Gagal mengunci saldo" | Race condition on hold | Yes (retry) |
| `TRANSACTION_ALREADY_PROCESSING` | 409 | "Transaksi sedang diproses" | Duplicate idempotency key | No |
| `TRANSACTION_CANCELLED` | 400 | "Transaksi dibatalkan" | User cancelled | No |
| `TRANSACTION_NOT_FOUND` | 404 | "Transaksi tidak ditemukan" | Invalid transaction_id | No |
| `TRANSACTION_EXPIRED` | 410 | "Transaksi sudah kadaluarsa" | Pending >15min | No (resubmit) |
| `TRANSACTION_INQUIRY_EXPIRED` | 400 | "Pengecekan tagihan kadaluarsa, silakan cek ulang" | Postpaid inquiry >24h old | No |

### 4.3 DIGIFLAZZ Errors (Mapped from RC)

| Code | RC | HTTP Status | Message (Indonesian) | Retry? | Backoff |
|---|---|---|---|---|---|
| `DIGIFLAZZ_SUCCESS` | 00 | 200 | "Transaksi berhasil" | No | — |
| `DIGIFLAZZ_TIMEOUT` | 01 | 502 | "Timeout dari provider, mencoba ulang..." | Yes (3x) | 2s,4s,8s |
| `DIGIFLAZZ_PENDING` | 03 | 202 | "Transaksi sedang diproses" | No (poll) | — |
| `DIGIFLAZZ_GENERAL_FAILURE` | 02 | 502 | "Transaksi gagal" | No | — |
| `DIGIFLAZZ_PAYLOAD_ERROR` | 40 | 400 | "Format request tidak valid" | No | — |
| `DIGIFLAZZ_INVALID_SIGNATURE` | 41 | 401 | "Signature tidak valid" | No (config fix) | — |
| `DIGIFLAZZ_SELLER_NOT_FOUND` | 42 | 502 | "Gagal memproses API seller" | No | — |
| `DIGIFLAZZ_SKU_NOT_FOUND` | 43 | 404 | "Produk tidak ditemukan" | No | — |
| `DIGIFLAZZ_INSUFFICIENT_DEPOSIT` | 44 | 500 | "Saldo platform tidak cukup — hubungi admin" | No | **ALERT ADMIN** |
| `DIGIFLAZZ_IP_BLOCKED` | 45 | 403 | "IP tidak dikenali" | No (whitelist fix) | — |
| `DIGIFLAZZ_REF_ID_DUPLICATE` | 49 | 409 | "ID transaksi sudah digunakan" | No (new ref_id) | — |
| `DIGIFLAZZ_NUMBER_BLOCKED` | 51 | 400 | "Nomor tujuan diblokir" | No | — |
| `DIGIFLAZZ_PREFIX_MISMATCH` | 52 | 400 | "Prefix nomor tidak sesuai operator" | No | — |
| `DIGIFLAZZ_PRODUCT_UNAVAILABLE` | 53 | 503 | "Produk tidak tersedia saat ini" | No | — |
| `DIGIFLAZZ_INVALID_NUMBER` | 54 | 400 | "Nomor tujuan salah" | No | — |
| `DIGIFLAZZ_PRODUCT_DISRUPTION` | 55 | 503 | "Produk sedang gangguan" | Yes (after delay) | 5s retry |
| `DIGIFLAZZ_CUT_OFF` | 58 | 503 | "Transaksi cut-off, coba lagi dalam 15 menit" | Yes (after wait) | 900s |
| `DIGIFLAZZ_BILL_NOT_AVAILABLE` | 60 | 404 | "Tagihan belum tersedia" | No (try later) | — |
| `DIGIFLAZZ_SELLER_NOT_VERIFIED` | 67 | 403 | "Seller belum ter-verifikasi" | No | — |
| `DIGIFLAZZ_OUT_OF_STOCK` | 68 | 503 | "Stok habis" | No | — |
| `DIGIFLAZZ_PRICE_MISMATCH` | 69 | 400 | "Harga telah berubah, silakan sinkronisasi ulang" | Yes (retry sync) | 10s |
| `DIGIFLAZZ_BILLER_TIMEOUT` | 70 | 502 | "Timeout dari biller, mencoba ulang..." | Yes (3x) | 2s,4s,8s |
| `DIGIFLAZZ_PRODUCT_UNSTABLE` | 71 | 502 | "Produk tidak stabil" | Yes (max 2) | 5s,10s |
| `DIGIFLAZZ_PRICELIST_LIMIT` | 83 | 429 | "Limit API pencarian harga, coba dalam 4 menit" | Yes (after wait) | 240s |
| `DIGIFLAZZ_RATE_LIMIT` | 85 | 429 | "Terlalu banyak permintaan, coba 1 menit lagi" | Yes (once) | 60s |
| `DIGIFLAZZ_PLN_INQUIRY_LIMIT` | 86 | 429 | "Limit cek PLN, coba lagi nanti" | Yes (after wait) | 60s |
| `DIGIFLAZZ_EMONEY_MULTIPLE` | 87 | 400 | "Nominal e-money harus kelipatan Rp 1.000" | No | — |
| `DIGIFLAZZ_ACCOUNT_BLOCKED` | 80 | 403 | "Akun diblokir oleh operator" | No | — |

**Implicit mapping:** Any other RC (88, 99, etc.) → `DIGIFLAZZ_UNKNOWN_ERROR`.

### 4.4 VALIDATION Errors

| Code | HTTP Status | Message (Indonesian) | When to Use |
|---|---|---|---|
| `VALIDATION_PHONE_FORMAT` | 400 | "Format nomor HP tidak valid (contoh: +6281234567890)" | Phone regex fail |
| `VALIDATION_PIN_FORMAT` | 400 | "PIN harus 6 digit angka" | PIN not 6 digits |
| `VALIDATION_PIN_SEQUENTIAL` | 400 | "PIN tidak boleh berurutan (123456, 654321)" | Weak PIN |
| `VALIDATION_PASSWORD_WEAK` | 400 | "Password minimal 8 karakter, mengandung huruf besar & angka" | Password fails complexity |
| `VALIDATION_CUSTOMER_NO_FORMAT` | 400 | "Format nomor tujuan tidak sesuai" | Customer number fails product-specific regex |
| `VALIDATION_AMOUNT_NEGATIVE` | 400 | "Jumlah tidak boleh negatif" | Amount < 0 |
| `VALIDATION_MISSING_FIELD` | 400 | "Field '{field}' wajib diisi" | Required field missing |
| `VALIDATION_JSON_INVALID` | 400 | "Format JSON tidak valid" | Parsing failure |

### 4.5 SYSTEM Errors

| Code | HTTP Status | Message (Indonesian) | When to Use | Auto-Retry? |
|---|---|---|---|---|
| `SYSTEM_DB_UNAVAILABLE` | 503 | "Database tidak tersedia, coba beberapa saat lagi" | DB connection pool exhausted | Yes (circuit breaker) |
| `SYSTEM_REDIS_UNAVAILABLE` | 503 | "Cache tidak tersedia" | Redis down | Yes |
| `SYSTEM_VAULT_UNAVAILABLE` | 503 | "Konfigurasi tidak dapat diakses" | Vault unreachable | Yes |
| `SYSTEM_INTERNAL` | 500 | "Terjadi kesalahan internal" | Unhandled panic/exception | No |
| `SYSTEM_TIMEOUT` | 504 | "Request timeout" | Request took > configured timeout | Yes (client may retry) |
| `SYSTEM_RATE_LIMIT` | 429 | "Server sedang sibuk, coba lagi nanti" | Rate limit internal (not endpoint-specific) | Yes (after Retry-After) |
| `SYSTEM_CIRCUIT_OPEN` | 503 | "Layanan sedang gangguan" | Circuit breaker open for downstream | Yes (after cooldown) |
| `SYSTEM_MAINTENANCE` | 503 | "Sistem dalam pemeliharaan" | Feature flag maintenance mode | No |

---

## 5. Retry Policy Matrix

**Client (Mobile App) Retry Logic:**

| Error Code | HTTP Status | Retry? | Max Attempts | Backoff Strategy | UI Behavior |
|---|---|---|---|---|---|
| `DIGIFLAZZ_TIMEOUT` | 502 | Yes | 3 | 2s → 4s → 8s | Show "Timeout, retrying..." toast |
| `DIGIFLAZZ_BILLER_TIMEOUT` | 502 | Yes | 3 | 2s → 4s → 8s | Same |
| `DIGIFLAZZ_PRODUCT_UNSTABLE` | 502 | Yes | 2 | 5s → 10s | Show "produk tidak stabil" |
| `DIGIFLAZZ_RATE_LIMIT` | 429 | Yes | 1 | Wait `Retry-After` header (or 60s) | Disable button, show timer |
| `DIGIFLAZZ_PRICELIST_LIMIT` | 429 | Yes | 1 | Wait 240s | Show "limit API, coba 4 menit" |
| `DIGIFLAZZ_CUT_OFF` | 503 | Yes | Infinite | Wait 900s (15min) | Show "cut-off, coba 15 menit" |
| `SYSTEM_DB_UNAVAILABLE` | 503 | Yes | 3 | 1s → 2s → 4s | Show "Server sibuk" |
| `SYSTEM_RATE_LIMIT` | 429 | Yes | 1 | Wait `Retry-After` (if present) else 30s | Show "Rate limit" |
| `SYSTEM_CIRCUIT_OPEN` | 503 | Yes (after 10s) | Infinite | 10s cooldown via circuit breaker | Show "Layanan gangguan" |
| `VALIDATION_*` | 400 | No | 0 | — | Show field errors inline |
| `AUTH_*` (except `OTP_RATE_LIMIT`) | 401/403 | No | 0 | — | Show specific message |
| `TRANSACTION_INSUFFICIENT_BALANCE` | 400 | No | 0 | — | Show balance, top-up prompt |
| `TRANSACTION_DAILY_LIMIT_EXCEEDED` | 403 | No | 0 | — | Show "limit harian" |
| `DIGIFLAZZ_SUCCESS` | 200 | No | — | — | Proceed |

**Implementation in Flutter (Dio interceptor):**
```dart
dio.interceptors.add(InterceptorsWrapper(
  onError: (error) async {
    final code = error.response?.data['error']['code'];
    final retryConfig = retryPolicy[code];
    
    if (retryConfig != null && retryConfig.shouldRetry) {
      if (retryConfig.maxAttempts > 0) {
        await Future.delayed(retryConfig.backoff);
        return dio.request(error.requestOptions);
      }
    }
    
    // Transform to user-friendly message
    final userMessage = error.response?.data['error']['message'] ?? 'Terjadi kesalahan';
    showToast(userMessage);
    return error;
  },
));
```

---

## 6. Error Mapping: Digiflazz RC → Internal Code

**Mapping Table** (used in Integration Service):

```go
var rcMapping = map[string]ErrorCode{
    // Success
    "00": "", // no error, success
    
    // Retryable
    "01": "DIGIFLAZZ_TIMEOUT",
    "70": "DIGIFLAZZ_BILLER_TIMEOUT",
    "71": "DIGIFLAZZ_PRODUCT_UNSTABLE",
    "99": "DIGIFLAZZ_ROUTER_ISSUE",
    
    // Pending — not error, but not success either
    "03": "DIGIFLAZZ_PENDING",
    
    // Non-retryable errors
    "40": "DIGIFLAZZ_PAYLOAD_ERROR",
    "41": "DIGIFLAZZ_INVALID_SIGNATURE",
    "42": "DIGIFLAZZ_SELLER_NOT_FOUND",
    "43": "DIGIFLAZZ_SKU_NOT_FOUND",
    "44": "DIGIFLAZZ_INSUFFICIENT_DEPOSIT", // ALERT ADMIN
    "45": "DIGIFLAZZ_IP_BLOCKED",
    "47": "DIGIFLAZZ_DUPLICATE_BUYER",
    "49": "DIGIFLAZZ_REF_ID_DUPLICATE",
    "50": "DIGIFLAZZ_NOT_FOUND",
    "51": "DIGIFLAZZ_NUMBER_BLOCKED",
    "52": "DIGIFLAZZ_PREFIX_MISMATCH",
    "53": "DIGIFLAZZ_PRODUCT_UNAVAILABLE",
    "54": "DIGIFLAZZ_INVALID_NUMBER",
    "55": "DIGIFLAZZ_PRODUCT_DISRUPTION",
    "58": "DIGIFLAZZ_CUT_OFF",
    "60": "DIGIFLAZZ_BILL_NOT_AVAILABLE",
    "67": "DIGIFLAZZ_SELLER_NOT_VERIFIED",
    "68": "DIGIFLAZZ_OUT_OF_STOCK",
    "69": "DIGIFLAZZ_PRICE_MISMATCH",
    "80": "DIGIFLAZZ_ACCOUNT_BLOCKED",
    "81": "DIGIFLAZZ_SELLER_BLOCKED_BY_USER",
    "82": "DIGIFLAZZ_ACCOUNT_UNVERIFIED",
    "83": "DIGIFLAZZ_PRICELIST_LIMIT",
    "85": "DIGIFLAZZ_RATE_LIMIT",
    "86": "DIGIFLAZZ_PLN_INQUIRY_LIMIT",
    "87": "DIGIFLAZZ_EMONEY_MULTIPLE",
    "88": "DIGIFLAZZ_FORBIDDEN_ACTION",
}
```

**Function:**
```go
func MapRCtoErrorCode(rc string) (string, bool) {
    code, ok := rcMapping[rc]
    return code, ok
}
```

If RC not in map → return `DIGIFLAZZ_UNKNOWN_ERROR` with raw RC in details.

---

## 7. Error Response Examples

### Example 1 — Validation Error (Missing Field)
```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": {
    "code": "VALIDATION_MISSING_FIELD",
    "message": "Field 'phone_number' wajib diisi",
    "details": {
      "missing_field": "phone_number"
    },
    "trace_id": "abc123def456",
    "timestamp": "2026-05-07T00:30:15.123Z"
  }
}
```

### Example 2 — Business Rule Violation (Daily Limit)
```http
HTTP/1.1 403 Forbidden
Content-Type: application/json

{
  "error": {
    "code": "TRANSACTION_DAILY_LIMIT_EXCEEDED",
    "message": "Limit harian transaksi tercapai. Batas: 50 transaksi / Rp 5.000.000 per hari",
    "details": {
      "limit_type": "count",
      "used_today": 50,
      "reset_at": "2026-05-08T00:00:00+07:00"
    },
    "trace_id": "xyz789",
    "timestamp": "2026-05-07T14:22:10.456Z"
  }
}
```

### Example 3 — Upstream Failure (Digiflazz Timeout)
```http
HTTP/1.1 502 Bad Gateway
Content-Type: application/json
Retry-After: 2

{
  "error": {
    "code": "DIGIFLAZZ_TIMEOUT",
    "message": "Timeout dari provider, mencoba ulang...",
    "details": {
      "rc": "01",
      "retry_attempt": 1,
      "max_retries": 3
    },
    "trace_id": "timeout123",
    "timestamp": "2026-05-07T00:31:45.789Z"
  }
}
```

### Example 4 — Internal System Error
```http
HTTP/1.1 500 Internal Server Error
Content-Type: application/json

{
  "error": {
    "code": "SYSTEM_INTERNAL",
    "message": "Terjadi kesalahan internal, tim kami sudah diberitahu",
    "details": {
      "error_id": "ERR-20260507-001234"  -- internal ticket reference
    },
    "trace_id": "internal456",
    "timestamp": "2026-05-07T00:32:00.000Z"
  }
}
```
**Never** expose stack traces or internal SQL errors to client in production.

---

## 8. Logging Errors

**Every error must be logged with full context:**
```go
log.Error("transaction failed",
    "error_code", errCode,
    "user_id", userID,
    "transaction_id", txID,
    "details", err.Details,
    "trace_id", traceID,
    "duration_ms", duration,
    "request_id", requestID,
)
```

**Structured log (JSON):**
```json
{
  "timestamp": "...",
  "level": "error",
  "service": "transaction-service",
  "action": "initiate_transaction",
  "error_code": "TRANSACTION_INSUFFICIENT_BALANCE",
  "user_id": "550e8400...",
  "transaction_id": null,
  "details": {"current_balance":25000,"required":50000},
  "trace_id": "abc123",
  "version": 1
}
```

**Error Aggregation:** Use Grafana Loki or similar to aggregate errors by code for alerting.

---

## 9. Client-Side Error Handling (Flutter)

**Error Display Patterns:**

| Severity | Display Method | Example |
|---|---|---|
| **Blocking** (cannot proceed) | Modal dialog with single "OK" button | Insufficient balance, daily limit |
| **Transient** (auto-retrying) | Snackbar with spinner + "Retrying..." text | Digiflazz timeout, network error |
| **Validation** (inline) | Text below input field in red | Invalid phone format |
| **Network** (no connectivity) | Full-screen retry UI with button | No internet connection |

**Error Translation:** Mobile app ships Indonesian string resources; all error messages from server already in Indonesian; no translation needed.

**Retry UI:**
```dart
if (error.code == 'DIGIFLAZZ_TIMEOUT' && retryCount < 3) {
  showSnackBar('Timeout, mencoba ulang... (${retryCount+1}/3)');
  await Future.delay(backoff);
  // retry request
} else {
  showErrorDialog(error.message);
}
```

---

## 10. Support & Traceability

**Each error includes `trace_id`.** Support workflow:
1. User reports issue with error screen showing trace ID (copy button)
2. Support asks for trace ID
3. Support looks up trace ID in Grafana Loki:
   ```json
   { trace_id="abc123" } | json
   ```
4. See full request/response logs across all services
5. Identify root cause (which service, which line)

**Enhancement:** Link trace_id to Sentry issue (Sentry groups by fingerprint but allows search by `trace_id`).

---

## 11. Error Budget & SLO Impact

Errors contribute to Service Level Objective (SLO) calculations:

- **Success criterion:** HTTP 2xx for success, 4xx for client error (not counted as service failure), 5xx counts as error
- **Error budget:** (1 − SLO) = 0.5% error budget for transaction service (99.5% success)
- **Burn rate:** If error budget consumed >50% in 28 days, prioritize bug fixes over features

**4xx vs 5xx:**
- `VALIDATION_*`, `AUTH_*`, `TRANSACTION_*` (business rule violations) → 4xx (client's fault)
- `SYSTEM_*`, `DIGIFLAZZ_*` (except 4xx mapped like RC 44/45) → 5xx (our system/providers fault)

---

## 12. Monitoring Error Rates

**Prometheus Alerts:**

```yaml
- alert: High5xxErrorRate
  expr: |
    sum(rate(http_requests_total{status=~"5.."}[5m])) 
    / 
    sum(rate(http_requests_total[5m])) > 0.005  -- 0.5%
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "High 5xx error rate detected ({{ $value | humanizePercentage }})"
    runbook_url: "https://wiki/runbooks/5xx-error-rate"

- alert: DigiflazzErrorRateHigh
  expr: |
    sum(rate(digiflazz_api_calls_total{status="error"}[5m])) 
    / 
    sum(rate(digiflazz_api_calls_total[5m])) > 0.10
  for: 3m
  labels:
    severity: critical
  annotations:
    summary: "Digiflazz error rate >10% ({{ $value | humanizePercentage }})"
```

**Dashboard Panel:**
- Top 10 error codes by frequency (bar chart)
- Error rate over time (line graph)
- Breakdown by service (pie chart)

---

## 13. Error Code Evolution & Versioning

**Error codes are immutable once released.** If message text needs change, keep code same and update message in translation file.

**Deprecation:** Mark codes as deprecated in `ERROR_HANDLING.md` after 6 months; avoid reuse.

**Versioning:** Add `api_version` field if breaking changes to error format needed. Initial version `1`.

---

## 14. Localization (i18n)

**Server sends Indonesian messages.** Future: support English locale via `Accept-Language` header.

```go
messages := map[string]map[string]string{
    "id": {
        "VALIDATION_PHONE_FORMAT": "Format nomor HP tidak valid",
        // ...
    },
    "en": {
        "VALIDATION_PHONE_FORMAT": "Invalid phone number format",
        // ...
    }
}
```

Mobile app currently only supports Indonesian (per PRD). Keep i18n enabled on server for future expansion.

---

## Appendix A — Go Error Wrapper

```go
type AppError struct {
    Code         string                 `json:"code"`
    Message      string                 `json:"message"`
    Details      map[string]interface{} `json:"details,omitempty"`
    TraceID      string                 `json:"trace_id"`
    StatusCode   int                    `json:"-"`
    Retryable    bool                   `json:"-"`
    RetryAfter   int                    `json:"-"` // seconds
}

func (e *AppError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewError(code, message string, details map[string]interface{}) *AppError {
    return &AppError{
        Code:      code,
        Message:   message,
        Details:   details,
        TraceID:   generateTraceID(),
        StatusCode: httpStatusForCode(code),
    }
}

func httpStatusForCode(code string) int {
    switch {
    case strings.HasPrefix(code, "AUTH_"):
        return http.StatusUnauthorized
    case strings.HasPrefix(code, "VALIDATION_"):
        return http.StatusBadRequest
    case strings.HasPrefix(code, "TRANSACTION_") || strings.HasPrefix(code, "DIGIFLAZZ_"):
        return http.StatusBadRequest // or 403/502 as appropriate; refine per code
    case strings.HasPrefix(code, "SYSTEM_"):
        return http.StatusInternalServerError
    default:
        return http.StatusInternalServerError
    }
}
```

---

## Appendix B — Error Code to HTTP Status Mapping

| Error Category | Default Status | Overrides |
|---|---|---|
| `AUTH_*` | 401 Unauthorized | `AUTH_ACCOUNT_LOCKED` → 403 |
| `VALIDATION_*` | 400 Bad Request | — |
| `TRANSACTION_*` | 400 Bad Request | `TRANSACTION_DAILY_LIMIT_EXCEEDED` → 403 |
| `DIGIFLAZZ_*` | 502 Bad Gateway | `DIGIFLAZZ_RATE_LIMIT` → 429, `DIGIFLAZZ_PENDING` → 202 |
| `SYSTEM_*` | 500 Internal Server Error | `SYSTEM_RATE_LIMIT` → 429, `SYSTEM_CIRCUIT_OPEN` → 503 |

---

## Appendix C — Retry-After Header

When returning 429 (rate limit) or 503 (circuit open), include `Retry-After` header:

```go
w.Header().Set("Retry-After", "60")  -- seconds
```

Client should respect and wait before retrying.

---

## Appendix D — Error Message UI Copy Guidelines

- **Be polite but concise:** "Saldo tidak cukup" not "Maaf, saldo Anda tidak mencukupi untuk melakukan transaksi ini"
- **Be actionable:** Include what user should do next: "Top up saldo di menu Wallet"
- **Include numbers:** Show expected vs actual: "Minimal Rp 25.000"
- **Avoid technical terms:** Say "provider" not "Digiflazz" in user messages; use "sistem" or "layanan"
- **Capitalization:** Sentence case, not ALL CAPS

---

## Open Questions

1. **Should we expose `Retry-After` for all retryable errors?** Yes, applicable to 429 and 503 responses. For 502 transient, we can include `Retry-After: 2` as suggestion.

2. **Should error codes be exposed to mobile?** Yes, so app can handle programmatically (e.g., show top-up prompt on `INSUFFICIENT_BALANCE` vs show "try again" on `DIGIFLAZZ_TIMEOUT`).

3. **Do we need localized messages (Bahasa Indonesia only)?** PRD says Indonesian UI; all messages in Indonesian. Future English possible via Accept-Language.

4. **How to handle errors from wallet concurrency deadlocks?** Should be retried automatically server-side (up to 3 times); if still failing, return `SYSTEM_INTERNAL` with trace ID and alert.

---

**Owner:** Backend Team Lead  
**Enforcement:** PR must include updated error codes if new ones added; update this document with every new error type introduced.

---

**Related Docs:**
- `DIGIFLAZZ_INTEGRATION_GUIDE.md` (RC → error code mapping details)
- `SECURITY_ARCHITECTURE.md` (error logging & PII redaction)
- `API_AUTHENTICATION_FLOW.md` (auth-specific errors)
