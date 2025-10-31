# Daybook Deployment Guide

## Overview

This deployment system ensures your git repository remains clean by using a separate build directory. The deployment follows this architecture:

```
/opt/daybook              → Source code (git repository - READ ONLY)
/opt/daybook-build        → Temporary build directory (created during deployment)
/opt/daybook-app          → Final deployment location (compiled artifacts)
```

## Prerequisites

1. Ubuntu 20.04 or 22.04 LTS
2. Git repository cloned at `/opt/daybook`
3. Sudo access
4. Internet connection

## Directory Structure

- **SOURCE_DIR** (`/opt/daybook`): Your git repository - never modified by deployment scripts
- **BUILD_DIR** (`/opt/daybook-build`): Temporary directory where code is copied and built
- **APP_DIR** (`/opt/daybook-app`): Where compiled applications run from

## Initial Setup

### 1. Clone Repository on Server

```bash
# SSH into your server
ssh ubuntu@your-server-ip

# Clone the repository
sudo mkdir -p /opt
cd /opt
sudo git clone <your-repo-url> daybook
sudo chown -R ubuntu:ubuntu /opt/daybook
```

### 2. Configure Deployment

```bash
cd /opt/daybook/deploy

# Edit the configuration file
nano config/deploy.conf

# Update these values:
# - SERVER_IP: Your server IP
# - DOMAIN_OR_IP: Your domain or IP address
# - DB_PASSWORD: Database password (leave empty to auto-generate)
```

### 3. Run Initial Deployment

```bash
cd /opt/daybook/deploy

# Fresh installation (installs all dependencies)
./deploy.sh --fresh

# Or regular deployment (if dependencies already installed)
./deploy.sh
```

## Deployment Process

The deployment script performs these steps:

1. **Copies** source code from `/opt/daybook` to `/opt/daybook-build`
2. **Builds** the application in the build directory
3. **Deploys** compiled artifacts to `/opt/daybook-app`
4. **Cleans up** the build directory
5. **No changes** are made to your git repository

## Common Deployment Commands

### Full Deployment

```bash
cd /opt/daybook/deploy
./deploy.sh
```

### Update Backend Only

```bash
cd /opt/daybook/deploy/scripts
./update_backend.sh
```

### Update Frontend Only

```bash
cd /opt/daybook/deploy/scripts
./update_frontend.sh
```

### Deploy with Options

```bash
# Skip dependency installation
./deploy.sh --skip-deps

# Skip database setup
./deploy.sh --skip-db

# Skip backend deployment
./deploy.sh --skip-backend

# Skip frontend deployment
./deploy.sh --skip-frontend

# Combine multiple options
./deploy.sh --skip-deps --skip-db
```

## Updating Application Code

When you have new code to deploy:

```bash
# 1. SSH into server
ssh ubuntu@your-server-ip

# 2. Pull latest changes from git
cd /opt/daybook
git pull origin master  # or your branch name

# 3. Run deployment
cd deploy
./deploy.sh
```

**The deployment will:**
- Copy fresh code from `/opt/daybook` to build directory
- Build the new version
- Deploy it to `/opt/daybook-app`
- Your git repository at `/opt/daybook` remains clean

## Service Management

### Backend Service

```bash
# Check status
sudo systemctl status daybook-backend

# View logs
sudo journalctl -u daybook-backend -f

# Restart service
sudo systemctl restart daybook-backend

# Stop service
sudo systemctl stop daybook-backend

# Start service
sudo systemctl start daybook-backend
```

### Frontend (Nginx)

```bash
# Check status
sudo systemctl status nginx

# View logs
sudo tail -f /var/log/nginx/daybook-access.log
sudo tail -f /var/log/nginx/daybook-error.log

# Restart Nginx
sudo systemctl restart nginx

# Reload Nginx (graceful)
sudo systemctl reload nginx
```

## Troubleshooting

### Build Fails

If the build fails, check:
```bash
# Check if source directory exists
ls -la /opt/daybook

# Check build logs
# Backend: Check Go installation
go version

# Frontend: Check Node.js installation
node --version
npm --version
```

### Service Won't Start

```bash
# Check backend logs
sudo journalctl -u daybook-backend -n 100

# Check if port is already in use
sudo netstat -tulpn | grep 8080

# Check environment file
sudo cat /opt/daybook-app/backend/.env
```

### Frontend Not Loading

```bash
# Check Nginx configuration
sudo nginx -t

# Check if frontend files exist
ls -la /opt/daybook-app/frontend/dist

# Check Nginx logs
sudo tail -f /var/log/nginx/daybook-error.log
```

### Database Connection Issues

```bash
# Test database connection
sudo -u postgres psql -l

# Check database credentials
PGPASSWORD='your-password' psql -h localhost -U daybook_user -d daybook_prod
```

## File Permissions

The deployment system automatically sets correct permissions:

- Application user: `daybook`
- Application group: `daybook`
- Backend directory: owned by `daybook:daybook`
- Frontend directory: owned by `daybook:daybook`
- Uploads directory: `775` permissions

## Backup and Restore

### Backup Database

```bash
cd /opt/daybook/deploy/scripts
./backup_database.sh
```

Backups are stored in: `/opt/daybook/backups/`

### Restore Database

```bash
cd /opt/daybook/deploy/scripts
./restore_database.sh /path/to/backup.sql
```

## Security Notes

1. **Never commit** `.env` files to git
2. The deployment creates `.env` files from templates
3. Database passwords are stored securely with `600` permissions
4. Backend `.env` is readable only by the `daybook` user
5. Your git repository remains clean and secure

## Advanced Configuration

### Custom Build Directory

Edit `deploy/config/deploy.conf`:
```bash
BUILD_DIR="/custom/build/path"
```

### Keep Build Directory for Debugging

Comment out the cleanup line in deployment scripts:
```bash
# In 06_deploy_backend.sh and 07_deploy_frontend.sh
# sudo rm -rf "$BUILD_BACKEND_DIR"
```

### SSL/HTTPS Setup

After initial deployment, configure SSL:
```bash
# Install certbot
sudo apt-get install -y certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d your-domain.com

# Certificate will auto-renew
```

## Getting Help

- Check logs: `sudo journalctl -u daybook-backend -f`
- Check deployment script output
- Verify configuration: `cat deploy/config/deploy.conf`
- Test connectivity: `curl http://localhost:8080/health`

## Summary

✅ **Clean git repository** - No build artifacts in source
✅ **Isolated builds** - Build in temporary directory
✅ **Simple updates** - Just `git pull` and `./deploy.sh`
✅ **Safe deployments** - No risk of modifying source code
✅ **Easy rollbacks** - Can always rebuild from git
