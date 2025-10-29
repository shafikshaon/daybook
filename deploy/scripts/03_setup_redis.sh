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

log_info "Setting up Redis..."

# Install Redis
log_info "Installing Redis..."
sudo apt-get install -y redis-server

# Configure Redis
log_info "Configuring Redis..."
REDIS_CONF="/etc/redis/redis.conf"

# Backup original configuration
sudo cp "$REDIS_CONF" "${REDIS_CONF}.backup"
log_info "Redis configuration backed up to ${REDIS_CONF}.backup"

# Configure Redis for systemd
sudo sed -i 's/^supervised no/supervised systemd/' "$REDIS_CONF"

# Bind to localhost only for security
sudo sed -i 's/^bind .*/bind 127.0.0.1/' "$REDIS_CONF"

# Set maxmemory policy (optional)
if ! grep -q "^maxmemory " "$REDIS_CONF"; then
    echo "maxmemory 256mb" | sudo tee -a "$REDIS_CONF" > /dev/null
    echo "maxmemory-policy allkeys-lru" | sudo tee -a "$REDIS_CONF" > /dev/null
    log_info "Set Redis maxmemory to 256MB with LRU eviction policy"
fi

# Enable persistence (optional)
if ! grep -q "^appendonly " "$REDIS_CONF"; then
    echo "appendonly yes" | sudo tee -a "$REDIS_CONF" > /dev/null
    log_info "Enabled Redis AOF persistence"
fi

# Start and enable Redis
log_info "Starting Redis service..."
sudo systemctl restart redis-server
sudo systemctl enable redis-server

# Wait for Redis to be ready
sleep 2

# Test Redis connection
log_info "Testing Redis connection..."
if redis-cli ping | grep -q "PONG"; then
    log_info "Redis connection successful!"
else
    log_error "Redis connection failed!"
    exit 1
fi

# Display Redis info
log_info "Redis setup completed successfully!"
redis-cli INFO server | grep "redis_version"
redis-cli INFO memory | grep "used_memory_human"

exit 0
