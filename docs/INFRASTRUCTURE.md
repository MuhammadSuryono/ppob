# Infrastructure & Deployment Specification for PPOB

**Audience:** DevOps, SRE, Platform Engineers  
**Last Updated:** 2026-05-07  
**Status:** Draft — cost estimates pending vendor quotes

---

## 1. Overview

This document defines the infrastructure architecture, deployment strategy, CI/CD pipeline, monitoring stack, and operational procedures for the PPOB multi-tenant application.

**Guiding Principles:**
- **Cloud-native:** Kubernetes-first, containerized microservices
- **Infrastructure as Code:** Terraform modules for all resources
- **Observability:** Centralized logging, metrics, tracing
- **Security:** mTLS, secrets in Vault, least-privilege IAM
- **Cost-Effective:** Autoscaling, spot instances for stateless workers, reserved instances for stateful

---

## 2. Cloud Provider & Region

**Provider:** Amazon Web Services (AWS) — or Google Cloud Platform (GCP) if preferred; AWS chosen for mature Jakarta region.

**Region:** `ap-southeast-3` (Jakarta, Indonesia) — ensures data residency compliance (Indonesian PDP law).

**Availability Zones:** 3 AZs (a, b, c) for HA.

**Account Strategy:** Single AWS account with IAM boundaries; separate VPC per environment (dev/staging/prod).

---

## 3. Kubernetes Cluster (EKS)

### 3.1 Cluster Configuration

| Setting | Value | Rationale |
|---|---|---|
| Kubernetes Version | 1.28 (latest stable at launch) | Security & features |
| Control Plane | EKS Managed (3 master nodes, multi-AZ) | High availability |
| Node Groups | 3 groups (see below) | Mixed workload optimization |
| VPC CNI | AWS VPC CNI (native) | Pods get ENI, IP per pod |
| Service Mesh | Istio (optional v1.20) | mTLS, traffic shifting, observability |
| Ingress Controller | Nginx Ingress Controller (nginx-ingress) | TLS termination, routing |
| Cluster Autoscaler | Enabled | Scale node groups based on pod demand |

### 3.2 Node Groups

| Node Group | Instance Type | Min | Max | Purpose | Taints |
|---|---|---|---|---|---|
| `system-nodes` | m5.large (2 vCPU, 8 GB) | 3 | 5 | System pods (Ingress, Metrics Server, Cluster Autoscaler) | No |
| `app-nodes` | m5.xlarge (4 vCPU, 16 GB) | 3 | 15 | All microservices (stateless) | No |
| `worker-nodes` | m5.2xlarge (8 vCPU, 32 GB) | 2 | 10 | Heavy workloads (sync jobs, reconciliation) | `dedicated=worker:NoSchedule` |

**Pricing:** On-Demand for system-nodes (always 3); 50% Spot Instances for app-nodes and worker-nodes to save ~60% cost (with fallback to on-demand).

**Amazon Machine Images (AMI):** Amazon Linux 2023 or Bottlerocket (for cost savings & security).

### 3.3 Pod Scheduling

- System-critical pods: `nodeSelector: {"node-role": "system"}` → pinned to system-nodes
- Regular services: no selector → distributed across app-nodes
- Heavy batch jobs (product sync, reconciliation): `nodeSelector: {"node-role": "worker"}`

---

## 4. Namespace & Environment Isolation

**Namespaces:**
- `ppob-dev`
- `ppob-staging`
- `ppob-prod`

**Resource Quotas per Namespace:**

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: ppob-quota
  namespace: ppob-prod
spec:
  hard:
    requests.cpu: "20"
    requests.memory: 40Gi
    limits.cpu: "40"
    limits.memory: 80Gi
    pods: "100"
    services: "50"
    secrets: "200"
    configmaps: "200"
```

---

## 5. Microservices Deployment

### 5.1 General Deployment Spec

Each service gets:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: ppob-prod
  labels:
    app: auth-service
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
        image: 123456789012.dkr.ecr.ap-southeast-3.amazonaws.com/ppob/auth-service:v1.2.3
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: auth-config
        - secretRef:
            name: auth-secrets
        resources:
          requests:
            cpu: "250m"
            memory: "512Mi"
          limits:
            cpu: "500m"
            memory: "1Gi"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
        securityContext:
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          capabilities:
            drop: ["ALL"]
```

### 5.2 Service (ClusterIP)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: ppob-prod
spec:
  selector:
    app: auth-service
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

### 5.3 Horizontal Pod Autoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
  namespace: ppob-prod
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

---

## 6. Database (PostgreSQL RDS)

### 6.1 Instance Configuration

| Setting | Value |
|---|---|
| Engine | PostgreSQL 15 |
| Instance Class | db.r6g.large (2 vCPU, 16 GB RAM) for prod; db.t4g.micro for dev |
| Storage | 100 GB General Purpose (SSD) |
| Multi-AZ | Enabled (production) |
| Backup Retention | 35 days (daily snapshots) |
| Point-in-Time Recovery | Enabled (WAL archiving to S3) |
| Storage Autoscaling | Enabled (max 1 TB) |
| Deletion Protection | Enabled (prod) |
| Encryption | KMS (AWS-managed or customer-managed CMK) |

### 6.2 Connection Pooling

**PgBouncer in Transaction Pooling Mode** (as sidecar or separate EC2):

```yaml
# pgbouncer.ini
[databases]
ppob = host=ppob-db.prod.ap-southeast-3.rds.amazonaws.com port=5432 dbname=ppob

[pgbouncer]
listen_port = 6432
listen_addr = *
auth_type = md5
auth_file = /etc/pgbouncer/userlist.txt
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20
reserve_pool_size = 10
```

**Services connect to `ppob-pool:6432` (ClusterIP service for PgBouncer).**

---

## 7. Caching (Redis)

### 7.1 ElastiCache Redis Cluster

| Setting | Value |
|---|---|
| Engine | Redis 7.2 |
| Node Type | cache.r6g.large (2 vCPU, 13 GB) |
| Number of Nodes | 3 (cluster mode enabled, 1 shard + 2 replicas) |
| Multi-AZ | Enabled |
| At-Rest Encryption | Enabled (KMS) |
| In-Transit Encryption | TLS enabled |
| Maintenance Window | Sun 02:00–04:00 UTC+7 |

### 7.2 Redis Usage

| Key Pattern | TTL | Purpose |
|---|---|---|
| `refresh_token:{jti}` | 7d (Mitra), 1d (Staff) | Token allowlist |
| `rate:auth:login:{phone}` | 15m sliding | Rate limiting |
| `product:cache:{sku}` | 3600s | Product catalog cache |
| `wallet:balance:{wallet_id}` | 60s | Wallet balance cache (write-through invalidate) |
| `idempotency:{key_hash}` | 24h | Idempotency check |
| `digiflazz:deposit` | 60s | Cached Digiflazz balance check |

---

## 8. Secrets Management (HashiCorp Vault)

### 8.1 Vault Cluster

**Deployment:** Vault OSS in HA mode (3 nodes) on EKS behind internal load balancer, using RDS for storage backend.

**Alternative:** AWS Secrets Manager (managed) — less operational overhead, higher cost.

**Choice:** **AWS Secrets Manager** for simplicity (unless advanced dynamic secrets needed). 

**Secrets Structure (Secrets Manager):**

```
/ppob/dev/digiflazz/api_key
/ppob/dev/digiflazz/webhook_secret
/ppob/dev/jwt/private_key
/ppob/dev/postgres/password
/ppob/prod/digiflazz/api_key
...
```

**Rotation:** Lambda functions for automatic rotation (30d for DB, 90d for JWT, quarterly for Digiflazz).

---

## 9. Networking

### 9.1 VPC Design

```
VPC: 10.0.0.0/16
  ├─ Public Subnets (AZ a,b,c): 10.0.1.0/24, 10.0.2.0/24, 10.0.3.0/24
  │   └─ NAT Gateways (one per AZ) for egress
  ├─ Private App Subnets (AZ a,b,c): 10.0.11.0/24, 10.0.12.0/24, 10.0.13.0/24
  │   └─ EKS worker nodes (no public IP)
  └─ Private Data Subnets (AZ a,b,c): 10.0.21.0/24, 10.0.22.0/24, 10.0.23.0/24
      └─ RDS, ElastiCache (no public IP)
```

**Route Tables:** Private subnets route through NAT Gateway for internet egress (to Digiflazz, Vault, etc.). No inbound from internet except through ALB (public subnets).

### 9.2 Ingress (ALB/Nginx)

**Option A:** AWS ALB (managed) with Ingress Controller (service type LoadBalancer).

**Option B:** Nginx Ingress Controller on EKS with Elastic IPs.

**Chosen:** Nginx Ingress Controller (more control, flexible annotations).

TLS via Let's Encrypt `cert-manager` or AWS Certificate Manager (ACM) with ` ingress` annotations.

---

## 10. CI/CD Pipeline

### 10.1 GitHub Actions Workflow

**File:** `.github/workflows/deploy.yml`

```yaml
name: Deploy PPOB

on:
  push:
    branches: [main, staging]
  pull_request:
    branches: [main]

env:
  REGISTRY: 123456789012.dkr.ecr.ap-southeast-3.amazonaws.com
  IMAGE_NAME: ppob/${{ github.event.repository.name }}

jobs:
  lint-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with: { go-version: '1.22' }
      - run: go mod download
      - run: make lint
      - run: make test
      - run: make build

  docker-build:
    needs: lint-test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main' || github.ref == 'refs/heads/staging'
    steps:
      - uses: actions/checkout@v3
      - uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-region: ap-southeast-3
          role-to-assume: arn:aws:iam::123456789012:role/GitHubActionsDeploy
      - run: aws ecr get-login-password | docker login --username AWS --password-stdin $REGISTRY
      - run: |
          TAG=${GITHUB_SHA:0:7}
          docker build -t $REGISTRY/$IMAGE_NAME:$TAG .
          docker push $REGISTRY/$IMAGE_NAME:$TAG

  deploy-dev:
    needs: docker-build
    if: github.ref == 'refs/heads/staging'
    runs-on: ubuntu-latest
    steps:
      - uses: azure/setup-kubectl@v3
      - run: |
          aws eks update-kubeconfig --name ppob-dev --region ap-southeast-3
          kubectl set image deployment/auth-service auth=$REGISTRY/$IMAGE_NAME:$TAG -n ppob-dev
          kubectl rollout status deployment/auth-service -n ppob-dev

  deploy-prod:
    needs: [docker-build, test-integration]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to production (manual approval)
        uses: peter-evans/slash-command-dispatch@v2
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          command: '/deploy-prod'
          permission: write
```

**Auto-deploy to staging on push to `staging` branch. Production requires manual approval via slash command `/deploy-prod` or manual GitHub Actions workflow dispatch.**

### 10.2 ArgoCD (GitOps Alternative)

Instead of GitHub Actions pushing, use **ArgoCD** for declarative continuous delivery:

- Application definitions stored in `argocd/applications/` as YAML
- ArgoCD syncs cluster state to Git state
- Manual or automatic sync (auto for staging, manual for prod)

**Chosen:** GitHub Actions for simplicity; can migrate to ArgoCD later for full GitOps.

---

## 11. Monitoring Stack

### 11.1 Metrics Collection (Prometheus)

**Deployment:** `kube-prometheus-stack` Helm chart (includes Prometheus Operator, Grafana, Alertmanager).

**Scrape Configs:**
```yaml
scrape_configs:
- job_name: 'prometheus'
  static_configs:
  - targets: ['localhost:9090']

- job_name: 'ppob-services'
  kubernetes_sd_configs:
  - role: pod
  relabel_configs:
  - source_labels: [__meta_kubernetes_pod_label_app]
    action: keep
    regex: .*-service
  metrics_path: /metrics
  scheme: http
```

### 11.2 Logging (Loki or OpenSearch)

**Option A (Lighter):** Grafana Loki + Promtail (logs as labels; cheaper)
**Option B (Full-text):** OpenSearch (Elasticsearch fork) + Fluent Bit

**Chosen:** **Loki** for cost efficiency; searchable by service, pod, trace_id.

**Fluent Bit DaemonSet** collects container logs from `/var/log/containers/*.log`, forwards to Loki.

### 11.3 Tracing (Jaeger)

**Jaeger Operator** on EKS with all-in-one deployment for start (scale later).

**OpenTelemetry Collector** as sidecar or DaemonSet to receive traces from services (OTLP) and export to Jaeger.

### 11.4 Alertmanager

**Routing:**
- Critical → PagerDuty (SRE on-call)
- Warning → Slack #alerts-ppob
- Info → #ppob-monitoring (optional)

**Routes Example:**
```yaml
routes:
- match:
    severity: critical
  receiver: 'pagerduty'
- match:
    severity: warning
  receiver: 'slack-alerts'
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
```

---

## 12. Cost Estimation (Monthly)

| Resource | Qty | Unit Cost (USD) | Monthly (USD) | Notes |
|---|---|---|---|---|
| EKS Control Plane | 1 | $0.10/hr | $73 | Fixed |
| EC2 m5.large (system) | 3 | $0.096/hr | $69 | Always on |
| EC2 m5.xlarge (app) | 3 on-demand + 3 spot | avg $0.06/hr | $140 | 50% spot |
| RDS PostgreSQL (db.r6g.large) | 1 | $0.352/hr | $253 | Multi-AZ doubles? Listed is per instance; Multi-AZ adds second instance but only one active? Pricing includes Multi-AZ |
| ElastiCache Redis (cache.r6g.large) | 3 nodes | $0.182/hr | $131 | Cluster |
| ALB | 1 | $0.025/hr + LCU | $40 | Est. 10 LCU |
| NAT Gateway | 3 | $0.045/hr | $97 | 1 per AZ |
| Data Transfer | — | — | $50 | Egress to internet |
| S3 (backups) | 1 TB | $0.023/GB | $24 | — |
| Secrets Manager (10 secret) | 10 | $0.40/secret/mo | $4 | — |
| CloudWatch Logs (ingest 10 GB/mo) | — | $0.50/GB | $5 | — |
| **Total** | — | — | **~$886** | ~Rp 14.5 juta/month |

**With Overprovisioning & Data Transfer:** ~$1,200/month (includes buffer, more nodes as load grows).

---

## 13. Backup & Disaster Recovery

### 13.1 Backup Strategy

- **RDS Automated Backups:** Daily snapshots retained 35 days; point-in-time recovery to any second within period
- **WAL Archiving:** Amazon RDS automatically stores WAL segments in S3 (for PITR)
- **Manual Export:** Weekly full dump to S3 (`pg_dump`) for long-term archive
- **Redis RDB:** Daily snapshot to S3 (if data needs persistence; usually cache only)
- **Config & Code:** Git repository + Terraform state in S3 with versioning + DynamoDB lock

**Recovery Objectives:**
- **RTO (Recovery Time Objective):** 4 hours (bring up new cluster from latest backup)
- **RPO (Recovery Point Objective):** 15 minutes (max data loss from last WAL archive)

### 13.2 DR Plan

**Scenario 1 — AZ failure:** EKS control plane & nodes spread across 3 AZs; cluster self-heals by moving pods to healthy AZs. RDS Multi-AZ promotes standby automatically. **RTO:** <30min.

**Scenario 2 — Region failure (Jakarta down):**
- Promote read replica in backup region (Singapore ap-southeast-1) to primary
- Update DNS/Route53 to point to new region
- **RTO:** 4 hours (manual failover, DNS propagation)

**Run DR drills quarterly.**

---

## 14. Security Hardening

- **Pod Security Standards:** `restricted` profile (non-root, read-only root filesystem, drop capabilities)
- **Network Policies:** Default deny all ingress/egress between namespaces; allow only necessary communication
- **Secrets in K8s:** Use External Secrets Operator to fetch from AWS Secrets Manager into K8s Secrets (encrypted at rest)
- **Pod Identity:** IAM Roles for Service Accounts (IRSA) — each service gets minimal IAM permissions (e.g., can only read its own secrets)
- **Image Scanning:** Trivy scans images on push to ECR; ImageScanningConfiguration in ECR
- **Runtime Security:** Falco (optional) for anomaly detection

---

## 15. Scaling Strategy

### 15.1 Horizontal Pod Autoscaling

Based on CPU/Memory and custom metrics (Prometheus):
- **Auth Service:** CPU > 70% → scale 3→5
- **Transaction Service:** Queue length (pending transactions) metric → scale 3→10
- **Product Sync Worker:** Based on `products_last_sync_duration_seconds`; keep single replica unless taking >10min, then spawn parallel workers

### 15.2 Database Scaling

- **Read Replicas:** 1 async replica for reporting queries; offload analytics
- **Connection Pooling:** PgBouncer in transaction pooling mode limits connections to 100; each service max 10 connections
- **Partitioning:** `transactions` table partition by month (monthly partitions after 6 months)

### 15.3 Caching Strategy

- **Product Catalog:** Cache-aside (lazy load) with TTL 1h; warm on sync job completion
- **Wallet Balance:** Write-through (update cache on every debit/credit via WalletService)
- **Rate Limiting:** In-memory Redis (fast)

---

## 16. Health Checks & Liveness

**Endpoints per service:**

| Endpoint | Purpose | Response |
|---|---|---|
| `GET /health/live` | K8s liveness (is process alive?) | 200 + `{"status":"ok"}` |
| `GET /health/ready` | K8s readiness (can accept traffic?) | 200 `{"db":true,"redis":true}` or 503 |
| `GET /metrics` | Prometheus scrape | `/metrics` format |

**Readiness checks:**
- Database connectivity (SELECT 1)
- Redis connectivity (PING)
- Vault connectivity (health endpoint)
- Circuit breaker state (if Digiflazz down, still ready — don't fail readiness, but may degrade)

---

## 17. Maintenance & Operations

### 17.1 kubectl Tips

```bash
# View logs with labels
kubectl logs -f -l app=auth-service -n ppob-prod

# Port-forward to DB
kubectl port-forward svc/ppob-db 5432:5432 -n ppob-prod --address 0.0.0.0

# Exec into pod
kubectl exec -it auth-service-xxxx -- /bin/sh

# Watch HPA
kubectl get hpa -n ppob-prod --watch
```

### 17.2 Database Maintenance

- **Vacuum & Analyze:** Run nightly `VACUUM ANALYZE` (RDS auto-vacuum sufficient for low-medium load)
- **Index Rebuild:** Quarterly `REINDEX` if bloat detected (`pg_stat_user_indexes`)
- **Backup Restoration Test:** Monthly restore from backup to test environment

---

## 18. Compliance & Data Residency

- All data stored in Jakarta AWS region (ap-southeast-3)
- No cross-border data transfer (except Digiflazz API call which routes to their servers; they may process outside Indonesia — review their data location)
- Backup S3 bucket with region-restricted access only
- Enable S3 Object Lock for compliance retention (WORM) if required by regulator

---

## 19. Disaster Recovery Runbook

**RTO: 4 hours**

1. **Assess:** Identify failed component (AZ, DB, entire cluster)
2. **Decide:** Failover to standby region or restore from backup?
3. **Restore DB:** If RDS lost, create new instance from latest snapshot, apply WAL replay
4. **Redeploy Cluster:** If EKS lost, `eksctl create cluster` from Terraform state (state stored in S3)
5. **Update DNS:** Point ALB to new cluster IPs/ELB
6. **Validate:** Smoke test all endpoints, verify balance reconciliation
7. **Post-mortem:** Document root cause, update runbook

---

## 20. Terraform Module Structure

```
terraform/
  modules/
    vpc/
      main.tf
      variables.tf
      outputs.tf
    eks/
      main.tf
      node_groups/
    rds/
      main.tf
    redis/
      main.tf
    secrets/
      main.tf (AWS Secrets Manager)
  environments/
    dev/
      terraform.tfvars
    staging/
      terraform.tfvars
    prod/
      terraform.tfvars
  backend.tf (S3 state store)
```

**State Locking:** DynamoDB table `terraform-lock`.

**Apply Order:** VPC → EKS → RDS → Redis → Secrets → Applications.

---

## Appendix A — Sample Terraform for EKS

```hcl
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = "ppob-prod"
  cluster_version = "1.28"
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnets

  node_security_group_additional_rules = {
    ingress_self_all = {
      description = "Node to node all ports/protocols"
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      type        = "ingress"
      self        = true
    }
    egress_all = {
      description = "Egress to anywhere"
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      type        = "egress"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  eks_managed_node_group_defaults = {
    ami_type       = "AL2_x86_64"
    instance_types = ["m5.xlarge"]
    capacity_type  = "SPOT"
  }

  eks_managed_node_groups = {
    app = {
      min_size = 3
      max_size = 15
      desired_size = 5
    }
    worker = {
      min_size = 2
      max_size = 10
      desired_size = 3
      instance_types = ["m5.2xlarge"]
      capacity_type  = "ON_DEMAND"  // worker uses on-demand for stability
    }
  }
}
```

---

## Appendix B — Deployment Rollout Strategy

**Strategy:** Rolling update with max surge 25% and max unavailable 25%.

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
```

**Zero-Downtime:** Readiness probes ensure old pods not terminated until new pods healthy.

**Canary Deployments (future):** Use Flagger + Istio to shift traffic 5% → 25% → 100% based on metrics.

---

## Appendix C — On-Call Runbook Snippet

**Alert:** `DigiflazzAPIErrorRateHigh`

1. Check Digiflazz status page: https://status.digiflazz.com
2. If their incident: inform customers, wait
3. If our circuit breaker open: check logs, consider manual reset
4. If no known Digiflazz issue: inspect Integration Service logs for `ERROR` level with `digiflazz` Tag
5. Restart Integration Service pods if stuck in error loop
6. If persists >30min, page on-call SRE lead

---

## Open Questions

1. **Database Multi-Tenant Isolation:** Single DB shared (current design) or separate DB per Mitra? Currently single DB with `user_id` scoping — simpler, but performance isolation weaker. **Decision:** Single DB; scale reads with replica if needed.

2. **Service Mesh Adoption:** Istio adds complexity; could use Linkerd instead. **Decision:** Defer; add in Phase 2 after MVP stable.

3. **Node Scaling:** Autoscaling based on CPU alone ok; or custom metric (queue length). **Decision:** HPA on CPU + manual HPA on custom metric for transaction service.

---

**Owner:** DevOps Team  
**Next Steps:**
1. Spin up dev VPC with Terraform
2. Deploy EKS cluster (dev)
3. Install ingress, metrics server, cert-manager
4. Deploy RDS + Redis
5. CI/CD: GitHub Actions → ECR → EKS dev namespace
6. Smoke test deploy of auth service

---

**Related Infrastructure Docs:**
- `SECURITY_ARCHITECTURE.md` (Vault, network policies)
- `OBSERVABILITY_SPEC.md` (Prometheus, Grafana dashboards, which this infra hosts)
