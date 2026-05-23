# 🔌 Digiflazz API Technical Reference

## 1. Authentication & Security

### 1.1 Request Signature
Digiflazz requires an MD5 signature for most API requests.
- **Formula:** `sign = md5(username + apiKey + command)`
- **`command` mapping:**
    - Cek Saldo: `"depo"`
    - Price List: `"pricelist"`
    - Prepaid Transaction: `<ref_id>`
    - Postpaid Inquiry: `"inq-pasca" + <ref_id>`
    - Postpaid Payment: `"pay-pasca" + <ref_id>`
    - Status Check (Postpaid): `"status-pasca" + <ref_id>`

### 1.2 Webhook Signature
Inbound webhooks use HMAC-SHA1 verification.
- **Header:** `X-Hub-Signature: sha1=<hmac_hex>`
- **Validation:** `hmac_sha1(raw_body, webhook_secret)`

## 2. Core Endpoints

### 2.1 Balance Check
- **Endpoint:** `POST https://api.digiflazz.com/v1/cek-saldo`
- **Request:** `{ "cmd": "deposit", "username": "...", "sign": "..." }`
- **Response:** `{ "data": { "deposit": 500000 } }`

### 2.2 Price List (Prepaid/Postpaid)
- **Endpoint:** `POST https://api.digiflazz.com/v1/price-list`
- **Rate Limit:** 1 request per 5 minutes per type.
- **Request:** `{ "cmd": "prepaid", "username": "...", "sign": "..." }`

### 2.3 Transaction (Prepaid)
- **Endpoint:** `POST https://api.digiflazz.com/v1/transaction`
- **Response Codes:**
    - `00`: Success (Terminal)
    - `03`: Pending (Wait for Webhook)
    - `01`, `70`: Timeout (Retryable)

### 2.4 Transaction (Postpaid)
- **Two-Step Flow:**
    1.  **Inquiry:** `commands: "inq-pasca"`. Returns bill details.
    2.  **Payment:** `commands: "pay-pasca"`. Executes final payment using the same `ref_id`.

## 3. Webhook Payloads

### 3.1 Delivery Headers
| Header | Value |
|---|---|
| `X-Digiflazz-Event` | `create` or `update` |
| `User-Agent` | `Digiflazz-Hookshot` (Prepaid) or `Digiflazz-Pasca-Hookshot` (Postpaid) |

### 3.2 Status Update Example
```json
{
  "data": {
    "ref_id": "unique-123",
    "status": "Sukses",
    "rc": "00",
    "sn": "REFERENCE-NUMBER-123",
    "buyer_last_saldo": 250000
  }
}
```

## 4. Response Code (RC) Mapping
| RC | Internal Status | Action |
|---|---|---|
| `00` | Success | terminal |
| `01` | Initiated | retry (3x) |
| `03` | Pending | wait webhook |
| `44` | Failed | ALERT ADMIN (insufficient platform deposit) |
| `68` | Failed | Show "Out of Stock" |
| `85` | Failed | Wait 60s, retry 1x |
| `99` | Pending | retry (3x) |
