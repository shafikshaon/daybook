# Kubernetes Deployment - Quick Start

## 🚀 **Deploy Daybook on Kubernetes in 3 Steps**

### **Step 1: Build Docker Images**

```bash
cd k8s
./build-images.sh

# Login to Docker Hub
docker login

# Push images
docker push shafikshaon/daybook-backend:latest
docker push shafikshaon/daybook-frontend:latest
```

### **Step 2: Configure (Optional)**

Edit secrets if needed:
```bash
nano k8s/01-secrets.yaml
```

Change your domain:
```bash
nano k8s/07-ingress.yaml
# Change: host: daybook.shafik.xyz
```

### **Step 3: Deploy**

```bash
cd k8s
./deploy.sh
```

**Done!** Your app is running on Kubernetes.

---

## 📊 **What You Get**

- ✅ **PostgreSQL** with persistent storage
- ✅ **Redis** for caching
- ✅ **Backend** API (2 replicas, auto-scaling ready)
- ✅ **Frontend** web app (2 replicas)
- ✅ **Ingress** for external access
- ✅ **Health checks** and **auto-restart**

---

## 🔍 **Quick Commands**

```bash
# Check status
kubectl get all -n daybook

# View logs
kubectl logs -f deployment/backend -n daybook

# Scale up
kubectl scale deployment/backend --replicas=5 -n daybook

# Update backend
docker build -t shafikshaon/daybook-backend:latest backend/
docker push shafikshaon/daybook-backend:latest
kubectl rollout restart deployment/backend -n daybook

# Delete everything
cd k8s
./delete.sh
```

---

## 🌐 **Access Your App**

```bash
# Get ingress IP
kubectl get ingress -n daybook

# Access
http://daybook.shafik.xyz
```

---

## 📚 **Full Documentation**

See [k8s/README.md](k8s/README.md) for:
- Detailed configuration
- Monitoring & debugging
- Backup & restore
- SSL/TLS setup
- Production checklist
- Troubleshooting

---

## ☸️ **Kubernetes Requirements**

**Cluster:**
- 2+ vCPUs
- 4GB+ RAM
- 20GB+ storage
- NGINX Ingress Controller

**Supported Platforms:**
- AWS EKS
- Google GKE
- Azure AKS
- Minikube
- Kind
- Docker Desktop

---

## 💡 **Why Kubernetes?**

- ✅ **Zero-downtime deployments** - Rolling updates
- ✅ **Auto-scaling** - Handles traffic spikes
- ✅ **Self-healing** - Auto-restart failed pods
- ✅ **Load balancing** - Distribute traffic
- ✅ **Declarative** - Infrastructure as code
- ✅ **Portable** - Run anywhere

---

## 🎯 **Architecture**

```
Internet
   ↓
Ingress (NGINX)
   ↓
   ├─► Backend Service → Backend Pods (2x)
   │                      ↓
   │                      ├─► PostgreSQL
   │                      └─► Redis
   │
   └─► Frontend Service → Frontend Pods (2x)
```

All isolated in the `daybook` namespace with persistent storage for database and uploads.
