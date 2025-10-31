# Daybook Kubernetes Deployment

Complete Kubernetes deployment for Daybook application.

## Prerequisites

1. **Kubernetes Cluster** (any of these):
   - AWS EKS
   - Google GKE
   - Azure AKS
   - Local: Minikube, Kind, or Docker Desktop

2. **kubectl** installed and configured

3. **Docker** installed (for building images)

4. **NGINX Ingress Controller** installed in your cluster

## Quick Start

### 1. Build Docker Images

```bash
cd k8s

# Build images
./build-images.sh

# Login to Docker Hub
docker login

# Push images
docker push shafikshaon/daybook-backend:latest
docker push shafikshaon/daybook-frontend:latest
```

### 2. Deploy to Kubernetes

```bash
# Deploy everything
./deploy.sh
```

That's it! Your application is now running on Kubernetes.

## What Gets Deployed

- **Namespace**: `daybook` - Isolated environment
- **PostgreSQL**: Database with persistent storage (10Gi)
- **Redis**: Cache service
- **Backend**: Go API (2 replicas)
- **Frontend**: React app (2 replicas)
- **Ingress**: NGINX ingress for routing

## Architecture

```
┌─────────────────────────────────────────────────┐
│              Ingress (NGINX)                    │
│         daybook.shafik.xyz                      │
└────┬────────────────────────────────────────────┘
     │
     ├─► /api/* ────────► Backend Service (8080)
     │                    ├─► Backend Pod 1
     │                    └─► Backend Pod 2
     │
     ├─► /health ────────► Backend Service
     │
     ├─► /uploads/* ─────► Backend Service
     │
     └─► /* ─────────────► Frontend Service (80)
                            ├─► Frontend Pod 1
                            └─► Frontend Pod 2

Backend connects to:
  ├─► PostgreSQL Service (5432)
  └─► Redis Service (6379)
```

## Configuration

### Edit Secrets

```bash
nano 01-secrets.yaml
```

Change these values:
```yaml
DB_PASSWORD: "your-db-password"
JWT_SECRET: "your-jwt-secret"
POSTGRES_PASSWORD: "your-db-password"
```

### Edit Domain

```bash
nano 07-ingress.yaml
```

Change:
```yaml
host: daybook.shafik.xyz  # Your domain
```

### Edit Backend Config

```bash
nano 04-backend-config.yaml
```

Adjust environment variables as needed.

## Access Your Application

After deployment:

```bash
# Get ingress IP/hostname
kubectl get ingress -n daybook

# Access the application
http://daybook.shafik.xyz
```

## Monitoring

### Check Status

```bash
# All resources
kubectl get all -n daybook

# Pods
kubectl get pods -n daybook

# Services
kubectl get svc -n daybook

# Ingress
kubectl get ingress -n daybook
```

### View Logs

```bash
# Backend logs
kubectl logs -f deployment/backend -n daybook

# Frontend logs
kubectl logs -f deployment/frontend -n daybook

# PostgreSQL logs
kubectl logs -f deployment/postgres -n daybook

# All pods
kubectl logs -f -l app=backend -n daybook
```

### Debug

```bash
# Describe a pod
kubectl describe pod <pod-name> -n daybook

# Shell into a pod
kubectl exec -it deployment/backend -n daybook -- /bin/sh

# Check events
kubectl get events -n daybook --sort-by='.lastTimestamp'
```

## Scaling

### Scale Backend

```bash
kubectl scale deployment/backend --replicas=3 -n daybook
```

### Scale Frontend

```bash
kubectl scale deployment/frontend --replicas=3 -n daybook
```

## Updates

### Update Backend

```bash
# Build new image
cd k8s
./build-images.sh

# Push to registry
docker push shafikshaon/daybook-backend:latest

# Rolling update
kubectl rollout restart deployment/backend -n daybook

# Check rollout status
kubectl rollout status deployment/backend -n daybook
```

### Update Frontend

```bash
# Build new image
cd k8s
./build-images.sh

# Push to registry
docker push shafikshaon/daybook-frontend:latest

# Rolling update
kubectl rollout restart deployment/frontend -n daybook
```

## Backup

### Backup Database

```bash
# Create a backup job
kubectl exec -it deployment/postgres -n daybook -- pg_dump -U postgres daybook > backup.sql
```

### Restore Database

```bash
kubectl exec -i deployment/postgres -n daybook -- psql -U postgres daybook < backup.sql
```

## Troubleshooting

### Pods Not Starting

```bash
# Check pod status
kubectl describe pod <pod-name> -n daybook

# Check logs
kubectl logs <pod-name> -n daybook
```

### Database Connection Issues

```bash
# Check if PostgreSQL is running
kubectl get pods -l app=postgres -n daybook

# Test connection
kubectl exec -it deployment/postgres -n daybook -- psql -U postgres -d daybook -c "SELECT 1;"
```

### Ingress Not Working

```bash
# Check ingress controller
kubectl get pods -n ingress-nginx

# Check ingress resource
kubectl describe ingress daybook-ingress -n daybook

# Check service endpoints
kubectl get endpoints -n daybook
```

## Clean Up

```bash
# Delete everything
./delete.sh

# Or manually
kubectl delete namespace daybook
```

## Resource Requirements

**Minimum Cluster Requirements:**
- 2 vCPUs
- 4GB RAM
- 20GB Storage

**Per Service:**

| Service    | CPU Request | CPU Limit | Memory Request | Memory Limit | Storage |
|------------|-------------|-----------|----------------|--------------|---------|
| PostgreSQL | 250m        | 500m      | 256Mi          | 512Mi        | 10Gi    |
| Redis      | 100m        | 200m      | 128Mi          | 256Mi        | -       |
| Backend    | 250m        | 500m      | 256Mi          | 512Mi        | 5Gi     |
| Frontend   | 100m        | 200m      | 128Mi          | 256Mi        | -       |

## Production Checklist

- [ ] Change default passwords in `01-secrets.yaml`
- [ ] Set up proper domain with DNS
- [ ] Configure SSL/TLS (add cert-manager)
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure autoscaling (HPA)
- [ ] Set up log aggregation (ELK/Loki)
- [ ] Configure backup strategy
- [ ] Set up CI/CD pipeline
- [ ] Review resource limits
- [ ] Set up network policies

## Advanced: SSL/TLS with cert-manager

```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Create ClusterIssuer
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF

# Update ingress with TLS
# Edit 07-ingress.yaml and add TLS section
```

## Support

For issues or questions:
- Check logs: `kubectl logs -f deployment/backend -n daybook`
- Check events: `kubectl get events -n daybook`
- Describe resources: `kubectl describe <resource> -n daybook`
