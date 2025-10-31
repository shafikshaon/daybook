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

log_info "Checking database migrations..."

# Note: This application uses GORM's AutoMigrate feature
# Migrations are automatically applied when the backend service starts
# See backend/database/database.go -> InitDatabase() for details

# Check if backend service is running
if systemctl is-active --quiet daybook-backend; then
    log_info "Backend service is running"
    log_info "Database migrations are automatically handled by the backend on startup"
    log_info "Check backend logs to verify migrations: sudo journalctl -u daybook-backend -n 50"
else
    log_warn "Backend service is not running"
    log_warn "Migrations will be applied when the backend service starts"
fi

# You can also verify database connection
log_info "Verifying database connection..."
if PGPASSWORD="${DB_PASSWORD}" psql -h localhost -U "${DB_USER}" -d "${DB_NAME}" -c "SELECT 1" > /dev/null 2>&1; then
    log_info "Database connection successful"

    # Show tables in database
    log_info "Current database tables:"
    PGPASSWORD="${DB_PASSWORD}" psql -h localhost -U "${DB_USER}" -d "${DB_NAME}" -c "\dt" 2>/dev/null || true
else
    log_error "Failed to connect to database"
    log_error "Please check database configuration and ensure PostgreSQL is running"
    exit 1
fi

log_info "Migration check completed"
log_info "Note: This application uses automatic migrations via GORM AutoMigrate"

exit 0
