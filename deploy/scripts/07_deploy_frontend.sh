#!/bin/bash
set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Load configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/../config/deploy.conf"

if [ -f "$CONFIG_FILE" ]; then
    source "$CONFIG_FILE"
else
    log_error "Configuration file not found: $CONFIG_FILE"
    exit 1
fi

log_info "Deploying frontend application..."

# Get the project root directory (2 levels up from scripts)
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
FRONTEND_SOURCE="${PROJECT_ROOT}/frontend"

if [ ! -d "$FRONTEND_SOURCE" ]; then
    log_error "Frontend source directory not found: $FRONTEND_SOURCE"
    exit 1
fi

log_info "Frontend source: $FRONTEND_SOURCE"
log_info "Frontend destination: $FRONTEND_DIR"

# Create a temporary directory for building
TEMP_BUILD_DIR="/tmp/daybook-frontend-build"
rm -rf "$TEMP_BUILD_DIR"
mkdir -p "$TEMP_BUILD_DIR"

# Copy frontend files to temp directory
log_info "Copying frontend files to temporary build directory..."
rsync -av \
    --exclude='.git' \
    --exclude='.idea' \
    --exclude='.claude' \
    --exclude='node_modules' \
    --exclude='dist' \
    --exclude='*.log' \
    "${FRONTEND_SOURCE}/" "${TEMP_BUILD_DIR}/"

# Copy the .env file if it exists in the destination
if [ -f "${FRONTEND_DIR}/.env" ]; then
    log_info "Using existing .env file"
    cp "${FRONTEND_DIR}/.env" "${TEMP_BUILD_DIR}/.env"
elif [ -f "${FRONTEND_DIR}/.env.production" ]; then
    log_info "Using existing .env.production file"
    cp "${FRONTEND_DIR}/.env.production" "${TEMP_BUILD_DIR}/.env"
else
    log_warn "No .env file found, using template"
    if [ -f "${TEMP_BUILD_DIR}/.env.example" ]; then
        cp "${TEMP_BUILD_DIR}/.env.example" "${TEMP_BUILD_DIR}/.env"

        # Update VITE_API_URL
        if [[ "$DOMAIN_OR_IP" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            API_URL="http://${DOMAIN_OR_IP}/api/v1"
        else
            API_URL="http://${DOMAIN_OR_IP}/api/v1"
        fi
        sed -i "s|^VITE_API_URL=.*|VITE_API_URL=${API_URL}|" "${TEMP_BUILD_DIR}/.env"
    fi
fi

# Build the frontend
cd "$TEMP_BUILD_DIR"

log_info "Installing npm dependencies..."
npm ci --production=false

log_info "Building frontend application..."
npm run build

# Verify dist directory was created
if [ ! -d "${TEMP_BUILD_DIR}/dist" ]; then
    log_error "Build failed: dist directory not found"
    exit 1
fi

log_info "Frontend build completed successfully"

# Copy dist to destination
log_info "Deploying built files..."
sudo mkdir -p "${FRONTEND_DIR}/dist"
sudo rsync -av --delete "${TEMP_BUILD_DIR}/dist/" "${FRONTEND_DIR}/dist/"

# Copy .env file to destination
if [ -f "${TEMP_BUILD_DIR}/.env" ]; then
    sudo cp "${TEMP_BUILD_DIR}/.env" "${FRONTEND_DIR}/.env"
fi

# Set ownership
log_info "Setting file ownership..."
sudo chown -R "$APP_USER:$APP_GROUP" "$FRONTEND_DIR"

# Set permissions
sudo chmod -R 755 "${FRONTEND_DIR}/dist"

# Clean up temp directory
log_info "Cleaning up temporary files..."
rm -rf "$TEMP_BUILD_DIR"

# Test Nginx configuration
log_info "Testing Nginx configuration..."
if sudo nginx -t; then
    log_info "Reloading Nginx..."
    sudo systemctl reload nginx
else
    log_error "Nginx configuration test failed!"
    exit 1
fi

log_info "Frontend deployment completed successfully!"
log_info "Frontend URL: http://${DOMAIN_OR_IP}"
log_info "Nginx logs: sudo tail -f /var/log/nginx/daybook-*.log"

exit 0
