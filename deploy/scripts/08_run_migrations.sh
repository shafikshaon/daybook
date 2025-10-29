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

log_info "Running database migrations..."

# Check if migrations directory exists
MIGRATIONS_DIR="${BACKEND_DIR}/migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
    log_warn "Migrations directory not found: $MIGRATIONS_DIR"
    log_info "Skipping migrations"
    exit 0
fi

# Load database credentials from .env file
if [ -f "${BACKEND_DIR}/.env" ]; then
    source "${BACKEND_DIR}/.env"
else
    log_error "Backend .env file not found: ${BACKEND_DIR}/.env"
    exit 1
fi

# Create migrations tracking table if it doesn't exist
log_info "Creating migrations tracking table..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME <<EOF
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    version VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
EOF

# Get list of applied migrations
APPLIED_MIGRATIONS=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT version FROM schema_migrations ORDER BY version;")

# Run pending migrations
log_info "Checking for pending migrations..."
MIGRATION_COUNT=0

for migration_file in $(ls -1 "$MIGRATIONS_DIR"/*.sql 2>/dev/null | sort); do
    migration_name=$(basename "$migration_file")
    migration_version=$(echo "$migration_name" | sed 's/\.sql$//')

    # Check if migration has been applied
    if echo "$APPLIED_MIGRATIONS" | grep -q "$migration_version"; then
        log_info "Skipping already applied migration: $migration_name"
        continue
    fi

    # Apply migration
    log_info "Applying migration: $migration_name"

    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f "$migration_file"; then
        # Record successful migration
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "INSERT INTO schema_migrations (version) VALUES ('$migration_version');"
        log_info "Successfully applied migration: $migration_name"
        ((MIGRATION_COUNT++))
    else
        log_error "Failed to apply migration: $migration_name"
        exit 1
    fi
done

if [ $MIGRATION_COUNT -eq 0 ]; then
    log_info "No pending migrations found. Database is up to date."
else
    log_info "Successfully applied $MIGRATION_COUNT migration(s)"
fi

# Display current schema version
log_info "Current database schema versions:"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT version, applied_at FROM schema_migrations ORDER BY applied_at DESC LIMIT 5;"

log_info "Database migrations completed successfully!"

exit 0
