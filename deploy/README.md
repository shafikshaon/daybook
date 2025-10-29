# Daybook Deployment Guide

This directory contains automated deployment scripts for deploying the Daybook application (frontend and backend) on an Ubuntu server running on AWS with PostgreSQL and Redis.

## Overview

The deployment system is modular and automated, requiring minimal manual intervention. Each script performs a specific task and can be run independently or as part of the main deployment process.

## Prerequisites

### Local Machine
- SSH access to your Ubuntu server
- Git repository cloned

### Ubuntu Server (AWS EC2)
- Ubuntu 20.04 LTS or later
- Minimum 2GB RAM recommended
- At least 10GB free disk space
- SSH access with sudo privileges
- Security group allowing:
  - Port 22 (SSH)
  - Port 80 (HTTP)
  - Port 443 (HTTPS) - optional, for SSL

## Quick Start

### 1. Configure Deployment

Edit the configuration file with your server details:

```bash
cd deploy
cp config/deploy.conf config/deploy.conf.local
nano config/deploy.conf
```

**Important:** Update these values in `deploy.conf`:
- `DOMAIN_OR_IP`: Your server's IP address or domain name
- `DB_PASSWORD`: A strong database password (leave empty to auto-generate)
- `GO_VERSION`: Go version to install (default: 1.21.5)
- `NODE_VERSION`: Node.js major version (default: 20)

### 2. Copy Files to Server

Copy the entire project to your server:

```bash
# From your local machine
rsync -avz --exclude 'node_modules' --exclude '.git' \
  /path/to/daybook ubuntu@YOUR_SERVER_IP:/tmp/daybook
```

### 3. Run Deployment

SSH into your server and run the deployment:

```bash
ssh ubuntu@YOUR_SERVER_IP
cd /tmp/daybook/deploy

# For fresh installation (first time)
./deploy.sh --fresh

# For updates (when infrastructure is already set up)
./deploy.sh --skip-deps --skip-db
```

### 4. Access Your Application

Once deployment completes:
- Frontend: `http://YOUR_SERVER_IP`
- Backend API: `http://YOUR_SERVER_IP/api/v1`
- Health Check: `http://YOUR_SERVER_IP/health`

## Directory Structure

```
deploy/
├── deploy.sh                    # Main deployment orchestration script
├── config/
│   ├── deploy.conf              # Main deployment configuration
│   ├── production.env.template  # Backend environment template
│   └── frontend.env.template    # Frontend environment template
├── scripts/
│   ├── 01_install_dependencies.sh   # Install Go, Node.js, system packages
│   ├── 02_setup_postgres.sh         # Install and configure PostgreSQL
│   ├── 03_setup_redis.sh            # Install and configure Redis
│   ├── 04_setup_nginx.sh            # Install and configure Nginx
│   ├── 05_setup_environment.sh      # Create app user, directories, .env files
│   ├── 06_deploy_backend.sh         # Build and deploy Go backend
│   ├── 07_deploy_frontend.sh        # Build and deploy Vue.js frontend
│   ├── 08_run_migrations.sh         # Run database migrations
│   ├── backup_database.sh           # Create database backup
│   ├── restore_database.sh          # Restore database from backup
│   ├── update_backend.sh            # Quick backend update
│   └── update_frontend.sh           # Quick frontend update
├── nginx/
│   └── daybook.conf.template    # Nginx configuration template
├── systemd/
│   └── daybook-backend.service.template  # Backend systemd service
└── README.md                    # This file
```

## Deployment Scripts

### Main Deployment Script

**`deploy.sh`** - Orchestrates the entire deployment process

```bash
# Full fresh installation
./deploy.sh --fresh

# Update deployment (skip infrastructure setup)
./deploy.sh --skip-deps --skip-db

# Deploy only backend
./deploy.sh --skip-deps --skip-db --skip-frontend

# Deploy only frontend
./deploy.sh --skip-deps --skip-db --skip-backend

# Show help
./deploy.sh --help
```

**Options:**
- `--fresh`: Fresh installation (install all dependencies and services)
- `--skip-deps`: Skip system dependencies installation
- `--skip-db`: Skip PostgreSQL and Redis setup
- `--skip-backend`: Skip backend deployment
- `--skip-frontend`: Skip frontend deployment

### Individual Scripts

All scripts can be run independently:

#### 1. Install Dependencies
```bash
./scripts/01_install_dependencies.sh
```
Installs:
- Go (version specified in config)
- Node.js and npm
- Build tools and utilities

#### 2. Setup PostgreSQL
```bash
./scripts/02_setup_postgres.sh
```
- Installs PostgreSQL
- Creates database and user
- Configures authentication
- Tests connection

#### 3. Setup Redis
```bash
./scripts/03_setup_redis.sh
```
- Installs Redis
- Configures for systemd
- Sets memory limits
- Enables persistence

#### 4. Setup Nginx
```bash
./scripts/04_setup_nginx.sh
```
- Installs Nginx
- Configures reverse proxy
- Sets up static file serving
- Configures firewall

#### 5. Setup Environment
```bash
./scripts/05_setup_environment.sh
```
- Creates application user
- Creates directory structure
- Generates secrets (JWT, DB password)
- Creates .env files

#### 6. Deploy Backend
```bash
./scripts/06_deploy_backend.sh
```
- Copies backend source
- Builds Go application
- Creates systemd service
- Starts backend service

#### 7. Deploy Frontend
```bash
./scripts/07_deploy_frontend.sh
```
- Copies frontend source
- Installs npm dependencies
- Builds production bundle
- Deploys to Nginx

#### 8. Run Migrations
```bash
./scripts/08_run_migrations.sh
```
- Tracks applied migrations
- Runs pending migrations
- Updates schema version

### Utility Scripts

#### Quick Update Scripts

Update backend only:
```bash
./scripts/update_backend.sh
```

Update frontend only:
```bash
./scripts/update_frontend.sh
```

#### Database Management

Create backup:
```bash
./scripts/backup_database.sh
```
- Creates timestamped backup
- Compresses with gzip
- Keeps last 7 backups

Restore from backup:
```bash
./scripts/restore_database.sh /path/to/backup.sql.gz
```
- Creates safety backup
- Restores from specified backup
- Restarts backend service

## Configuration Details

### Deploy Configuration (`config/deploy.conf`)

Key configuration options:

```bash
# Server Details
SERVER_IP="YOUR_SERVER_IP"
DOMAIN_OR_IP="YOUR_DOMAIN_OR_IP"

# Application
APP_NAME="daybook"
APP_USER="daybook"
APP_DIR="/opt/daybook"

# Database
DB_NAME="daybook_prod"
DB_USER="daybook_user"
DB_PASSWORD=""  # Auto-generated if empty

# Ports
BACKEND_PORT=8080

# Versions
GO_VERSION="1.21.5"
NODE_VERSION="20"
```

### Environment Files

Backend environment (`.env`) includes:
- Server configuration
- Database credentials
- Redis configuration
- JWT secret

Frontend environment (`.env`) includes:
- API URL
- App metadata

## Post-Deployment

### Service Management

#### Backend Service
```bash
# Status
sudo systemctl status daybook-backend

# Restart
sudo systemctl restart daybook-backend

# Stop
sudo systemctl stop daybook-backend

# Start
sudo systemctl start daybook-backend

# View logs
sudo journalctl -u daybook-backend -f
```

#### Nginx
```bash
# Status
sudo systemctl status nginx

# Reload configuration
sudo systemctl reload nginx

# Restart
sudo systemctl restart nginx

# Test configuration
sudo nginx -t

# View logs
sudo tail -f /var/log/nginx/daybook-access.log
sudo tail -f /var/log/nginx/daybook-error.log
```

#### PostgreSQL
```bash
# Status
sudo systemctl status postgresql

# Connect to database
PGPASSWORD=your_password psql -h localhost -U daybook_user -d daybook_prod

# View connections
sudo -u postgres psql -c "SELECT * FROM pg_stat_activity WHERE datname='daybook_prod';"
```

#### Redis
```bash
# Status
sudo systemctl status redis-server

# Connect
redis-cli

# Monitor commands
redis-cli monitor

# Get info
redis-cli INFO
```

### Log Files

Application logs:
- Backend: `/var/log/daybook/backend.log`
- Backend errors: `/var/log/daybook/backend-error.log`
- Nginx access: `/var/log/nginx/daybook-access.log`
- Nginx error: `/var/log/nginx/daybook-error.log`

### File Locations

```
/opt/daybook/
├── backend/
│   ├── daybook-backend     # Binary
│   ├── .env                # Backend environment
│   ├── migrations/         # Database migrations
│   └── uploads/            # Uploaded files
├── frontend/
│   ├── dist/               # Built frontend files
│   └── .env                # Frontend environment
└── backups/                # Database backups
```

## SSL/TLS Setup (Optional but Recommended)

### Using Let's Encrypt (Certbot)

1. Install Certbot:
```bash
sudo apt-get install -y certbot python3-certbot-nginx
```

2. Obtain certificate:
```bash
sudo certbot --nginx -d your-domain.com
```

3. Auto-renewal is configured automatically

4. Update Nginx configuration to uncomment SSL directives in `/etc/nginx/sites-available/daybook`

## Monitoring and Maintenance

### Regular Maintenance Tasks

1. **Database Backups** (Recommended: Daily)
```bash
./scripts/backup_database.sh
```

2. **System Updates**
```bash
sudo apt-get update && sudo apt-get upgrade -y
```

3. **Log Rotation** (Automatic via systemd/logrotate)

4. **Disk Space Monitoring**
```bash
df -h
du -sh /opt/daybook/*
```

### Monitoring

Check application health:
```bash
curl http://localhost/health
```

Monitor system resources:
```bash
htop
```

Database size:
```bash
sudo -u postgres psql -c "SELECT pg_size_pretty(pg_database_size('daybook_prod'));"
```

## Troubleshooting

### Backend Service Won't Start

1. Check logs:
```bash
sudo journalctl -u daybook-backend -n 100 --no-pager
```

2. Check configuration:
```bash
cat /opt/daybook/backend/.env
```

3. Test database connection:
```bash
cd /opt/daybook/backend
sudo -u daybook ./daybook-backend
```

### Frontend Not Loading

1. Check Nginx:
```bash
sudo nginx -t
sudo systemctl status nginx
```

2. Check files exist:
```bash
ls -la /opt/daybook/frontend/dist/
```

3. Check Nginx logs:
```bash
sudo tail -f /var/log/nginx/daybook-error.log
```

### Database Connection Issues

1. Verify PostgreSQL is running:
```bash
sudo systemctl status postgresql
```

2. Test connection:
```bash
PGPASSWORD=your_password psql -h localhost -U daybook_user -d daybook_prod
```

3. Check pg_hba.conf:
```bash
sudo nano /etc/postgresql/*/main/pg_hba.conf
```

### Permission Issues

Fix ownership:
```bash
sudo chown -R daybook:daybook /opt/daybook
sudo chown -R daybook:daybook /var/log/daybook
sudo chmod 775 /opt/daybook/backend/uploads
```

## Security Recommendations

1. **Firewall**: Ensure only necessary ports are open
```bash
sudo ufw status
```

2. **SSH**: Disable password authentication, use SSH keys only

3. **Database**: Ensure PostgreSQL only listens on localhost

4. **Secrets**: Keep `.env` files secure (600 permissions)

5. **Updates**: Regularly update system packages

6. **Backups**: Automate database backups and test restores

7. **SSL/TLS**: Use HTTPS in production

8. **Monitoring**: Set up monitoring and alerting

## CI/CD Integration

To automate deployments from CI/CD:

1. Set up SSH key for CI/CD runner
2. Configure secrets in your CI/CD platform
3. Add deployment step:

```yaml
# Example GitHub Actions
- name: Deploy to Production
  run: |
    ssh ubuntu@${{ secrets.SERVER_IP }} "
      cd /opt/daybook &&
      git pull origin master &&
      cd deploy &&
      ./deploy.sh --skip-deps --skip-db
    "
```

## Rollback

To rollback to a previous version:

1. Restore database backup:
```bash
./scripts/restore_database.sh /opt/daybook/backups/daybook_backup_TIMESTAMP.sql.gz
```

2. Deploy previous code version:
```bash
cd /opt/daybook
git checkout <previous-commit>
cd deploy
./deploy.sh --skip-deps --skip-db
```

## Support

For issues or questions:
1. Check logs in `/var/log/daybook/` and `/var/log/nginx/`
2. Review service status: `sudo systemctl status daybook-backend`
3. Verify configurations in `/opt/daybook/`

## License

This deployment system is part of the Daybook project.
