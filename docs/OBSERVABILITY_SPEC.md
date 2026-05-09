# 📊 Observability Specification for PPOB Application

**Audience:** DevOps, SRE, Backend developers  
**Last Updated:** 2026-05-07  
**Status:** Draft — metric thresholds to be calibrated in staging

---

## 1. Overview

This document defines logging, metrics, tracing, and alerting standards across all microservices. Observability enables rapid diagnosis of issues, performance monitoring, and business insight.

**Three Pillars:**
1. **Logs** — detailed event records (structured JSON)
2. **Metrics** — numerical measurements over time (Prometheus)
3. **Traces** — request flow across services (OpenTelemetry)

---

## 2. Structured Logging

### 2.1 Log Format (JSON)

All services must output logs in JSON format (one log entry per line). No plain text.

**Required Fields:**

| Field | Type | Description | Example |
|---|---|---|---|
| `timestamp` | string (RFC3339) | When event occurred | `"2026-05-07T00:30:15.123Z"` |
| `level` | string | Log level | `"info"`, `"error"`, `"warn"`, `"debug"` |
| `service` | string | Service name | `"auth-service"` |
| `trace_id` | string | OpenTelemetry trace ID (hex) | `"4bf92f3577b34da6a3ce929d0e0e4736"` |
| `span_id` | string | OpenTelemetry span ID (hex) | `"00f067aa0ba902b7"` |
| `user_id` | string (optional) | User performing action | `"550e8400-e29b-41d4-a716-446655440000"` |
| `action` | string | Logical operation | `"login"`, `"transaction_initiate"`, `"wallet_debit"` |
| `duration_ms` | number | Duration of operation (if applicable) | `145.2` |
| `error` | object (optional) | Error details if level=error | `{"code":"INSUFFICIENT_BALANCE","message":"..."}` |
| `details` | object (optional) | Additional context (key-values) | `{"phone":"+6281234567890","device_id":"abc123"}` |
| `version` | int | Log schema version (for future migrations) | `1` |

**Example Log Entry:**
```json
{
  "timestamp": "2026-05-07T00:30:15.123Z",
  "level": "info",
  "service": "transaction-service",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "action": "transaction_initiate",
  "details": {
    "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "product_id": "prod-123",
    "product_name": "XL Data 25GB",
    "selling_price": 27000,
    "customer_no": "+6281234567890",
    "device_fingerprint": "sha256:abc...",
    "idempotency_key": "idx-456"
  },
  "duration_ms": 145.2,
  "version": 1
}
```

### 2.2 PII Redaction

**Never log in plain text:**
- `password`
- `pin`
- `phone_number` → mask to `+62 812-345-***-****` or `+62***XXXX`
- `customer_no` → mask to last 4 digits: `****9012`
- `account_no` → mask last 4 digits: `****5678`
- `ref_id` → keep full (not PII, but treat as sensitive; hash in logs? Keep full for debugging)

**Redaction Middleware (Go):**
```go
func RedactLogFields(entry map[string]interface{}) map[string]interface{} {
    redactPatterns := []struct{
        Key string
        Mask func(string) string
    }{
        {"password", func(s string) string { return "***" }},
        {"pin", func(s string) string { return "***" }},
        {"phone_number", maskPhone},
        {"customer_no", maskCustomerNo},
        {"account_no", maskAccount},
    }

    for _, pattern := range redactPatterns {
        if val, ok := entry[pattern.Key]; ok {
            if str, ok := val.(string); ok {
                entry[pattern.Key] = pattern.Mask(str)
            }
        }
    }
    return entry
}
```

### 2.3 Log Levels

| Level | Use Case | Should Alert? |
|---|---|---|
| `debug` | Detailed development info, request/response bodies (dev only) | No |
| `info` | Normal operations: transaction initiated, success, staff added | No |
| `warn` | Recoverable errors: retry, rate limit hit, stale data, timeout | No (monitor trends) |
| `error` | Failed operation: transaction failed, wallet debit error, webhook processing error | Yes (if > threshold) |
| `fatal` | Service crashing: panic, DB connection lost | Yes (immediate) |

**Production:** Set log level to `info`. Enable `debug` only in dev or on-demand via config reload.

---

## 3. Metrics (Prometheus Format)

All services expose `/metrics` endpoint in Prometheus text format.

### 3.1 Service-Level Metrics

**HTTP Metrics (use Prometheus client middleware):**
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"service", "method", "path", "status"},
    )
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request latency",
            Buckets: prometheus.DefBuckets,
        },
        []string{"service", "method", "path"},
    )
)
```

**Business Metrics:**

| Metric | Type | Labels | Description |
|---|---|---|---|
| `transactions_total` | Counter | `service, status` (success/failed/pending/expired) | Total transactions processed |
| `transactions_amount_total` | Counter | `service, currency` | Total monetary value of transactions (IDR) |
| `wallet_debits_total` | Counter | `service, role` (Mitra/Staff) | Number of wallet debits |
| `wallet_credits_total` | Counter | `service, source` (topup/refund/commission) | Number of wallet credits |
| `wallet_balance_usd`{wallet_id} | Gauge | `wallet_id, role` | Current wallet balance (IDR) |
| `staff_active_count` | Gauge | `mitra_id` | Number of active staff per Mitra |
| `products_active_count` | Gauge | — | Number of active products in catalog |
| `digiflazz_api_calls_total` | Counter | `service, endpoint, status` | Calls to Digiflazz API |
| `digiflazz_api_duration_seconds` | Histogram | `service, endpoint` | Latency of Digiflazz API calls |
| `webhook_received_total` | Counter | `service, event_type` (create/update) | Webhooks from Digiflazz |
| `webhook_processing_duration_seconds` | Histogram | `service` | Time to process webhook |
| `rate_limit_exceeded_total` | Counter | `service, endpoint` | Number of 429 responses |
| `compensation_jobs_retried_total` | Counter | `job_type` | Compensation retry count |
| `audit_log_entries_total` | Counter | `service, action` | Audit log entries created |

### 3.2 SLI/SLO Definitions

| Service | SLI (Service Level Indicator) | SLO Target | Measurement |
|---|---|---|---|
| Auth Service | Login success rate (2xx / total) | 99.9% | `rate(http_requests_total{status=~"2.."}[5m])` |
| Transaction Service | Transaction success rate (final status Success) | 99.5% within 5min | `rate(transactions_total{status="success"}[5m]) / rate(transactions_total[5m])` |
| Wallet Service | Wallet debit/credit latency | p99 < 100ms | `histogram_quantile(0.99, wallet_operation_duration_seconds)` |
| Product Service | Product list API latency | p95 < 200ms | `histogram_quantile(0.95, http_request_duration_seconds{path="/products"}[5m])` |
| Integration Service | Webhook processing latency | p99 < 500ms | `histogram_quantile(0.99, webhook_processing_duration_seconds)` |
| System | End-to-end transaction time (initiate → success) | p95 < 30s | `transaction_init_to_success_duration_seconds` |

---

## 4. Distributed Tracing (OpenTelemetry)

### 4.1 Trace Context Propagation

**Trace ID:** 32-hex-character UUID (16 bytes)  
**Span ID:** 16-hex-character (8 bytes)

**Propagation via HTTP Headers:**
```
Traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
Tracestate: (optional)
```

**Instrumentation:**
- Auto-instrument Go HTTP servers/clients using `go.opentelemetry.io/otel`
- Add span attributes:
  - `user.id`
  - `user.role`
  - `transaction.id`
  - `product.id`
  - `wallet.id`
- No sensitive PII in spans (redact phone numbers)

### 4.2 Trace Sampling

**Strategy:** Parent-based trace sampling (keep all traces for error requests, sample 1% of success).

```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

exporters:
  jaeger:
    endpoint: jaeger-collector:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [jaeger]
```

**Sampling Rule:**  
- Errors (span status = error): 100%  
- Success: 1% (probabilistic)  
- Health checks: 0%

Estimated traces: 100 TPS × 3600 × 24 × 0.01 = ~86,400 spans/day (1%) — manageable.

---

## 5. Alerting Rules (Prometheus Alertmanager)

**File:** `alerting_rules.yml`

```yaml
groups:
- name: ppob_critical
  rules:
  - alert: HighErrorRate
    expr: |
      sum(rate(http_requests_total{status=~"5.."}[5m])) 
      / 
      sum(rate(http_requests_total[5m])) > 0.05
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "High 5xx error rate ({{ $value | humanizePercentage }})"
      description: "Service {{ $labels.service }} has >5% 5xx errors"

  - alert: DigiflazzAPIUnavailable
    expr: |
      rate(digiflazz_api_calls_total{status="error"}[5m]) 
      / 
      rate(digiflazz_api_calls_total[5m]) > 0.1
    for: 3m
    labels:
      severity: critical
    annotations:
      summary: "Digiflazz API error rate high ({{ $value | humanizePercentage }})"
      description: "Check Digiflazz status page"

  - alert: PendingTransactionsStuck
    expr: |
      sum(transactions_total{status="pending"}) 
      by (service) > 100
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Many pending transactions ({{ $value }})"
      description: "Webhook processing backlog detected"

  - alert: WalletDebitFailure
    expr: |
      rate(wallet_operation_errors_total{operation="debit"}[5m]) > 0.01
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "Wallet debit failures elevated"
      description: "Check wallet service logs for errors"

  - alert: DailyLimitHitRateHigh
    expr: |
      rate(rate_limit_exceeded_total{endpoint="/transactions/initiate"}[5m]) > 0.1
    for: 3m
    labels:
      severity: warning
    annotations:
      summary: "Many users hitting daily transaction limit"
      description: "Consider raising limits or investigate abuse"

  - alert: DatabaseConnectionsHigh
    expr: |
      pg_stat_activity_count{datname="ppob"} > 150
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "High DB connection count ({{ $value }})"
      description: "Connection pool size may need adjustment"

  - alert: CompensationJobFailed
    expr: |
      increase(compensation_jobs_retried_total[10m]) > 0
    labels:
      severity: critical
    annotations:
      summary: "Compensation job failed ({{ $value }})"
      description: "Manual intervention required"

  - alert: ReconciliationDrift
    expr: |
      wallet_reconciliation_drift_amount{>10000}
    for: 1h
    labels:
      severity: warning
    annotations:
      summary: "Wallet balance drift detected >Rp10,000"
      description: "Run manual reconciliation"

- name: ppob_warnings
  rules:
  - alert: SlowDigiflazzLatency
    expr: |
      histogram_quantile(0.95, rate(digiflazz_api_duration_seconds_bucket[5m])) > 2
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Digiflazz p95 latency >2s ({{ $value | humanizeDuration }})"

  - alert: HighHoldBalance
    expr: |
      sum(wallet_balance_held_amount) > 100000000
    for: 10m
    labels:
      severity: info
    annotations:
      summary: "Large amount locked in holds (Rp {{ $value | humanize }})"
```

**Alert Routing (Alertmanager config):**
- `critical` → PagerDuty (immediate page)
- `warning` → Slack #alerts-ppob
- `info` → Slack #ppob-monitoring

---

## 6. Dashboards (Grafana)

### 6.1. Service Overview Dashboard

**Panels:**
1. **Request Rate** (graph) — all services, by status code
2. **Error Rate** (stat) — 5xx % last 5min
3. **p50/p95/p99 Latency** (graph) — per service
4. **Active Connections** (gauge) — DB, Redis
5. **JWT Cache Hit Rate** (stat) — Redis hit/miss

### 6.2. Transaction Dashboard

**Panels:**
1. **Transactions per minute** (time series) — by status (success/failed/pending)
2. **Success Rate** (gauge) — last 1h
3. **Average Transaction Value** (stat) — IDR
4. **Top Products** (bar) — by count
5. **Pending Age** (histogram) — age distribution of pending txns
6. **Commission Earned Today** (stat) — total staff commission

### 6.3. Wallet Dashboard

**Panels:**
1. **Total Wallet Balance** (stat) — sum of all wallets
2. **Held Balance** (stat) — amount currently on hold (pending transactions)
3. **Top-up Volume** (graph) — daily top-up total
4. **Wallet Debit Errors** (alert list) — last 10 failures
5. **Reconciliation Drift** (gauge) — current mismatch amount

### 6.4. Digiflazz Health Dashboard

**Panels:**
1. **API Response Time** (graph) — p50/p95/p99
2. **Success Rate** (stat) — last 1h
3. **RC Breakdown** (pie) — distribution of response codes
4. **Webhook Delivery Lag** (graph) — time from transaction to webhook (should be <30s)
5. **Pending Webhooks Queue** (gauge) — count in DLQ

### 6.5. Infrastructure Dashboard

**Panels:**
1. **CPU / Memory** (graphs) — per node, per service
2. **Pod Restarts** (alert list) — recent crash loops
3. **DB Connections** (gauge) — used vs max
4. **Redis Memory** (gauge) — used vs max
5. **Disk Usage** (graph) — PVC usage

---

## 7. Correlation & Debugging

**Trace ID in Logs:** Every structured log entry includes `trace_id`. In Grafana, you can:
- Search logs by trace_id across all services
- Jump from a log line to the corresponding trace in Jaeger
- Correlate metric spikes with specific trace examples

**Correlation Example:**
1. Alert: `HighErrorRate` on transaction-service
2. Query logs: `{service="transaction-service", level="error"} | json | trace_id="abc123"`
3. Find trace_id `abc123` → search in Jaeger → see full trace: auth → wallet → digiflazz
4. Identify failing span (Digiflazz call) → root cause

---

## 8. Log Aggregation & Retention

**Stack:** Fluent Bit → OpenSearch (or Loki)

**Fluent Bit Config:**
```conf
[INPUT]
    Name tail
    Path /var/log/containers/*.log
    Parser docker
    Tag kubernetes.*

[FILTER]
    Name    kubernetes
    Match   kubernetes.*
    Keep_Log   true
    Annotations  On

[FILTER]
    Name    modify
    Match   *
    Add level $${?record["level"]}
    Remove level

[OUTPUT]
    Name   es
    Match  *
    Host   opensearch
    Port   9200
    Index  ppob-logs-%Y.%m.%d
    Type   _doc
```

**Retention:**
- **Index per day:** `ppob-logs-2026.05.07`
- **Hot tier (30 days):** SSD storage, fast search
- **Warm tier (90 days):** HDD storage, slower queries
- **Cold tier (5 years):** Glacier/S3 for compliance; queryable via S3 Select (rare)

---

## 9. Health Checks

Each service exposes `/health/live` (liveness) and `/health/ready` (readiness).

### 9.1 Liveness Probe (`/health/live`)
Checks if service process is alive. Never checks dependencies.

```go
func livenessHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]Byte(`{"status":"alive"}`))
}
```

### 9.2 Readiness Probe (`/health/ready`)
Checks if service can accept traffic. Verifies:
- Database connection
- Redis connection
- Vault connectivity (for secret fetch)
- Circuit breakers (if Digiflazz down, still ready? Yes, but degraded)

```go
func readinessHandler(w http.ResponseWriter, r *http.Request) {
    checks := map[string]bool{
        "database":   db.PingContext(r.Context()) == nil,
        "redis":      redis.Ping() == nil,
        "vault":      vault.HealthCheck() == nil,
    }

    allOK := true
    for _, ok := range checks {
        if !ok {
            allOK = false
            break
        }
    }

    if allOK {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(checks)
    }
}
```

### 9.3 Startup Probe (K8s)
Delay readiness check for 30s after container start (allows time for Vault connection, DB pool init).

---

## 10. Synthetic Monitoring

**Uptime Checks (external):**
- External service (UptimeRobot, Pingdom) hits `/health/ready` every 1min from 3 global regions
- Alert if >2 consecutive failures
- Also test full transaction flow via synthetic transaction:
  ```
  GET /health/synthetic
  → triggers: register test user → login → initiate test transaction → verify success
  ```
  (runs in isolated test environment, not prod)

---

## 11. Error Tracking (Sentry)

**Integration:** Sentry SDK in all services.

**Before Send Hook:** Redact PII from error events:
```go
sentry.CurrentHub().ConfigureScope(func(scope *sentry.Scope) {
    scope.SetExtra("phone_number", mask(original))
    scope.SetExtra("customer_no", mask(original))
})
```

**Fingerprinting:** Group similar errors (e.g., `INSUFFICIENT_BALANCE` errors same fingerprint).

**Alerting:** Sentry sends to Slack #sentry-ppob on new error spikes (>10 occurrences in 1min).

---

## 12. Log Query Examples (for Developers)

**Find transaction by ID:**
```json
{
  "query": "transaction_id:a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "sort": "@timestamp desc"
}
```

**Find all wallet debits for a user:**
```json
{
  "query": "user_id:550e8400-e29b-41d4-a716-446655440000 AND action:wallet_debit",
  "sort": "@timestamp desc"
}
```

**Find errors in last 1 hour:**
```json
{
  "query": "level:error AND timestamp:[now-1h TO now]",
  "sort": "@timestamp desc"
}
```

**Find all traces for a user's transaction:**
```json
{
  "query": "user_id:550e8400-e29b-41d4-a716-446655440000 AND trace_id:abc123",
  "sort": "@timestamp asc"
}
```

---

## 13. Metrics Collection Code Snippet

```go
package metrics

import (
    "time"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    TransactionInitiateDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "transaction_initiate_duration_seconds",
            Help:    "Time to process transaction initiation",
            Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
        },
        []string{"status"},
    )

    WalletDebitDuration = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "wallet_debit_duration_seconds",
            Help:    "Time to debit wallet",
            Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1},
        },
    )
)

func RecordTransactionInitiate(start time.Time, status string) {
    duration := time.Since(start).Seconds()
    TransactionInitiateDuration.WithLabelValues(status).Observe(duration)
}
```

---

## 14. Alert Escalation Policy

| Severity | Response Time | Notification Channel | Escalation |
|---|---|---|---|
| Critical | 5 minutes | PagerDuty (SRE on-call) | If no ack in 15min → secondary on-call |
| Warning | 1 hour | Slack #alerts-ppob | Daily digest if unresolved 24h |
| Info | 24 hours | Slack #ppob-monitoring | Weekly review |

---

## 15. SLO Reporting

**Monthly Report to Stakeholders:**
- Availability (uptime) per service
- Error budget burn rate
- Top 5 incident post-mortems (if any)
- Performance trends (latency, throughput)
- Business metrics: transaction volume, success rate, avg transaction value

---

## Appendix A — OpenTelemetry Instrumentation Go Example

```go
package main

import (
    "go.opentelemetry.io/otel"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/trace"
)

func initTracer() {
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIdRatioBased(0.01))),
        sdktrace.WithBatcher(sdktrace.NewBatchSpanProcessor(exporter)),
    )
    otel.SetTracerProvider(tp)
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
    ctx, span := otel.Tracer("transaction").Start(r.Context(), "handleTransaction")
    defer span.End()

    // attach user_id to span
    span.SetAttributes(
        attribute.String("user.id", userID),
        attribute.String("transaction.id", txID),
    )

    // business logic
}
```

---

## Appendix B — Grafana Dashboard JSON Templates

Dashboard definitions in JSON to import are available in `dashboards/` folder (to be generated separately).

---

## Appendix C — Retention & Compliance

- **Operational logs (INFO/ERROR):** 90 days in hot storage, 2 years in warm
- **Audit logs:** 5 years (legal requirement) in cold storage (immutable, WORM if possible)
- **Metrics:** 30 days at 1s resolution, 1 year at 1min resolution, 5 years at 1h resolution (downsampled)

---

**Next Actions:**
1. Deploy Prometheus + Grafana + Alertmanager to staging
2. Instrument all services with middleware for logs/metrics/traces
3. Create dashboard templates (import JSON)
4. Configure alert routing (PagerDuty, Slack)
5. Set up log aggregation pipeline (Fluent Bit → OpenSearch)
6. Define runbooks for each alert (linked from alert annotation `runbook_url`)

---

**Owner:** SRE Team  
**Runbooks:** See `RUNBOOKS.md` for detailed remediation steps per alert
