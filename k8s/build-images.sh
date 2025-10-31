#!/bin/bash
set -e

echo "🔨 Building Docker Images for Kubernetes"
echo ""

# Your Docker Hub username
DOCKER_USERNAME="shafikshaon"

# Build backend
echo "📦 Building Backend Image..."
cd ../backend
docker build -t ${DOCKER_USERNAME}/daybook-backend:latest .
echo "✅ Backend image built"

# Build frontend
echo "📦 Building Frontend Image..."
cd ../frontend
docker build -t ${DOCKER_USERNAME}/daybook-frontend:latest .
echo "✅ Frontend image built"

echo ""
echo "✅ All images built successfully!"
echo ""
echo "📤 Next steps:"
echo "   1. docker login"
echo "   2. docker push ${DOCKER_USERNAME}/daybook-backend:latest"
echo "   3. docker push ${DOCKER_USERNAME}/daybook-frontend:latest"
echo "   4. ./deploy.sh"
