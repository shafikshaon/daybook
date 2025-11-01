# Troubleshooting Guide

## Backend Connection Refused (Port 8082)

**Error**: `connect() failed (111: Connection refused) while connecting to upstream`

This means the backend service isn't running.

### Step 1: Check Backend Service Status

```bash
sudo systemctl status daybook-backend
```

**Possible outputs:**

#### If service doesn't exist:
```
Unit daybook-backend.service could not be found
```
**Fix**: Create the systemd service (see Part 7 of EC2_DEPLOYMENT_GUIDE.md)

#### If service failed to start:
```
● daybook-backend.service - Daybook Backend API
   Loaded: loaded
   Active: failed (Result: exit-code)
```
**Fix**: Check the logs (see Step 2)

#### If service is inactive:
```
   Active: inactive (dead)
```
**Fix**: Start the service
```bash
sudo systemctl start daybook-backend
```

### Step 2: Check Backend Logs

```bash
# Check systemd journal
sudo journalctl -u daybook-backend -n 50 --no-pager

# Check application logs
sudo tail -100 /var/log/daybook-backend.log
sudo tail -100 /var/log/daybook-backend-error.log
```

### Common Issues and Fixes

#### Issue 1: Binary Not Found

**Error in logs:**
```
Failed to execute command: No such file or directory
```

**Check:**
```bash
ls -la /var/www/daybook/backend/daybook-backend
```

**Fix:**
```bash
# Copy from build directory
cp ~/daybook/build/daybook-backend /var/www/daybook/backend/daybook-backend
chmod +x /var/www/daybook/backend/daybook-backend
```

#### Issue 2: Database Connection Failed

**Error in logs:**
```
failed to connect to database
FATAL: password authentication failed
```

**Fix:**
```bash
# Check .env file
cat /var/www/daybook/backend/.env | grep DB_

# Test database connection
psql -h localhost -U daybook_user -d daybook
```

If password is wrong, update `/var/www/daybook/backend/.env`:
```bash
nano /var/www/daybook/backend/.env
# Update DB_PASSWORD
```

#### Issue 3: Redis Connection Failed

**Error in logs:**
```
failed to connect to redis
NOAUTH Authentication required
```

**Fix:**
```bash
# Check Redis is running
sudo systemctl status redis-server

# Test Redis connection
redis-cli
AUTH your_redis_password
PING
```

If password is wrong, update `/var/www/daybook/backend/.env`:
```bash
nano /var/www/daybook/backend/.env
# Update REDIS_PASSWORD
```

#### Issue 4: Port Already in Use

**Error in logs:**
```
bind: address already in use
```

**Check what's using port 8082:**
```bash
sudo netstat -tulpn | grep 8082
# or
sudo lsof -i :8082
```

**Fix:**
```bash
# Kill the process
sudo kill -9 <PID>

# Restart backend
sudo systemctl restart daybook-backend
```

#### Issue 5: Permission Denied

**Error in logs:**
```
permission denied
```

**Fix:**
```bash
# Check binary permissions
ls -la /var/www/daybook/backend/daybook-backend

# Make executable
chmod +x /var/www/daybook/backend/daybook-backend

# Check directory permissions
ls -la /var/www/daybook/backend/
sudo chown -R ubuntu:ubuntu /var/www/daybook
```

#### Issue 6: Missing .env File

**Error in logs:**
```
unable to load .env file
```

**Fix:**
```bash
# Create .env file
nano /var/www/daybook/backend/.env
```

See EC2_DEPLOYMENT_GUIDE.md Part 6.3.2 for .env template.

### Step 3: Test Backend Manually

```bash
# Try running the backend directly
cd /var/www/daybook/backend
./daybook-backend

# You should see startup logs
# Press Ctrl+C to stop
```

If it starts successfully, the issue is with the systemd service.

### Step 4: Restart Services

```bash
# Reload systemd
sudo systemctl daemon-reload

# Start backend
sudo systemctl start daybook-backend

# Check status
sudo systemctl status daybook-backend

# Enable auto-start on boot
sudo systemctl enable daybook-backend
```

### Step 5: Verify Backend is Running

```bash
# Check if port 8082 is listening
sudo netstat -tulpn | grep 8082

# Test backend directly
curl http://localhost:8082

# Test through Nginx
curl https://api.daybook.shafik.xyz
```

---

## Frontend Issues

### Issue 1: 404 Not Found

**Error**: Files not found

**Fix:**
```bash
# Check files exist
ls -la /var/www/daybook/frontend/index.html

# If missing, extract again
cd /var/www/daybook/frontend
tar -xzf ~/daybook/build/frontend-dist.tar.gz
```

### Issue 2: Blank Page

**Check browser console for errors (F12)**

Common causes:
- API URL misconfigured (check build-time .env.production)
- CORS errors (check backend CORS_ALLOWED_ORIGINS)
- Backend not running

---

## Database Issues

### Check PostgreSQL Status

```bash
sudo systemctl status postgresql

# View logs
sudo tail -100 /var/log/postgresql/postgresql-15-main.log
```

### Test Connection

```bash
psql -h localhost -U daybook_user -d daybook
```

### Reset Password

```bash
sudo -u postgres psql
ALTER USER daybook_user WITH PASSWORD 'new_password';
\q
```

Then update `/var/www/daybook/backend/.env`

---

## Redis Issues

### Check Redis Status

```bash
sudo systemctl status redis-server
```

### Test Connection

```bash
redis-cli
AUTH your_password
PING
```

### Reset Password

```bash
sudo nano /etc/redis/redis.conf
# Update: requirepass your_new_password

sudo systemctl restart redis-server
```

Then update `/var/www/daybook/backend/.env`

---

## SSL/HTTPS Issues

### Certificate Errors

```bash
# Check certificates
sudo certbot certificates

# Renew if expired
sudo certbot renew --force-renewal
```

### Nginx Configuration Errors

```bash
# Test configuration
sudo nginx -t

# View errors
sudo tail -50 /var/log/nginx/error.log
```

---

## Monitoring Commands

### Check All Services

```bash
sudo systemctl status daybook-backend postgresql redis-server nginx
```

### View All Logs

```bash
# Backend
sudo journalctl -u daybook-backend -f

# Nginx
sudo tail -f /var/log/nginx/daybook-frontend-access.log
sudo tail -f /var/log/nginx/daybook-backend-access.log
sudo tail -f /var/log/nginx/daybook-backend-error.log

# PostgreSQL
sudo tail -f /var/log/postgresql/postgresql-15-main.log
```

### Check Resource Usage

```bash
# Memory
free -h

# Disk
df -h

# Processes
htop

# Network
sudo netstat -tulpn
```

---

## Complete Reset

If everything is broken and you want to start fresh:

```bash
# Stop all services
sudo systemctl stop daybook-backend nginx

# Remove deployment
sudo rm -rf /var/www/daybook

# Re-run setup
cd ~/daybook
git pull origin master
./scripts/ec2-setup.sh

# Start services
sudo systemctl start daybook-backend
sudo systemctl reload nginx
```

---

## Getting Help

When asking for help, provide:

1. **Service status:**
   ```bash
   sudo systemctl status daybook-backend
   ```

2. **Recent logs:**
   ```bash
   sudo journalctl -u daybook-backend -n 50 --no-pager
   ```

3. **Configuration:**
   ```bash
   cat /var/www/daybook/backend/.env | grep -v PASSWORD
   ```

4. **Port status:**
   ```bash
   sudo netstat -tulpn | grep 8082
   ```

5. **File permissions:**
   ```bash
   ls -la /var/www/daybook/backend/daybook-backend
   ```
