# EC2 Deployment Guide for Daybook

## System Requirements
- EC2 Instance: 10GB storage, 0.5GB RAM
- OS: Ubuntu 22.04 LTS or later
- Domains configured:
  - Frontend: `daybook.shafik.xyz` → EC2 IP
  - Backend: `api.daybook.shafik.xyz` → EC2 IP

## Architecture
- **Frontend**: Vue.js app served by Nginx on port 8081 (daybook.shafik.xyz)
- **Backend**: Go API served by Nginx on port 8082 (api.daybook.shafik.xyz)
- **Database**: PostgreSQL 15 (local)
- **Cache**: Redis 7 (local)

---

## Part 1: Local Build (On Your Development Machine)

⚠️ **Important**: Due to limited RAM (0.5GB), we'll build the backend binary and frontend assets locally to avoid OOM errors on EC2.

### 1.1 Use the Automated Build Script

The easiest way is to use the provided build script:

```bash
cd /path/to/daybook

# Run the build script
./scripts/local-build.sh
```

This will:
- Build the Go backend binary for Linux
- Build the Vue.js frontend production assets
- Package everything into the `build/` directory

### 1.2 Commit and Push Build Artifacts

After building, commit the artifacts to git:

```bash
# Add build artifacts
git add build/

# Commit
git commit -m "Add deployment build artifacts"

# Push to remote
git push origin master
```

### 1.3 Manual Build (Alternative)

If you prefer to build manually:

```bash
# Build backend
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../build/daybook-backend main.go

# Build frontend
cd ../frontend
npm install
cat > .env.production << EOF
VITE_API_URL=https://api.daybook.shafik.xyz
VITE_APP_NAME=Daybook
VITE_APP_VERSION=1.0.0
EOF
npm run build
cd dist
tar -czf ../../build/frontend-dist.tar.gz .

# Commit and push
cd ../..
git add build/
git commit -m "Add deployment build artifacts"
git push origin master
```

---

## Part 2: EC2 Initial Setup

### 2.1 Connect to EC2

```bash
ssh -i your-key.pem ubuntu@your-ec2-ip
```

### 2.2 Update System

```bash
sudo apt update
sudo apt upgrade -y
```

### 2.3 Install Essential Tools

```bash
sudo apt install -y \
  curl \
  wget \
  git \
  nano \
  ufw \
  software-properties-common \
  ca-certificates \
  gnupg \
  lsb-release
```

### 2.4 Configure Firewall

```bash
# Allow SSH, HTTP, HTTPS
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8081/tcp
sudo ufw allow 8082/tcp

# Enable firewall
sudo ufw --force enable
sudo ufw status
```

---

## Part 3: Install PostgreSQL

### 3.1 Install PostgreSQL 15

```bash
# Add PostgreSQL repository
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

# Install PostgreSQL
sudo apt update
sudo apt install -y postgresql-15 postgresql-contrib-15

# Start and enable PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql
sudo systemctl status postgresql
```

### 3.2 Configure PostgreSQL for Low Memory

```bash
# Optimize for limited RAM
sudo nano /etc/postgresql/15/main/postgresql.conf
```

Add/modify these settings:

```conf
# Memory Settings (optimized for 0.5GB RAM)
shared_buffers = 64MB
effective_cache_size = 128MB
maintenance_work_mem = 16MB
work_mem = 2MB
max_connections = 20

# Performance
random_page_cost = 1.1
effective_io_concurrency = 200
```

Save and restart:

```bash
sudo systemctl restart postgresql
```

### 3.3 Create Database and User

```bash
# Switch to postgres user
sudo -u postgres psql

# In PostgreSQL shell, run:
CREATE DATABASE daybook;
CREATE USER daybook_user WITH ENCRYPTED PASSWORD 'your_strong_password_here';
GRANT ALL PRIVILEGES ON DATABASE daybook TO daybook_user;
ALTER DATABASE daybook OWNER TO daybook_user;

# Grant schema privileges
\c daybook
GRANT ALL ON SCHEMA public TO daybook_user;

# Exit
\q
```

### 3.4 Test Connection

```bash
psql -h localhost -U daybook_user -d daybook -W
# Enter password when prompted, then \q to exit
```

---

## Part 4: Install Redis

### 4.1 Install Redis 7

```bash
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list

sudo apt update
sudo apt install -y redis-server
```

### 4.2 Configure Redis for Low Memory

```bash
sudo nano /etc/redis/redis.conf
```

Modify these settings:

```conf
# Memory limit (50MB max)
maxmemory 50mb
maxmemory-policy allkeys-lru

# Save to disk less frequently to save memory
save 900 1
save 300 10

# Bind to localhost only
bind 127.0.0.1

# Disable protected mode if needed
protected-mode yes

# Set password
requirepass your_redis_password_here
```

Save and restart:

```bash
sudo systemctl restart redis-server
sudo systemctl enable redis-server
sudo systemctl status redis-server
```

### 4.3 Test Redis

```bash
redis-cli
# In Redis CLI:
AUTH your_redis_password_here
PING
# Should return PONG
exit
```

---

## Part 5: Install Nginx

```bash
sudo apt install -y nginx

# Start and enable Nginx
sudo systemctl start nginx
sudo systemctl enable nginx
sudo systemctl status nginx
```

---

## Part 6: Deploy Application

### 6.1 Clone Repository

```bash
cd /home/ubuntu
git clone https://github.com/shafikshaon/daybook.git
cd daybook
```

**Note**: If you have SSH keys configured on EC2 for GitHub, you can use SSH:

```bash
git clone git@github.com:shafikshaon/daybook.git
```

### 6.2 Deploy Using Automated Script (Recommended)

The easiest way to deploy is using the provided setup script:

```bash
cd ~/daybook
./scripts/ec2-setup.sh
```

This script will:
- Create directory structure
- Deploy backend binary from `build/daybook-backend`
- Deploy frontend from `build/frontend-dist.tar.gz`
- Configure backend environment (.env)
- Create systemd service
- Configure Nginx

The script will prompt you for configuration details (database password, Redis password, etc.).

**After the script completes, skip to Part 7 to start the services.**

### 6.3 Manual Deployment (Alternative)

If you prefer manual deployment:

#### 6.3.1 Setup Application Directory Structure

```bash
sudo mkdir -p /var/www/daybook
sudo mkdir -p /var/www/daybook/frontend
sudo mkdir -p /var/www/daybook/backend
sudo mkdir -p /var/www/daybook/backend/uploads
sudo chown -R ubuntu:ubuntu /var/www/daybook
```

#### 6.3.2 Deploy Backend

```bash
# Copy binary from build directory
cp ~/daybook/build/daybook-backend /var/www/daybook/backend/daybook-backend
chmod +x /var/www/daybook/backend/daybook-backend
```

#### Create backend .env file:

```bash
nano /var/www/daybook/backend/.env
```

Add the following configuration:

```env
# Server Configuration
SERVER_PORT=8082
SERVER_MODE=release

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=daybook_user
DB_PASSWORD=your_strong_password_here
DB_NAME=daybook
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Dhaka

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password_here
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_min_32_characters_long_123456
JWT_EXPIRATION=24h

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://daybook.shafik.xyz
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_MAX_AGE=12h
CORS_ALLOW_CREDENTIALS=true

# Upload Configuration
UPLOAD_PATH=/var/www/daybook/backend/uploads
MAX_UPLOAD_SIZE=10485760

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m
```

Save the file (Ctrl+X, Y, Enter).

#### 6.3.3 Deploy Frontend

```bash
# Extract frontend build
cd /var/www/daybook/frontend
tar -xzf ~/daybook/build/frontend-dist.tar.gz
```

---

## Part 7: Create Systemd Services

### 7.1 Backend Service

```bash
sudo nano /etc/systemd/system/daybook-backend.service
```

Add:

```ini
[Unit]
Description=Daybook Backend API
After=network.target postgresql.service redis.service
Requires=postgresql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/var/www/daybook/backend
Environment="GIN_MODE=release"
ExecStart=/var/www/daybook/backend/daybook-backend
Restart=on-failure
RestartSec=5s
StandardOutput=append:/var/log/daybook-backend.log
StandardError=append:/var/log/daybook-backend-error.log

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### 7.2 Create Log Files

```bash
sudo touch /var/log/daybook-backend.log
sudo touch /var/log/daybook-backend-error.log
sudo chown ubuntu:ubuntu /var/log/daybook-backend*.log
```

### 7.3 Enable and Start Backend Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable daybook-backend

# Start the service
sudo systemctl start daybook-backend

# Check status
sudo systemctl status daybook-backend

# View logs
sudo journalctl -u daybook-backend -f
```

---

## Part 8: Configure Nginx

### 8.1 Create Nginx Configuration for Frontend

```bash
sudo nano /etc/nginx/sites-available/daybook-frontend
```

Add:

```nginx
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
```

### 8.2 Create Nginx Configuration for Backend

```bash
sudo nano /etc/nginx/sites-available/daybook-backend
```

Add:

```nginx
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
```

### 8.3 Enable Sites

```bash
# Create symbolic links
sudo ln -s /etc/nginx/sites-available/daybook-frontend /etc/nginx/sites-enabled/
sudo ln -s /etc/nginx/sites-available/daybook-backend /etc/nginx/sites-enabled/

# Test Nginx configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

---

## Part 9: SSL/TLS Setup with Let's Encrypt

### 9.1 Install Certbot

```bash
sudo apt install -y certbot python3-certbot-nginx
```

### 9.2 Obtain SSL Certificates

```bash
# For frontend
sudo certbot --nginx -d daybook.shafik.xyz

# For backend
sudo certbot --nginx -d api.daybook.shafik.xyz
```

Follow the prompts and select option 2 (redirect HTTP to HTTPS).

### 9.3 Auto-renewal Test

```bash
sudo certbot renew --dry-run
```

The certificates will auto-renew via cron job.

---

## Part 10: DNS Configuration

Before accessing the application, configure your DNS:

1. Go to your domain registrar (e.g., Namecheap, GoDaddy)
2. Add A records:
   ```
   daybook.shafik.xyz     → Your EC2 Public IP
   api.daybook.shafik.xyz → Your EC2 Public IP
   ```
3. Wait for DNS propagation (5-30 minutes)

---

## Part 11: Verify Deployment

### 11.1 Check Services

```bash
# Check backend
sudo systemctl status daybook-backend
curl http://localhost:8082/health  # If you have a health endpoint

# Check PostgreSQL
sudo systemctl status postgresql

# Check Redis
sudo systemctl status redis-server

# Check Nginx
sudo systemctl status nginx
```

### 11.2 Check Logs

```bash
# Backend logs
sudo tail -f /var/log/daybook-backend.log
sudo journalctl -u daybook-backend -f

# Nginx logs
sudo tail -f /var/log/nginx/daybook-frontend-access.log
sudo tail -f /var/log/nginx/daybook-backend-access.log
```

### 11.3 Test Endpoints

```bash
# Test backend (after SSL)
curl https://api.daybook.shafik.xyz/api/health

# Test frontend
curl https://daybook.shafik.xyz
```

### 11.4 Access Application

Open browser and visit:
- Frontend: https://daybook.shafik.xyz
- Backend API: https://api.daybook.shafik.xyz

---

## Part 12: Maintenance & Updates

### 12.1 Update Application

To deploy new changes:

**On local machine:**

```bash
# Rebuild the application
cd /path/to/daybook
./scripts/local-build.sh

# Commit and push changes
git add build/
git commit -m "Update deployment build"
git push origin master
```

**On EC2:**

```bash
# Pull latest changes
cd ~/daybook
git pull origin master

# Stop backend service
sudo systemctl stop daybook-backend

# Update backend
cp ~/daybook/build/daybook-backend /var/www/daybook/backend/daybook-backend
chmod +x /var/www/daybook/backend/daybook-backend

# Update frontend
cd /var/www/daybook/frontend
rm -rf *
tar -xzf ~/daybook/build/frontend-dist.tar.gz

# Restart services
sudo systemctl start daybook-backend
sudo systemctl reload nginx

# Verify
sudo systemctl status daybook-backend
```

### 12.2 Backup Database

Create automated backup script:

```bash
nano /home/ubuntu/backup-db.sh
```

Add:

```bash
#!/bin/bash
BACKUP_DIR="/home/ubuntu/backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

pg_dump -h localhost -U daybook_user -d daybook > $BACKUP_DIR/daybook_$DATE.sql

# Keep only last 7 days
find $BACKUP_DIR -name "daybook_*.sql" -mtime +7 -delete
```

Make executable and add to cron:

```bash
chmod +x /home/ubuntu/backup-db.sh

# Add to crontab (daily at 2 AM)
crontab -e
# Add line:
0 2 * * * /home/ubuntu/backup-db.sh
```

### 12.3 Monitor Resources

```bash
# Check memory usage
free -h

# Check disk usage
df -h

# Check running processes
htop  # Install with: sudo apt install htop

# Check backend process
ps aux | grep daybook-backend
```

### 12.4 View Application Logs

```bash
# Backend application logs
sudo tail -100 /var/log/daybook-backend.log
sudo tail -100 /var/log/daybook-backend-error.log

# Systemd logs
sudo journalctl -u daybook-backend --since "1 hour ago"

# Nginx access logs
sudo tail -100 /var/log/nginx/daybook-frontend-access.log
sudo tail -100 /var/log/nginx/daybook-backend-access.log

# Nginx error logs
sudo tail -100 /var/log/nginx/daybook-frontend-error.log
sudo tail -100 /var/log/nginx/daybook-backend-error.log
```

---

## Part 13: Troubleshooting

### Backend Won't Start

```bash
# Check service status
sudo systemctl status daybook-backend

# Check logs for errors
sudo journalctl -u daybook-backend -n 50

# Verify binary permissions
ls -l /var/www/daybook/backend/daybook-backend

# Test manual start
cd /var/www/daybook/backend
./daybook-backend
```

### Database Connection Issues

```bash
# Verify PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -h localhost -U daybook_user -d daybook

# Check PostgreSQL logs
sudo tail -100 /var/log/postgresql/postgresql-15-main.log
```

### Nginx Issues

```bash
# Test configuration
sudo nginx -t

# Check for port conflicts
sudo netstat -tulpn | grep :8081
sudo netstat -tulpn | grep :8082

# Restart Nginx
sudo systemctl restart nginx
```

### Out of Memory

```bash
# Check memory usage
free -h

# Check swap
swapon --show

# Add swap if needed (1GB)
sudo fallocate -l 1G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

### SSL Certificate Issues

```bash
# Check certificate status
sudo certbot certificates

# Renew manually
sudo certbot renew

# Check Nginx SSL configuration
sudo nginx -t
```

---

## Part 14: Security Checklist

- [ ] Changed default PostgreSQL password
- [ ] Set strong Redis password
- [ ] Generated unique JWT secret (32+ characters)
- [ ] Configured UFW firewall
- [ ] Enabled HTTPS with Let's Encrypt
- [ ] Set proper file permissions (backend binary, .env files)
- [ ] Disabled PostgreSQL remote access (only localhost)
- [ ] Configured Nginx security headers
- [ ] Set up automated database backups
- [ ] Reviewed and secured .env files (not world-readable)

```bash
# Secure .env file
chmod 600 /var/www/daybook/backend/.env
```

---

## Quick Reference Commands

### Service Management
```bash
# Backend
sudo systemctl start daybook-backend
sudo systemctl stop daybook-backend
sudo systemctl restart daybook-backend
sudo systemctl status daybook-backend

# Nginx
sudo systemctl reload nginx
sudo systemctl restart nginx

# PostgreSQL
sudo systemctl restart postgresql

# Redis
sudo systemctl restart redis-server
```

### Log Viewing
```bash
# Real-time backend logs
sudo journalctl -u daybook-backend -f

# Recent errors
sudo journalctl -u daybook-backend -p err -n 50
```

### Database Operations
```bash
# Connect to database
psql -h localhost -U daybook_user -d daybook

# Backup database
pg_dump -h localhost -U daybook_user -d daybook > backup.sql

# Restore database
psql -h localhost -U daybook_user -d daybook < backup.sql
```

---

## Support

For issues or questions:
1. Check application logs first
2. Review this deployment guide
3. Consult the main project README
4. Check backend API documentation

---

**Deployment Complete!** Your Daybook application should now be running at:
- Frontend: https://daybook.shafik.xyz
- Backend API: https://api.daybook.shafik.xyz
