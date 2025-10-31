# Daybook Deployment - Quick Start Guide

Get your Daybook application deployed in minutes!

## Prerequisites

- Ubuntu 20.04+ server on AWS
- SSH access with sudo privileges
- Your server's IP address or domain name

## Step-by-Step Deployment

### 1. Prepare Configuration (On Local Machine)

```bash
cd deploy
nano config/deploy.conf
```

**Update these required fields:**
```bash
DOMAIN_OR_IP="YOUR_SERVER_IP_OR_DOMAIN"  # e.g., "1.2.3.4" or "daybook.example.com"
```

Optional fields (will auto-generate if left empty):
```bash
DB_PASSWORD=""  # Leave empty for auto-generation
JWT_SECRET=""   # Leave empty for auto-generation
```

### 2. Copy Project to Server

```bash
# From project root on your local machine
rsync -avz --exclude 'node_modules' --exclude '.git' \
  ./ ubuntu@YOUR_SERVER_IP:/tmp/daybook
```

### 3. SSH into Server

```bash
ssh ubuntu@YOUR_SERVER_IP
```

### 4. Run Deployment

```bash
cd /tmp/daybook/deploy

# Make deploy script executable (if not already)
chmod +x deploy.sh

# Run fresh installation
./deploy.sh --fresh
```

**What this does:**
1. Installs Go, Node.js, and system dependencies
2. Installs and configures PostgreSQL
3. Installs and configures Redis
4. Installs and configures Nginx
5. Creates application user and directories
6. Builds and deploys backend
7. Builds and deploys frontend
8. Runs database migrations
9. Starts all services

### 5. Verify Deployment

After deployment completes, test your application:

```bash
# Check backend health
curl http://localhost/health

# Check services
sudo systemctl status daybook-backend
sudo systemctl status nginx
```

Access your application:
- Frontend: `http://YOUR_SERVER_IP`
- Backend API: `http://YOUR_SERVER_IP/api/v1`

## Common Commands

### View Logs
```bash
# Backend logs
sudo journalctl -u daybook-backend -f

# Nginx logs
sudo tail -f /var/log/nginx/daybook-access.log
sudo tail -f /var/log/nginx/daybook-error.log
```

### Restart Services
```bash
# Restart backend
sudo systemctl restart daybook-backend

# Reload Nginx
sudo systemctl reload nginx
```

### Update Application
```bash
# Update backend only
cd /home/ubuntu/projects/opt/daybook/deploy
./scripts/update_backend.sh

# Update frontend only
cd /home/ubuntu/projects/opt/daybook/deploy
./scripts/update_frontend.sh

# Update both
cd /home/ubuntu/projects/opt/daybook/deploy
./deploy.sh --skip-deps --skip-db
```

### Backup Database
```bash
cd /home/ubuntu/projects/opt/daybook/deploy
./scripts/backup_database.sh
```

## Troubleshooting

### Backend won't start
```bash
# Check logs
sudo journalctl -u daybook-backend -n 50

# Check environment
cat /home/ubuntu/projects/opt/daybook/backend/.env

# Test manually
cd /home/ubuntu/projects/opt/daybook/backend
sudo -u daybook ./daybook-backend
```

### Can't access frontend
```bash
# Test Nginx config
sudo nginx -t

# Check Nginx status
sudo systemctl status nginx

# Check frontend files exist
ls -la /home/ubuntu/projects/opt/daybook/frontend/dist/
```

### Database connection error
```bash
# Check PostgreSQL
sudo systemctl status postgresql

# Test connection
PGPASSWORD=your_password psql -h localhost -U daybook_user -d daybook_prod
```

## Security Hardening (Recommended)

### 1. Configure Firewall
```bash
sudo ufw allow OpenSSH
sudo ufw allow 'Nginx Full'
sudo ufw enable
```

### 2. Set up SSL/TLS with Let's Encrypt
```bash
sudo apt-get install -y certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 3. Secure SSH
Edit `/etc/ssh/sshd_config`:
```bash
PasswordAuthentication no
PermitRootLogin no
```

Then restart SSH:
```bash
sudo systemctl restart sshd
```

## Next Steps

1. ✅ Test all application features
2. ✅ Set up automated backups (add to crontab)
3. ✅ Configure SSL/TLS for HTTPS
4. ✅ Set up monitoring and alerts
5. ✅ Review and update security settings
6. ✅ Configure domain name (if using one)

## Important Files and Locations

```
/home/ubuntu/projects/opt/daybook/          - Application root
/home/ubuntu/projects/opt/daybook/backend/  - Backend application and binary
/home/ubuntu/projects/opt/daybook/frontend/ - Frontend built files
/var/log/daybook/      - Application logs
/home/ubuntu/projects/opt/daybook/backups/  - Database backups
```

## Need More Details?

See the full [README.md](README.md) for:
- Detailed documentation
- Individual script usage
- Advanced configuration
- CI/CD integration
- Monitoring setup

## Quick Reference

| Task | Command |
|------|---------|
| Deploy everything | `./deploy.sh --fresh` |
| Update only | `./deploy.sh --skip-deps --skip-db` |
| Backend logs | `sudo journalctl -u daybook-backend -f` |
| Restart backend | `sudo systemctl restart daybook-backend` |
| Backup database | `./scripts/backup_database.sh` |
| Update backend | `./scripts/update_backend.sh` |
| Update frontend | `./scripts/update_frontend.sh` |

---

**Estimated Time:** 10-15 minutes for fresh installation

**Support:** Check logs in `/var/log/daybook/` and `/var/log/nginx/` for issues
