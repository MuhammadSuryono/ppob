# 🏗️ Infrastructure & Deployment

## 1. Overview
The PPOB system is deployed on **AWS (Jakarta Region - ap-southeast-3)** to ensure data residency compliance and low latency for Indonesian users. The architecture is cloud-native, utilizing Kubernetes for orchestration.

## 2. Compute: AWS EKS
We use Amazon Elastic Kubernetes Service (EKS) to manage our microservices.
- **Node Groups:**
    - `system-nodes`: For cluster utilities (Ingress, Monitoring).
    - `app-nodes`: m5.xlarge instances for microservices (mix of On-Demand and Spot).
    - `worker-nodes`: Heavy batch jobs (Sync, Reconciliation).
- **Service Mesh:** Istio for mTLS, traffic splitting, and observability.

## 3. Storage & Caching
- **Database:** Amazon RDS PostgreSQL 15 (Multi-AZ for HA).
- **Caching:** Amazon ElastiCache Redis 7.2 (Cluster mode enabled).
- **Connection Pooling:** PgBouncer (Transaction pooling mode) to manage DB connections efficiently.

## 4. Networking
- **VPC Design:** 3-AZ layout with public subnets for Load Balancers and private subnets for Apps and Databases.
- **Ingress:** Nginx Ingress Controller handles SSL termination and internal routing.
- **NAT Gateways:** One per AZ for secure egress.

## 5. CI/CD Pipeline
We use **GitHub Actions** for our automated delivery pipeline:
1.  **Continuous Integration:** Linting, Unit Tests, and Docker Image building.
2.  **Image Registry:** Amazon ECR (Elastic Container Registry).
3.  **Deployment:** `kubectl` or `Helm` to update EKS deployments.
4.  **Strategy:** Rolling updates for stateless services; manual approval for production.

## 6. Observability Stack
- **Metrics:** Prometheus Operator + Grafana dashboards.
- **Logging:** Fluent Bit forwarding to Grafana Loki.
- **Tracing:** OpenTelemetry Collector + Jaeger for distributed request tracking.
- **Alerting:** Alertmanager routing critical issues to PagerDuty/Slack.

## 7. Disaster Recovery
- **RPO:** 15 minutes (WAL archiving to S3).
- **RTO:** 4 hours (Cluster recreation from Terraform and latest DB snapshot).
- **Backup:** Daily automated RDS snapshots; weekly S3 exports.
