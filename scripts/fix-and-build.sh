#!/bin/bash
#
# Fix npm permissions and build for EC2 deployment
#

set -e

echo "=========================================="
echo "Fixing npm permissions and building..."
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

echo -e "${YELLOW}Step 1: Fixing npm permissions...${NC}"
echo "You may be prompted for your password..."
sudo chown -R $(whoami) "$HOME/.npm"
echo -e "${GREEN}✓ npm permissions fixed${NC}"
echo ""

echo -e "${YELLOW}Step 2: Cleaning frontend...${NC}"
cd "$PROJECT_DIR/frontend"
rm -rf node_modules package-lock.json
npm cache clean --force
echo -e "${GREEN}✓ Frontend cleaned${NC}"
echo ""

echo -e "${YELLOW}Step 3: Installing dependencies...${NC}"
npm install
echo -e "${GREEN}✓ Dependencies installed${NC}"
echo ""

echo -e "${YELLOW}Step 4: Building backend...${NC}"
cd "$PROJECT_DIR/backend"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$PROJECT_DIR/build/daybook-backend" main.go
echo -e "${GREEN}✓ Backend built${NC}"
echo ""

echo -e "${YELLOW}Step 5: Building frontend...${NC}"
cd "$PROJECT_DIR/frontend"

# Create production .env
cat > .env.production << EOF
VITE_API_URL=https://api.daybook.shafik.xyz
VITE_APP_NAME=Daybook
VITE_APP_VERSION=1.0.0
EOF

npm run build
cd dist
# Exclude macOS extended attributes to avoid warnings on Linux
COPYFILE_DISABLE=1 tar --no-xattrs -czf "$PROJECT_DIR/build/frontend-dist.tar.gz" . 2>/dev/null || \
tar -czf "$PROJECT_DIR/build/frontend-dist.tar.gz" .
echo -e "${GREEN}✓ Frontend built${NC}"
echo ""

# Copy .env.example files
cp "$PROJECT_DIR/backend/.env.example" "$PROJECT_DIR/build/backend.env.example"
cp "$PROJECT_DIR/frontend/.env.example" "$PROJECT_DIR/build/frontend.env.example"

echo -e "${GREEN}=========================================="
echo "Build Complete!"
echo "==========================================${NC}"
echo ""
echo "Build artifacts:"
ls -lh "$PROJECT_DIR/build"
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
