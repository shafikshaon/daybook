# Production Deployment Checklist

## Issue: /api/v1/auth/me Returns 404 in Production

The endpoint is fully implemented. Follow this checklist to fix your production deployment.

## Quick Diagnosis

Run this command on your **production server**:

```bash
curl http://localhost:8080/health
```

**If you get a response** → Backend is running, skip to [Step 3](#step-3-check-nginx)
**If connection refused** → Backend is not running, follow from [Step 1](#step-1-verify-backend-files)

---

## Step 1: Verify Backend Files

SSH into your production server:

```bash
ssh user@your-production-server
```

Check if backend files exist:

```bash
cd /path/to/daybook/backend
ls -la daybook-backend
ls -la .env
```

**Expected**: Both files should exist

**If missing**:
```bash
# Pull latest code
git pull origin master

# Build backend
go build -o daybook-backend

# Copy .env.example to .env and configure
cp .env.example .env
nano .env
```

---

## Step 2: Start Backend Service

### Option A: Using systemd (Recommended)

Create service file:
```bash
sudo nano /etc/systemd/system/daybook-backend.service
```

Paste this configuration:
```ini
[Unit]
Description=Daybook Backend API
After=network.target postgresql.service

[Service]
Type=simple
User=your-username
WorkingDirectory=/path/to/daybook/backend
Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"
ExecStart=/path/to/daybook/backend/daybook-backend
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**Replace**:
- `your-username` with your actual username
- `/path/to/daybook` with actual path

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable daybook-backend
sudo systemctl start daybook-backend
sudo systemctl status daybook-backend
```

### Option B: Using PM2

```bash
cd /path/to/daybook/backend
pm2 start ./daybook-backend --name daybook-backend
pm2 save
pm2 startup
```

### Option C: Direct Execution (Testing Only)

```bash
cd /path/to/daybook/backend
./daybook-backend
```

**Verify backend started**:
```bash
# Check process
ps aux | grep daybook-backend

# Check port
sudo lsof -i :8080

# Test locally
curl http://localhost:8080/health
```

**Expected**: `{"status":"ok","message":"Daybook API is running"}`

---

## Step 3: Check Nginx Configuration

### Find your nginx config:
```bash
ls /etc/nginx/sites-available/
ls /etc/nginx/conf.d/
```

Look for `api.daybook.shafik.xyz` or similar.

### Verify configuration:

```bash
cat /etc/nginx/sites-available/api.daybook.shafik.xyz
# OR
cat /etc/nginx/conf.d/api.daybook.shafik.xyz.conf
```

**Required configuration**:
```nginx
server {
    listen 443 ssl http2;
    server_name api.daybook.shafik.xyz;

    ssl_certificate /etc/letsencrypt/live/api.daybook.shafik.xyz/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.daybook.shafik.xyz/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**If config is missing or wrong**:
```bash
sudo nano /etc/nginx/sites-available/api.daybook.shafik.xyz
# Paste the configuration above

# Enable it
sudo ln -s /etc/nginx/sites-available/api.daybook.shafik.xyz /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

---

## Step 4: Verify .env Configuration

```bash
cd /path/to/daybook/backend
cat .env
```

**Required settings**:
```bash
SERVER_PORT=8080
SERVER_MODE=release

DB_HOST=your-db-host
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=daybook

CORS_ALLOWED_ORIGINS=http://localhost:3000,https://daybook.shafik.xyz,https://api.daybook.shafik.xyz
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_ALLOW_CREDENTIALS=true

JWT_SECRET=your-secret-key
JWT_EXPIRATION=168
```

**If CORS_ALLOWED_ORIGINS is missing or wrong**:
```bash
nano .env
# Add or update the CORS settings
# Restart backend after changes
```

---

## Step 5: Test Everything

### Test Backend Locally (on server)
```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/auth/me
```

**Expected**:
- Health: HTTP 200 - `{"status":"ok"}`
- /auth/me: HTTP 401 - `{"error":"Authorization header required"}`

### Test Public API (from anywhere)
```bash
curl https://api.daybook.shafik.xyz/health
curl https://api.daybook.shafik.xyz/api/v1/auth/me
```

**Expected**:
- Health: HTTP 200 - `{"status":"ok"}`
- /auth/me: HTTP 401 (NOT 404!)

### Test Frontend
1. Open https://daybook.shafik.xyz in browser
2. Open DevTools (F12) → Console
3. Login with your credentials
4. Check for errors in console

**Expected**: No 404 errors, successful login

---

## Step 6: Check Firewall

```bash
# If using ufw
sudo ufw status
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw reload

# If using iptables
sudo iptables -L -n | grep -E "80|443"
```

---

## Step 7: Check Logs

### Backend Logs
```bash
# If using systemd
sudo journalctl -u daybook-backend -f

# If using pm2
pm2 logs daybook-backend

# If using log file
tail -f /var/log/daybook-backend.log
```

### Nginx Logs
```bash
sudo tail -f /var/log/nginx/error.log
sudo tail -f /var/log/nginx/access.log
```

---

## Automated Diagnostic

Run this script on your production server:

```bash
cd /path/to/daybook
chmod +x diagnose-production.sh
sudo ./diagnose-production.sh
```

It will automatically check all of the above and tell you exactly what's wrong.

---

## Common Issues & Solutions

### Issue: "Connection refused"
**Cause**: Backend not running
**Fix**: Follow Step 2 to start backend

### Issue: "502 Bad Gateway"
**Cause**: Nginx can't connect to backend
**Fix**:
- Check backend is running on port 8080
- Check nginx proxy_pass is `http://localhost:8080`

### Issue: "404 Not Found"
**Cause**: Backend not running OR routes not registered
**Fix**:
- Start backend (Step 2)
- Verify latest code is deployed (Step 1)

### Issue: "CORS error in browser"
**Cause**: CORS_ALLOWED_ORIGINS missing or wrong
**Fix**: Update .env (Step 4), restart backend

### Issue: "SSL certificate error"
**Cause**: Certificate missing or expired
**Fix**:
```bash
sudo certbot --nginx -d api.daybook.shafik.xyz
```

---

## Success Verification

When everything is working, you should see:

### Backend Process Running
```bash
$ ps aux | grep daybook-backend
user  12345  0.0  1.2  daybook-backend
```

### Port 8080 Listening
```bash
$ sudo lsof -i :8080
COMMAND   PID  USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
daybook  12345 user   3u  IPv6  12345      0t0  TCP *:8080 (LISTEN)
```

### Health Check OK
```bash
$ curl https://api.daybook.shafik.xyz/health
{"status":"ok","message":"Daybook API is running"}
```

### Auth Endpoint Returns 401 (Not 404)
```bash
$ curl https://api.daybook.shafik.xyz/api/v1/auth/me
{"success":false,"error":"Authorization header required"}
```

### Frontend Working
- No console errors
- Login works
- Export features work
- No 404 on /auth/me

---

## Quick Commands Cheat Sheet

```bash
# Check if backend is running
ps aux | grep daybook-backend
sudo lsof -i :8080

# Start backend (systemd)
sudo systemctl start daybook-backend
sudo systemctl status daybook-backend

# Start backend (pm2)
pm2 start daybook-backend
pm2 status

# Test backend
curl http://localhost:8080/health
curl https://api.daybook.shafik.xyz/health

# Check nginx
sudo nginx -t
sudo systemctl status nginx
sudo systemctl reload nginx

# View logs
sudo journalctl -u daybook-backend -f
sudo tail -f /var/log/nginx/error.log

# Restart everything
sudo systemctl restart daybook-backend
sudo systemctl reload nginx
```

---

## Still Not Working?

1. Run the diagnostic script: `./diagnose-production.sh`
2. Share the output
3. Check `PRODUCTION_DEPLOYMENT_GUIDE.md` for detailed troubleshooting

---

## Files Reference

- ✅ `AUTH_ME_ENDPOINT_STATUS.md` - Confirms endpoint is implemented
- ✅ `PRODUCTION_QUICK_FIX.md` - Quick fixes for common issues
- ✅ `PRODUCTION_DEPLOYMENT_GUIDE.md` - Detailed deployment guide
- ✅ `diagnose-production.sh` - Automated diagnostics
- ✅ `test-auth-me-endpoint.sh` - Test endpoint locally
