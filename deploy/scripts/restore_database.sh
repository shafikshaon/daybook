#!/bin/bash
set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

# Load database credentials
if [ -f "${BACKEND_DIR}/.env" ]; then
    source "${BACKEND_DIR}/.env"
else
    log_error "Backend .env file not found"
    exit 1
fi

# Check if backup file is provided
if [ -z "$1" ]; then
    log_error "Usage: $0 <backup_file>"
    log_info "Available backups:"
    ls -lh "${APP_DIR}/backups"/daybook_backup_*.sql.gz 2>/dev/null || echo "No backups found"
    exit 1
fi

BACKUP_FILE="$1"

if [ ! -f "$BACKUP_FILE" ]; then
    log_error "Backup file not found: $BACKUP_FILE"
    exit 1
fi

log_warn "WARNING: This will restore the database from backup!"
log_warn "Current database will be dropped and recreated!"
log_warn "Backup file: $BACKUP_FILE"
echo ""
read -p "Are you sure you want to continue? (yes/no) " -r
echo
if [ "$REPLY" != "yes" ]; then
    log_info "Restore cancelled"
    exit 0
fi

# Stop backend service
log_info "Stopping backend service..."
sudo systemctl stop daybook-backend

# Create a backup of current database before restore
log_info "Creating backup of current database..."
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
SAFETY_BACKUP="${APP_DIR}/backups/pre_restore_backup_${TIMESTAMP}.sql"
PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > "$SAFETY_BACKUP"
gzip "$SAFETY_BACKUP"
log_info "Safety backup created: ${SAFETY_BACKUP}.gz"

# Drop and recreate database
log_info "Recreating database..."
sudo -u postgres psql <<EOF
DROP DATABASE IF EXISTS $DB_NAME;
CREATE DATABASE $DB_NAME OWNER $DB_USER;
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
EOF

# Restore from backup
log_info "Restoring database from backup..."
if [[ "$BACKUP_FILE" == *.gz ]]; then
    gunzip -c "$BACKUP_FILE" | PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME
else
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME < "$BACKUP_FILE"
fi

# Start backend service
log_info "Starting backend service..."
sudo systemctl start daybook-backend

# Wait for service to start
sleep 3

if systemctl is-active --quiet daybook-backend; then
    log_info "Database restored successfully!"
    log_info "Backend service is running"
else
    log_error "Backend service failed to start after restore"
    log_error "Check logs: sudo journalctl -u daybook-backend -n 50"
    exit 1
fi

exit 0
