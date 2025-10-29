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

# Create backup directory
BACKUP_DIR="${APP_DIR}/backups"
sudo mkdir -p "$BACKUP_DIR"

# Generate backup filename with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/daybook_backup_${TIMESTAMP}.sql"

log_info "Creating database backup..."
log_info "Database: $DB_NAME"
log_info "Backup file: $BACKUP_FILE"

# Create backup
PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > "$BACKUP_FILE"

# Compress backup
log_info "Compressing backup..."
gzip "$BACKUP_FILE"
BACKUP_FILE="${BACKUP_FILE}.gz"

# Set permissions
sudo chown "$APP_USER:$APP_GROUP" "$BACKUP_FILE"
sudo chmod 600 "$BACKUP_FILE"

# Display backup info
BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
log_info "Backup completed successfully!"
log_info "File: $BACKUP_FILE"
log_info "Size: $BACKUP_SIZE"

# Optional: Keep only last 7 backups
BACKUP_COUNT=$(ls -1 "${BACKUP_DIR}"/daybook_backup_*.sql.gz 2>/dev/null | wc -l)
if [ "$BACKUP_COUNT" -gt 7 ]; then
    log_info "Removing old backups (keeping last 7)..."
    ls -1t "${BACKUP_DIR}"/daybook_backup_*.sql.gz | tail -n +8 | xargs rm -f
fi

exit 0
