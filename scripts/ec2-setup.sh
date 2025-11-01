#!/bin/bash
#
# EC2 Setup Script
# Run this script on your EC2 instance after transferring build artifacts
#

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Daybook EC2 Setup Script"
echo "=========================================="
echo ""

# Check if running as ubuntu user
if [ "$USER" != "ubuntu" ]; then
    echo -e "${RED}Please run this script as ubuntu user${NC}"
    exit 1
fi

# Variables
DEPLOY_DIR="/var/www/daybook"
BACKEND_DIR="$DEPLOY_DIR/backend"
FRONTEND_DIR="$DEPLOY_DIR/frontend"
HOME_DIR="/home/ubuntu"
REPO_DIR="$HOME_DIR/daybook"
BUILD_DIR="$REPO_DIR/build"

# Function to prompt for input
prompt_input() {
    local prompt_text=$1
    local var_name=$2
    local default_value=$3

    if [ -n "$default_value" ]; then
        read -p "$prompt_text [$default_value]: " input
        input=${input:-$default_value}
    else
        read -p "$prompt_text: " input
    fi

    eval "$var_name='$input'"
}

# Function to prompt for password
prompt_password() {
    local prompt_text=$1
    local var_name=$2

    read -s -p "$prompt_text: " password
    echo ""
    eval "$var_name='$password'"
}

echo -e "${GREEN}[1/6] Creating Directory Structure...${NC}"
sudo mkdir -p $BACKEND_DIR
sudo mkdir -p $FRONTEND_DIR
sudo mkdir -p $BACKEND_DIR/uploads
sudo chown -R ubuntu:ubuntu $DEPLOY_DIR
echo -e "${GREEN}✓ Directories created${NC}"
echo ""

echo -e "${GREEN}[2/6] Deploying Backend...${NC}"
if [ ! -f "$BUILD_DIR/daybook-backend" ]; then
    echo -e "${RED}Error: daybook-backend not found in $BUILD_DIR${NC}"
    echo "Please ensure you have:"
    echo "  1. Built the application locally (./scripts/local-build.sh)"
    echo "  2. Committed and pushed build artifacts to git"
    echo "  3. Pulled the latest changes (git pull origin master)"
    exit 1
fi

cp "$BUILD_DIR/daybook-backend" "$BACKEND_DIR/daybook-backend"
chmod +x "$BACKEND_DIR/daybook-backend"
echo -e "${GREEN}✓ Backend deployed${NC}"
echo ""

echo -e "${GREEN}[3/6] Deploying Frontend...${NC}"
if [ ! -f "$BUILD_DIR/frontend-dist.tar.gz" ]; then
    echo -e "${RED}Error: frontend-dist.tar.gz not found in $BUILD_DIR${NC}"
    echo "Please ensure you have:"
    echo "  1. Built the application locally (./scripts/local-build.sh)"
    echo "  2. Committed and pushed build artifacts to git"
    echo "  3. Pulled the latest changes (git pull origin master)"
    exit 1
fi

cd $FRONTEND_DIR
tar -xzf "$BUILD_DIR/frontend-dist.tar.gz"
echo -e "${GREEN}✓ Frontend deployed${NC}"
echo ""

echo -e "${GREEN}[4/6] Configuring Backend Environment...${NC}"
echo ""
echo "Please provide the following information for backend configuration:"
echo ""

# Database configuration
prompt_input "Database User" DB_USER "daybook_user"
prompt_password "Database Password" DB_PASSWORD
prompt_input "Database Name" DB_NAME "daybook"

# Redis configuration
prompt_password "Redis Password" REDIS_PASSWORD

# JWT configuration
echo ""
echo "Generating JWT Secret..."
JWT_SECRET=$(openssl rand -base64 48)

# Create .env file
cat > "$BACKEND_DIR/.env" << EOF
# Server Configuration
SERVER_PORT=8082
SERVER_MODE=release

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Dhaka

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=$REDIS_PASSWORD
REDIS_DB=0

# JWT Configuration
JWT_SECRET=$JWT_SECRET
JWT_EXPIRATION=24h

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://daybook.shafik.xyz
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_MAX_AGE=12h
CORS_ALLOW_CREDENTIALS=true

# Upload Configuration
UPLOAD_PATH=$BACKEND_DIR/uploads
MAX_UPLOAD_SIZE=10485760

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m
EOF

chmod 600 "$BACKEND_DIR/.env"
echo -e "${GREEN}✓ Backend .env configured${NC}"
echo ""

echo -e "${GREEN}[5/6] Creating Systemd Service...${NC}"

sudo tee /etc/systemd/system/daybook-backend.service > /dev/null << EOF
[Unit]
Description=Daybook Backend API
After=network.target postgresql.service redis.service
Requires=postgresql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=$BACKEND_DIR
Environment="GIN_MODE=release"
ExecStart=$BACKEND_DIR/daybook-backend
Restart=on-failure
RestartSec=5s
StandardOutput=append:/var/log/daybook-backend.log
StandardError=append:/var/log/daybook-backend-error.log

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

# Create log files
sudo touch /var/log/daybook-backend.log
sudo touch /var/log/daybook-backend-error.log
sudo chown ubuntu:ubuntu /var/log/daybook-backend*.log

echo -e "${GREEN}✓ Systemd service created${NC}"
echo ""

echo -e "${GREEN}[6/6] Creating Nginx Configurations...${NC}"

# Frontend Nginx config
sudo tee /etc/nginx/sites-available/daybook-frontend > /dev/null << 'EOF'
server {
    listen 8081;
    server_name daybook.shafik.xyz;

    root /var/www/daybook/frontend;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Disable access to hidden files
    location ~ /\. {
        deny all;
    }

    access_log /var/log/nginx/daybook-frontend-access.log;
    error_log /var/log/nginx/daybook-frontend-error.log;
}
EOF

# Backend Nginx config
sudo tee /etc/nginx/sites-available/daybook-backend > /dev/null << 'EOF'
upstream backend_app {
    server 127.0.0.1:8082;
    keepalive 32;
}

server {
    listen 8081;
    server_name api.daybook.shafik.xyz;

    client_max_body_size 10M;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    location / {
        proxy_pass http://backend_app;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Connection "";

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    access_log /var/log/nginx/daybook-backend-access.log;
    error_log /var/log/nginx/daybook-backend-error.log;
}
EOF

# Enable sites
sudo ln -sf /etc/nginx/sites-available/daybook-frontend /etc/nginx/sites-enabled/
sudo ln -sf /etc/nginx/sites-available/daybook-backend /etc/nginx/sites-enabled/

# Test Nginx config
sudo nginx -t

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Nginx configured${NC}"
else
    echo -e "${RED}✗ Nginx configuration error${NC}"
    exit 1
fi
echo ""

echo -e "${GREEN}=========================================="
echo "Setup Complete!"
echo "==========================================${NC}"
echo ""
echo "Next steps:"
echo "1. Start the backend service:"
echo "   sudo systemctl daemon-reload"
echo "   sudo systemctl enable daybook-backend"
echo "   sudo systemctl start daybook-backend"
echo ""
echo "2. Reload Nginx:"
echo "   sudo systemctl reload nginx"
echo ""
echo "3. Check service status:"
echo "   sudo systemctl status daybook-backend"
echo "   sudo journalctl -u daybook-backend -f"
echo ""
echo "4. Configure SSL with Let's Encrypt:"
echo "   sudo certbot --nginx -d daybook.shafik.xyz"
echo "   sudo certbot --nginx -d api.daybook.shafik.xyz"
echo ""
echo -e "${YELLOW}Important: Make sure DNS records are configured before setting up SSL${NC}"
