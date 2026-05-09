# PPOB Monolith to Microservices Migration Guide

## Overview

This guide documents the step-by-step migration from a monolithic Go application to a microservices architecture following the Strangler Fig Pattern.

## Current Architecture (Monolith)

```
Monolithic Application (v1/backend/)
├── Auth Module
├── User Module
├── Wallet Module
├── Transaction Module
├── Product Module
└── Integration Module (Digiflazz)
    └── Single PostgreSQL Database
    └── Single Codebase
    └── Shared Dependencies
```

## Target Architecture (Microservices)

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │    │   User Service  │    │ Wallet Service  │
│   Port: 8081    │    │   Port: 8082    │    │   Port: 8083    │
└────────┬────────┘    └────────┬────────┘    └────────┬────────┘
         │                       │                       │
    ┌────┴────┐             ┌─────┴─────┐           ┌─────┴─────┐
    │  gRPC   │             │   gRPC    │           │   gRPC    │
    │  Events │◄────────────┤  Events   │◄──────────┤  Events   │
    └────┬────┘             └─────┬─────┘           └─────┬─────┘
         │                       │                       │
┌────────┴─────────┐    ┌────────┴─────────┐    ┌────────┴─────────┐
│ Transaction Svc  │    │ Product Service  │    │ Integration Svc  │
│     Port: 8084    │    │     Port: 8085    │    │     Port: 8086    │
└──────────────────┘    └──────────────────┘    └──────────────────┘
         │                       │                       │
         └───────────┬───────────┼───────────┬───────────┘
                     │           │
              ┌──────┴──────┐   │
              │   PostgreSQL│   │
              │   (Shared)  │◄──┘
              └─────────────┘
```

## Migration Strategy

### Phase 1: Preparation (Week 1)

**Objectives:**
- [x] Document current monolith structure
- [x] Define service boundaries
- [x] Create migration plan
- [x] Setup infrastructure (K8s, Redis, PostgreSQL)
- [x] Create service skeletons

**Deliverables:**
- Architecture documentation
- Service ownership matrix
- Database schema mapping
- Infrastructure as Code (Terraform/K8s)

### Phase 2: Service Extraction (Weeks 2-4)

#### Step 1: Extract Auth Service (Week 2)

**Scope:**
- User registration and authentication
- JWT token management
- OTP generation and verification
- Device fingerprinting

**Implementation:**

1. Create new Auth Service directory structure:
```
v1/auth-service/
├── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   └── service/
└── pkg/
    ├── redis/
    └── utils/
```

2. Copy relevant code from monolith:
```bash
cp v1/backend/internal/handler/auth_handler.go v1/auth-service/internal/handler/
cp v1/backend/internal/service/auth_service.go v1/auth-service/internal/service/
cp v1/backend/internal/repository/*_repository.go v1/auth-service/internal/repository/
```

3. Refactor dependencies:
- Remove cross-service dependencies
- Add event publishing to Redis
- Implement gRPC client for User Service

4. Test independently:
```bash
cd v1/auth-service
go test ./...
go run main.go
```

**API Gateway Configuration:**
```yaml
services:
  - name: auth-service
    url: http://auth-service:8081
    routes:
      - name: auth-route
        paths:
          - /api/v1/auth
```

**Verification:**
- [x] Auth service runs independently
- [x] Registration endpoint works
- [x] Login endpoint works
- [x] OTP verification works
- [x] JWT tokens are valid

#### Step 2: Extract User Service (Week 2-3)

**Scope:**
- User profile management
- Role and permission management
- Staff management

**Implementation:**

1. Create User Service structure
2. Extract user-related handlers, services, repositories
3. Implement role management
4. Add gRPC server for profile lookups
5. Subscribe to user.registered events

**Events:**
- Consumes: `user.registered` (from Auth)
- Publishes: `user.updated`, `user.role_changed`

#### Step 3: Extract Wallet Service (Week 3)

**Scope:**
- Wallet balance management
- Balance holds and releases
- Transaction event processing

**Implementation:**

1. Create Wallet Service structure
2. Extract wallet handlers, services, repositories
3. Implement event sourcing for wallet events
4. Add gRPC server for balance operations
5. Subscribe to transaction events

**Concurrency Control:**
```go
func (s *WalletService) PlaceHold(ctx context.Context, req *HoldRequest) error {
    // Use database transaction with SELECT FOR UPDATE
    tx := s.db.WithContext(ctx).Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Lock wallet row
    var wallet model.Wallet
    if err := tx.Set("gorm:query_option", "FOR UPDATE").
        Where("wallet_id = ? AND owner_id = ?", req.WalletID, req.UserID).
        First(&wallet).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // Check balance
    if wallet.BalanceAvailable < req.Amount {
        tx.Rollback()
        return ErrInsufficientBalance
    }
    
    // Update balance
    wallet.BalanceAvailable -= req.Amount
    wallet.BalanceHold += req.Amount
    
    if err := tx.Save(&wallet).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // Create wallet event
    event := model.WalletEvent{
        WalletID:      wallet.WalletID,
        EventType:     "hold_placed",
        Amount:        req.Amount,
        BalanceBefore: wallet.BalanceAvailable + req.Amount,
        BalanceAfter:  wallet.BalanceAvailable,
    }
    
    if err := tx.Create(&event).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit().Error
}
```

#### Step 4: Extract Transaction Service (Week 3-4)

**Scope:**
- Transaction initiation and processing
- Digiflazz integration
- Saga orchestration

**Implementation:**

1. Create Transaction Service structure
2. Implement Saga orchestration pattern
3. Add Digiflazz client with retry logic
4. Subscribe to wallet and product events

**Saga Implementation:**
```go
type TransactionSaga struct {
    steps []SagaStep
}

type SagaStep struct {
    Execute    func(ctx context.Context) error
    Compensate func(ctx context.Context) error
}

func (s *TransactionSaga) Execute(ctx context.Context) error {
    for i, step := range s.steps {
        if err := step.Execute(ctx); err != nil {
            // Compensate in reverse order
            for j := i - 1; j >= 0; j-- {
                s.steps[j].Compensate(ctx)
            }
            return err
        }
    }
    return nil
}

// Usage
saga := NewTransactionSaga(req)
saga.AddStep(
    func(ctx context.Context) error {
        return productClient.Validate(ctx, req.ProductID)
    },
    func(ctx context.Context) error {
        return walletClient.ReleaseHold(ctx, holdID)
    },
)
```

### Phase 3: Communication Layer (Week 4-5)

**gRPC Definitions:**

```protobuf
// proto/wallet.proto
syntax = "proto3";

package wallet;

service WalletService {
    rpc PlaceHold (HoldRequest) returns (HoldResponse);
    rpc ReleaseHold (HoldRequest) returns (ReleaseResponse);
    rpc ConvertHoldToDebit (HoldRequest) returns (DebitResponse);
    rpc GetBalance (BalanceRequest) returns (BalanceResponse);
}

message HoldRequest {
    string user_id = 1;
    string wallet_id = 2;
    double amount = 3;
    string reference = 4;
}

message HoldResponse {
    string hold_id = 1;
    double balance_available = 2;
    double balance_hold = 3;
}
```

**Event-Driven Architecture:**

```go
// Publish event
func (s *Service) publishEvent(ctx context.Context, eventType string, data interface{}) error {
    payload, _ := json.Marshal(data)
    return s.redis.XAdd(ctx, &redis.XAddArgs{
        Stream: "events:" + eventType,
        Values: map[string]interface{}{"data": string(payload)},
    }).Err()
}

// Consume events
func (s *Service) consumeEvents(ctx context.Context, eventType string) {
    for {
        streams, err := s.redis.XRead(ctx, &redis.XReadArgs{
            Streams: []string{"events:" + eventType, lastID},
            Block:   0,
        }).Result()
        
        // Process events
    }
}
```

### Phase 4: Data Migration (Week 6)

**Database Refactoring:**

1. **Remove Foreign Keys:**
```sql
-- Remove FK between transactions and wallets
ALTER TABLE transactions 
DROP CONSTRAINT transactions_wallet_id_fkey;

-- Make references optional
ALTER TABLE transactions 
ALTER COLUMN wallet_id DROP NOT NULL;
```

2. **Add Audit Columns:**
```sql
ALTER TABLE transactions 
ADD COLUMN created_by UUID,
ADD COLUMN updated_by UUID,
ADD COLUMN created_from_service VARCHAR(50);
```

3. **Data Duplication (Read Models):**
```sql
-- Create denormalized transaction summary
CREATE TABLE transaction_summaries AS
SELECT 
    t.transaction_id,
    t.user_id,
    u.phone_number,
    t.amount,
    t.status,
    t.created_at
FROM transactions t
JOIN users u ON t.user_id = u.user_id;
```

**Migration Script:**
```go
func migrateData(ctx context.Context) error {
    // 1. Enable dual-write
    enableDualWrite()
    
    // 2. Backfill historical data
    if err := backfillData(ctx); err != nil {
        return err
    }
    
    // 3. Verify data consistency
    if err := verifyData(ctx); err != nil {
        return err
    }
    
    // 4. Switch traffic gradually
    if err := switchTraffic(ctx, 0.05); err != nil { // 5% to new services
        return err
    }
    
    // 5. Monitor and increase
    time.Sleep(24 * time.Hour)
    if err := verifyMetrics(); err != nil {
        rollback()
        return err
    }
    
    // 6. Complete cutover
    return switchTraffic(ctx, 1.0) // 100% to new services
}
```

### Phase 5: Deployment (Week 7)

**Kubernetes Manifests:**

```yaml
# auth-service-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: ppob
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth
        image: ppob/auth-service:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: auth-service-config
        - secretRef:
            name: auth-service-secrets
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
```

**Service Mesh (Istio):**
```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: wallet-service
spec:
  host: wallet-service
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 100
        maxRequestsPerConnection: 2
    outlierDetection:
      consecutive5xxErrors: 5
      interval: 30s
      baseEjectionTime: 30s
```

### Phase 6: Cutover (Week 8)

**Blue-Green Deployment:**

1. Deploy new microservices alongside monolith
2. Configure API Gateway for dual-write
3. Route 5% traffic to new services
4. Monitor for 24 hours
5. Gradually increase to 100%
6. Decommission monolith

**Rollback Plan:**
```bash
#!/bin/bash
# rollback.sh

# Route all traffic to monolith
kubectl patch ingress ppob-ingress -p '
{
  "spec": {
    "rules": [{
      "http": {
        "paths": [{
          "path": "/",
          "backend": {
            "serviceName": "monolith",
            "servicePort": 8080
          }
        }]
      }
    }]
  }
}'

# Scale down microservices
kubectl scale deployment --all -n ppob --replicas=0

echo "Rolled back to monolith"
```

## Verification Checklist

### Pre-Migration
- [x] All services have test coverage > 80%
- [x] Database migrations tested in staging
- [x] Load testing completed
- [x] Disaster recovery plan documented

### Post-Migration
- [ ] Zero data loss verified
- [ ] All endpoints return correct responses
- [ ] Latency within acceptable limits (< 500ms p95)
- [ ] Error rate < 0.1%
- [ ] All events published and consumed correctly
- [ ] Distributed tracing working end-to-end

## Monitoring During Migration

### Key Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| Error Rate | < 0.1% | > 1% |
| P95 Latency | < 500ms | > 1s |
| Database Connections | < 80% | > 90% |
| Queue Depth | < 100 | > 1000 |
| CPU Usage | < 70% | > 85% |

### Dashboards

1. **Service Health**: Uptime, error rates, response times
2. **Database**: Query performance, connection pool, replication lag
3. **Events**: Throughput, processing time, dead letter queue size
4. **Business Metrics**: Transaction success rate, wallet balances

## Common Issues & Solutions

### Issue 1: Distributed Transaction Failures
**Solution:**
- Implement Saga pattern with compensation
- Use idempotent operations
- Add retry logic with exponential backoff

### Issue 2: Event Ordering
**Solution:**
- Use Kafka instead of Redis Streams if ordering critical
- Add sequence numbers to events
- Implement event deduplication

### Issue 3: Network Partitions
**Solution:**
- Circuit breakers to prevent cascading failures
- Timeout configurations
- Fallback to cached data when possible

### Issue 4: Data Inconsistency
**Solution:**
- Regular reconciliation jobs
- Event sourcing for audit trail
- Alerts on data drift

## Performance Optimization

1. **Connection Pooling**: Optimize database connection pools per service
2. **Caching**: Implement Redis caching for read-heavy operations
3. **Async Processing**: Move non-critical operations to background workers
4. **Database Indexes**: Add indexes for common queries
5. **Query Optimization**: Profile and optimize slow queries

## Security Considerations

1. **Service-to-Service Auth**: mTLS with service mesh
2. **API Gateway**: JWT validation, rate limiting
3. **Secrets Management**: Kubernetes Secrets or Vault
4. **Network Policies**: Restrict inter-service communication
5. **Audit Logging**: Log all sensitive operations

## Post-Migration Tasks

1. **Remove Technical Debt**:
   - Delete duplicated code in monolith
   - Remove migration scaffolding
   - Clean up deprecated endpoints

2. **Optimize Resources**:
   - Right-size service instances
   - Optimize database indexes
   - Tune connection pools

3. **Document Everything**:
   - Update architecture diagrams
   - Document runbooks
   - Create troubleshooting guides

4. **Team Training**:
   - Microservices best practices
   - Distributed systems patterns
   - New deployment processes

## Success Criteria

- [ ] All services deployed and running
- [ ] Zero downtime during migration
- [ ] All existing functionality preserved
- [ ] Performance improved or maintained
- [ ] Team comfortable with new architecture
- [ ] Monitoring and alerting in place
- [ ] Documentation complete
- [ ] Disaster recovery tested

## Lessons Learned

1. **Start Simple**: Begin with least critical service
2. **Automate Everything**: CI/CD, testing, deployment
3. **Monitor Aggressively**: You can't fix what you can't see
4. **Communicate Constantly**: Keep stakeholders informed
5. **Plan for Rollback**: Always have an escape route
6. **Document Decisions**: Why you chose certain approaches

## References

- [Strangler Fig Pattern](https://martinfowler.com/bliki/StranglerFigApplication.html)
- [Microservices Patterns](https://microservices.io/patterns/)
- [Saga Pattern](https://microservices.io/patterns/data/saga.html)
- [Event-Driven Architecture](https://www.confluent.io/event-driven-architecture/)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)

## Support

For questions or issues during migration:
- Slack: #microservices-migration
- Email: tech-leads@yontech.com
- On-call: Check PagerDuty rotation