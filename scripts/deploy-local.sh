#!/bin/bash
#
# Local Deployment Script
# Detects changes in frontend/backend and rebuilds/pushes to git
#

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Directories
BACKEND_DIR="$PROJECT_DIR/backend"
FRONTEND_DIR="$PROJECT_DIR/frontend"
BUILD_DIR="$PROJECT_DIR/build"

echo -e "${BLUE}=========================================="
echo "Daybook Auto-Deploy (Local)"
echo "==========================================${NC}"
echo ""

# Check for uncommitted changes
if ! git diff-index --quiet HEAD -- 2>/dev/null; then
    echo -e "${YELLOW}⚠ You have uncommitted changes${NC}"
    git status --short
    echo ""
    read -p "Do you want to continue? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Detect changes
BACKEND_CHANGED=false
FRONTEND_CHANGED=false

# Check if backend files changed
if git diff --name-only HEAD~1 HEAD | grep -q "^backend/"; then
    BACKEND_CHANGED=true
    echo -e "${YELLOW}✓ Backend changes detected${NC}"
fi

# Check if frontend files changed
if git diff --name-only HEAD~1 HEAD | grep -q "^frontend/"; then
    FRONTEND_CHANGED=true
    echo -e "${YELLOW}✓ Frontend changes detected${NC}"
fi

# Force build option
if [ "$1" == "--force" ] || [ "$1" == "-f" ]; then
    BACKEND_CHANGED=true
    FRONTEND_CHANGED=true
    echo -e "${YELLOW}✓ Force build enabled (building both)${NC}"
fi

# If no changes detected, check manually
if [ "$BACKEND_CHANGED" = false ] && [ "$FRONTEND_CHANGED" = false ]; then
    echo -e "${BLUE}No changes detected in backend or frontend${NC}"
    echo ""
    read -p "Build backend anyway? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        BACKEND_CHANGED=true
    fi

    read -p "Build frontend anyway? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        FRONTEND_CHANGED=true
    fi
fi

# If still nothing to build, exit
if [ "$BACKEND_CHANGED" = false ] && [ "$FRONTEND_CHANGED" = false ]; then
    echo -e "${RED}Nothing to build. Exiting.${NC}"
    exit 0
fi

echo ""
echo -e "${GREEN}Build Plan:${NC}"
[ "$BACKEND_CHANGED" = true ] && echo "  ✓ Backend will be rebuilt"
[ "$FRONTEND_CHANGED" = true ] && echo "  ✓ Frontend will be rebuilt"
echo ""

# Create build directory
mkdir -p "$BUILD_DIR"

# Build Backend
if [ "$BACKEND_CHANGED" = true ]; then
    echo -e "${GREEN}[1/2] Building Backend...${NC}"
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
fi

# Build Frontend
if [ "$FRONTEND_CHANGED" = true ]; then
    echo -e "${GREEN}[2/2] Building Frontend...${NC}"
    cd "$FRONTEND_DIR"

    if [ ! -f "package.json" ]; then
        echo -e "${RED}Error: package.json not found in $FRONTEND_DIR${NC}"
        exit 1
    fi

    # Check if node_modules exists
    if [ ! -d "node_modules" ]; then
        echo "Installing dependencies..."
        npm install
    fi

    # Create production .env if it doesn't exist
    if [ ! -f ".env.production" ]; then
        echo "Creating .env.production..."
        cat > .env.production << EOF
VITE_API_URL=https://api.daybook.shafik.xyz/api/v1
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
        # Exclude macOS extended attributes to avoid warnings on Linux
        COPYFILE_DISABLE=1 tar --no-xattrs -czf "$BUILD_DIR/frontend-dist.tar.gz" . 2>/dev/null || \
        tar -czf "$BUILD_DIR/frontend-dist.tar.gz" .
        echo -e "${GREEN}✓ Frontend tarball created${NC}"
        ls -lh "$BUILD_DIR/frontend-dist.tar.gz"
    else
        echo -e "${RED}✗ Frontend build failed${NC}"
        exit 1
    fi
    echo ""
fi

# Copy .env.example files
cp "$BACKEND_DIR/.env.example" "$BUILD_DIR/backend.env.example" 2>/dev/null || true
cp "$FRONTEND_DIR/.env.example" "$BUILD_DIR/frontend.env.example" 2>/dev/null || true

# Create deployment info
cat > "$BUILD_DIR/DEPLOYMENT_INFO.txt" << EOF
Daybook Deployment Package
Generated: $(date)
Hostname: $(hostname)
User: $(whoami)

Files in this package:
$(ls -lh "$BUILD_DIR" | grep -v "^total" | grep -v "DEPLOYMENT_INFO")

Backend Built: $BACKEND_CHANGED
Frontend Built: $FRONTEND_CHANGED

Git Information:
Branch: $(git branch --show-current)
Commit: $(git rev-parse --short HEAD)
Message: $(git log -1 --pretty=%B | head -1)
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

# Ask to commit and push
echo -e "${YELLOW}Next Steps:${NC}"
echo ""
read -p "Do you want to commit and push to git? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    cd "$PROJECT_DIR"

    # Add build directory
    git add build/

    # Add source changes if any
    [ "$BACKEND_CHANGED" = true ] && git add backend/ || true
    [ "$FRONTEND_CHANGED" = true ] && git add frontend/ || true

    # Get commit message
    echo ""
    echo "Enter commit message (or press Enter for default):"
    read -r COMMIT_MSG

    if [ -z "$COMMIT_MSG" ]; then
        COMMIT_MSG="Update deployment build - $(date +%Y-%m-%d)"
        if [ "$BACKEND_CHANGED" = true ] && [ "$FRONTEND_CHANGED" = true ]; then
            COMMIT_MSG="$COMMIT_MSG (backend + frontend)"
        elif [ "$BACKEND_CHANGED" = true ]; then
            COMMIT_MSG="$COMMIT_MSG (backend)"
        else
            COMMIT_MSG="$COMMIT_MSG (frontend)"
        fi
    fi

    # Commit
    git commit -m "$COMMIT_MSG" || echo "Nothing to commit"

    # Push
    echo ""
    echo -e "${YELLOW}Pushing to remote...${NC}"
    git push origin master

    if [ $? -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✓ Successfully pushed to remote${NC}"
        echo ""
        echo -e "${BLUE}=========================================="
        echo "Deploy on EC2 with:"
        echo "==========================================${NC}"
        echo "cd ~/daybook"
        echo "./scripts/deploy-ec2.sh"
        echo ""
    else
        echo -e "${RED}✗ Failed to push to remote${NC}"
        exit 1
    fi
else
    echo ""
    echo -e "${YELLOW}Skipped git commit/push${NC}"
    echo ""
    echo "To push manually:"
    echo "  git add build/"
    echo "  git commit -m 'Update deployment build'"
    echo "  git push origin master"
fi
