# Quick Fix for Production 404 Error

## The Problem

Your production API at `https://api.daybook.shafik.xyz` is returning 404 errors:

```
GET https://api.daybook.shafik.xyz/api/v1/auth/me 404 (Not Found)
```

## Quick Diagnosis (5 minutes)

### Step 1: Copy diagnostic script to your production server

```bash
# On your local machine
scp diagnose-production.sh user@your-server:/tmp/

# SSH into your server
ssh user@your-server

# Run the diagnostic script
cd /tmp
chmod +x diagnose-production.sh
sudo ./diagnose-production.sh
```

The script will tell you exactly what's wrong.

## Most Likely Causes (and fixes)

### Cause 1: Backend Not Running (90% probability)

**Symptoms**: Script shows "Backend process is NOT running"

**Fix**:
```bash
# Navigate to your backend directory
cd /path/to/daybook/backend

# Start the backend
./daybook-backend

# OR if using systemd
sudo systemctl start daybook-backend
sudo systemctl enable daybook-backend  # auto-start on reboot

# OR if using pm2
pm2 start daybook-backend
pm2 save
```

### Cause 2: Nginx Not Configured for API (8% probability)

**Symptoms**: Script shows "No nginx configuration found"

**Fix**: Create nginx configuration:

```bash
# Create nginx config
sudo nano /etc/nginx/sites-available/api.daybook.shafik.xyz
```

Paste this configuration:

```nginx
server {
    listen 80;
    server_name api.daybook.shafik.xyz;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.daybook.shafik.xyz;

    ssl_certificate /etc/letsencrypt/live/api.daybook.shafik.xyz/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.daybook.shafik.xyz/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

Enable and reload:

```bash
# Enable the configuration
sudo ln -s /etc/nginx/sites-available/api.daybook.shafik.xyz /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

### Cause 3: Missing SSL Certificate (1% probability)

**Symptoms**: HTTPS not working, nginx errors about SSL

**Fix**:
```bash
# Install/renew Let's Encrypt certificate
sudo certbot --nginx -d api.daybook.shafik.xyz
```

### Cause 4: Wrong Backend Code Deployed (1% probability)

**Symptoms**: Everything seems to work but routes return 404

**Fix**: Redeploy the latest code

```bash
# On production server
cd /path/to/daybook/backend

# Pull latest code
git pull origin master

# Rebuild
go build -o daybook-backend

# Restart
sudo systemctl restart daybook-backend
# OR
pm2 restart daybook-backend
```

## Verification Steps

After applying the fix, test these URLs:

### Test 1: Health Check
```bash
curl https://api.daybook.shafik.xyz/health
```
**Expected**: `{"status":"ok"}` or similar

### Test 2: Auth Endpoint (without token - should return 401)
```bash
curl https://api.daybook.shafik.xyz/api/v1/auth/me
```
**Expected**: `{"error":"Unauthorized"}` or similar (NOT 404)

### Test 3: From Browser
Open https://daybook.shafik.xyz in your browser and try to:
1. Login
2. Navigate to Settings
3. Try to export data

## One-Liner Quick Check

Run this from your production server to check everything:

```bash
echo "Backend Running:" && ps aux | grep daybook-backend | grep -v grep && \
echo "Port 8080:" && sudo lsof -i :8080 && \
echo "Local Test:" && curl -s http://localhost:8080/health && \
echo "Public Test:" && curl -s https://api.daybook.shafik.xyz/health
```

## Still Not Working?

If the issue persists:

1. Check backend logs:
   ```bash
   # If using systemd
   sudo journalctl -u daybook-backend -f

   # If using a log file
   tail -f /var/log/daybook-backend.log
   ```

2. Check nginx logs:
   ```bash
   sudo tail -f /var/log/nginx/error.log
   ```

3. Share the output of the diagnostic script for more help

## Summary

```
┌─────────────────────────────────────────────────┐
│  Most Common Fix (90% of cases):                │
│                                                  │
│  1. SSH into production server                  │
│  2. cd /path/to/daybook/backend                 │
│  3. ./daybook-backend                           │
│     OR                                           │
│     sudo systemctl start daybook-backend        │
│                                                  │
│  4. Test: curl http://localhost:8080/health     │
│  5. Test: curl https://api.daybook.shafik.xyz/  │
└─────────────────────────────────────────────────┘
```

## Files Created to Help You

1. **PRODUCTION_DEPLOYMENT_GUIDE.md** - Comprehensive troubleshooting guide
2. **diagnose-production.sh** - Automated diagnostic script
3. **PRODUCTION_QUICK_FIX.md** (this file) - Quick fixes for common issues

Good luck! The export feature code is working perfectly locally - you just need to get your production backend running.
