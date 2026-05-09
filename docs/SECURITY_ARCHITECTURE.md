# 🔐 Security Architecture for PPOB Mobile Application

**Audience:** Backend developers, security reviewers, DevOps  
**Last Updated:** 2026-05-07  
**Status:** Draft — requires security review before implementation

---

## 1. Threat Model (STRIDE)

We apply STRIDE per element to identify threats and necessary mitigations.

| Element | Spoofing | Tampering | Repudiation | Information Disclosure | Denial of Service | Elevation of Privilege | Mitigations |
|---|---|---|---|---|---|---|---|
| Mobile App → Auth | ✅ Device fingerprint + JWT | ❌ Request body in transit | ✅ Audit log | ✅ TLS 1.3 | ✅ Rate limiting | ❌ — | TLS everywhere, HMAC signatures, mTLS internal |
| Auth Service | ✅ JWT signature | ✅ DB writes audited | ✅ Full audit trail | ✅ Secrets in Vault | ✅ Rate limiting | ✅ Role-based access control | JWT RS256, signed logs, RBAC |
| Wallet Service | ✅ Operator authentication | ✅ Pessimistic locking | ✅ Wallet event log | ✅ Field-level encryption for PII | ✅ Circuit breaker | ✅ Ownership checks | SELECT FOR UPDATE, event sourcing |
| Digiflazz API | ✅ HMAC signature | ✅ Signature verification | ✅ Webhook audit log | ✅ Secrets from Vault | ✅ Timeout + retry | ✅ — | MD5 request sig, HMAC-SHA1 webhook |

---

## 2. Authentication & Authorization

### 2.1 JWT Configuration

**Algorithm:** RS256 (asymmetric) — private key in service, public key for verification

**Keys:**
- Key ID: `kid` header indicates which key (supports key rotation)
- Key pair generation: 2048-bit RSA, rotated every 90 days
- Storage: Private key in HashiCorp Vault at `secret/jwt/private_key`, public key published via `/.well-known/jwks.json`

**Token Structure:**
```json
{
  "sub": "user-uuid",
  "role": "active_role_uuid",
  "roles": ["role_uuid_1", "role_uuid_2"],
  "wallet_id": "wallet_uuid",
  "device_id": "device_fingerprint_uuid",
  "iat": 1704556800,
  "exp": 1704571200,
  "jti": "jwt-id-unique"
}
```

**Refresh Token:**
- Type: Random 256-bit entropy, stored hashed (bcrypt) in Redis
- TTL: 7 days (configurable per role: Mitra 7d, Staff 1d)
- One-time use: consumed on refresh, new refresh token issued
- Revocation: added to Redis blocklist until expiry

**Access Token:**
- TTL: 15 minutes
- Rotation: New token issued on refresh; old token valid until expiry

### 2.2 Password & PIN Hashing

**Passwords (login):**
- Algorithm: bcrypt
- Cost factor: 12 (takes ~350ms on modern CPU)
- Salt: 16-byte random per user
- Minimum length: 8 characters
- Requirements: 1 uppercase, 1 lowercase, 1 digit

**PINs (transaction authorization):**
- Algorithm: Argon2id (reduced parameters for mobile UX)
  - Memory: 64 MB
  - Iterations: 2
  - Parallelism: 1
  - Salt: 16-byte random per PIN change
  - Hash length: 32 bytes
- Rationale: PIN is only 6 digits (1M entropy). Argon2id memory-hard slows GPU brute-force.
- Salt stored server-side in `users.pin_salt` column (new column needed)
- PIN change requires re-authentication (password + OTP)

### 2.3 Device Fingerprinting & Trust

**Fingerprint Signals (collect from mobile app):**
1. `device_id` — UUID generated on first launch, stored in secure storage
2. `user_agent` — Flutter version + OS + device model
3. `ip_subnet` — First 3 octets of IP (masked)
4. `app_version` — Semantic version string
5. `install_ts` — App installation timestamp (days since install)
6. `last_login_ts` — Previous successful login timestamp

**Trust Score Calculation (score 0–100):**
- Device seen before on this account: +30
- Same IP subnet as recent logins: +20
- App version current (±1 minor): +10
- Install age > 7 days: +10
- Last login < 24h ago: +20
- ≥3 previous successful logins on this device: +10

**Thresholds:**
- `score ≥ 70` → Trusted → Skip OTP (direct PIN login)
- `30 ≤ score < 70` → Semi-trusted → Password + OTP required
- `score < 30` → Untrusted → Full password + OTP + possible additional verification

**Storage Table:** `device_fingerprints`
```sql
CREATE TABLE device_fingerprints (
    device_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    fingerprint_hash VARCHAR(64),
    user_agent TEXT,
    ip_address INET,
    trust_score INT DEFAULT 0,
    is_trusted BOOLEAN DEFAULT FALSE,
    first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
CREATE INDEX idx_device_user ON device_fingerprints(user_id);
```

**Trust Updates:** On each login, update `last_seen` and recalculate `trust_score`. If score crosses threshold, update `is_trusted`.

### 2.4 Rate Limiting

**Implementation:** Token bucket algorithm using Redis (sliding window with fixed increments)

**Limits per Identity (phone_number / user_id / IP):**

| Endpoint | Key | Rate Limit | Burst | Purpose |
|---|---|---|---|---|
| POST `/auth/register` | `rate:auth:register:{phone}` | 3 / hour | 1 | Prevent OTP spamming |
| POST `/auth/verify-otp` | `rate:auth:verify:{phone}` | 5 / 15 min | 2 | Brute-force OTP protection |
| POST `/auth/login` | `rate:auth:login:{phone}` | 10 / 15 min | 3 | Password/PIN brute-force |
| POST `/wallets/topup-staff` | `rate:topup:{mitra_id}` | 30 / min | 5 | Mitra topup flood protection |
| POST `/transactions/initiate` | `rate:txn:{user_id}` | 30 / min | 5 | Transaction flooding |
| POST `/users/add-staff` | `rate:staffadd:{mitra_id}` | 20 / hour | 3 | Staff creation spam |

**Action on Exceed:** Return `429 Too Many Requests` with:
```json
{
  "error": "RATE_LIMIT_EXCEEDED",
  "message": "Too many requests, please try again later",
  "retry_after": 60
}
```

### 2.5 Session & Token Revocation

**Refresh Token Allowlist (Redis):**
```
key: refresh_token:{jti} → { user_id, device_id, expires_at }
TTL: 7 days (matches refresh token TTL)
```

**On Refresh:** Old token `jti` deleted from allowlist (consumed); new token added.  
**On Logout:** Access token added to blocklist (Redis set `blocked_jti` with TTL = remaining `exp`); refresh token deleted from allowlist.  
**On Privilege Escalation (Role Switch):** All existing refresh tokens invalidated (remove all `refresh_token:{user_id}:*` keys); force re-login.

### 2.6 PIN Security

- Max 5 wrong attempts → account locked for 1 hour
- After lock: require password + OTP to unlock and reset PIN
- PIN change requires current PIN OR password + OTP
- New PIN cannot be same as last 3 historical PINs (enforced via `pin_history` table if implemented)

---

## 3. Secrets Management

**Provider:** HashiCorp Vault (recommended) or AWS Secrets Manager

**Secrets Required:**

| Secret | Vault Path | Rotation | Purpose |
|---|---|---|---|
| Digiflazz API Key | `secret/digiflazz/api_key` | Quarterly | Request signature |
| Digiflazz Webhook Secret | `secret/digiflazz/webhook_secret` | Quarterly | HMAC verification |
| JWT Private Key | `secret/jwt/private_key` | 90 days | Token signing |
| Database Credentials | `secret/postgres/ppob_prod` | 30 days | DB connection |
| Redis Password | `secret/redis/password` | 30 days | Cache auth |
| Encryption Master Key | `secret/encryption/master_key` | Yearly | Hive encryption key |

**Access Control:** Services authenticate via K8S service account token; least-privilege ACLs per service.

**Rotation Strategy:**
- Automated via Vault for database passwords (30d)
- Manual rotation for Digiflazz credentials (coordinate with provider)
- JWT keys: Generate new keypair, update public key endpoint, gradually phase out old key by keeping it valid for 7 days while services rotate

---

## 4. Network Security

### 4.1 Transport Layer Security

- **External Traffic:** TLS 1.3 minimum, strong cipher suites only
- **Internal Traffic:** Mutual TLS (mTLS) via service mesh (Istio automatic)
- **Certificates:** Let's Encrypt (staging/prod) or internal CA for service mesh

### 4.2 Firewall & Egress

**Egress Rules:**
- Allow outbound to Digiflazz IP ranges (need their published list; else allow all but monitor)
- Allow outbound to Vault, S3, CloudWatch (if AWS)
- Deny all other external egress from app pods (except via NAT gateway)

**Ingress Rules:**
- Public: API Gateway / Ingress (Nginx) on ports 80/443 only
- Internal: Service-to-service communication via mesh; no NodePort/LoadBalancer for individual services

### 4.3 API Gateway

**Recommended:** Kong or AWS API Gateway

**Features:**
- Rate limiting (enforced at edge before reaching cluster)
- IP allowlisting for Digiflazz webhook source IPs (if they provide fixed ranges)
- Request size limits (max 1MB)
- CORS configured for mobile app origins (if hybrid app)

---

## 5. Data Protection

### 5.1 PII Handling

**Personally Identifiable Information (PII):**
- `phone_number` — stored as provided, but display masked: `+62 812-345-***-****`
- `name` — stored plain (needed for receipts)
- `customer_no` (in transactions) — stored plain (needed for inquiry), but masked in logs: last 4 digits visible

**Redaction in Logs:**
```go
func RedactSensitive(data map[string]interface{}) map[string]interface{} {
    redactKeys := []string{"password", "pin", "phone_number", "customer_no", "account_no"}
    for _, key := range redactKeys {
        if val, ok := data[key]; ok {
            if str, ok := val.(string); ok {
                data[key] = maskString(str) // e.g. "0812345***"
            }
        }
    }
    return data
}
```

### 5.2 Encryption at Rest

- PostgreSQL: Transparent Data Encryption (TDE) via AWS RDS (storage-level encryption with KMS)
- Additional column-level encryption for sensitive fields? Not required if DB encryption sufficient
- Backups encrypted with KMS key

### 5.3 Backup Security

- Backup files stored in private S3 bucket, no public access
- Encryption: SSE-KMS with separate key from DB
- Retention: 35 days, automatic deletion
- Access logs: S3 access logs and CloudTrail enabled

---

## 6. Input Validation & Injection Prevention

**All endpoints must:**
- Validate request body against JSON Schema (using `github.com/xeipuuv/gojsonschema`)
- Sanitize strings to prevent XSS if ever displayed in HTML admin panel (use text-only)
- Use parameterized queries everywhere (never string-concatenate SQL)
- For full-text search (if added later), use prepared statements

**Specific Checks:**

**Phone Number Format:**
```go
// Indonesian phone regex: +62 followed by 8-12 digits, no leading zero after +62
var phoneRegex = regexp.MustCompile(`^\+62[0-9]{8,12}$`)
```

**PIN Format:**
```go
// Exactly 6 digits, not sequential, not all same
var pinRegex = regexp.MustCompile(`^[0-9]{6}$`)
// Reject: 123456, 654321, 111111, 000000
```

**Customer Number Validation by Product:**
- Prepaid mobile: regex per operator stored in `product_validation_rules` table (e.g., Telkomsel: `^8[0-9]{8,10}$`)
- PLN: numeric, 8-12 digits
- PDAM: varies by region — stored in product metadata

---

## 7. OWASP Top 10 Mitigations

| OWASP Risk | Mitigation |
|---|---|
| **A01 — Broken Access Control** | Enforce ownership checks in every service; never trust client-provided IDs alone; always verify `user_id` owns resource before action |
| **A02 — Cryptographic Failures** | TLS everywhere; Argon2id/bcrypt for passwords; no MD5 except Digiflazz signature (external requirement) |
| **A03 — Injection** | Parameterized queries; ORM (GORM) with safe query building; input validation |
| **A04 — Insecure Design** | Threat modeling done; security considered from start; code reviews with security checklist |
| **A05 — Security Misconfiguration** | Infrastructure as Code (Terraform) ensures consistent config; CIS Docker/K8s benchmarks; secrets never in config files |
| **A06 — Vulnerable Components** | Dependabot/renovate for dependency updates; weekly security scans (trivy, grype) |
| **A07 — Identity & Auth Failures** | Strong hashing, rate limiting, device trust scoring, MFA future-ready (OTP) |
| **A08 — Software & Data Integrity** | CI/CD pipeline with signed artifacts; image signing (cosign); SBOM generation |
| **A09 — Security Logging & Monitoring** | Structured logs with correlation IDs; audit log immutable; alerts on anomalies |
| **A10 — Server-Side Request Forgery** | Outbound requests to Digiflazz only from Integration Service; URL allowlist; no user-controlled URLs |

---

## 8. Audit Logging

**Table:** `audit_logs` (already defined in database schema, expand with additional columns)

**Additional Columns:**
```sql
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS trace_id VARCHAR(32);
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS ip_address INET;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS user_agent TEXT;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS severity VARCHAR(20) DEFAULT 'INFO'; -- INFO, WARNING, CRITICAL
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS schema_version INT DEFAULT 1;
```

**Required Audit Events (non-exhaustive):**
- User login success/failure
- PIN change
- Role switch
- Wallet debit/credit (amount, balance_before, balance_after)
- Transaction initiation (product_id, customer_no masked, amount)
- Transaction status change (webhook received)
- Staff creation/modification
- Product price update (sync)
- Refund issuance
- Rate limit exceeded (potential attack)

**Immutable:** Never UPDATE or DELETE audit_logs. If correction needed, insert new record with `action = 'correction'` and `details` explaining.

---

## 9. Secrets in Code & Config

**Never commit:**
- Private keys
- API keys
- Database passwords
- Vault tokens

**Use:**
- Environment variables for service configs (references to Vault paths)
- Kubernetes Secrets (encrypted at rest with etcd encryption)
- `.env.local` for local development (gitignored)

**Example config structure:**
```yaml
database:
  host: ${DB_HOST}
  port: 5432
  name: ppob
  credentialsSecret: "secret/postgres/ppob_dev" # fetched from Vault at startup

digiflazz:
  baseURL: "https://api.digiflazz.com/v1"
  apiKeySecret: "secret/digiflazz/api_key"
  webhookSecretSecret: "secret/digiflazz/webhook_secret"

jwt:
  privateKeySecret: "secret/jwt/private_key"
  publicKeyEndpoint: "https://api.ppob.co.id/.well-known/jwks.json"
```

---

## 10. Input Validation & Sanitization

**Validation Strategy:**
1. **Schema Validation:** All incoming JSON validated against OpenAPI/JSON Schema spec before handler invocation
2. **Field-Level Validation:** Business rule validation (e.g., selling_price ≥ platform_price)
3. **Sanitization:** Strip HTML/script tags from any user-provided text fields (name, product search queries)

**Validation Libraries:**
- Go: `github.com/go-playground/validator/v10` for struct tags
- JSON Schema: `github.com/xeipuuv/gojsonschema` for request body validation

**Critical Validations:**
- `phone_number`: Must match Indonesian format `^\+62[0-9]{8,12}$`
- `pin`: Exactly 6 digits, numeric only, not sequential, not all same
- `password`: Min 8 chars, 1 upper, 1 lower, 1 digit
- `customer_no`: Format varies by product — validate against `product_validation_patterns` table
- `amount`/`selling_price`: Positive decimal, max 100,000,000 (avoid overflow)
- `otp_code`: Exactly 6 digits numeric, matches stored hash

---

## 11. Penetration Testing & Compliance

**Before Launch:**
- External pentest by licensed firm (check Indonesian regulations — possibly need OJK/FSS approval for PPOB aggregator)
- Internal security review (static analysis: `golangci-lint --enable=ethnicity`, `bandit` for Python scripts if any)
- Dependency scanning: `trivy fs .`, `syft` for SBOM generation

**Compliance:**
- **Data Residency:** All data stored in Indonesia (AWS Jakarta ap-southeast-3)
- **PDP Law (Personal Data Protection):** Consent for phone number use; right to be forgotten (anonymize user while preserving transaction integrity for accounting)
- **Financial Record Retention:** 5 years (Indonesia tax law) — ensure backups cover this period
- **Tax:** Prices shown inclusive of all taxes; platform responsible for tax reporting

---

## 12. Security Incident Response Playbook

### Scenario 1 — Compromised Mitra Account (stolen PIN)
1. User reports unauthorized transactions to support
2. Support revokes all refresh tokens for user: `redis-cli DEL refresh_token:{user_id}:*`
3. Force password + PIN reset via support portal
4. Review audit logs for suspicious activity (IP change, device change, unusual transaction patterns)
5. Refund policy: if fraud confirmed, refund from platform reserve fund (not from Digiflazz; they do not refund once transaction success)

### Scenario 2 — Digiflazz API Downtime
1. Circuit breaker opens after 5 consecutive failures (see INFRASTRUCTURE.md for circuit breaker config)
2. All transaction attempts fail fast with `DIGIFLAZZ_UNAVAILABLE` error (no retry storm)
3. Mobile app shows "Layanan sedang gangguan, silakan coba beberapa saat lagi"
4. Alert sent to Slack #digiflazz-outage and PagerDuty if >5 min
5. Wait for Digiflazz status page or manually reset circuit breaker after 5 min cooldown

### Scenario 3 — Database Corruption / Data Loss
1. Immediately stop all write operations (enable maintenance mode via feature flag)
2. Restore from latest backup (daily snapshot) to new isolated instance
3. Apply WAL replay to point-in-time just before corruption detected
4. Validate data integrity (checksum verification)
5. Resume service, notify users of brief outage via push notification

### Scenario 4 — JWT Key Compromise
1. Immediately publish new public key via `/.well-known/jwks.json` (mark old key as revoked)
2. Deploy updated services with new private key from Vault
3. All existing access tokens expire naturally within 15 min; force logout critical accounts (Mitra with high balance) via refresh token revocation
4. Rotate all active refresh tokens (bulk invalidation and forced re-login)
5. Rotate Digiflazz API keys if they were also stored in same compromised vault

---

## 13. Additional Security Controls (Future Phase)

- **TOTP MFA** optional for Mitra accounts with balance > Rp 10,000,000
- **Hardware Security Modules (HSM)** for JWT key storage in production (AWS CloudHSM)
- **Anomaly Detection:** Machine learning on transaction patterns (unusual time, amount, customer_no) — alert on deviation
- **Bug Bounty Program** after 6 months of stable operation
- **WAF (Web Application Firewall)** in front of API Gateway for additional layer
- **DPA (Data Processing Agreement)** required for any third-party analytics services

---

## 14. Developer Security Checklist

Every developer must run this checklist before merging PR:

- [ ] All DB queries use parameterized statements or ORM (no string concatenation)
- [ ] Passwords hashed with bcrypt (cost 12) before storage
- [ ] PINs hashed with argon2id (memory=64MB, iter=2, parallelism=1) before storage
- [ ] JWT signed with RSA private key, verified with public key from JWKS endpoint
- [ ] Secrets fetched from Vault at service startup, never hardcoded or in config files
- [ ] Rate limiting middleware enabled on all public endpoints
- [ ] Audit log entry written for every state-changing operation (wallet, transaction, staff assignment)
- [ ] PII (password, pin, phone, customer_no) masked/redacted in structured logs
- [ ] Idempotency keys checked to prevent duplicate transactions (duplicate ref_id rejected)
- [ ] Ownership verified: user can only access their own resources (wallet, staff, transactions) — always query with `user_id = ? AND ...`
- [ ] Input validated against JSON schema before processing (reject with 400 if invalid)
- [ ] Errors returned without stack traces or internal details in production (`debug=false`)
- [ ] TLS used for all outbound calls (including Digiflazz)
- [ ] mTLS enabled for internal service-to-service communication (if using service mesh)
- [ ] Circuit breaker configured for all outbound calls (Digiflazz, internal services)

---

## Appendix A — Algorithm Parameters Reference

### Argon2id (PIN Hashing)
```go
argon2.IDKey([]byte(pin), salt, 2, 64*1024, 1, 32)
// Memory: 64 MB, Iterations: 2, Parallelism: 1, Hash length: 32 bytes
// Expected verification time: ~200ms on mid-range mobile device
```

### bcrypt (Password Hashing)
```go
bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // cost=10 initially
// After load testing, increase to cost=12 if verification <500ms acceptable
// Cost factor doubling ≈ 2x time; recommended: 10–12
```

### MD5 (Digiflazz Signature)
External API requirement; cannot change. Used only for outbound request signature generation, not for security elsewhere.

---

## Appendix B — Vault Policy Example

```hcl
# Policy for Auth Service
path "secret/jwt/*" {
  capabilities = ["read"]
}

path "secret/postgres/ppob_prod" {
  capabilities = ["read"]
}

# Policy for Integration Service
path "secret/digiflazz/*" {
  capabilities = ["read"]
}

path "secret/postgres/ppob_prod" {
  capabilities = ["read"]
}

# Policy for Wallet Service
path "secret/postgres/ppob_prod" {
  capabilities = ["read"]
}

path "secret/redis/password" {
  capabilities = ["read"]
}
```

---

## Appendix C — Compliance Cross-Check

| Regulation | Requirement | Implementation | Status |
|---|---|---|---|
| Indonesian PDP Law | Right to be forgotten | User data anonymization procedure in `BUSINESS_LOGIC_SPEC.md` | ✅ Covered |
| Tax Law (5-year retention) | Financial records kept 5 years | Backup retention 35 days + quarterly tape archive to Glacier | ⚠️ Needs implementation |
| OJK PPOB Licensing | PPOB operator license | Legal review pending | ❌ TBD |
| PCI DSS (if storing cards) | Card data protection | Not applicable (no card storage) | N/A |

---

**Next Actions:**
1. Security review meeting with infosec team (schedule within 1 week)
2. Pen test planning (vendor selection)
3. Implement vault integration and secret rotation automation
4. Deploy rate limiting middleware to staging environment for load testing
