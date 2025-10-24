#!/bin/bash

# Daybook Backend Quick Start Script

set -e

echo "🚀 Starting Daybook Backend..."

# Check if docker-compose is available
if command -v docker-compose &> /dev/null; then
    echo "📦 Using Docker Compose..."
    docker-compose up -d
    echo ""
    echo "✅ Services started successfully!"
    echo ""
    echo "📍 Backend API: http://localhost:8080"
    echo "📍 PostgreSQL: localhost:5432"
    echo "📍 Redis: localhost:6379"
    echo ""
    echo "📋 Useful commands:"
    echo "  - View logs: docker-compose logs -f"
    echo "  - Stop services: docker-compose down"
    echo "  - Check health: curl http://localhost:8080/health"
    echo ""
else
    echo "🔧 Docker Compose not found. Starting locally..."

    # Check if PostgreSQL is running
    if ! pg_isready -h localhost -p 5432 &> /dev/null; then
        echo "⚠️  PostgreSQL is not running!"
        echo "Please start PostgreSQL and ensure it's accessible on localhost:5432"
        exit 1
    fi

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo "❌ Go is not installed!"
        echo "Please install Go 1.21 or higher"
        exit 1
    fi

    # Install dependencies
    echo "📦 Installing dependencies..."
    go mod download

    # Run the application
    echo "🚀 Starting backend..."
    go run main.go
fi
