#!/bin/bash
set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_step() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="${SCRIPT_DIR}/scripts"
CONFIG_FILE="${SCRIPT_DIR}/config/deploy.conf"

# Display banner
cat << "EOF"
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║           Daybook Application Deployment             ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝
EOF

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    log_error "Do not run this script as root. It will use sudo when needed."
    exit 1
fi

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    log_error "Configuration file not found: $CONFIG_FILE"
    log_error "Please copy deploy/config/deploy.conf and configure it"
    exit 1
fi

# Source configuration
source "$CONFIG_FILE"

# Validate configuration
log_info "Validating configuration..."
if [ "$DOMAIN_OR_IP" == "YOUR_DOMAIN_OR_IP" ]; then
    log_error "Please configure DOMAIN_OR_IP in $CONFIG_FILE"
    exit 1
fi

# Display configuration
log_info "Deployment Configuration:"
echo "  App Name: $APP_NAME"
echo "  App Directory: $APP_DIR"
echo "  Domain/IP: $DOMAIN_OR_IP"
echo "  Backend Port: $BACKEND_PORT"
echo "  Database: $DB_NAME"
echo ""

# Parse command line arguments
FRESH_INSTALL=false
SKIP_DEPS=false
SKIP_DB=false
SKIP_BACKEND=false
SKIP_FRONTEND=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --fresh)
            FRESH_INSTALL=true
            shift
            ;;
        --skip-deps)
            SKIP_DEPS=true
            shift
            ;;
        --skip-db)
            SKIP_DB=true
            shift
            ;;
        --skip-backend)
            SKIP_BACKEND=true
            shift
            ;;
        --skip-frontend)
            SKIP_FRONTEND=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --fresh          Fresh installation (install all dependencies)"
            echo "  --skip-deps      Skip dependency installation"
            echo "  --skip-db        Skip database setup"
            echo "  --skip-backend   Skip backend deployment"
            echo "  --skip-frontend  Skip frontend deployment"
            echo "  --help           Show this help message"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            log_info "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Confirm deployment
echo ""
read -p "Do you want to proceed with deployment? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log_info "Deployment cancelled"
    exit 0
fi

# Start deployment
START_TIME=$(date +%s)

# Step 1: Install dependencies
if [ "$FRESH_INSTALL" = true ] || [ "$SKIP_DEPS" = false ]; then
    log_step "Step 1: Installing System Dependencies"
    bash "${SCRIPTS_DIR}/01_install_dependencies.sh"
else
    log_warn "Skipping dependency installation"
fi

# Step 2: Setup PostgreSQL
if [ "$FRESH_INSTALL" = true ] || [ "$SKIP_DB" = false ]; then
    log_step "Step 2: Setting up PostgreSQL"
    bash "${SCRIPTS_DIR}/02_setup_postgres.sh"
else
    log_warn "Skipping PostgreSQL setup"
fi

# Step 3: Setup Redis
if [ "$FRESH_INSTALL" = true ] || [ "$SKIP_DB" = false ]; then
    log_step "Step 3: Setting up Redis"
    bash "${SCRIPTS_DIR}/03_setup_redis.sh"
else
    log_warn "Skipping Redis setup"
fi

# Step 4: Setup Nginx
if [ "$FRESH_INSTALL" = true ]; then
    log_step "Step 4: Setting up Nginx"
    bash "${SCRIPTS_DIR}/04_setup_nginx.sh"
fi

# Step 5: Setup environment
log_step "Step 5: Setting up Application Environment"
bash "${SCRIPTS_DIR}/05_setup_environment.sh"

# Step 6: Deploy backend
if [ "$SKIP_BACKEND" = false ]; then
    log_step "Step 6: Deploying Backend"
    bash "${SCRIPTS_DIR}/06_deploy_backend.sh"
else
    log_warn "Skipping backend deployment"
fi

# Step 7: Run migrations
if [ "$SKIP_DB" = false ]; then
    log_step "Step 7: Running Database Migrations"
    bash "${SCRIPTS_DIR}/08_run_migrations.sh"
else
    log_warn "Skipping database migrations"
fi

# Step 8: Deploy frontend
if [ "$SKIP_FRONTEND" = false ]; then
    log_step "Step 8: Deploying Frontend"
    bash "${SCRIPTS_DIR}/07_deploy_frontend.sh"
else
    log_warn "Skipping frontend deployment"
fi

# Restart Nginx if needed
if [ "$FRESH_INSTALL" = false ]; then
    log_step "Reloading Nginx"
    bash "${SCRIPTS_DIR}/04_setup_nginx.sh"
fi

# Calculate deployment time
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))
MINUTES=$((DURATION / 60))
SECONDS=$((DURATION % 60))

# Display success message
cat << EOF

╔═══════════════════════════════════════════════════════╗
║                                                       ║
║          Deployment Completed Successfully!          ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝

Deployment took: ${MINUTES}m ${SECONDS}s

Application URLs:
  Frontend: http://${DOMAIN_OR_IP}
  Backend API: http://${DOMAIN_OR_IP}/api/v1
  Health Check: http://${DOMAIN_OR_IP}/health

Useful Commands:
  Backend Service:
    - Status: sudo systemctl status daybook-backend
    - Logs: sudo journalctl -u daybook-backend -f
    - Restart: sudo systemctl restart daybook-backend

  Nginx:
    - Status: sudo systemctl status nginx
    - Logs: sudo tail -f /var/log/nginx/daybook-*.log
    - Reload: sudo systemctl reload nginx

  Database:
    - Connect: PGPASSWORD=${DB_PASSWORD} psql -h localhost -U ${DB_USER} -d ${DB_NAME}

  Log Files:
    - Backend: ${LOG_DIR}/backend.log
    - Nginx Access: /var/log/nginx/daybook-access.log
    - Nginx Error: /var/log/nginx/daybook-error.log

Next Steps:
  1. Test the application by visiting http://${DOMAIN_OR_IP}
  2. Configure SSL/TLS for HTTPS (recommended for production)
  3. Set up backups for database and uploads
  4. Configure monitoring and alerts

EOF

exit 0
