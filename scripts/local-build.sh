#!/bin/bash
#
# Local Build Script for EC2 Deployment
# Run this on your local machine before deploying to EC2
#

set -e

echo "=========================================="
echo "Daybook Local Build Script"
echo "=========================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Build directories
BACKEND_DIR="$PROJECT_DIR/backend"
FRONTEND_DIR="$PROJECT_DIR/frontend"
BUILD_DIR="$PROJECT_DIR/build"

echo -e "${YELLOW}Project Directory: $PROJECT_DIR${NC}"
echo ""

# Create build directory
mkdir -p "$BUILD_DIR"

# Build Backend
echo -e "${GREEN}[1/3] Building Backend...${NC}"
cd "$BACKEND_DIR"

if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found in $BACKEND_DIR${NC}"
    exit 1
fi

echo "Building Go binary for Linux..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/daybook-backend" main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Backend build successful${NC}"
    ls -lh "$BUILD_DIR/daybook-backend"
else
    echo -e "${RED}✗ Backend build failed${NC}"
    exit 1
fi
echo ""

# Build Frontend
echo -e "${GREEN}[2/3] Building Frontend...${NC}"
cd "$FRONTEND_DIR"

if [ ! -f "package.json" ]; then
    echo -e "${RED}Error: package.json not found in $FRONTEND_DIR${NC}"
    exit 1
fi

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install
else
    echo "Dependencies already installed, skipping npm install..."
fi

# Create production .env if it doesn't exist
if [ ! -f ".env.production" ]; then
    echo "Creating .env.production..."
    cat > .env.production << EOF
VITE_API_URL=https://api.daybook.shafik.xyz
VITE_APP_NAME=Daybook
VITE_APP_VERSION=1.0.0
EOF
fi

echo "Building Vue.js application..."
npm run build

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Frontend build successful${NC}"
    echo "Creating tarball..."
    cd dist
    tar -czf "$BUILD_DIR/frontend-dist.tar.gz" .
    echo -e "${GREEN}✓ Frontend tarball created${NC}"
    ls -lh "$BUILD_DIR/frontend-dist.tar.gz"
else
    echo -e "${RED}✗ Frontend build failed${NC}"
    exit 1
fi
echo ""

# Create deployment package
echo -e "${GREEN}[3/3] Creating Deployment Package...${NC}"

# Copy .env.example files
cp "$BACKEND_DIR/.env.example" "$BUILD_DIR/backend.env.example"
cp "$FRONTEND_DIR/.env.example" "$BUILD_DIR/frontend.env.example"

# Create deployment info
cat > "$BUILD_DIR/DEPLOYMENT_INFO.txt" << EOF
Daybook Deployment Package
Generated: $(date)

Files in this package:
- daybook-backend          : Backend binary (Linux x64)
- frontend-dist.tar.gz     : Frontend build (gzipped tarball)
- backend.env.example      : Backend environment template
- frontend.env.example     : Frontend environment template

Next Steps:
1. Transfer files to EC2:
   scp -i your-key.pem daybook-backend ubuntu@your-ec2-ip:/home/ubuntu/
   scp -i your-key.pem frontend-dist.tar.gz ubuntu@your-ec2-ip:/home/ubuntu/

2. Follow the EC2_DEPLOYMENT_GUIDE.md for deployment instructions

Build Information:
- Backend: Go (CGO_ENABLED=0, GOOS=linux, GOARCH=amd64)
- Frontend: Vue.js + Vite production build
EOF

echo ""
echo -e "${GREEN}=========================================="
echo "Build Complete!"
echo "==========================================${NC}"
echo ""
echo "Build artifacts location: $BUILD_DIR"
echo ""
ls -lh "$BUILD_DIR"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo ""
echo "1. Commit and push the build artifacts:"
echo "   git add build/"
echo "   git commit -m 'Add deployment build artifacts'"
echo "   git push origin master"
echo ""
echo "2. On EC2, pull the latest changes:"
echo "   cd ~/daybook"
echo "   git pull origin master"
echo "   ./scripts/ec2-setup.sh"
echo ""
echo -e "${YELLOW}Or follow the EC2_DEPLOYMENT_GUIDE.md for detailed steps.${NC}"
