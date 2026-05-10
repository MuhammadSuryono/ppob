# 🚀 AI-Assisted Execution Timeline for PPOB Mobile Application

**Project:** PPOB Multi-Tenant Mobile Application  
**Based on:** `v1/docs/` specifications (API contracts, business logic, architecture)  
**Status:** Ready for Execution  
**Created:** 2026-05-07  
**Last Updated:** 2026-05-09  

---

## 📋 Executive Summary

This document provides a phased execution timeline for implementing the complete PPOB mobile application using AI-assisted development. The timeline covers 8 major phases across **~20-24 weeks** (5-6 months) with clear deliverables, AI augmentation opportunities, and risk mitigation strategies.

**AI Execution Approach:**
- Use AI for boilerplate code generation, documentation, test creation, and architecture validation
- Human oversight for business logic, financial calculations, security-critical code, and integration points
- Iterative delivery with 2-week sprints

---

## 🎯 Project Phases Overview

| Phase | Name | Duration | Key Deliverables | AI Role |
|---|---|---|---|---|
| **1** | Foundation & Core Services | 3 weeks | Auth Service, User/Role Service, Database schema | High (boilerplate, tests) |
| **2** | Wallet & Financial Engine | 3 weeks | Wallet Service, wallet_events, accounting logic | Medium-High (business rules) |
| **3** | Product & Transaction Service | 2 weeks | Product Service, Transaction state machine | Medium (business logic validation) |
| **4** | Digiflazz Integration | 2 weeks | Integration Service, sync jobs, webhook handling | Medium (API mapping) |
| **5** | Mobile App (Android Native) | 4 weeks | Auth flow, transaction UI, offline support | Very High (UI generation) |
| **6** | Infrastructure & CI/CD | 2 weeks | EKS cluster, RDS, Redis, CI/CD pipelines | Low-Medium (Terraform IaC) |
| **7** | Observability & Security | 2 weeks | Prometheus/Grafana, logging, security hardening | Low (config, rules) |
| **8** | Testing & Go-Live | 2 weeks | E2E tests, load testing, deployment, cutover | Medium (test scenarios) |

**Total Estimated Duration: 20-24 weeks** (concurrent development possible with team splitting)

---

## 📅 Detailed Timeline & Tasks

### **Phase 1: Foundation & Core Services (Weeks 1-3)**

**Objective:** Establish authentication, user management, and database foundation.

**API Reference:** All endpoint specifications are in `MICROSERVICES_API_DOC.md` (base URL: `https://fedora.sinauplatform.id/api/v1/{api_name}`).

#### Week 1: Project Setup & Auth Service
- **Task 1.1:** Initialize monorepo structure, CI/CD pipeline skeleton
  - AI: Generate GitHub Actions workflows, Dockerfile templates, Makefile, Gradle build config
  - Deliverable: Working dev environment with lint/test/format commands
- **Task 1.2:** Database schema migration setup (`golang-migrate`)
  - AI: Generate all migration files from existing schema docs (`DATA_MODEL_VALIDATION.md`)
  - Deliverable: SQL migration files (001-012), seed data scripts
- **Task 1.3:** Implement Auth Service (registration, OTP, login)
  - AI: Generate JWT token service, bcrypt/argon2id utilities, rate limiting middleware
  - Human: Implement trust score calculation logic, device fingerprinting
  - Reference: `MICROSERVICES_API_DOC.md` → Auth Service (Base URL: `https://fedora.sinauplatform.id/api/v1/auth`)
  - Deliverable: Auth endpoints working with Postman collection against `https://fedora.sinauplatform.id/api/v1/auth/*`
- **Task 1.4:** Write unit tests for Auth Service (70%+ coverage)
  - AI: Generate test cases for all auth flows, edge cases, error scenarios
  - Deliverable: `auth_service_test.go` with mocked dependencies

#### Week 2: User & Role Service
- **Task 2.1:** User CRUD with role assignment
  - AI: Generate service layer, repository interfaces, API handlers
  - Human: Implement role switching logic, wallet auto-creation trigger
  - Reference: `MICROSERVICES_API_DOC.md` → User Service (Base URL: `https://fedora.sinauplatform.id/api/v1/user`)
  - Deliverable: `/api/v1/user/*` endpoints working (Postman collection)
- **Task 2.2:** Implement multi-role wallet linkage
  - AI: Generate SQL queries for wallet lookup by role, database trigger code
  - Deliverable: Wallet resolution working, role switching updates JWT claims
- **Task 2.3:** Write integration tests (Auth + User Service interaction)
  - AI: Generate test scenarios for multi-tenant data isolation
  - Deliverable: `integration_test_user_role.go`

#### Week 3: Database & Validation Layer
- **Task 3.1:** Validate all CHECK constraints, triggers
  - AI: Generate SQL validation scripts, test data violating constraints
  - Human: Review business rule constraints (selling_price ≥ platform_price, etc.)
  - Deliverable: All constraints enforced, constraint violation tests passing
- **Task 3.2:** Implement structured logging middleware
  - AI: Generate JSON logger setup with trace_id propagation
  - Deliverable: All services log in JSON format, PII redaction enabled
- **Task 3.3:** API documentation (OpenAPI spec generation)
  - AI: Generate OpenAPI 3.0 YAML from code annotations (Swaggo or oapi-codegen)
  - Deliverable: `openapi.yaml`, interactive Swagger UI at `/docs`

**Phase 1 Milestone:** ✅ User can register, login, switch roles; database schema validated; 70%+ test coverage on core services.

---

### **Phase 2: Wallet & Financial Engine (Weeks 4-6)**

**Objective:** Implement wallet operations with event sourcing, concurrency control, and financial integrity.

#### Week 4: Wallet Service Core
- **Task 4.1:** Wallet CRUD + balance queries
  - AI: Generate repository methods, service layer with pessimistic locking (`SELECT FOR UPDATE`)
  - Deliverable: `/wallets/v1/balance` endpoint, real-time balance display
- **Task 4.2:** Implement `wallet_events` table & event sourcing logic
  - AI: Generate event structs, append-only repository pattern, balance reconstruction algorithm
  - Human: Ensure atomicity of event + balance update in single DB transaction
  - Deliverable: All wallet changes create immutable events; balance computed correctly
- **Task 4.3:** Hold & release mechanism for pending transactions
  - AI: Generate hold placement/release SQL with atomic `INSERT ... ON CONFLICT`
  - Deliverable: `balance_available` and `balance_held` update correctly

#### Week 5: Financial Operations
- **Task 5.1:** Mitra → Staff top-up implementation
  - AI: Generate atomic transfer logic (debit Mitra, credit Staff, linked events)
  - Human: Implement compensation job for partial failures
  - Deliverable: `/wallets/v1/staff/{id}/topup` with audit trail
- **Task 5.2:** Daily limit enforcement (atomic, race-free)
  - AI: Generate `daily_limits` table queries with `INSERT ... ON CONFLICT ... DO UPDATE`
  - Human: Verify limits cannot be bypassed with concurrent requests
  - Deliverable: Staff transaction count/amount limits enforced correctly
- **Task 5.3:** Wallet reconciliation job (hourly)
  - AI: Generate cron job (K8s CronJob) that compares event sum to cached balance
  - Deliverable: Drift detection + alert if mismatch > Rp 1,000

#### Week 6: Margin & Commission
- **Task 6.1:** Implement margin calculation engine
  - AI: Generate service that computes `Margin = SellingPrice - PlatformPrice`
  - Human: Implement `staff_global_margin_settings` + `staff_product_margin_overrides` lookup
  - Deliverable: Commission calculated correctly for both schemes (FixedAllowance, MarginShare)
- **Task 6.2:** Commission crediting (real-time to staff wallet)
  - AI: Generate transaction success handler that creates `commissions` record + credit event
  - Deliverable: Staff wallet balance updates immediately upon transaction success
- **Task 6.3:** Financial integrity invariants tests
  - AI: Generate SQL queries that verify invariants daily (see `BUSINESS_LOGIC_SPEC.md`)
  - Deliverable: Daily reconciliation job with CRITICAL alerting on violations

**Phase 2 Milestone:** ✅ Wallet operations thread-safe, financial calculations accurate, daily limits enforced, audit-compliant event sourcing.

---

### **Phase 3: Product & Transaction Service (Weeks 7-8)**

**Objective:** Product catalog management and transaction lifecycle with state machine.

#### Week 7: Product Service
- **Task 7.1:** Product CRUD + category management
  - AI: Generate CRUD endpoints, pagination, filtering logic
  - Deliverable: `/products/v1` endpoints (list by category, search)
- **Task 7.2:** Product sync from Digiflazz (mock mode)
  - AI: Generate sync job that fetches `pricelist` (prepaid/postpaid), upserts products
  - Human: Implement Redlock distributed locking to prevent concurrent sync
  - Deliverable: Hourly prepaid sync (simulated), every-5-min postpaid sync
- **Task 7.3:** Price validation (selling_price ≥ platform_price)
  - AI: Generate validation middleware, error response for below-cost attempts
  - Deliverable: Transactions reject if Mitra sets price below platform cost

#### Week 8: Transaction Service Core
- **Task 8.1:** Transaction initiation with idempotency
  - AI: Generate handler with `Idempotency-Key` check, wallet hold placement
  - Deliverable: `/transactions/v1/initiate` returns transaction_id, holds balance
- **Task 8.2:** State machine implementation (see `TRANSACTION_STATE_MACHINE.md`)
  - AI: Generate state transition guard functions, status update logic
  - Human: Implement webhook processing with idempotent deduplication
  - Deliverable: Proper state transitions (Initiated → Pending → Success/Failed/Expired)
- **Task 8.3:** Transaction history & detail endpoints
  - AI: Generate paginated queries, filtering by status/date/category
  - Deliverable: `/transactions/v1/history` with cursor-based pagination

**Phase 3 Milestone:** ✅ Products synced, transactions initiated with holds, state machine functional, history queryable.

---

### **Phase 4: Digiflazz Integration (Weeks 9-10)**

**Objective:** Real integration with Digiflazz API via the Integration Service, webhook handling, and error mapping.

**API Reference:** All endpoint specifications are defined in `MICROSERVICES_API_DOC.md` under "Integration Service" (Base URL: `https://fedora.sinauplatform.id/api/v1/integration`).

#### Week 9: Integration Service Setup
- **Task 9.1:** Digiflazz API client with HMAC signature generation
  - AI: Generate signature generator (MD5 as required), HTTP client with retry
  - Human: Store credentials in Vault, implement circuit breaker (`gobreaker`)
  - Reference: `MICROSERVICES_API_DOC.md` → `/integration/digiflazz/transaction` endpoint spec
  - Deliverable: Connection to Digiflazz sandbox successful via `POST /integration/digiflazz/transaction`
- **Task 9.2:** Product sync against real Digiflazz sandbox
  - AI: Generate sync job using real API, RC error handling
  - Deliverable: Live product data flowing into DB, prices updated (ref: `GET /integration/providers`, `GET /integration/errors`)
- **Task 9.3:** Transaction initiation to Digiflazz
  - AI: Generate request builder, response parser, RC → internal status mapping
  - Deliverable: Successful transaction in sandbox (RC 00), pending handling (RC 03)

#### Week 10: Webhook & Error Handling
- **Task 10.1:** Webhook endpoint with HMAC verification
  - AI: Generate signature verification (SHA1), Redis lock to prevent duplicate processing
  - Reference: `MICROSERVICES_API_DOC.md` → `/integration/webhook/digiflazz` endpoint spec
  - Deliverable: Webhook endpoint verified at `https://fedora.sinauplatform.id/api/v1/integration/webhook/digiflazz`
- **Task 10.2:** Complete error mapping from `ERROR_HANDLING.md` + `MICROSERVICES_API_DOC.md` → `/integration/errors`
  - AI: Generate error code catalog, user-friendly Indonesian messages
  - Human: Validate error messages with native speaker
  - Deliverable: All RC codes mapped per `MICROSERVICES_API_DOC.md` error catalog, error response format standardized
- **Task 10.3:** Retry policy & dead letter queue
  - AI: Generate compensation job table (`compensation_jobs`) and worker
  - Deliverable: Failed multi-step operations auto-retry with exponential backoff (monitor via `GET /integration/compensation/dead-letter`)

**Phase 4 Milestone:** ✅ End-to-end transaction flow: Initiate → Digiflazz → Webhook → Wallet debit → Commission credited.

---

### **Phase 5: Mobile App (Android Native) (Weeks 11-14)**

**Objective:** Build complete native Android app with offline-first architecture using Kotlin, Jetpack Compose, Hilt, and Room.

#### Week 11: Foundation & Auth Screens
- **Task 11.1:** Project setup, dependency injection (Hilt), networking (Retrofit + OkHttp)
  - AI: Generate project structure (modular Gradle setup), Hilt DI modules, OkHttp client with interceptors
  - Deliverable: App runs on emulator, API connectivity established
- **Task 11.2:** Login/Register/OTP screens with form validation
  - AI: Generate all auth UI screens (Jetpack Compose) with Indonesian copy, custom PIN pad composable
  - Deliverable: User can register, verify OTP, set password/PIN, login
- **Task 11.3:** Device fingerprint collection & trust logic UI
  - AI: Generate device info collection utility, trust score display
  - Deliverable: Login flow adapts based on trust (OTP required if low score)

#### Week 12: Core Transaction Flow
- **Task 12.1:** Home screen, category grid, product listing
  - AI: Generate product catalog UI with Coil image loading, search functionality
  - Deliverable: Grid of PPOB categories, product list with Indonesian names
- **Task 12.2:** Transaction initiation flow (customer input, confirmation, PIN entry)
  - AI: Generate multi-step transaction wizard with Compose Navigation
  - Deliverable: Full transaction flow from product selection to PIN confirmation
- **Task 12.3:** Transaction history screen with filters
  - AI: Generate paged list with Paging 3 library, status tabs (All/Success/Pending/Failed)
  - Deliverable: Transaction history with offline cache (Room)

#### Week 13: Wallet & Staff Management (Mitra)
- **Task 13.1:** Wallet screen (balance display, held balance explanation)
  - AI: Generate balance card UI, transaction list linked to wallet ViewModel
  - Deliverable: Real-time balance (with caching refresh)
- **Task 13.2:** Staff list, add/edit staff screens
  - AI: Generate CRUD UI for staff management, margin scheme selector (Fixed/Share)
  - Deliverable: Mitra can add staff, set margin/daily limits
- **Task 13.3:** Staff top-up modal
  - AI: Generate top-up dialog, confirmation with balance check
  - Deliverable: Mitra can top-up staff wallet from own wallet

#### Week 14: Offline, Push & Polish
- **Task 14.1:** Offline queue & sync mechanism
  - AI: Generate Room entities for pending sync, `PendingSyncItem` model, `WorkManager` sync service
  - Deliverable: Transactions queue when offline, auto-sync on reconnect
- **Task 14.2:** Push notifications (FCM) integration
  - AI: Generate `FirebaseMessagingService`, notification handlers for success/failure/low-balance events, notification channel setup
  - Steps: Add `google-services.json` (from Firebase project ppob-uat) to `app/` directory, apply `com.google.gms.google-services` plugin in `app/build.gradle.kts`, add Firebase BoM (`com.google.firebase:firebase-bom:34.13.0`) and `firebase-analytics` dependency
  - Deliverable: Transaction result push, deep linking to detail screen, Firebase Analytics event tracking enabled
- **Task 14.3:** Biometric authentication (optional)
  - AI: Generate `BiometricPrompt` wrapper, settings toggle in Compose
  - Deliverable: Fingerprint/Face ID can replace PIN entry on trusted devices
- **Task 14.4:** UI polish & performance optimization
  - AI: Generate splash screen, error states, empty states, loading skeletons, baseline profile
  - Deliverable: APK <30MB (split per ABI), cold start <2s, 60fps scrollable lists

**Phase 5 Milestone:** ✅ Complete native Android app with auth, transaction flow, wallet/staff management, offline support, push notifications.

---

### **Phase 6: Infrastructure & CI/CD (Weeks 15-16)**

**Objective:** Deploy production-grade Kubernetes infrastructure and automated deployment.

#### Week 15: Kubernetes & Database
- **Task 15.1:** Terraform VPC + EKS cluster (Jakarta region)
  - AI: Generate Terraform modules for VPC, EKS with 3 node groups (system/app/worker)
  - Human: Review AWS resource sizing, cost optimization (spot instances)
  - Deliverable: Running EKS cluster in `ap-southeast-3` (Jakarta)
- **Task 15.2:** RDS PostgreSQL + ElastiCache Redis deployment
  - AI: Generate Terraform for RDS (Multi-AZ), Redis cluster (3 nodes)
  - Deliverable: Database accessible via private endpoints, PgBouncer configured
- **Task 15.3:** Secrets management (AWS Secrets Manager or Vault)
  - AI: Generate secret rotation lambdas, External Secrets Operator config
  - Deliverable: All secrets fetched at runtime, no hardcoded credentials

#### Week 16: CI/CD & Deployment
- **Task 16.1:** Docker image builds + ECR repository
  - AI: Generate multi-stage Dockerfile (Gradle build → distroless), ECR repo creation scripts
  - Deliverable: Images pushed to ECR on every commit
- **Task 16.2:** GitHub Actions workflows (build/test/deploy)
  - AI: Generate workflows for dev (auto-deploy) and prod (manual approval)
  - Deliverable: Push to `main` triggers prod deploy via slash command
- **Task 16.3:** Kubernetes manifests (Helm charts or Kustomize)
  - AI: Generate Helm charts for all services with HPA, resource limits, probes
  - Deliverable: `helm install ppob ./chart` deploys full stack to cluster
- **Task 16.4:** Ingress + TLS (cert-manager + Let's Encrypt)
  - AI: Generate Nginx Ingress config, TLS certificate requests
  - Deliverable: `https://api.ppob.co.id` serving all services

**Phase 6 Milestone:** ✅ Infrastructure fully automated, CI/CD pipeline functional, staging environment live.

---

### **Phase 7: Observability & Security (Weeks 17-18)**

**Objective:** Comprehensive monitoring, logging, tracing, and security hardening.

#### Week 17: Monitoring Stack
- **Task 17.1:** Prometheus + Grafana deployment
  - AI: Generate `kube-prometheus-stack` Helm values, custom dashboards JSON
  - Deliverable: Metrics scraped from all services, Grafana dashboards (Service Overview, Transaction, Wallet, Digiflazz)
- **Task 17.2:** Alertmanager rules + routing
  - AI: Generate alert rules from `OBSERVABILITY_SPEC.md` (error rate, pending backlog, wallet drift)
  - Deliverable: Critical alerts → PagerDuty, warnings → Slack
- **Task 17.3:** Loki (logs) + Jaeger (traces) setup
  - AI: Generate Fluent Bit DaemonSet config, OpenTelemetry Collector config
  - Deliverable: Logs searchable by `trace_id`, traces viewable in Jaeger UI

#### Week 18: Security Hardening
- **Task 18.1:** Network policies + pod security standards
  - AI: Generate K8s NetworkPolicy resources (default-deny), PodSecurityContext
  - Deliverable: Pods non-root, drop-all capabilities, network segmentation
- **Task 18.2:** mTLS internal communication (service mesh setup)
  - AI: Generate Istio PeerAuthentication, DestinationRule manifests
  - Deliverable: All service-to-service traffic encrypted
- **Task 18.3:** Security scanning integration
  - AI: Generate Trivy scanner config, GitHub CodeQL workflow
  - Deliverable: PRs blocked if critical vulnerabilities detected
- **Task 18.4:** Penetration test preparation
  - Human: Create test environment, document API surface, prepare test accounts
  - Deliverable: Ready for external security audit

**Phase 7 Milestone:** ✅ Full observability stack operational, security controls in place, compliance checks passed.

---

### **Phase 8: Testing & Go-Live (Weeks 19-20+)**

**Objective:** Validate system under load, conduct UAT, and deploy to production.

#### Week 19: Testing & Validation
- **Task 19.1:** Load testing (k6 or Locust)
  - AI: Generate load test scripts simulating 1000 concurrent transactions
  - Human: Run tests, identify bottlenecks, tune DB/Redis connection pools
  - Deliverable: System handles 100 TPS with p99 < 1s, no errors
- **Task 19.2:** End-to-end test automation
  - AI: Generate Maestro or Appium E2E tests for critical user journeys on Android
  - Deliverable: 90%+ E2E coverage, CI runs on every PR
- **Task 19.3:** User Acceptance Testing (UAT) with pilot Mitra
  - Human: Onboard 3-5 pilot Mitras, collect feedback, fix bugs
  - Deliverable: UAT sign-off from product owner
- **Task 19.4:** Disaster recovery drill
  - Human: Simulate DB failure, restore from backup, measure RTO/RPO
  - Deliverable: Recovery procedure documented and validated (<4h RTO)

#### Week 20: Production Deployment
- **Task 20.1:** Production environment provisioning (Terraform apply)
  - AI: Generate production-specific Terraform vars (larger instance sizes)
  - Deliverable: Production EKS, RDS, Redis live
- **Task 20.2:** Blue-green deployment or canary rollout
  - AI: Generate ArgoCD manifests for gradual traffic shifting
  - Deliverable: Zero-downtime deployment strategy ready
- **Task 20.3:** Cutover & monitoring
  - Human: Deploy to production, monitor alerts, verify transaction flow end-to-end
  - Deliverable: **System LIVE** with successful transactions flowing

**Post-Launch (Weeks 21-24):**
- Week 21-22: Hypercare — 24/7 monitoring, rapid bug fixes
- Week 23-24: Post-mortem, performance tuning, phase 2 planning

---

## 🤖 AI-Assistance Strategy

### Where to Use AI (High Leverage)

| Task | AI Tool | Prompt Example |
|---|---|---|
| **Boilerplate CRUD** | Cursor/Claude Code | "Generate Go structs and GORM repositories for `wallets` and `wallet_events` tables with all fields from schema" |
| **Unit Tests** | AI test generator | "Write unit tests for `CalculateCommission` function covering FixedAllowance and MarginShare schemes" |
| **API Handlers** | AI code completion | "Create Gin handler for `POST /users/v1/staff` with validation, error handling, and audit logging" |
| **Database Migrations** | AI SQL generation | "Write PostgreSQL migration to split `staff_margin_settings` into `staff_global_margin_settings` and `staff_product_margin_overrides` with data backfill" |
| **Android UI (Compose)** | AI UI generator | "Generate Jetpack Compose screen for transaction history with LazyColumn, pull-to-refresh, and status tabs" |
| **Terraform IaC** | AI IaC generator | "Create Terraform for EKS cluster with 3 node groups (system, app, worker) in AWS Jakarta region" |
| **Dockerfiles** | AI container config | "Write multi-stage Dockerfile for Go service with distroless final image" |
| **Documentation** | AI doc writer | "Generate OpenAPI spec for Auth Service endpoints from Go code comments" |
| **Test Data** | AI data factory | "Create factory functions for random Indonesian phone numbers and valid OTP codes" |
| **Room DAO & Entities** | AI code generation | "Generate Room entity, DAO, and TypeConverter for Transaction data class with Parcelize" |

### Where Humans Must Lead (Critical)

- Financial calculations (commission, margin, tax implications)
- Security-critical code (JWT signing, PIN hashing, signature verification)
- Database transaction boundaries & concurrency control
- Integration with external API (Digiflazz signature, webhook verification)
- Compliance & legal (data retention, PDP anonymization)
- Production deployment decisions & rollback plans
- Android app store compliance review (Google Play policies)

---

## 📊 Dependencies & Critical Path

```
Phase 1 (Auth/User) → Phase 2 (Wallet) → Phase 3 (Transaction) → Phase 4 (Integration) → Phase 5 (Mobile)
        ↓                      ↓                    ↓                        ↓
     Database            Event Sourcing       State Machine         Webhook Handling
```

**Critical Path:** Database → Auth → Wallet → Transaction → Integration → Mobile

**Parallelizable Paths:**
- Mobile UI development can start after Phase 1 (API contracts defined) using mock servers
- Infrastructure (Phase 6) can start after Phase 1 (services defined)
- Observability (Phase 7) can start after Phase 3 (metrics instrumented)

---

## ⚠️ Risk Mitigation & Blockers

| Risk | Impact | Mitigation | Owner |
|---|---|---|---|
| **Digiflazz API unstable in sandbox** | High — integration delayed | Contact Digiflazz support early, request stable sandbox credentials, build API client with feature flags to toggle to mock mode | Integration Lead |
| **Team lacks Golang/Android Native expertise** | Medium — velocity reduced | Allocate 1 week training per phase, use AI pair-programming, hire consultant for review | Tech Lead |
| **Database performance at scale** | High — system bottlenecks | Load test early (Phase 2), add indexes proactively, partition `transactions` by month after 1M rows | DBA |
| **Security audit findings** | Critical — may require rework | Conduct internal threat modeling (STRIDE) before Phase 7, fix issues pre-external audit | Security Architect |
| **Google Play Store rejection** | Medium — launch delayed | Review Play Store guidelines early, ensure no policy violations (payments, data usage, SMS restrictions) | Product Manager |
| **Android device fragmentation** | Medium — UI/compat issues | Use AndroidX libraries, test on multiple API levels (26, 28, 30, 33, 34), use Firebase Test Lab | Mobile Lead |
| **Keystore compatibility issues** | High — security feature broken | Test Keystore on multiple OEMs (Samsung, Xiaomi, etc.), implement fallback to encrypted SharedPreferences | Security Lead |
| **Team burnout from 6-month timeline** | Medium — quality degrades | Plan 1-week buffer between phases, rotate developers across services, celebrate milestones | Project Manager |

---

## 📈 Success Metrics & KPIs

**Technical KPIs:**
- Test coverage: ≥80% unit, ≥50% integration
- API response time: p95 < 300ms (without Digiflazz)
- System availability: 99.5% (excluding Digiflazz downtime)
- Transaction success rate: ≥99% (final status Success)
- Pending-to-success latency: p95 < 30s
- Android app startup time: <2s cold start
- Crash-free rate: ≥99.5% (Firebase Crashlytics)
- ANR rate: <0.5%

**Business KPIs (Post-Launch):**
- Daily active Mitra: Target 100 within 1 month
- Transactions per day: Target 1,000 within 1 month
- Average transaction value: Rp 25,000–50,000
- Staff per Mitra avg: ≥3

---

## 🛠️ Toolchain & AI Assistants

**Recommended Stack for AI-Assisted Development:**
- **Code Generation:** Claude Code, Cursor, GitHub Copilot X
- **Android Development:** Android Studio with Gemini / Copilot integration
- **Documentation:** AI markdown generator (for API specs, READMEs)
- **Testing:** AI test case generator (TestGPT, Ponicode), Firebase Test Lab
- **Review:** AI code reviewer (CodeRabbit, Sweep)
- **Architecture:** AI diagram generator (from D2 or Mermaid text)

**Prompt Templates:** Store in `docs/ai-prompts/` for reuse:
```
PROMPT: Generate Room entity, DAO, and Repository for Transaction with offline sync queue
CONTEXT: Transaction schema from DATA_MODEL_VALIDATION.md, sync pattern from MOBILE_APP_SPEC.md section 4.2
OUTPUT: Transaction.kt, TransactionDao.kt, PendingSyncItem.kt, SyncWorker.kt
```

```
PROMPT: Generate Retrofit service interface for {service} API with base URL https://fedora.sinauplatform.id/api/v1/{api_name}
CONTEXT: Endpoint definitions from MICROSERVICES_API_DOC.md, API name mapping: auth, user, wallet, transaction, integration, product
OUTPUT: {ServiceName}Service.kt with suspend functions, DTOs, error handling
```

```
PROMPT: Add Google Services plugin and Firebase dependencies to app/build.gradle.kts
CONTEXT: google-services.json located at app/google-services.json (Firebase project: ppob-uat)
OUTPUT: Updated build.gradle.kts files (project-level and module-level) with google-services plugin and Firebase BoM dependencies
```

---

## 📚 Reference Documentation Index

All AI prompts should reference these source documents:

```
v1/docs/
├── API_AUTHENTICATION_FLOW.md          → Auth flow details
├── BUSINESS_LOGIC_SPEC.md              → Pricing, margins, limits
├── CONCURRENCY_CONTROL.md              → Wallet locking, event sourcing
├── DATA_MODEL_VALIDATION.md            → DB constraints, validation rules
├── Database Schema for PPOB.md         → Full table definitions
├── DIGIFLAZZ_INTEGRATION_GUIDE.md      → API mapping, RC codes
├── ERROR_HANDLING.md                   → Error responses, retry policy
├── MICROSERVICES_API_DOC.md           → Full endpoint reference with base URL https://fedora.sinauplatform.id/api/v1/{api_name}
├── MOBILE_APP_SPEC.md                 → Android Native architecture, state management
├── OBSERVABILITY_SPEC.md              → Metrics, logs, traces standards
├── SECURITY_ARCHITECTURE.md           → JWT, hashing, rate limiting
├── TRANSACTION_STATE_MACHINE.md       → State transitions, webhooks
├── Product Requirement Document...md  → Feature requirements
└── google-services.json               → Firebase config (project: ppob-uat, package: com.yonotech.ppob.mobile)
```

---

## 🔄 Iteration & Adaptation

**Review Cadence:**
- **Daily:** Standup, progress vs timeline
- **Weekly:** Demo completed features, adjust next week's tasks
- **Phase End:** Retrospective, phase gate review before starting next

**If Behind Schedule:**
- Prioritize MVP features (prepaid only, skip postpaid initially)
- Reduce polish (basic UI instead of refined)
- Defer non-critical (biometrics, advanced reporting)
- Add developer resources (contractors for specific phases)

**If Ahead Schedule:**
- Add stretch goals: TOTP MFA, advanced fraud detection, admin dashboard
- Conduct more thorough load testing, security audit prep

---

## 🎉 Go-Live Checklist

Before marking project complete:

- [ ] All API endpoints functional with >80% test coverage
- [ ] Android app published to Google Play Store (internal → closed → production track)
- [ ] Load test passed: 100 TPS sustained, p99 < 1s
- [ ] Security audit completed, all critical issues resolved
- [ ] Digiflazz production credentials obtained and tested
- [ ] Backup & DR drill performed (RTO <4h achieved)
- [ ] Support team trained, runbooks documented
- [ ] Legal/compliance sign-off (PDP, tax retention)
- [ ] Pilot Mitras onboarded and transacting successfully
- [ ] Monitoring dashboards green, alerts routed correctly
- [ ] Documentation complete: API spec, architecture diagrams, runbooks
- [ ] Crash-free rate ≥99.5% over 72-hour period post-launch

---

## 📞 Emergency Contacts

**During Active Development:**
- **Tech Lead:** [Contact]
- **DevOps / Infra:** [Contact]
- **Product Owner:** [Contact]

**Post-Launch (24/7 On-Call):**
- **PagerDuty:** `ppob-sre@company.com`
- **Slack:** `#ppob-incidents`
- **Digiflazz Support:** [Hotline/Email pre-saved]

---

## 🏁 Conclusion

This timeline provides a structured path from zero to production-ready PPOB application using AI-assisted development. The key to success is **iterative delivery** (2-week sprints), **continuous validation** (tests + load), and **early integration** (Phase 4 connects real Digiflazz, not mocks).

**Technology shift:** This timeline has been updated to reflect the move from Flutter cross-platform to **Android Native (Kotlin + Jetpack Compose + Hilt + Room + Retrofit)**. Android Native offers:
- ~50% smaller APK size (no Flutter engine overhead)
- Direct access to Android platform APIs (Keystore, BiometricPrompt, FCM)
- Better performance on low-end Indonesian devices
- Native Material Design 3 compliance
- Simplified dependency graph (no Dart/Flutter SDK dependencies)

**Next Action:** Begin Phase 1, Task 1.1 — generate CI/CD pipeline skeleton with AI assistant.

---

**Document Owner:** Project Manager / Tech Lead  
**Last Updated:** 2026-05-09  
**Next Review:** After Phase 1 completion (Week 3)