# 🛠️ Microservices Architecture

## 1. Introduction
The PPOB backend is architected as a set of six core microservices, each owning a specific business domain. This decoupling allows for independent scaling, deployment, and technology evolution.

## 2. Core Services

### 2.1 Auth Service (Port 8081)
**Domain:** Identity & Access Management.
- **Responsibilities:**
    - User registration and adaptive login flow.
    - JWT token issuance and revocation.
    - OTP generation/verification (SMS/WhatsApp).
    - Device fingerprinting and trust scoring.
- **Primary Data:** `users`, `otp_codes`, `device_fingerprints`.

### 2.2 User Service (Port 8082)
**Domain:** Profile & Role Management.
- **Responsibilities:**
    - User profile CRUD.
    - RBAC (Role-Based Access Control) management.
    - Staff management (Mitra-Staff hierarchy).
    - Notification preferences and history.
- **Primary Data:** `roles`, `user_roles`, `notifications`.

### 2.3 Wallet Service (Port 8083)
**Domain:** Financial Ledger.
- **Responsibilities:**
    - Real-time balance management.
    - Atomic hold/release/debit operations.
    - Event-sourced transaction log.
    - Mitra-to-Staff fund transfers (top-ups).
- **Primary Data:** `wallets`, `wallet_events`.

### 2.4 Transaction Service (Port 8084)
**Domain:** PPOB Operations.
- **Responsibilities:**
    - Transaction lifecycle orchestration.
    - Idempotency enforcement.
    - Commission and margin distribution logic.
    - Reporting and KPI aggregation.
- **Primary Data:** `transactions`, `daily_limits`, `commissions`.

### 2.5 Product Service (Port 8085)
**Domain:** Catalog & Pricing.
- **Responsibilities:**
    - Product and category management.
    - Platform markup calculation.
    - Digiflazz pricelist synchronization (Prepaid/Postpaid).
    - Mitra-specific price overrides.
- **Primary Data:** `products`, `categories`, `mitra_product_prices`.

### 2.6 Integration Service (Port 8086)
**Domain:** External Gateway.
- **Responsibilities:**
    - Digiflazz API client with HMAC signatures.
    - Webhook reception and HMAC verification.
    - Circuit breaker and retry logic.
    - Provider error mapping.
- **Primary Data:** `compensation_jobs`, `postpaid_inquiries`.

## 3. Resilience Strategy
- **Circuit Breakers:** Implemented using `gobreaker` to prevent cascading failures when external providers (Digiflazz) or internal services are down.
- **Retries with Backoff:** Used for transient network errors.
- **Timeouts:** Aggressive context-based timeouts (e.g., 30s for Digiflazz calls).

## 4. Monitoring & Health
- **Endpoints:** `/health/live` (liveness) and `/health/ready` (readiness).
- **Metrics:** Each service exports `/metrics` for Prometheus scraping.
