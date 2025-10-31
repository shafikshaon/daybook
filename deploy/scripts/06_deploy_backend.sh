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

# Get the project root directory (2 levels up from scripts)
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
BACKEND_SOURCE="${PROJECT_ROOT}/backend"

if [ ! -d "$BACKEND_SOURCE" ]; then
    log_error "Backend source directory not found: $BACKEND_SOURCE"
    exit 1
fi

log_info "Backend source: $BACKEND_SOURCE"
log_info "Backend destination: $BACKEND_DIR"

# Stop the service if it's running
if systemctl is-active --quiet daybook-backend; then
    log_info "Stopping daybook-backend service..."
    sudo systemctl stop daybook-backend
fi

# Copy backend files
log_info "Copying backend files..."
sudo rsync -av --delete \
    --exclude='.git' \
    --exclude='.idea' \
    --exclude='.claude' \
    --exclude='node_modules' \
    --exclude='*.log' \
    --exclude='.env' \
    --exclude='daybook-backend' \
    "${BACKEND_SOURCE}/" "${BACKEND_DIR}/"

# Ensure uploads directory exists and has correct permissions
sudo mkdir -p "${BACKEND_DIR}/uploads"

# Build the backend
log_info "Building backend application..."
cd "$BACKEND_DIR"

# Set Go environment for the app user
APP_USER_HOME=$(eval echo ~$APP_USER)
GOPATH_DIR="${APP_USER_HOME}/go"

# Create Go cache directory for app user
sudo mkdir -p "${GOPATH_DIR}"
sudo chown -R "$APP_USER:$APP_GROUP" "${GOPATH_DIR}"

# Download dependencies
log_info "Downloading Go dependencies..."
sudo -u "$APP_USER" env PATH=/usr/local/go/bin:$PATH HOME=${APP_USER_HOME} GOPATH=${GOPATH_DIR} go mod download

# Build the application
log_info "Compiling Go application..."
sudo -u "$APP_USER" env PATH=/usr/local/go/bin:$PATH HOME=${APP_USER_HOME} GOPATH=${GOPATH_DIR} go build -o daybook-backend -ldflags="-s -w" main.go

# Verify the binary was created
if [ ! -f "${BACKEND_DIR}/daybook-backend" ]; then
    log_error "Failed to build backend application"
    exit 1
fi

log_info "Backend application built successfully"

# Make the binary executable
sudo chmod +x "${BACKEND_DIR}/daybook-backend"

# Set ownership
log_info "Setting file ownership..."
sudo chown -R "$APP_USER:$APP_GROUP" "$BACKEND_DIR"

# Set proper permissions for uploads directory
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

log_info "Backend deployment completed successfully!"
log_info "Service: daybook-backend"
log_info "Status: sudo systemctl status daybook-backend"
log_info "Logs: sudo journalctl -u daybook-backend -f"
log_info "Log file: ${LOG_DIR}/backend.log"

exit 0
