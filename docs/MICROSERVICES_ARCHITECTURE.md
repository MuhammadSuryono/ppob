# PPOB Microservices Architecture

## Overview

This document describes the microservices architecture for the PPOB (Pulsa, Pasca Bayar, dan Pembayaran) system. The monolithic application has been refactored into 6 independent microservices following domain-driven design principles.

## Architecture Principles

1. **Domain-Driven Design**: Each service owns its data and business logic
2. **Database per Service**: Each service has its own database (shared initially, will be split)
3. **Event-Driven Communication**: Services communicate via events for eventual consistency
4. **API Gateway**: Single entry point for external clients
5. **Resilience Patterns**: Circuit breakers, retries, and timeouts
6. **Observability**: Distributed tracing, centralized logging, metrics

## Service Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   API Gateway   │────│  Auth Service    │────│  User Service   │
│   (Kong/NGINX)  │    │  Port: 8081      │    │  Port: 8082     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         │                       ▼                       ▼
         │                ┌─────────────────┐    ┌─────────────────┐
         │                │ Wallet Service  │────│ Transaction     │
         │                │ Port: 8083      │    │ Service         │
         │                └─────────────────┘    │ Port: 8084      │
         │                       │               └─────────────────┘
         ▼                       ▼                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Product Service │    │ Integration     │    │ Notification    │
│ Port: 8085      │    │ Service         │    │ Service         │
└─────────────────┘    │ Port: 8086      │    │ (Future)        │
                        └─────────────────┘    └─────────────────┘
```

## Services

### 1. Auth Service (Port 8081)
**Responsibilities:**
- User registration and authentication
- JWT token generation and validation
- OTP generation and verification
- Device fingerprinting
- Session management

**Endpoints:**
- `POST /auth/register` - Register new user
- `POST /auth/login` - User login
- `POST /auth/verify-otp` - Verify OTP
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout

**Database:** Owns `users`, `otp_codes`, `device_fingerprints` tables

**Events Published:**
- `user.registered` - When new user registers
- `user.logged_in` - When user logs in
- `user.logged_out` - When user logs out

### 2. User Service (Port 8082)
**Responsibilities:**
- User profile management
- Role and permission management
- Staff management
- User status management

**Endpoints:**
- `GET /users/{id}` - Get user profile
- `PUT /users/{id}` - Update user profile
- `GET /users/{id}/roles` - Get user roles
- `POST /users/{id}/roles` - Assign role to user

**Database:** Owns `user_roles`, `roles` tables

**Events Published:**
- `user.updated` - When user profile is updated
- `user.role_changed` - When user role changes

**Events Consumed:**
- `user.registered` - From Auth Service

### 3. Wallet Service (Port 8083)
**Responsibilities:**
- Wallet balance management
- Balance holds and releases
- Transaction processing
- Wallet event sourcing
- Commission calculations

**Endpoints:**
- `GET /wallets/{id}/balance` - Get wallet balance
- `POST /wallets/{id}/hold` - Place hold on balance
- `POST /wallets/{id}/debit` - Debit wallet
- `POST /wallets/{id}/credit` - Credit wallet
- `POST /wallets/{id}/release-hold` - Release hold

**Database:** Owns `wallets`, `wallet_events` tables

**Events Published:**
- `wallet.credited` - When wallet is credited
- `wallet.debited` - When wallet is debited
- `wallet.hold_placed` - When hold is placed
- `wallet.hold_released` - When hold is released

**Events Consumed:**
- `transaction.created` - From Transaction Service
- `transaction.completed` - From Transaction Service

### 4. Transaction Service (Port 8084)
**Responsibilities:**
- Transaction initiation and processing
- Digiflazz integration
- Transaction status management
- Saga orchestration
- Margin calculations

**Endpoints:**
- `POST /transactions` - Create new transaction
- `GET /transactions/{id}` - Get transaction details
- `POST /transactions/{id}/status` - Update transaction status
- `GET /transactions` - List transactions with filters

**Database:** Owns `transactions`, `daily_limits` tables

**Events Published:**
- `transaction.created` - When transaction is created
- `transaction.completed` - When transaction completes
- `transaction.failed` - When transaction fails

**Events Consumed:**
- `wallet.hold_placed` - From Wallet Service
- `product.updated` - From Product Service

**Saga Pattern:** Orchestration-based saga for transaction processing:
1. Validate product (Product Service)
2. Place hold (Wallet Service)
3. Execute transaction (Integration Service)
4. Debit wallet (Wallet Service)
5. Create commission (Transaction Service)

### 5. Product Service (Port 8085)
**Responsibilities:**
- Product catalog management
- Category management
- Pricing management
- Digiflazz sync
- Product search

**Endpoints:**
- `GET /products` - List products
- `GET /products/{id}` - Get product details
- `POST /products` - Create product
- `PUT /products/{id}` - Update product
- `GET /products/search` - Search products

**Database:** Owns `products`, `categories` tables

**Events Published:**
- `product.created` - When product is created
- `product.updated` - When product is updated
- `product.deleted` - When product is deleted

**Events Consumed:** None (read-only for other services)

### 6. Integration Service (Port 8086)
**Responsibilities:**
- Digiflazz API integration
- Webhook processing
- External provider management
- Retry and compensation logic
- Rate limiting for external APIs

**Endpoints:**
- `POST /integrations/digiflazz/transaction` - Initiate Digiflazz transaction
- `POST /integrations/digiflazz/callback` - Digiflazz webhook
- `GET /integrations/providers` - List available providers

**Database:** Owns `integration_logs`, `provider_configs` tables

**Events Published:**
- `integration.completed` - When integration completes
- `integration.failed` - When integration fails

**Events Consumed:**
- `transaction.created` - From Transaction Service

## Communication Patterns

### 1. Synchronous Communication (gRPC)
Used for request-response patterns where immediate response is needed.

```protobuf
// wallet_service.proto
service WalletService {
    rpc PlaceHold (HoldRequest) returns (HoldResponse);
    rpc ReleaseHold (HoldRequest) returns (ReleaseResponse);
    rpc GetBalance (BalanceRequest) returns (BalanceResponse);
}
```

**Benefits:**
- Strong consistency
- Immediate feedback
- Type-safe contracts

**Use Cases:**
- Balance checks
- Hold placement
- User validation

### 2. Asynchronous Communication (Events)
Used for eventual consistency and decoupled processing.

**Event Bus:** Redis Streams

```
Transaction Service --"TransactionCreated"--> Redis Stream
                                                      |
                                                      v
                                               Wallet Service
                                               (Consumer)
```

**Benefits:**
- Loose coupling
- Scalability
- Resilience (services can be down)

**Use Cases:**
- Transaction processing
- Commission calculations
- Notifications
- Audit logging

## Data Management

### Database Strategy (Phased Approach)

**Phase 1 (Current): Shared Database, Separate Schemas**
- All services use same PostgreSQL instance
- Each service owns its tables
- No cross-service foreign keys
- Services access other data via APIs, not DB joins

**Phase 2 (Future): Database per Service**
- Each service gets its own database
- Data replication via CDC
- Event-driven synchronization

### Handling Shared Data

**Problem:** `user_roles` used by Auth and User services

**Solution:** Event-driven cache invalidation
```
User Service (owns user_roles)
    │
    └── publishes "RoleChanged" event
        │
        ▼
Auth Service (has user_roles_cache)
    └── updates local cache
```

### Transaction Coordination

**Pattern:** Saga with Orchestration

```go
type TransactionSaga struct {
    steps []SagaStep
}

func (s *TransactionSaga) Execute(ctx context.Context) error {
    for i, step := range s.steps {
        if err := step.Execute(ctx); err != nil {
            // Execute compensation
            for j := i - 1; j >= 0; j-- {
                s.steps[j].Compensate(ctx)
            }
            return err
        }
    }
    return nil
}
```

**Saga Steps:**
1. `ValidateProductStep` - Call Product Service
2. `PlaceHoldStep` - Call Wallet Service
3. `ExecuteTransactionStep` - Call Integration Service
4. `DebitWalletStep` - Call Wallet Service (compensation: ReleaseHold)
5. `CreateCommissionStep` - Internal

## Resilience Patterns

### 1. Circuit Breaker
```go
var cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "WalletService",
    MaxRequests: 5,
    Interval:    1 * time.Minute,
    Timeout:     30 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.Requests >= 20 && 
               float64(counts.TotalFailures)/float64(counts.Requests) > 0.5
    },
})
```

**Benefits:**
- Prevents cascading failures
- Fast failure when service is down
- Automatic recovery

### 2. Retry with Backoff
```go
func withRetry[T any](op func() (T, error), maxRetries int) (T, error) {
    delay := 100 * time.Millisecond
    for i := 0; i <= maxRetries; i++ {
        result, err := op()
        if err == nil {
            return result, nil
        }
        time.Sleep(delay)
        delay *= 2 // exponential backoff
    }
    return zero, err
}
```

**Benefits:**
- Handles transient failures
- Configurable retry policy

### 3. Timeout
All service calls have context-based timeouts:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

balance, err := walletClient.GetBalance(ctx, userID)
```

## Observability

### 1. Distributed Tracing
**Implementation:** OpenTelemetry with Jaeger

```go
func initTracer() {
    exp, _ := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ))
    
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
    )
    otel.SetTracerProvider(tp)
}
```

**Trace Propagation:**
```http
GET /api/v1/transactions
Traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
```

### 2. Structured Logging
JSON format with correlation IDs:
```json
{
    "timestamp": "2026-05-07T08:47:27Z",
    "level": "INFO",
    "service": "transaction-service",
    "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
    "message": "Transaction completed",
    "transaction_id": "abc-123"
}
```

### 3. Metrics
Prometheus metrics:
- Request count by endpoint
- Latency histograms
- Error rates
- Circuit breaker state

## Deployment

### Kubernetes Architecture
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    spec:
      containers:
      - name: auth
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
```

### Service Mesh (Istio)
```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: wallet-service
spec:
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        maxRequestsPerConnection: 2
    outlierDetection:
      consecutive5xxErrors: 5
      interval: 30s
```

## Security

### 1. Authentication
- JWT with RS256 signing
- Token expiration: 15 minutes (access), 7 days (refresh)
- Token introspection endpoint

### 2. Authorization
- Role-based access control (RBAC)
- Service-to-service authentication via mTLS

### 3. Network Security
- Service mesh with Istio
- mTLS for internal communication
- Network policies

### 4. Data Protection
- Password hashing with bcrypt
- PIN encryption
- Database encryption at rest

## Monitoring & Alerting

### Key Metrics
- Error rate < 0.1%
- P95 latency < 500ms
- Availability > 99.9%
- Transaction success rate > 99.5%

### Alerting Rules
```yaml
groups:
- name: transaction-alerts
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 5m
```

## Rollback Strategy

### Partial Rollback (per service)
1. Update API Gateway routing
2. Scale down failed service
3. Analyze logs
4. Fix and redeploy

### Full Rollback
1. Route all traffic to monolith
2. Disable all microservices
3. Stop dual-write
4. Restore DB from backup if needed

## Migration Plan

### Phase 1: Preparation (Week 1)
- Document current state
- Define service boundaries
- Setup infrastructure

### Phase 2: Service Extraction (Weeks 2-4)
- Extract Auth Service
- Extract User Service
- Extract Wallet Service
- Extract Transaction Service

### Phase 3: Integration (Weeks 4-5)
- Implement gRPC communication
- Setup event bus
- Add circuit breakers

### Phase 4: Data Migration (Week 6)
- Split shared database
- Remove foreign keys
- Add audit columns

### Phase 5: Deployment (Week 7)
- Kubernetes deployment
- Service mesh setup
- Monitoring setup

### Phase 6: Cutover (Week 8)
- Blue-green deployment
- Gradual traffic shift
- Final verification

## Benefits

1. **Scalability**: Scale services independently based on load
2. **Resilience**: Failure isolation prevents cascading failures
3. **Deployment**: Independent deployment and rollback
4. **Technology**: Freedom to choose best tech per service
5. **Team Autonomy**: Teams can own services end-to-end
6. **Performance**: Optimize each service for its workload

## Trade-offs

1. **Complexity**: Distributed systems are harder to debug
2. **Consistency**: Eventual consistency instead of strong consistency
3. **Latency**: Network overhead for inter-service calls
4. **Data**: Distributed transactions are challenging
5. **Testing**: Integration testing is more complex
