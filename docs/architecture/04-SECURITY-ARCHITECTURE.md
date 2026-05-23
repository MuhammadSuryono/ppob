# 🔐 Security Architecture

## 1. Introduction
Security is integrated at every layer of the PPOB architecture, from device fingerprinting on the client to mTLS within the service mesh.

## 2. Authentication & Authorization

### 2.1 Adaptive Authentication
The system uses an adaptive flow based on **Device Trust Scores (0-100)**:
- **Trusted (score ≥ 70):** PIN-only login for fast access.
- **Semi-trusted (30-69):** Password + OTP required.
- **Untrusted (score < 30):** Full authentication (Password + OTP + Challenge).

### 2.2 JWT Policy
- **Algorithm:** RS256 (asymmetric signing).
- **Access Token TTL:** 15 minutes.
- **Refresh Token TTL:** 7 days (Mitra), 1 day (Staff).
- **Rotation:** One-time use refresh tokens with Redis-backed allowlist.
- **Revocation:** Centralized blocklist in Redis for compromised sessions.

### 2.3 Credential Hashing
- **Passwords:** `bcrypt` with a cost factor of 12.
- **Transaction PINs:** `Argon2id` (64MB, 2 iterations, 1 parallelism) to resist GPU brute-force attacks on the 6-digit entropy space.

## 3. Data Protection

### 3.1 PII Redaction
Personally Identifiable Information is never logged in plain text.
- **Phone Numbers:** Masked to `+62 812-345-***-****`.
- **Customer IDs:** Masked to the last 4 digits in displays and logs.

### 3.2 Encryption at Rest
- **Database:** AWS RDS Transparent Data Encryption (TDE) via KMS.
- **Local Cache:** Mobile app uses SQLCipher for Room DB and EncryptedSharedPreferences for tokens.

## 4. Network Security
- **External:** TLS 1.3 only, strictly enforced.
- **Internal:** mTLS (Mutual TLS) for all service-to-service communication via the service mesh (Istio).
- **Firewall:** Egress filtering allows outbound calls only to Digiflazz and essential cloud services.

## 5. Fraud Detection & Limits
- **Daily Limits:** Enforced atomically in the `daily_limits` table (Transaction Count & Turnover Amount).
- **Rate Limiting:** Sliding window limits on OTP requests (3/hr) and login attempts (10/15min).
- **Idempotency:** SHA256-hashed idempotency keys stored for 24 hours to prevent duplicate billing.

## 6. Secrets Management
All secrets (API keys, private keys, DB passwords) are stored in **AWS Secrets Manager** or **HashiCorp Vault**.
- Injected at runtime via environment variables or sidecars.
- No hardcoded secrets allowed in codebase or CI/CD pipelines.
- Automated rotation for database credentials.
