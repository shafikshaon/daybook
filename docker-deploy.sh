#!/bin/bash
set -e

echo "🚀 Deploying Daybook with Docker"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop first."
    exit 1
fi

# Check if PostgreSQL is running on host
echo "🔍 Checking PostgreSQL on host..."
if ! nc -z localhost 5432 2>/dev/null; then
    echo "❌ PostgreSQL is not running on localhost:5432"
    echo "   Please start PostgreSQL on your host machine"
    exit 1
fi
echo "✅ PostgreSQL is running"

# Check if Redis is running on host
echo "🔍 Checking Redis on host..."
if ! nc -z localhost 6379 2>/dev/null; then
    echo "⚠️  Redis is not running on localhost:6379"
    echo "   Continuing anyway, but Redis caching will not work"
else
    echo "✅ Redis is running"
fi

echo ""
echo "📦 Building Docker images..."
docker-compose build

echo ""
echo "🚀 Starting containers..."
docker-compose up -d

echo ""
echo "⏳ Waiting for services to start..."
sleep 5

# Check if containers are running
if docker ps | grep -q daybook-backend && docker ps | grep -q daybook-frontend; then
    echo ""
    echo "✅ Deployment successful!"
    echo ""
    echo "🌐 Application URLs:"
    echo "   Frontend: http://localhost:3000"
    echo "   Backend:  http://localhost:8080"
    echo "   Health:   http://localhost:8080/health"
    echo ""
    echo "📊 View logs:"
    echo "   docker-compose logs -f backend"
    echo "   docker-compose logs -f frontend"
    echo ""
    echo "🛑 Stop containers:"
    echo "   docker-compose down"
else
    echo "❌ Deployment failed!"
    echo ""
    echo "Check logs:"
    echo "   docker-compose logs"
    exit 1
fi
