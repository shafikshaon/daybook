# EC2 Deployment Checklist

Quick reference checklist for deploying Daybook to EC2.

## Pre-Deployment (Local Machine)

### 1. Build Application
- [ ] Navigate to project directory
- [ ] Run build script: `./scripts/local-build.sh`
- [ ] Verify build artifacts in `build/` directory:
  - [ ] `daybook-backend` (Linux binary)
  - [ ] `frontend-dist.tar.gz`

### 2. Prepare EC2 Instance
- [ ] EC2 instance is running (Ubuntu 22.04 LTS)
- [ ] Security Group allows:
  - [ ] SSH (port 22)
  - [ ] HTTP (port 80)
  - [ ] HTTPS (port 443)
  - [ ] Custom TCP (ports 8081, 8082) - optional, for testing
- [ ] You have SSH key (.pem file)
- [ ] You know the EC2 public IP address

### 3. Configure DNS
- [ ] Create A record: `daybook.shafik.xyz` → EC2 IP
- [ ] Create A record: `api.daybook.shafik.xyz` → EC2 IP
- [ ] Wait for DNS propagation (verify with `nslookup`)

---

## EC2 Setup

### 4. Initial Server Setup
- [ ] Connect to EC2: `ssh -i your-key.pem ubuntu@your-ec2-ip`
- [ ] Update system: `sudo apt update && sudo apt upgrade -y`
- [ ] Install essential tools:
  ```bash
  sudo apt install -y curl wget git nano ufw software-properties-common
  ```

### 5. Configure Firewall
- [ ] Enable UFW rules:
  ```bash
  sudo ufw allow 22/tcp
  sudo ufw allow 80/tcp
  sudo ufw allow 443/tcp
  sudo ufw allow 8081/tcp
  sudo ufw allow 8082/tcp
  sudo ufw --force enable
  ```
- [ ] Verify: `sudo ufw status`

### 6. Install PostgreSQL
- [ ] Add PostgreSQL repository
- [ ] Install PostgreSQL 15
- [ ] Optimize for low memory (edit `/etc/postgresql/15/main/postgresql.conf`)
- [ ] Create database and user:
  ```sql
  CREATE DATABASE daybook;
  CREATE USER daybook_user WITH ENCRYPTED PASSWORD 'your_password';
  GRANT ALL PRIVILEGES ON DATABASE daybook TO daybook_user;
  ```
- [ ] Test connection: `psql -h localhost -U daybook_user -d daybook`

### 7. Install Redis
- [ ] Add Redis repository
- [ ] Install Redis 7
- [ ] Configure for low memory (edit `/etc/redis/redis.conf`):
  - [ ] Set `maxmemory 50mb`
  - [ ] Set `requirepass your_redis_password`
- [ ] Restart and test: `redis-cli` then `AUTH password` then `PING`

### 8. Install Nginx
- [ ] Install: `sudo apt install -y nginx`
- [ ] Start: `sudo systemctl start nginx`
- [ ] Enable: `sudo systemctl enable nginx`

---

## Application Deployment

### 9. Deploy Application (Option A: Manual)
On **EC2**:
- [ ] Create directories:
  ```bash
  sudo mkdir -p /var/www/daybook/{backend,frontend}
  sudo chown -R ubuntu:ubuntu /var/www/daybook
  ```
- [ ] Deploy backend:
  ```bash
  cp ~/daybook/build/daybook-backend /var/www/daybook/backend/
  cp ~/daybook/build/frontend-dist.tar.gz /var/www/daybook/frontend/
  chmod +x /var/www/daybook/backend/daybook-backend
  ```
- [ ] Deploy frontend:
  ```bash
  cd /var/www/daybook/frontend
  tar -xzf frontend-dist.tar.gz
  ```

### 12. Configure Backend Environment
- [ ] Create `/var/www/daybook/backend/.env` with:
  - [ ] Database credentials
  - [ ] Redis password
  - [ ] JWT secret (generate with `openssl rand -base64 48`)
  - [ ] CORS origins: `https://daybook.shafik.xyz`
  - [ ] Upload path
- [ ] Secure permissions: `chmod 600 /var/www/daybook/backend/.env`

### 13. Create Systemd Service
- [ ] Create `/etc/systemd/system/daybook-backend.service`
- [ ] Create log files:
  ```bash
  sudo touch /var/log/daybook-backend{,-error}.log
  sudo chown ubuntu:ubuntu /var/log/daybook-backend*.log
  ```
- [ ] Enable and start:
  ```bash
  sudo systemctl daemon-reload
  sudo systemctl enable daybook-backend
  sudo systemctl start daybook-backend
  ```
- [ ] Verify: `sudo systemctl status daybook-backend`

### 14. Configure Nginx
- [ ] Create `/etc/nginx/sites-available/daybook-frontend`
- [ ] Create `/etc/nginx/sites-available/daybook-backend`
- [ ] Enable sites:
  ```bash
  sudo ln -s /etc/nginx/sites-available/daybook-frontend /etc/nginx/sites-enabled/
  sudo ln -s /etc/nginx/sites-available/daybook-backend /etc/nginx/sites-enabled/
  ```
- [ ] Test config: `sudo nginx -t`
- [ ] Reload: `sudo systemctl reload nginx`

---

## SSL/TLS Setup

### 15. Install Certbot
- [ ] Install: `sudo apt install -y certbot python3-certbot-nginx`

### 16. Obtain SSL Certificates
- [ ] For frontend:
  ```bash
  sudo certbot --nginx -d daybook.shafik.xyz
  ```
- [ ] For backend:
  ```bash
  sudo certbot --nginx -d api.daybook.shafik.xyz
  ```
- [ ] Select option 2 (redirect HTTP to HTTPS)
- [ ] Test auto-renewal: `sudo certbot renew --dry-run`

---

## Verification

### 17. Test Services
- [ ] Backend status: `sudo systemctl status daybook-backend`
- [ ] PostgreSQL status: `sudo systemctl status postgresql`
- [ ] Redis status: `sudo systemctl status redis-server`
- [ ] Nginx status: `sudo systemctl status nginx`

### 18. Check Logs
- [ ] Backend logs: `sudo journalctl -u daybook-backend -f`
- [ ] Nginx access: `sudo tail -f /var/log/nginx/daybook-frontend-access.log`
- [ ] Nginx errors: `sudo tail -f /var/log/nginx/daybook-backend-error.log`

### 19. Test Endpoints
- [ ] Frontend: `curl https://daybook.shafik.xyz`
- [ ] Backend: `curl https://api.daybook.shafik.xyz` (should see API response)
- [ ] Open in browser: https://daybook.shafik.xyz
- [ ] Test signup/login functionality

### 20. Performance Check
- [ ] Memory usage: `free -h`
- [ ] Disk usage: `df -h`
- [ ] Process list: `htop`

---

## Post-Deployment

### 21. Security
- [ ] Verify firewall: `sudo ufw status`
- [ ] Check .env permissions: `ls -la /var/www/daybook/backend/.env` (should be 600)
- [ ] Review PostgreSQL config (local access only)
- [ ] Review Redis config (password protected)
- [ ] Test HTTPS redirect (HTTP should redirect to HTTPS)

### 22. Backups
- [ ] Create backup script: `~/backup-db.sh`
- [ ] Make executable: `chmod +x ~/backup-db.sh`
- [ ] Add to crontab: `crontab -e`
  ```
  0 2 * * * /home/ubuntu/backup-db.sh
  ```
- [ ] Test backup: `./backup-db.sh`

### 23. Monitoring
- [ ] Set up log rotation (if needed)
- [ ] Document application URLs
- [ ] Create runbook for common operations
- [ ] Test backup restoration procedure

### 24. Documentation
- [ ] Document database credentials (secure location)
- [ ] Document Redis password
- [ ] Document JWT secret
- [ ] Save .env file backup (secure location)
- [ ] Note EC2 instance details
- [ ] Update team documentation

---

## Common Issues

### Backend Won't Start
1. Check logs: `sudo journalctl -u daybook-backend -n 50`
2. Verify binary: `ls -l /var/www/daybook/backend/daybook-backend`
3. Check .env file: `cat /var/www/daybook/backend/.env`
4. Test database connection: `psql -h localhost -U daybook_user -d daybook`

### 502 Bad Gateway
1. Check if backend is running: `sudo systemctl status daybook-backend`
2. Check Nginx config: `sudo nginx -t`
3. Check port binding: `sudo netstat -tulpn | grep 8082`

### Out of Memory
1. Check usage: `free -h`
2. Add swap:
   ```bash
   sudo fallocate -l 1G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```
3. Make permanent: `echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab`

### SSL Issues
1. Check certificates: `sudo certbot certificates`
2. Renew manually: `sudo certbot renew`
3. Check Nginx config: `sudo nginx -t`

---

## Quick Commands

```bash
# Restart everything
sudo systemctl restart daybook-backend
sudo systemctl reload nginx

# View logs
sudo journalctl -u daybook-backend -f

# Check all services
sudo systemctl status daybook-backend postgresql redis-server nginx

# Database backup
pg_dump -h localhost -U daybook_user -d daybook > backup_$(date +%Y%m%d).sql

# Monitor resources
htop
```

---

## Contact & Support

- Main Guide: See `EC2_DEPLOYMENT_GUIDE.md` for detailed instructions
- Build Script: `scripts/local-build.sh`
- Setup Script: `scripts/ec2-setup.sh`

**Deployment Date:** _______________

**Deployed By:** _______________

**Notes:**
_______________________________________________
_______________________________________________
_______________________________________________
