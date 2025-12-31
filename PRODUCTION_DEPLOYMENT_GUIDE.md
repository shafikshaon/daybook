# Production Deployment Guide

## Current Issue

Your production backend at `https://api.daybook.shafik.xyz` is returning 404 errors:

```
GET https://api.daybook.shafik.xyz/api/v1/auth/me 404 (Not Found)
```

This guide will help you diagnose and fix the issue.

## Step 1: Verify Backend is Running on Production Server

SSH into your production server and check if the backend process is running:

```bash
# Check if the backend process is running
ps aux | grep daybook-backend

# Or if you're using systemd
sudo systemctl status daybook-backend

# Or if you're using pm2
pm2 list
```

**Expected**: You should see a running process
**If not running**: Start the backend service

## Step 2: Check Backend Logs

View the backend logs to see what's happening:

```bash
# If using systemd
sudo journalctl -u daybook-backend -f

# If using a log file
tail -f /var/log/daybook-backend.log

# Or check the log location from your .env
# Check where logs are being written
```

**Look for**:
- Server startup messages
- What port it's listening on
- Any error messages
- Route registration logs

## Step 3: Test Backend Directly on Server

SSH into your production server and test the backend locally:

```bash
# Test if backend is responding on localhost
curl http://localhost:8080/health

# Test the auth endpoint
curl http://localhost:8080/api/v1/auth/me

# Check what's actually listening on port 8080
sudo netstat -tulpn | grep 8080
# OR
sudo lsof -i :8080
```

**Expected**: Backend should respond with JSON
**If not responding**: Backend is not running or crashed

## Step 4: Check Nginx/Reverse Proxy Configuration

Your domain `https://api.daybook.shafik.xyz` needs to be configured to proxy requests to your backend.

Find your nginx configuration:

```bash
# Common locations
ls /etc/nginx/sites-available/
ls /etc/nginx/conf.d/

# Look for daybook or api.daybook configuration
cat /etc/nginx/sites-available/api.daybook.shafik.xyz
# OR
cat /etc/nginx/conf.d/api.daybook.shafik.xyz.conf
```

### Correct Nginx Configuration

Your nginx config should look like this:

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name api.daybook.shafik.xyz;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name api.daybook.shafik.xyz;

    # SSL certificates (Let's Encrypt or your provider)
    ssl_certificate /etc/letsencrypt/live/api.daybook.shafik.xyz/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.daybook.shafik.xyz/privkey.pem;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Proxy to backend
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;

        # Important headers
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_cache_bypass $http_upgrade;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint (optional)
    location /health {
        proxy_pass http://localhost:8080/health;
        access_log off;
    }
}
```

### Apply Nginx Changes

```bash
# Test nginx configuration
sudo nginx -t

# If test passes, reload nginx
sudo systemctl reload nginx

# Check nginx status
sudo systemctl status nginx

# View nginx error logs if there are issues
sudo tail -f /var/log/nginx/error.log
```

## Step 5: Verify Backend Environment Variables

Make sure your production backend has the correct .env file:

```bash
# SSH into production server
cd /path/to/daybook/backend

# Check .env file exists
ls -la .env

# Verify critical settings
cat .env | grep -E "SERVER_PORT|DB_HOST|CORS_ALLOWED_ORIGINS"
```

**Required settings for production**:

```bash
SERVER_PORT=8080
SERVER_MODE=release

# Database
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=daybook
DB_SSLMODE=require

# CORS - CRITICAL!
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://daybook.shafik.xyz,https://api.daybook.shafik.xyz
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_EXPOSE_HEADERS=
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=12

# JWT
JWT_SECRET=your-production-secret
JWT_EXPIRATION=168
```

## Step 6: Check Firewall Rules

Ensure the firewall allows traffic:

```bash
# Check if firewall is blocking
sudo ufw status

# If using ufw, allow nginx
sudo ufw allow 'Nginx Full'
sudo ufw allow 'Nginx HTTPS'

# Or allow specific ports
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Reload firewall
sudo ufw reload
```

## Step 7: Test API Endpoints

Once everything is configured, test the API:

```bash
# From your local machine

# Test health endpoint
curl https://api.daybook.shafik.xyz/health

# Test auth endpoint (should return 401 without token)
curl https://api.daybook.shafik.xyz/api/v1/auth/me

# Test with a valid token
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://api.daybook.shafik.xyz/api/v1/auth/me
```

## Common Issues and Solutions

### Issue 1: "Connection refused"
**Cause**: Backend not running
**Solution**: Start the backend service

### Issue 2: "502 Bad Gateway"
**Cause**: Nginx can't connect to backend
**Solution**:
- Check backend is running on port 8080
- Check nginx proxy_pass is correct

### Issue 3: "404 Not Found"
**Cause**: Routes not properly registered or wrong URL path
**Solution**:
- Verify routes are registered in `routes/routes.go`
- Check nginx is not rewriting URLs incorrectly
- Ensure backend is running the latest code

### Issue 4: CORS errors
**Cause**: Backend CORS not allowing frontend domain
**Solution**: Update CORS_ALLOWED_ORIGINS in backend .env

### Issue 5: SSL Certificate errors
**Cause**: Missing or expired SSL certificate
**Solution**:
```bash
# Install/renew Let's Encrypt certificate
sudo certbot --nginx -d api.daybook.shafik.xyz
```

## Deployment Checklist

- [ ] Backend code deployed to production server
- [ ] .env file configured with production settings
- [ ] Database accessible from production server
- [ ] Backend service running (systemd/pm2)
- [ ] Nginx configured as reverse proxy
- [ ] SSL certificate installed for api.daybook.shafik.xyz
- [ ] Firewall allows HTTP/HTTPS traffic
- [ ] CORS configured for daybook.shafik.xyz
- [ ] Frontend .env.production has correct VITE_API_URL
- [ ] Frontend built with production environment
- [ ] Frontend deployed to web server

## Quick Diagnosis Commands

Run these commands on your production server to get a complete picture:

```bash
#!/bin/bash
echo "=== Backend Process ==="
ps aux | grep daybook-backend

echo -e "\n=== Port 8080 Listener ==="
sudo lsof -i :8080

echo -e "\n=== Nginx Status ==="
sudo systemctl status nginx | head -20

echo -e "\n=== Nginx Config Test ==="
sudo nginx -t

echo -e "\n=== Backend Local Test ==="
curl -s http://localhost:8080/health

echo -e "\n=== Public API Test ==="
curl -s https://api.daybook.shafik.xyz/health

echo -e "\n=== Nginx Error Logs (last 10 lines) ==="
sudo tail -10 /var/log/nginx/error.log

echo -e "\n=== CORS Settings ==="
grep CORS_ALLOWED_ORIGINS /path/to/daybook/backend/.env
```

## Next Steps

1. Run the diagnostic commands above
2. Share the output to identify the exact issue
3. Follow the relevant solution based on the diagnosis
4. Test the frontend at https://daybook.shafik.xyz once backend is working

## Contact

If you're still stuck, provide:
1. Output of the diagnostic commands
2. Nginx configuration file
3. Backend startup logs
4. Any error messages from browser console
