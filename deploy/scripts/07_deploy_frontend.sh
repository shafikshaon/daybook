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

# Verify source directory exists
if [ ! -d "$SOURCE_DIR/frontend" ]; then
    log_error "Source frontend directory not found: $SOURCE_DIR/frontend"
    log_error "Please ensure the repository is cloned at $SOURCE_DIR"
    exit 1
fi

log_info "Source directory: $SOURCE_DIR/frontend"
log_info "Build directory: $BUILD_FRONTEND_DIR"
log_info "Deploy directory: $FRONTEND_DIR"

# Create build directory
log_info "Preparing build directory..."
sudo rm -rf "$BUILD_FRONTEND_DIR"
sudo mkdir -p "$BUILD_FRONTEND_DIR"

# Copy frontend source to build directory
log_info "Copying source code to build directory..."
sudo cp -r "$SOURCE_DIR/frontend"/* "$BUILD_FRONTEND_DIR/"

# Set ownership
sudo chown -R "$APP_USER:$APP_GROUP" "$BUILD_FRONTEND_DIR"

cd "$BUILD_FRONTEND_DIR"

# Setup .env file if needed
if [ ! -f "${BUILD_FRONTEND_DIR}/.env" ] && [ ! -f "${BUILD_FRONTEND_DIR}/.env.production" ]; then
    log_warn "No .env file found, creating from template"
    if [ -f "${BUILD_FRONTEND_DIR}/.env.example" ]; then
        sudo -u "$APP_USER" cp "${BUILD_FRONTEND_DIR}/.env.example" "${BUILD_FRONTEND_DIR}/.env"

        # Update VITE_API_URL
        if [[ "$DOMAIN_OR_IP" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            API_URL="http://${DOMAIN_OR_IP}/api/v1"
        else
            API_URL="http://${DOMAIN_OR_IP}/api/v1"
        fi
        sudo sed -i "s|^VITE_API_URL=.*|VITE_API_URL=${API_URL}|" "${BUILD_FRONTEND_DIR}/.env"
    fi
else
    log_info "Using existing .env file"
fi

# Install npm dependencies
log_info "Installing npm dependencies..."
# Use npm install to handle missing package-lock.json gracefully
sudo -u "$APP_USER" npm install --production=false

# Build the frontend
log_info "Building frontend application..."
sudo -u "$APP_USER" npm run build

# Verify dist directory was created
if [ ! -d "${BUILD_FRONTEND_DIR}/dist" ]; then
    log_error "Build failed: dist directory not found"
    exit 1
fi

log_info "Build successful!"

# Create deployment directory
log_info "Preparing deployment directory..."
sudo mkdir -p "$FRONTEND_DIR"

# Deploy the built files
log_info "Deploying built files..."
sudo rm -rf "${FRONTEND_DIR}/dist"
sudo cp -r "${BUILD_FRONTEND_DIR}/dist" "$FRONTEND_DIR/"

# Set ownership and permissions
log_info "Setting ownership and permissions..."
sudo chown -R "$APP_USER:$APP_GROUP" "$FRONTEND_DIR"
sudo chmod -R 755 "${FRONTEND_DIR}/dist"

# Test Nginx configuration
log_info "Testing Nginx configuration..."
if sudo nginx -t; then
    log_info "Reloading Nginx..."
    sudo systemctl reload nginx
else
    log_error "Nginx configuration test failed!"
    exit 1
fi

# Clean up build directory (optional, comment out if you want to keep it for debugging)
log_info "Cleaning up build directory..."
sudo rm -rf "$BUILD_FRONTEND_DIR"

log_info "Frontend deployment completed successfully!"
log_info "Frontend URL: http://${DOMAIN_OR_IP}"
log_info "Nginx logs: sudo tail -f /var/log/nginx/daybook-*.log"

exit 0
