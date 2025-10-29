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

log_info "Setting up Nginx..."

# Install Nginx
log_info "Installing Nginx..."
sudo apt-get install -y nginx

# Create Nginx configuration from template
log_info "Creating Nginx configuration..."
NGINX_TEMPLATE="${SCRIPT_DIR}/../nginx/daybook.conf.template"
NGINX_CONF="/etc/nginx/sites-available/daybook"

if [ ! -f "$NGINX_TEMPLATE" ]; then
    log_error "Nginx template not found: $NGINX_TEMPLATE"
    exit 1
fi

# Replace placeholders in template
sudo cp "$NGINX_TEMPLATE" "$NGINX_CONF"
sudo sed -i "s|DOMAIN_OR_IP|${DOMAIN_OR_IP}|g" "$NGINX_CONF"
sudo sed -i "s|FRONTEND_DIR|${FRONTEND_DIR}|g" "$NGINX_CONF"
sudo sed -i "s|BACKEND_DIR|${BACKEND_DIR}|g" "$NGINX_CONF"

# Enable the site
log_info "Enabling Nginx site..."
sudo ln -sf "$NGINX_CONF" /etc/nginx/sites-enabled/daybook

# Remove default site if exists
if [ -f "/etc/nginx/sites-enabled/default" ]; then
    sudo rm -f /etc/nginx/sites-enabled/default
    log_info "Removed default Nginx site"
fi

# Test Nginx configuration
log_info "Testing Nginx configuration..."
if sudo nginx -t; then
    log_info "Nginx configuration is valid"
else
    log_error "Nginx configuration is invalid!"
    exit 1
fi

# Start and enable Nginx
log_info "Starting Nginx service..."
sudo systemctl restart nginx
sudo systemctl enable nginx

# Configure firewall
log_info "Configuring firewall..."
sudo ufw allow 'Nginx Full'
sudo ufw allow OpenSSH
echo "y" | sudo ufw enable || true

log_info "Nginx setup completed successfully!"
log_info "Nginx is listening on port 80"
log_info "Site configuration: $NGINX_CONF"

exit 0
