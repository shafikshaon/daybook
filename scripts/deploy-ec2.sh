#!/bin/bash
#
# EC2 Deployment Script
# Pulls latest changes and deploys backend/frontend
#

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=========================================="
echo "Daybook Auto-Deploy (EC2)"
echo "==========================================${NC}"
echo ""

# Check if running as ubuntu user
if [ "$USER" != "ubuntu" ]; then
    echo -e "${YELLOW}⚠ Warning: Not running as ubuntu user${NC}"
fi

# Variables
REPO_DIR="/home/ubuntu/daybook"
DEPLOY_DIR="/var/www/daybook"
BACKEND_DIR="$DEPLOY_DIR/backend"
FRONTEND_DIR="$DEPLOY_DIR/frontend"
BUILD_DIR="$REPO_DIR/build"

# Check if repo exists
if [ ! -d "$REPO_DIR" ]; then
    echo -e "${RED}Error: Repository not found at $REPO_DIR${NC}"
    echo "Please clone the repository first:"
    echo "  cd ~"
    echo "  git clone https://github.com/shafikshaon/daybook.git"
    exit 1
fi

# Pull latest changes
echo -e "${GREEN}[1/5] Pulling latest changes from git...${NC}"
cd "$REPO_DIR"

# Show current commit
echo "Current commit: $(git rev-parse --short HEAD)"

# Pull
git pull origin master

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Failed to pull changes${NC}"
    exit 1
fi

echo "New commit: $(git rev-parse --short HEAD)"
echo -e "${GREEN}✓ Successfully pulled latest changes${NC}"
echo ""

# Detect what changed
BACKEND_CHANGED=false
FRONTEND_CHANGED=false

if [ -f "$BUILD_DIR/daybook-backend" ]; then
    # Check if backend binary is newer than deployed one
    if [ ! -f "$BACKEND_DIR/daybook-backend" ] || \
       [ "$BUILD_DIR/daybook-backend" -nt "$BACKEND_DIR/daybook-backend" ]; then
        BACKEND_CHANGED=true
    fi
fi

if [ -f "$BUILD_DIR/frontend-dist.tar.gz" ]; then
    FRONTEND_CHANGED=true
fi

# Force deploy option
if [ "$1" == "--force" ] || [ "$1" == "-f" ]; then
    BACKEND_CHANGED=true
    FRONTEND_CHANGED=true
    echo -e "${YELLOW}✓ Force deploy enabled${NC}"
fi

echo -e "${BLUE}Deployment Plan:${NC}"
[ "$BACKEND_CHANGED" = true ] && echo "  ✓ Backend will be deployed" || echo "  - Backend unchanged"
[ "$FRONTEND_CHANGED" = true ] && echo "  ✓ Frontend will be deployed" || echo "  - Frontend unchanged"
echo ""

if [ "$BACKEND_CHANGED" = false ] && [ "$FRONTEND_CHANGED" = false ]; then
    echo -e "${YELLOW}No changes detected. Nothing to deploy.${NC}"
    exit 0
fi

# Deploy Backend
if [ "$BACKEND_CHANGED" = true ]; then
    echo -e "${GREEN}[2/5] Deploying Backend...${NC}"

    # Check if binary exists
    if [ ! -f "$BUILD_DIR/daybook-backend" ]; then
        echo -e "${RED}Error: Backend binary not found in $BUILD_DIR${NC}"
        exit 1
    fi

    # Stop backend service
    echo "Stopping backend service..."
    sudo systemctl stop daybook-backend

    # Deploy binary
    echo "Copying backend binary..."
    sudo cp "$BUILD_DIR/daybook-backend" "$BACKEND_DIR/daybook-backend"
    sudo chmod +x "$BACKEND_DIR/daybook-backend"
    sudo chown ubuntu:ubuntu "$BACKEND_DIR/daybook-backend"

    echo -e "${GREEN}✓ Backend deployed${NC}"
    echo ""
fi

# Deploy Frontend
if [ "$FRONTEND_CHANGED" = true ]; then
    echo -e "${GREEN}[3/5] Deploying Frontend...${NC}"

    # Check if tarball exists
    if [ ! -f "$BUILD_DIR/frontend-dist.tar.gz" ]; then
        echo -e "${RED}Error: Frontend tarball not found in $BUILD_DIR${NC}"
        exit 1
    fi

    # Extract frontend
    echo "Extracting frontend files..."
    cd "$FRONTEND_DIR"
    sudo rm -rf *
    sudo tar -xzf "$BUILD_DIR/frontend-dist.tar.gz"
    sudo chown -R www-data:www-data "$FRONTEND_DIR"

    echo -e "${GREEN}✓ Frontend deployed${NC}"
    echo ""
fi

# Restart Services
echo -e "${GREEN}[4/5] Restarting Services...${NC}"

# Start backend if it was deployed
if [ "$BACKEND_CHANGED" = true ]; then
    echo "Starting backend service..."
    sudo systemctl start daybook-backend
    sleep 2

    # Check if backend is running
    if sudo systemctl is-active --quiet daybook-backend; then
        echo -e "${GREEN}✓ Backend service started${NC}"
    else
        echo -e "${RED}✗ Backend service failed to start${NC}"
        echo "Check logs with: sudo journalctl -u daybook-backend -n 50"
        exit 1
    fi
fi

# Reload Nginx
echo "Reloading Nginx..."
sudo systemctl reload nginx

echo -e "${GREEN}✓ Services restarted${NC}"
echo ""

# Verify Deployment
echo -e "${GREEN}[5/5] Verifying Deployment...${NC}"

# Check backend
if [ "$BACKEND_CHANGED" = true ]; then
    echo -n "Backend service: "
    if sudo systemctl is-active --quiet daybook-backend; then
        echo -e "${GREEN}✓ Running${NC}"

        # Check port
        if sudo netstat -tulpn | grep -q ":8082"; then
            echo -e "  Port 8082: ${GREEN}✓ Listening${NC}"
        else
            echo -e "  Port 8082: ${RED}✗ Not listening${NC}"
        fi
    else
        echo -e "${RED}✗ Not running${NC}"
    fi
fi

# Check frontend
if [ "$FRONTEND_CHANGED" = true ]; then
    echo -n "Frontend files: "
    if [ -f "$FRONTEND_DIR/index.html" ]; then
        echo -e "${GREEN}✓ Deployed${NC}"
    else
        echo -e "${RED}✗ Missing${NC}"
    fi
fi

# Check Nginx
echo -n "Nginx: "
if sudo systemctl is-active --quiet nginx; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Not running${NC}"
fi

echo ""
echo -e "${GREEN}=========================================="
echo "Deployment Complete!"
echo "==========================================${NC}"
echo ""

# Show URLs
echo -e "${BLUE}Your application is available at:${NC}"
echo "  Frontend: https://daybook.shafik.xyz"
echo "  Backend:  https://api.daybook.shafik.xyz"
echo ""

# Show useful commands
echo -e "${YELLOW}Useful commands:${NC}"
echo "  View backend logs:  sudo journalctl -u daybook-backend -f"
echo "  View nginx logs:    sudo tail -f /var/log/nginx/daybook-backend-error.log"
echo "  Restart backend:    sudo systemctl restart daybook-backend"
echo "  Check services:     sudo systemctl status daybook-backend nginx"
echo ""

# Show deployment info if available
if [ -f "$BUILD_DIR/DEPLOYMENT_INFO.txt" ]; then
    echo -e "${BLUE}Deployment Info:${NC}"
    cat "$BUILD_DIR/DEPLOYMENT_INFO.txt" | grep -E "Generated|Commit|Message" | sed 's/^/  /'
    echo ""
fi
