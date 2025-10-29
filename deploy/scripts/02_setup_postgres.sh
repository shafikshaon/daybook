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

log_info "Setting up PostgreSQL..."

# Install PostgreSQL
log_info "Installing PostgreSQL..."
sudo apt-get install -y postgresql postgresql-contrib

# Start and enable PostgreSQL
log_info "Starting PostgreSQL service..."
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Wait for PostgreSQL to be ready
sleep 3

# Generate random password if not set
if [ -z "$DB_PASSWORD" ]; then
    DB_PASSWORD=$(openssl rand -base64 32)
    log_info "Generated database password: $DB_PASSWORD"
    log_warn "IMPORTANT: Save this password! It will be written to the backend .env file"

    # Update deploy.conf with the generated password
    sed -i "s/^DB_PASSWORD=.*/DB_PASSWORD=\"$DB_PASSWORD\"/" "$CONFIG_FILE"
fi

# Create database user and database
log_info "Creating database user and database..."
sudo -u postgres psql <<EOF
-- Create user if not exists
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = '$DB_USER') THEN
        CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
    END IF;
END
\$\$;

-- Create database if not exists
SELECT 'CREATE DATABASE $DB_NAME OWNER $DB_USER'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME')\gexec

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;

-- Connect to the database and grant schema privileges
\c $DB_NAME
GRANT ALL ON SCHEMA public TO $DB_USER;
EOF

# Configure PostgreSQL for local connections (optional)
log_info "Configuring PostgreSQL..."

# Backup original pg_hba.conf
PG_VERSION=$(psql --version | awk '{print $3}' | cut -d. -f1)
PG_HBA_CONF="/etc/postgresql/${PG_VERSION}/main/pg_hba.conf"

if [ -f "$PG_HBA_CONF" ]; then
    sudo cp "$PG_HBA_CONF" "${PG_HBA_CONF}.backup"
    log_info "PostgreSQL configuration backed up to ${PG_HBA_CONF}.backup"
fi

# Restart PostgreSQL to apply changes
log_info "Restarting PostgreSQL..."
sudo systemctl restart postgresql

# Test connection
log_info "Testing database connection..."
if PGPASSWORD=$DB_PASSWORD psql -h localhost -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
    log_info "Database connection successful!"
else
    log_error "Database connection failed!"
    exit 1
fi

log_info "PostgreSQL setup completed successfully!"
log_info "Database: $DB_NAME"
log_info "User: $DB_USER"
log_info "Password: $DB_PASSWORD"

exit 0
