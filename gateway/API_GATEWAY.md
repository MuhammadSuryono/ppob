# API Gateway

The API Gateway serves as the single entry point for all client requests. It handles routing, authentication, rate limiting, and load balancing across microservices.

## Gateway Responsibilities

1. **Request Routing**: Routes incoming requests to appropriate microservices
2. **Authentication & Authorization**: Validates JWT tokens and checks permissions
3. **Rate Limiting**: Prevents abuse and ensures fair usage
4. **Load Balancing**: Distributes requests across service instances
5. **Caching**: Caches responses for improved performance
6. **SSL Termination**: Handles HTTPS/TLS decryption
7. **Request/Response Transformation**: Modifies requests and responses as needed

## Technology Stack

- **Kong API Gateway** (Recommended)
  - Open-source, scalable, and cloud-native
  - Built on top of NGINX
  - Plugin architecture for extensibility
  
- **Alternative: NGINX Plus**
  - High-performance load balancer
  - Advanced routing capabilities
  - Commercial support available

## Kong Configuration

### Basic Setup

```yaml
_format_version: "2.1"
_transform: true

services:
  - name: auth-service
    url: http://auth-service:8081
    routes:
      - name: auth-route
        paths:
          - /api/v1/auth
        strip_path: true
        
  - name: user-service
    url: http://user-service:8082
    routes:
      - name: user-route
        paths:
          - /api/v1/users
        strip_path: true
        plugins:
          - name: jwt
          - name: rate-limiting
            config:
              minute: 100
              hour: 1000
              
  - name: wallet-service
    url: http://wallet-service:8083
    routes:
      - name: wallet-route
        paths:
          - /api/v1/wallets
        strip_path: true
        plugins:
          - name: jwt
          - name: rate-limiting
            config:
              minute: 200
              hour: 2000
              
  - name: transaction-service
    url: http://transaction-service:8084
    routes:
      - name: transaction-route
        paths:
          - /api/v1/transactions
        strip_path: true
        plugins:
          - name: jwt
          - name: rate-limiting
            config:
              minute: 50
              hour: 500
              
  - name: product-service
    url: http://product-service:8085
    routes:
      - name: product-route
        paths:
          - /api/v1/products
          - /api/v1/categories
        strip_path: true
        
  - name: integration-service
    url: http://integration-service:8086
    routes:
      - name: integration-route
        paths:
          - /api/v1/integrations
          - /api/v1/webhooks
        strip_path: true
        plugins:
          - name: key-auth

plugins:
  - name: cors
    config:
      origins: "*"
      methods: "GET,POST,PUT,DELETE,OPTIONS"
      headers: "Content-Type,Authorization"
      exposed_headers: "X-Authenticated-User"
      
  - name: prometheus
    config:
      per_consumer: true
      
  - name: correlation-id
    config:
      header_name: "X-Request-ID"
      generator: "uuid#counter"
      echo_downstream: true
```

### JWT Authentication Plugin

```yaml
- name: jwt
  config:
    key_claim_name: "kid"
    uri_param_names:
      - "jwt"
    claims_to_verify:
      - "exp"
    maximum_expiration: 3600
```

### Rate Limiting Plugin

```yaml
- name: rate-limiting
  config:
    minute: 100
    hour: 1000
    day: 10000
    policy: local
    limit_by: consumer
    fault_tolerant: true
```

### Circuit Breaker Configuration

```yaml
- name: proxy-cache
  config:
    strategy: redis
    content_type:
      - application/json
    cache_ttl: 30
    cache_control: true
```

## NGINX Configuration (Alternative)

```nginx
upstream auth_service {
    server auth-service-1:8081;
    server auth-service-2:8081;
    server auth-service-3:8081;
}

upstream user_service {
    server user-service-1:8082;
    server user-service-2:8082;
    server user-service-3:8082;
}

upstream wallet_service {
    server wallet-service-1:8083;
    server wallet-service-2:8083;
    server wallet-service-3:8083;
}

upstream transaction_service {
    server transaction-service-1:8084;
    server transaction-service-2:8084;
    server transaction-service-3:8084;
}

map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

server {
    listen 80;
    server_name api.ppob.yontech.com;
    
    # SSL Configuration
    listen 443 ssl http2;
    ssl_certificate /etc/nginx/ssl/ppob.crt;
    ssl_certificate_key /etc/nginx/ssl/ppob.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # Logging
    access_log /var/log/nginx/gateway_access.log;
    error_log /var/log/nginx/gateway_error.log;
    
    # Request ID for tracing
    set $request_id $http_x_request_id;
    if ($request_id = "") {
        set $request_id $request_id;
    }
    
    add_header X-Request-ID $request_id always;
    
    # Rate Limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=100r/s;
    limit_conn_zone $binary_remote_addr zone=addr_limit:10m;
    
    # Auth Service Routes
    location /api/v1/auth {
        limit_req zone=api_limit burst=20 nodelay;
        limit_conn addr_limit 10;
        
        proxy_pass http://auth_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
    }
    
    # User Service Routes
    location /api/v1/users {
        limit_req zone=api_limit burst=20 nodelay;
        limit_conn addr_limit 10;
        
        # JWT Authentication
        auth_jwt "Restricted";
        auth_jwt_key_file /etc/nginx/jwt_public_key.pem;
        
        proxy_pass http://user_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        proxy_set_header X-Authenticated-User $remote_user;
        
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
    }
    
    # Wallet Service Routes
    location /api/v1/wallets {
        limit_req zone=api_limit burst=30 nodelay;
        limit_conn addr_limit 10;
        
        auth_jwt "Restricted";
        auth_jwt_key_file /etc/nginx/jwt_public_key.pem;
        
        proxy_pass http://wallet_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        proxy_set_header X-Authenticated-User $remote_user;
        
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
    }
    
    # Transaction Service Routes
    location /api/v1/transactions {
        limit_req zone=api_limit burst=20 nodelay;
        limit_conn addr_limit 10;
        
        auth_jwt "Restricted";
        auth_jwt_key_file /etc/nginx/jwt_public_key.pem;
        
        proxy_pass http://transaction_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        proxy_set_header X-Authenticated-User $remote_user;
        
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
    }
    
    # Product Service Routes
    location /api/v1/products {
        proxy_pass http://product_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
    }
    
    location /api/v1/categories {
        proxy_pass http://product_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
    }
    
    # Integration Service Routes
    location /api/v1/integrations {
        limit_conn addr_limit 5;
        
        auth_jwt "Restricted";
        auth_jwt_key_file /etc/nginx/jwt_public_key.pem;
        
        proxy_pass http://integration_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        proxy_set_header X-Authenticated-User $remote_user;
        
        proxy_connect_timeout 10s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
    
    location /api/v1/webhooks {
        proxy_pass http://integration_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        
        proxy_connect_timeout 30s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Health Check Endpoint
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
    
    # Default route
    location / {
        return 404 '{"error": "Not found"}';
        add_header Content-Type application/json;
    }
}
```

## Load Balancing Strategies

### 1. Round Robin (Default)
- Distributes requests evenly across all instances
- Simple and predictable
- Best for: Homogeneous instances with similar capacity

### 2. Least Connections
- Routes to instance with fewest active connections
- Better for: Long-lived connections or variable request processing times

### 3. IP Hash
- Routes based on client IP address
- Ensures same client goes to same instance
- Best for: Session persistence without sticky sessions

### 4. Weighted Round Robin
- Assigns weights based on instance capacity
- Routes more traffic to powerful instances
- Best for: Heterogeneous infrastructure

## Caching Strategy

### Gateway-Level Caching

```nginx
# Cache product catalog
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=product_cache:10m max_size=1g inactive=60m use_temp_path=off;

location /api/v1/products {
    proxy_cache product_cache;
    proxy_cache_valid 200 5m;
    proxy_cache_valid 404 1m;
    proxy_cache_use_stale error timeout updating;
    proxy_cache_background_update on;
    add_header X-Cache-Status $upstream_cache_status;
    
    proxy_pass http://product_service;
}
```

### Cache Invalidation

1. **Time-based**: TTL expiration
2. **Event-based**: Invalidate on product update events
3. **Manual**: Admin-triggered cache purge

## Security Features

### 1. SSL/TLS Termination
- Offloads encryption/decryption from microservices
- Centralized certificate management
- Supports multiple certificates (SNI)

### 2. WAF Integration
```yaml
- name: modsecurity
  config:
    config:
      rules:
        - /etc/nginx/modsec/main.conf
    mode: "ON"
```

### 3. IP Whitelisting/Blacklisting
```nginx
location /api/v1/admin {
    allow 10.0.0.0/8;
    allow 172.16.0.0/12;
    deny all;
    
    proxy_pass http://admin_service;
}
```

### 4. Bot Detection
- Rate limiting per IP
- User agent validation
- Request pattern analysis

## Monitoring & Observability

### Gateway Metrics

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| Request Rate | Requests per second | < 1000/s |
| Error Rate | 4xx + 5xx responses | > 1% |
| Latency P95 | 95th percentile latency | > 500ms |
| Active Connections | Current connections | > 10000 |
| Cache Hit Rate | Cache effectiveness | < 80% |

### Logging

```nginx
log_format gateway_log '
    $time_iso8601 $remote_addr
    $request_method $request_uri $server_protocol
    $status $body_bytes_sent
    "$http_referer" "$http_user_agent"
    $request_time $upstream_response_time
    $request_id $remote_user
    $http_x_forwarded_for
';

access_log /var/log/gateway.log gateway_log;
```

## High Availability

### Gateway Cluster

```
         ┌─────────────────┐
         │  Load Balancer  │
         │   (HAProxy)     │
         └────────┬────────┘
                  │
    ┌─────────────┼─────────────┐
    │             │             │
┌───▼───┐   ┌───▼───┐   ┌───▼───┐
│Gateway│   │Gateway│   │Gateway│
│   1   │   │   2   │   │   3   │
└───────┘   └───────┘   └───────┘
```

### Health Checks

```yaml
health_checks:
  active:
    http_path: /health
    healthy_threshold: 2
    unhealthy_threshold: 3
    timeout: 2s
    interval: 10s
```

## Performance Tuning

### NGINX Optimization

```nginx
worker_processes auto;
worker_rlimit_nofile 65535;

events {
    worker_connections 4096;
    use epoll;
    multi_accept on;
}

http {
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    keepalive_requests 1000;
    reset_timedout_connection on;
    
    # Buffers
    client_body_buffer_size 128k;
    client_max_body_size 10m;
    client_header_buffer_size 1k;
    large_client_header_buffers 4 8k;
    
    # Timeouts
    client_body_timeout 12;
    client_header_timeout 12;
    send_timeout 10;
}
```

## Deployment Strategy

### Blue-Green Deployment

1. Deploy new gateway version (green)
2. Run smoke tests
3. Switch load balancer to green
4. Monitor for issues
5. Rollback if needed (switch back to blue)

### Canary Release

1. Deploy new version to 10% of traffic
2. Monitor metrics
3. Gradually increase to 100%
4. Automatic rollback on errors

## Cost Optimization

1. **Right-sizing**: Choose appropriate instance types
2. **Auto-scaling**: Scale based on traffic patterns
3. **Reserved Instances**: Commit for discounted pricing
4. **Spot Instances**: Use for non-critical workloads
5. **Efficient Caching**: Reduce backend load

## Disaster Recovery

### RTO (Recovery Time Objective): 5 minutes
### RPO (Recovery Point Objective):  1 minute

### DR Plan

1. **Multi-AZ Deployment**: Gateways in multiple availability zones
2. **Multi-Region**: Active-passive setup for critical regions
3. **Configuration Management**: Infrastructure as Code (Terraform)
4. **Automated Backups**: Regular configuration backups
5. **Failover Procedures**: Documented and tested

## Summary

The API Gateway is a critical component that:
- Provides unified entry point for all clients
- Handles cross-cutting concerns (auth, rate limiting, caching)
- Enables service discovery and load balancing
- Ensures security and compliance
- Provides observability and monitoring
- Optimizes performance and cost

By centralizing these concerns at the gateway level, microservices can focus on business logic while maintaining consistency and security across the entire system.