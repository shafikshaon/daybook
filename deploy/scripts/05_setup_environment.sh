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

log_info "Setting up application environment..."

# Create application user if not exists
if ! id "$APP_USER" &>/dev/null; then
    log_info "Creating application user: $APP_USER"
    sudo useradd -r -s /bin/bash -d "$APP_DIR" -m "$APP_USER"
    log_info "User $APP_USER created"
else
    log_warn "User $APP_USER already exists"
fi

# Create application directories
log_info "Creating application directories..."
sudo mkdir -p "$APP_DIR"
sudo mkdir -p "$BACKEND_DIR"
sudo mkdir -p "$FRONTEND_DIR"
sudo mkdir -p "$LOG_DIR"
sudo mkdir -p "$BACKEND_DIR/uploads"

# Set ownership
log_info "Setting directory ownership..."
sudo chown -R "$APP_USER:$APP_GROUP" "$APP_DIR"
sudo chown -R "$APP_USER:$APP_GROUP" "$LOG_DIR"

# Set permissions
sudo chmod 755 "$APP_DIR"
sudo chmod 755 "$BACKEND_DIR"
sudo chmod 755 "$FRONTEND_DIR"
sudo chmod 755 "$LOG_DIR"
sudo chmod 775 "$BACKEND_DIR/uploads"

log_info "Application directories created and configured:"
log_info "  - App directory: $APP_DIR"
log_info "  - Backend directory: $BACKEND_DIR"
log_info "  - Frontend directory: $FRONTEND_DIR"
log_info "  - Log directory: $LOG_DIR"

# Generate JWT secret if not set
if [ -z "$JWT_SECRET" ] || [ "$JWT_SECRET" == "CHANGE_THIS_SECRET_KEY_TO_RANDOM_STRING" ]; then
    JWT_SECRET=$(openssl rand -base64 64 | tr -d '\n')
    log_info "Generated JWT secret"
fi

# Create backend .env file
log_info "Creating backend environment file..."
BACKEND_ENV_TEMPLATE="${SCRIPT_DIR}/../config/production.env.template"
BACKEND_ENV_FILE="${BACKEND_DIR}/.env"

if [ ! -f "$BACKEND_ENV_TEMPLATE" ]; then
    log_error "Backend env template not found: $BACKEND_ENV_TEMPLATE"
    exit 1
fi

# Replace placeholders and create .env file
sudo cp "$BACKEND_ENV_TEMPLATE" "$BACKEND_ENV_FILE"
sudo sed -i "s/^SERVER_PORT=.*/SERVER_PORT=${BACKEND_PORT}/" "$BACKEND_ENV_FILE"
sudo sed -i "s/^DB_NAME=.*/DB_NAME=${DB_NAME}/" "$BACKEND_ENV_FILE"
sudo sed -i "s/^DB_USER=.*/DB_USER=${DB_USER}/" "$BACKEND_ENV_FILE"
sudo sed -i "s/^DB_PASSWORD=.*/DB_PASSWORD=${DB_PASSWORD}/" "$BACKEND_ENV_FILE"
sudo sed -i "s|^JWT_SECRET=.*|JWT_SECRET=${JWT_SECRET}|" "$BACKEND_ENV_FILE"

# Set secure permissions for .env file
sudo chown "$APP_USER:$APP_GROUP" "$BACKEND_ENV_FILE"
sudo chmod 600 "$BACKEND_ENV_FILE"

log_info "Backend environment file created: $BACKEND_ENV_FILE"

# Create frontend .env file
log_info "Creating frontend environment file..."
FRONTEND_ENV_TEMPLATE="${SCRIPT_DIR}/../config/frontend.env.template"
FRONTEND_ENV_FILE="${FRONTEND_DIR}/.env"

if [ ! -f "$FRONTEND_ENV_TEMPLATE" ]; then
    log_error "Frontend env template not found: $FRONTEND_ENV_TEMPLATE"
    exit 1
fi

# Replace placeholders
sudo cp "$FRONTEND_ENV_TEMPLATE" "$FRONTEND_ENV_FILE"

# Determine the API URL based on DOMAIN_OR_IP
if [[ "$DOMAIN_OR_IP" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    # It's an IP address
    API_URL="http://${DOMAIN_OR_IP}/api/v1"
else
    # It's a domain name
    API_URL="http://${DOMAIN_OR_IP}/api/v1"
fi

sudo sed -i "s|^VITE_API_URL=.*|VITE_API_URL=${API_URL}|" "$FRONTEND_ENV_FILE"

sudo chown "$APP_USER:$APP_GROUP" "$FRONTEND_ENV_FILE"
sudo chmod 644 "$FRONTEND_ENV_FILE"

log_info "Frontend environment file created: $FRONTEND_ENV_FILE"

log_info "Environment setup completed successfully!"

exit 0
