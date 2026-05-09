# PPOB Deployment Guide

This directory contains separate Docker Compose files for infrastructure and each microservice.

## Files

- `infrastructure.yml` - Contains shared services (PostgreSQL, Redis, Jaeger, Prometheus, Grafana, Kong)
- `auth-service.yml` - Authentication service
- `user-service.yml` - User management service
- `wallet-service.yml` - Wallet and balance service
- `transaction-service.yml` - Transaction processing service
- `product-service.yml` - Product catalog service
- `integration-service.yml` - External provider integrations

## Networking

All services communicate over a shared Docker network called `ppob-internal`. The infrastructure creates this network, and each service references it as an external network.

## Usage

### 1. Start Infrastructure First
```bash
docker-compose -f infrastructure.yml up -d
```

This will start:
- PostgreSQL (ppob-postgres:5432)
- Redis (ppob-redis:6379)
- Jaeger (ppob-jaeger:16686 for UI)
- Prometheus (ppob-prometheus:9090)
- Grafana (ppob-grafana:3000)
- Kong API Gateway (ppob-kong:8000/8443 for proxy, 8001/8444 for admin)

### 2. Start Individual Services
Start only the services you need:
```bash
docker-compose -f auth-service.yml up -d
docker-compose -f user-service.yml up -d
# ... etc for other services
```

### 3. Service Ports
Each service exposes port 8080 internally, mapped to:
- Auth Service: 8081
- User Service: 8082
- Wallet Service: 8083
- Transaction Service: 8084
- Product Service: 8085
- Integration Service: 8086

### 4. Kong Configuration
After starting Kong, you need to configure routes via the Admin API (localhost:8001):

Example for auth service:
```bash
# Create service
curl -i -X POST http://localhost:8001/services/auth-service \
  --data name=auth-service \
  --data url=http://auth-service:8080

# Create route
curl -i -X POST http://localhost:8001/services/auth-service/routes \
  --data paths[]=/api/v1/auth
```

Repeat for other services with appropriate paths.

### 5. Environment Variables
Service-specific environment variables can be added to each service's compose file. Common variables:
- Database connection: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
- Redis connection: REDIS_HOST, REDIS_PORT
- Service mode: GIN_MODE (debug/release)
- Tracing: JaegerEndpoint (for services that use OpenTelemetry)

### 6. Stopping Services
```bash
# Stop specific service
docker-compose -f auth-service.yml down

# Stop all services and infrastructure
docker-compose -f infrastructure.yml down
docker-compose -f auth-service.yml down
# ... etc for all service files
```

## Benefits of This Approach

1. **Independent Scaling**: Scale services based on individual demand
2. **Fault Isolation**: Issues in one service don't affect others
3. **Selective Deployment**: Only start needed services for development/testing
4. **Clear Separation**: Infrastructure concerns vs. business logic
5. **Optimized Networking**: Direct service-to-service communication over Docker's internal network
6. **Standard Entry Point**: All external traffic through Kong API Gateway

## Notes

- Ensure you're running these commands from the `deployment` directory
- The `context` paths in service compose files point to `../service-name` assuming this deployment folder is at the same level as service folders
- Update environment variables (especially secrets) as needed for your environment
- For production, consider using Docker secrets or external secret management