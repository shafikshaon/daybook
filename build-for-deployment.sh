#!/bin/bash
# Build backend binary for deployment to server

set -e

echo "🔨 Building Daybook Backend for Production Deployment"
echo ""

cd backend

echo "📦 Downloading dependencies..."
go mod download

echo "🔧 Building Linux binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o daybook-backend \
    main.go

if [ -f "daybook-backend" ]; then
    SIZE=$(du -h daybook-backend | cut -f1)
    echo ""
    echo "✅ Build successful!"
    echo "   Binary: backend/daybook-backend"
    echo "   Size: $SIZE"
    echo ""
    echo "📤 Next steps:"
    echo "   1. git add backend/daybook-backend"
    echo "   2. git commit -m 'Add pre-built binary for deployment'"
    echo "   3. git push"
    echo ""
    echo "🚀 On server, just run: ./deploy.sh"
    echo "   (It will use the pre-built binary and skip compilation!)"
else
    echo "❌ Build failed!"
    exit 1
fi
