# Quick Start - Daybook Deployment

## What Changed?

✅ **Fixed Issues:**
- No more modifications to your git-tracked code
- Separate build directory for compilation
- Fixed package-lock.json issues
- Fixed path mismatches (now uses `/opt/daybook`)
- Created missing migrations script

## On Your Server

### First Time Setup

```bash
# 1. Ensure repository is at /opt/daybook
cd /opt
ls -la daybook  # Should show your repository

# 2. If not, clone it
sudo git clone <your-repo-url> /opt/daybook
sudo chown -R ubuntu:ubuntu /opt/daybook

# 3. Run deployment
cd /opt/daybook/deploy
./deploy.sh --fresh
```

### Regular Updates

When you have new code:

```bash
# 1. Pull latest code
cd /opt/daybook
git pull origin master

# 2. Deploy
cd deploy
./deploy.sh
```

## How It Works

```
┌─────────────────┐
│  /opt/daybook   │ ← Your git repo (NEVER modified)
│  (read-only)    │
└────────┬────────┘
         │ Copy
         ↓
┌──────────────────────┐
│  /opt/daybook-build  │ ← Temporary build directory
│  (deleted after)     │
└──────────┬───────────┘
           │ Deploy
           ↓
┌─────────────────────┐
│  /opt/daybook-app   │ ← Running application
│  (production)       │
└─────────────────────┘
```

## Key Commands

```bash
# Full deployment
cd /opt/daybook/deploy
./deploy.sh

# Update backend only
./scripts/update_backend.sh

# Update frontend only
./scripts/update_frontend.sh

# Check backend status
sudo systemctl status daybook-backend

# Check logs
sudo journalctl -u daybook-backend -f

# Restart services
sudo systemctl restart daybook-backend
sudo systemctl reload nginx
```

## Deployment Options

```bash
# Skip dependency installation (faster)
./deploy.sh --skip-deps

# Skip database setup (when DB is already configured)
./deploy.sh --skip-db

# Deploy only backend
./deploy.sh --skip-frontend

# Deploy only frontend
./deploy.sh --skip-backend
```

## Configuration

Edit `deploy/config/deploy.conf`:

```bash
# Important settings
SOURCE_DIR="/opt/daybook"           # Your git repo
BUILD_DIR="/opt/daybook-build"      # Temporary build dir
APP_DIR="/opt/daybook-app"          # Deployment location
DOMAIN_OR_IP="your-domain.com"      # Your domain/IP
```

## Verification

After deployment, check:

```bash
# 1. Backend health
curl http://localhost:8080/health

# 2. Frontend files
ls -la /opt/daybook-app/frontend/dist

# 3. Services running
sudo systemctl status daybook-backend
sudo systemctl status nginx

# 4. Git repo is clean
cd /opt/daybook
git status  # Should show: "working tree clean"
```

## Troubleshooting

### Deployment fails?

```bash
# Check source directory
ls -la /opt/daybook

# Check configuration
cat /opt/daybook/deploy/config/deploy.conf

# View full logs during deployment
./deploy.sh 2>&1 | tee deployment.log
```

### Service not starting?

```bash
# Backend logs
sudo journalctl -u daybook-backend -n 100

# Nginx logs
sudo tail -f /var/log/nginx/daybook-error.log
```

### Need to rebuild?

```bash
# Just run deploy again - it rebuilds from scratch
cd /opt/daybook/deploy
./deploy.sh
```

## Benefits

✅ **Git stays clean** - No build artifacts in source code
✅ **Easy updates** - Just `git pull` and `./deploy.sh`
✅ **Fast rebuilds** - Build in isolated directory
✅ **Safe deployments** - Original code never modified
✅ **Simple rollback** - Checkout different git commit and redeploy

## Next Steps

1. **Run your first deployment**
   ```bash
   cd /opt/daybook/deploy
   ./deploy.sh --fresh
   ```

2. **Test the application**
   ```bash
   curl http://your-domain-or-ip
   ```

3. **Check logs**
   ```bash
   sudo journalctl -u daybook-backend -f
   ```

4. **Setup SSL** (optional but recommended)
   ```bash
   sudo certbot --nginx -d your-domain.com
   ```

## Support

For detailed information, see `DEPLOYMENT_GUIDE.md`
