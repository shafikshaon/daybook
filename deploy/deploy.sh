#!/bin/bash
set -e

# ============================================================================
# Daybook Simple Deployment Script
# ============================================================================
# Just run: ./deploy.sh
# That's it!
# ============================================================================

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_step() { echo -e "\n${BLUE}==== $1 ====${NC}\n"; }

# ============================================================================
# CONFIGURATION - Edit these values
# ============================================================================

DOMAIN="daybook.shafik.xyz"              # Your domain or server IP
DB_PASSWORD="VtPuqR1er16VFI5F5CPrYtlX4bmsBavZ"  # Auto-generated secure password
APP_PORT=8080                            # Backend port

# ============================================================================
# DO NOT EDIT BELOW THIS LINE
# ============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    log_error "Do not run as root. Script will use sudo when needed."
    exit 1
fi

log_step "Starting Daybook Deployment"

# ============================================================================
# Step 1: Install System Dependencies
# ============================================================================

log_step "Installing System Dependencies"

# Update system
log_info "Updating system packages..."
sudo apt-get update

# Install basic tools
log_info "Installing basic tools (this may take 5-10 minutes on first run)..."
sudo apt-get install -y curl wget git build-essential openssl postgresql redis-server nginx

# Install Go
log_info "Installing Go 1.25..."
GO_VERSION="1.25.3"
if ! command -v go &> /dev/null || ! go version | grep -q "go1.25"; then
    log_info "Downloading Go ${GO_VERSION}..."
    wget --progress=bar:force https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
    log_info "Installing Go..."
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    rm go${GO_VERSION}.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh > /dev/null
fi
log_info "Go version: $(go version)"

# Install Node.js
log_info "Installing Node.js 22..."
if ! command -v node &> /dev/null || ! node --version | grep -q "v22"; then
    log_info "Adding Node.js repository..."
    curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
    log_info "Installing Node.js..."
    sudo apt-get install -y nodejs
fi
log_info "Node version: $(node --version)"
log_info "NPM version: $(npm --version)"

# ============================================================================
# Step 2: Setup PostgreSQL Database
# ============================================================================

log_step "Setting up PostgreSQL"

sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database and user
log_info "Creating database and user..."
sudo -u postgres psql -c "CREATE DATABASE daybook_prod;" 2>/dev/null || log_warn "Database already exists"
sudo -u postgres psql -c "CREATE USER daybook_user WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null || log_warn "User already exists"
sudo -u postgres psql -c "ALTER DATABASE daybook_prod OWNER TO daybook_user;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE daybook_prod TO daybook_user;"

log_info "Database setup complete"

# ============================================================================
# Step 3: Setup Redis
# ============================================================================

log_step "Setting up Redis"

sudo systemctl start redis-server
sudo systemctl enable redis-server
log_info "Redis is running"

# ============================================================================
# Step 4: Build Backend
# ============================================================================

log_step "Building Backend"

cd "$PROJECT_ROOT/backend"

# Create .env file
log_info "Creating backend .env file..."
cat > .env <<EOF
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=$APP_PORT
ENVIRONMENT=production

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=daybook_prod
DB_USER=daybook_user
DB_PASSWORD=$DB_PASSWORD
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=$(openssl rand -base64 64 | tr -d '\n')
JWT_EXPIRATION=24h

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://$DOMAIN,http://localhost:3000

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m

# File Upload
MAX_UPLOAD_SIZE=10485760
UPLOAD_PATH=./uploads
EOF

chmod 600 .env

# Build backend (skip if pre-built binary exists)
if [ -f "daybook-backend" ]; then
    log_info "Using pre-built backend binary (skipping compilation)"
    chmod +x daybook-backend
else
    log_info "No pre-built binary found, building from source..."
    log_info "Downloading Go dependencies (this may take a few minutes)..."
    go mod download -x

    log_info "Compiling backend application (this may take 1-2 minutes)..."
    CGO_ENABLED=0 go build -v -ldflags="-s -w" -o daybook-backend main.go

    if [ ! -f "daybook-backend" ]; then
        log_error "Backend build failed!"
        exit 1
    fi

    chmod +x daybook-backend
    log_info "Backend built successfully"
fi

mkdir -p uploads

# ============================================================================
# Step 5: Create Backend Service
# ============================================================================

log_step "Setting up Backend Service"

sudo tee /etc/systemd/system/daybook-backend.service > /dev/null <<EOF
[Unit]
Description=Daybook Backend Service
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=$USER
WorkingDirectory=$PROJECT_ROOT/backend
ExecStart=$PROJECT_ROOT/backend/daybook-backend
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable daybook-backend

log_info "Starting backend service..."
sudo systemctl restart daybook-backend

log_info "Waiting for backend to start..."
sleep 3

if ! systemctl is-active --quiet daybook-backend; then
    log_error "Backend service failed to start!"
    log_error "Check logs: sudo journalctl -u daybook-backend -n 50"
    exit 1
fi

log_info "Backend service is running"

# ============================================================================
# Step 6: Build Frontend
# ============================================================================

log_step "Building Frontend"

cd "$PROJECT_ROOT/frontend"

# Create .env file
log_info "Creating frontend .env file..."
cat > .env <<EOF
VITE_API_URL=http://$DOMAIN/api/v1
EOF

# Install dependencies and build
log_info "Installing npm dependencies (this may take 3-5 minutes)..."
npm install --progress=true

log_info "Building frontend application (this may take 2-3 minutes)..."
npm run build

if [ ! -d "dist" ]; then
    log_error "Frontend build failed!"
    exit 1
fi

log_info "Frontend built successfully"

# ============================================================================
# Step 7: Setup Nginx
# ============================================================================

log_step "Setting up Nginx"

sudo tee /etc/nginx/sites-available/daybook > /dev/null <<EOF
upstream backend {
    server 127.0.0.1:$APP_PORT;
    keepalive 32;
}

server {
    listen 80;
    server_name $DOMAIN;

    client_max_body_size 10M;

    # Frontend
    location / {
        root $PROJECT_ROOT/frontend/dist;
        try_files \$uri \$uri/ /index.html;
    }

    # Backend API
    location /api/ {
        proxy_pass http://backend;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header Connection "";
    }

    # Health check
    location /health {
        proxy_pass http://backend;
        access_log off;
    }

    # Uploads
    location /uploads/ {
        alias $PROJECT_ROOT/backend/uploads/;
    }
}
EOF

# Enable site
sudo ln -sf /etc/nginx/sites-available/daybook /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default

# Test and reload nginx
log_info "Testing Nginx configuration..."
if sudo nginx -t 2>&1 | grep -q "successful"; then
    log_info "Starting Nginx..."
    sudo systemctl restart nginx
    sudo systemctl enable nginx
    log_info "Nginx configured and running"
else
    log_error "Nginx configuration failed!"
    exit 1
fi

# ============================================================================
# DONE!
# ============================================================================

log_step "Deployment Complete!"

echo ""
echo "=========================================="
echo "  ✅ Daybook is now deployed!"
echo "=========================================="
echo ""
echo "🌐 Frontend: http://$DOMAIN"
echo "🔌 Backend API: http://$DOMAIN/api/v1"
echo "❤️  Health: http://$DOMAIN/health"
echo ""
echo "📊 Services:"
echo "   Backend: sudo systemctl status daybook-backend"
echo "   Logs: sudo journalctl -u daybook-backend -f"
echo ""
echo "🔒 Database:"
echo "   Name: daybook_prod"
echo "   User: daybook_user"
echo "   Password: $DB_PASSWORD"
echo "   (Password saved in: $PROJECT_ROOT/backend/.env)"
echo ""
echo "🔄 To update:"
echo "   cd $PROJECT_ROOT"
echo "   git pull"
echo "   cd deploy"
echo "   ./deploy.sh"
echo ""
echo "=========================================="

exit 0
