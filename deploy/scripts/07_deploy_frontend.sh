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

# Verify frontend directory exists (should be cloned from git)
if [ ! -d "$FRONTEND_DIR" ]; then
    log_error "Frontend directory not found: $FRONTEND_DIR"
    log_error "Please clone the repository to $FRONTEND_DIR first"
    exit 1
fi

log_info "Frontend directory: $FRONTEND_DIR"

# Setup .env file if needed
if [ ! -f "${FRONTEND_DIR}/.env" ] && [ ! -f "${FRONTEND_DIR}/.env.production" ]; then
    log_warn "No .env file found, creating from template"
    if [ -f "${FRONTEND_DIR}/.env.example" ]; then
        cp "${FRONTEND_DIR}/.env.example" "${FRONTEND_DIR}/.env"

        # Update VITE_API_URL
        if [[ "$DOMAIN_OR_IP" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            API_URL="http://${DOMAIN_OR_IP}/api/v1"
        else
            API_URL="http://${DOMAIN_OR_IP}/api/v1"
        fi
        sed -i "s|^VITE_API_URL=.*|VITE_API_URL=${API_URL}|" "${FRONTEND_DIR}/.env"
    fi
else
    log_info "Using existing .env file"
fi

# Build the frontend
cd "$FRONTEND_DIR"

log_info "Installing npm dependencies..."
if [ -f "package-lock.json" ]; then
    npm ci --production=false
else
    log_warn "No package-lock.json found, running npm install (not recommended for production)"
    npm install --production=false
    log_warn "Consider committing the generated package-lock.json to your repository"
fi

log_info "Building frontend application..."
npm run build

# Verify dist directory was created
if [ ! -d "${FRONTEND_DIR}/dist" ]; then
    log_error "Build failed: dist directory not found"
    exit 1
fi

log_info "Frontend build completed successfully"

# Set ownership
log_info "Setting file ownership..."
sudo chown -R "$APP_USER:$APP_GROUP" "$FRONTEND_DIR"

# Set permissions
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

log_info "Frontend deployment completed successfully!"
log_info "Frontend URL: http://${DOMAIN_OR_IP}"
log_info "Nginx logs: sudo tail -f /var/log/nginx/daybook-*.log"

exit 0
