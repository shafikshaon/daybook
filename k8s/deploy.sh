#!/bin/bash
set -e

echo "🚀 Deploying Daybook to Kubernetes"
echo ""

# Apply manifests in order
echo "📋 Applying Kubernetes manifests..."

kubectl apply -f 00-namespace.yaml
echo "✅ Namespace created"

kubectl apply -f 01-secrets.yaml
echo "✅ Secrets created"

kubectl apply -f 02-postgres.yaml
echo "✅ PostgreSQL deployed"

kubectl apply -f 03-redis.yaml
echo "✅ Redis deployed"

kubectl apply -f 04-backend-config.yaml
echo "✅ Backend config created"

kubectl apply -f 05-backend.yaml
echo "✅ Backend deployed"

kubectl apply -f 06-frontend.yaml
echo "✅ Frontend deployed"

kubectl apply -f 07-ingress.yaml
echo "✅ Ingress created"

echo ""
echo "⏳ Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres -n daybook --timeout=120s
kubectl wait --for=condition=ready pod -l app=redis -n daybook --timeout=120s
kubectl wait --for=condition=ready pod -l app=backend -n daybook --timeout=120s
kubectl wait --for=condition=ready pod -l app=frontend -n daybook --timeout=120s

echo ""
echo "✅ Deployment Complete!"
echo ""
echo "📊 Status:"
kubectl get pods -n daybook
echo ""
echo "🌐 Access your application:"
echo "   http://daybook.shafik.xyz"
echo ""
echo "🔍 Useful commands:"
echo "   kubectl get all -n daybook"
echo "   kubectl logs -f deployment/backend -n daybook"
echo "   kubectl logs -f deployment/frontend -n daybook"
echo "   kubectl describe ingress daybook-ingress -n daybook"
