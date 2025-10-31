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

# Verify source directory exists
if [ ! -d "$SOURCE_DIR/backend" ]; then
    log_error "Source backend directory not found: $SOURCE_DIR/backend"
    log_error "Please ensure the repository is cloned at $SOURCE_DIR"
    exit 1
fi

log_info "Source directory: $SOURCE_DIR/backend"
log_info "Build directory: $BUILD_BACKEND_DIR"
log_info "Deploy directory: $BACKEND_DIR"

# Create build directory
log_info "Preparing build directory..."
sudo rm -rf "$BUILD_BACKEND_DIR"
sudo mkdir -p "$BUILD_BACKEND_DIR"

# Copy backend source to build directory
log_info "Copying source code to build directory..."
sudo cp -r "$SOURCE_DIR/backend"/* "$BUILD_BACKEND_DIR/"

# Set Go environment for the app user
APP_USER_HOME=$(eval echo ~$APP_USER)
GOPATH_DIR="${APP_USER_HOME}/go"
GOCACHE_DIR="${GOPATH_DIR}/cache"

# Create Go cache directory for app user
sudo mkdir -p "${GOPATH_DIR}" "${GOCACHE_DIR}"
sudo chown -R "$APP_USER:$APP_GROUP" "${GOPATH_DIR}"
sudo chown -R "$APP_USER:$APP_GROUP" "$BUILD_BACKEND_DIR"

cd "$BUILD_BACKEND_DIR"

# Download dependencies
log_info "Downloading Go dependencies..."
sudo -u "$APP_USER" env \
    PATH=/usr/local/go/bin:$PATH \
    HOME=${APP_USER_HOME} \
    GOPATH=${GOPATH_DIR} \
    GOCACHE=${GOCACHE_DIR} \
    go mod download

# Build the application
log_info "Building Go application..."

# Stop the service before building (to avoid binary in-use issues on deploy dir)
if systemctl is-active --quiet daybook-backend; then
    log_info "Stopping daybook-backend service..."
    sudo systemctl stop daybook-backend
fi

# Build with optimizations
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
if [ ! -f "${BUILD_BACKEND_DIR}/daybook-backend" ]; then
    log_error "Failed to build backend application"
    exit 1
fi

log_info "Build successful!"

# Create deployment directory
log_info "Preparing deployment directory..."
sudo mkdir -p "$BACKEND_DIR"
sudo mkdir -p "${BACKEND_DIR}/uploads"

# Copy binary to deployment directory
log_info "Deploying binary..."
sudo cp "${BUILD_BACKEND_DIR}/daybook-backend" "$BACKEND_DIR/"
sudo chmod +x "${BACKEND_DIR}/daybook-backend"

# Copy .env file if it exists in build directory
if [ -f "${BUILD_BACKEND_DIR}/.env" ]; then
    sudo cp "${BUILD_BACKEND_DIR}/.env" "$BACKEND_DIR/"
fi

# Set ownership and permissions
log_info "Setting ownership and permissions..."
sudo chown -R "$APP_USER:$APP_GROUP" "$BACKEND_DIR"
sudo chmod 775 "${BACKEND_DIR}/uploads"

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

# Clean up build directory (optional, comment out if you want to keep it for debugging)
log_info "Cleaning up build directory..."
sudo rm -rf "$BUILD_BACKEND_DIR"

log_info "Backend deployment completed successfully!"
log_info "Service: daybook-backend"
log_info "Status: sudo systemctl status daybook-backend"
log_info "Logs: sudo journalctl -u daybook-backend -f"
log_info "Log file: ${LOG_DIR}/backend.log"

exit 0
