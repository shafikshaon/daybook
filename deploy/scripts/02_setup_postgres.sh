#!/bin/bash
set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/../config/deploy.conf"

if [ -f "$CONFIG_FILE" ]; then
    source "$CONFIG_FILE"
else
    log_error "Configuration file not found: $CONFIG_FILE"
    exit 1
fi

log_info "Installing PostgreSQL 18..."

# Add PostgreSQL APT repository
sudo apt-get install -y wget ca-certificates
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'

# Update and install PostgreSQL 18
sudo apt-get update
sudo apt-get install -y postgresql-18 postgresql-contrib-18

# Start and enable PostgreSQL
log_info "Starting PostgreSQL service..."
sudo systemctl start postgresql
sudo systemctl enable postgresql
sleep 3

# Use the specified password
DB_PASSWORD='$hafik$ha0n'

# Create database user and database
log_info "Creating database user and database..."
sudo -u postgres psql <<EOF
-- Create user if not exists
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = '$DB_USER') THEN
        CREATE USER $DB_USER WITH PASSWORD E'${DB_PASSWORD//\'/\'\'}';
    ELSE
        ALTER USER $DB_USER WITH PASSWORD E'${DB_PASSWORD//\'/\'\'}';
    END IF;
END
\$\$;

-- Create database if not exists
SELECT 'CREATE DATABASE $DB_NAME OWNER $DB_USER'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME')\gexec

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
EOF

# Grant schema privileges
sudo -u postgres psql -d $DB_NAME <<EOF
GRANT ALL ON SCHEMA public TO $DB_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $DB_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO $DB_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO $DB_USER;
EOF

# Create a system user matching the DB user for peer authentication
log_info "Creating system user for peer authentication..."
if ! id "$DB_USER" &>/dev/null; then
    sudo useradd -r -s /bin/bash "$DB_USER"
    log_info "System user $DB_USER created"
fi

# Test connection using peer authentication
log_info "Testing database connection (peer authentication)..."
if sudo -u "$DB_USER" psql -d "$DB_NAME" -c '\q' 2>/dev/null; then
    log_info "Peer authentication successful!"
else
    log_warn "Peer authentication test skipped (user may need shell access)"
fi

# Test connection as postgres user
log_info "Testing database connection as postgres..."
if sudo -u postgres psql -d "$DB_NAME" -c '\q' 2>/dev/null; then
    log_info "Database connection successful!"
else
    log_error "Database connection failed!"
    exit 1
fi

log_info "PostgreSQL 18 setup completed successfully!"
log_info "Database: $DB_NAME"
log_info "User: $DB_USER"
log_info "Password: $DB_PASSWORD"
log_info ""
log_info "Connection methods:"
log_info "  1. Peer authentication: sudo -u $DB_USER psql -d $DB_NAME"
log_info "  2. Unix socket with password: psql postgresql://$DB_USER:PASSWORD@/localhost/$DB_NAME"
log_info "  3. From your app: Use unix socket /var/run/postgresql or localhost without host parameter"

exit 0