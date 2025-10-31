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

log_info "Deploying backend application..."

# Verify backend directory exists (should be cloned from git)
if [ ! -d "$BACKEND_DIR" ]; then
    log_error "Backend directory not found: $BACKEND_DIR"
    log_error "Please clone the repository to $BACKEND_DIR first"
    exit 1
fi

log_info "Backend directory: $BACKEND_DIR"

# Set Go environment for the app user
APP_USER_HOME=$(eval echo ~$APP_USER)
GOPATH_DIR="${APP_USER_HOME}/go"
GOCACHE_DIR="${GOPATH_DIR}/cache"

# Create Go cache directory for app user
sudo mkdir -p "${GOPATH_DIR}" "${GOCACHE_DIR}"
sudo chown -R "$APP_USER:$APP_GROUP" "${GOPATH_DIR}"

cd "$BACKEND_DIR"

# Check if rebuild is needed (skip if no source changes)
LAST_BUILD_HASH_FILE="${BACKEND_DIR}/.last_build_hash"
CURRENT_HASH=$(git rev-parse HEAD 2>/dev/null || echo "no-git")

if [ -f "$LAST_BUILD_HASH_FILE" ] && [ -f "${BACKEND_DIR}/daybook-backend" ]; then
    LAST_HASH=$(cat "$LAST_BUILD_HASH_FILE")
    if [ "$LAST_HASH" = "$CURRENT_HASH" ]; then
        log_info "No source changes detected (hash: ${CURRENT_HASH}), skipping rebuild"
        log_info "To force rebuild, delete: $LAST_BUILD_HASH_FILE"
        SKIP_BUILD=true
    fi
fi

if [ "$SKIP_BUILD" != "true" ]; then
    # Download dependencies (only if go.mod changed)
    log_info "Checking dependencies..."
    sudo -u "$APP_USER" env \
        PATH=/usr/local/go/bin:$PATH \
        HOME=${APP_USER_HOME} \
        GOPATH=${GOPATH_DIR} \
        GOCACHE=${GOCACHE_DIR} \
        go mod download -x

    # Build the application with optimizations
    log_info "Compiling Go application..."

    # Stop the service before building (to avoid binary in-use issues)
    if systemctl is-active --quiet daybook-backend; then
        log_info "Stopping daybook-backend service..."
        sudo systemctl stop daybook-backend
    fi

    # Optimized build flags:
    # -p N: Use N parallel builds (defaults to GOMAXPROCS)
    # -ldflags="-s -w": Strip debug info and symbol table
    # -trimpath: Remove file system paths for reproducible builds
    # -buildvcs=false: Skip VCS info (faster)
    # CGO_ENABLED=0: Static binary (faster, no CGO overhead)
    sudo -u "$APP_USER" env \
        PATH=/usr/local/go/bin:$PATH \
        HOME=${APP_USER_HOME} \
        GOPATH=${GOPATH_DIR} \
        GOCACHE=${GOCACHE_DIR} \
        CGO_ENABLED=0 \
        go build \
        -p 4 \
        -trimpath \
        -buildvcs=false \
        -ldflags="-s -w -extldflags '-static'" \
        -o daybook-backend \
        main.go

    # Verify the binary was created
    if [ ! -f "${BACKEND_DIR}/daybook-backend" ]; then
        log_error "Failed to build backend application"
        exit 1
    fi

    # Save current hash
    echo "$CURRENT_HASH" > "$LAST_BUILD_HASH_FILE"

    log_info "Backend application built successfully (hash: ${CURRENT_HASH})"
else
    # Still stop the service for restart
    if systemctl is-active --quiet daybook-backend; then
        log_info "Stopping daybook-backend service..."
        sudo systemctl stop daybook-backend
    fi
fi

# Make the binary executable
sudo chmod +x "${BACKEND_DIR}/daybook-backend"

# Ensure uploads directory exists and has correct permissions
sudo mkdir -p "${BACKEND_DIR}/uploads"
sudo chmod 775 "${BACKEND_DIR}/uploads"

# Set ownership
log_info "Setting file ownership..."
sudo chown -R "$APP_USER:$APP_GROUP" "$BACKEND_DIR"

# Create systemd service from template
log_info "Creating systemd service..."
SERVICE_TEMPLATE="${SCRIPT_DIR}/../systemd/daybook-backend.service.template"
SERVICE_FILE="/etc/systemd/system/daybook-backend.service"

if [ ! -f "$SERVICE_TEMPLATE" ]; then
    log_error "Service template not found: $SERVICE_TEMPLATE"
    exit 1
fi

# Replace placeholders in service template
sudo cp "$SERVICE_TEMPLATE" "$SERVICE_FILE"
sudo sed -i "s|APP_USER|${APP_USER}|g" "$SERVICE_FILE"
sudo sed -i "s|APP_GROUP|${APP_GROUP}|g" "$SERVICE_FILE"
sudo sed -i "s|BACKEND_DIR|${BACKEND_DIR}|g" "$SERVICE_FILE"
sudo sed -i "s|LOG_DIR|${LOG_DIR}|g" "$SERVICE_FILE"

# Reload systemd
log_info "Reloading systemd..."
sudo systemctl daemon-reload

# Enable and start the service
log_info "Starting daybook-backend service..."
sudo systemctl enable daybook-backend
sudo systemctl start daybook-backend

# Wait a moment for the service to start
sleep 3

# Check service status
if systemctl is-active --quiet daybook-backend; then
    log_info "Backend service is running!"
    sudo systemctl status daybook-backend --no-pager -l
else
    log_error "Backend service failed to start!"
    log_error "Check logs with: sudo journalctl -u daybook-backend -n 50"
    exit 1
fi

log_info "Backend deployment completed successfully!"
log_info "Service: daybook-backend"
log_info "Status: sudo systemctl status daybook-backend"
log_info "Logs: sudo journalctl -u daybook-backend -f"
log_info "Log file: ${LOG_DIR}/backend.log"

exit 0