# Quick Deploy Commands

Your deployment setup:
- **Git Repository**: `/home/ubuntu/daybook` (on EC2)
- **Application Files**: `/var/www/daybook` (on EC2)
- **Backend Binary**: `/var/www/daybook/backend/daybook-backend`
- **Backend .env**: `/var/www/daybook/backend/.env`
- **Frontend Files**: `/var/www/daybook/frontend/`

---

## Complete Deployment Workflow

### Step 1: Build on Local Machine

```bash
cd /Users/shafikshaon/workplace/development/projects/daybook

# Auto-detect changes and build
./scripts/deploy-local.sh

# OR force build everything
./scripts/deploy-local.sh --force
```

This will:
1. Build backend → `build/daybook-backend`
2. Build frontend → `build/frontend-dist.tar.gz`
3. Ask to commit and push to git

### Step 2: Deploy on EC2

```bash
# SSH to EC2
ssh -i your-key.pem ubuntu@your-ec2-ip

# Navigate to repo
cd /home/ubuntu/daybook

# Run deployment
./scripts/deploy-ec2.sh

# OR force deploy everything
./scripts/deploy-ec2.sh --force
```

This will:
1. Pull from GitHub
2. Copy backend binary to `/var/www/daybook/backend/`
3. Extract frontend to `/var/www/daybook/frontend/`
4. Restart services

---

## Quick Fixes

### Fix Frontend API URL Issue Right Now

**On Local:**
```bash
cd /Users/shafikshaon/workplace/development/projects/daybook

# Ensure .env.production is correct
cat frontend/.env.production
# Should show: VITE_API_URL=https://api.daybook.shafik.xyz/api/v1

# Force rebuild frontend
./scripts/deploy-local.sh --force

# Commit and push when prompted
```

**On EC2:**
```bash
cd /home/ubuntu/daybook

# Deploy
./scripts/deploy-ec2.sh

# Verify frontend deployed
ls -la /var/www/daybook/frontend/index.html
```

**In Browser:**
- Hard refresh: `Ctrl+Shift+R` (Windows/Linux) or `Cmd+Shift+R` (Mac)
- Open DevTools (F12) → Network tab
- Check API calls go to `https://api.daybook.shafik.xyz/api/v1/...`

---

## Manual Deployment (Without Scripts)

### Backend Only

**Local:**
```bash
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../build/daybook-backend main.go
git add build/
git commit -m "Update backend"
git push origin master
```

**EC2:**
```bash
cd /home/ubuntu/daybook
git pull origin master
sudo systemctl stop daybook-backend
sudo cp build/daybook-backend /var/www/daybook/backend/daybook-backend
sudo chmod +x /var/www/daybook/backend/daybook-backend
sudo systemctl start daybook-backend
sudo systemctl status daybook-backend
```

### Frontend Only

**Local:**
```bash
cd frontend
npm run build
cd dist
tar -czf ../../build/frontend-dist.tar.gz .
cd ../..
git add build/
git commit -m "Update frontend"
git push origin master
```

**EC2:**
```bash
cd /home/ubuntu/daybook
git pull origin master
cd /var/www/daybook/frontend
sudo rm -rf *
sudo tar -xzf /home/ubuntu/daybook/build/frontend-dist.tar.gz
sudo systemctl reload nginx
```

---

## Service Management

### Backend Service

```bash
# Start
sudo systemctl start daybook-backend

# Stop
sudo systemctl stop daybook-backend

# Restart
sudo systemctl restart daybook-backend

# Status
sudo systemctl status daybook-backend

# Logs (real-time)
sudo journalctl -u daybook-backend -f

# Recent logs
sudo journalctl -u daybook-backend -n 50 --no-pager
```

### Nginx

```bash
# Reload (no downtime)
sudo systemctl reload nginx

# Restart
sudo systemctl restart nginx

# Status
sudo systemctl status nginx

# Test config
sudo nginx -t
```

---

## Configuration Files

### Backend .env

**Location:** `/var/www/daybook/backend/.env`

```bash
# Edit
nano /var/www/daybook/backend/.env

# View (hide passwords)
cat /var/www/daybook/backend/.env | grep -v PASSWORD

# After editing, restart backend
sudo systemctl restart daybook-backend
```

**Key variables:**
- `DB_PASSWORD` - PostgreSQL password
- `REDIS_PASSWORD` - Redis password
- `JWT_SECRET` - JWT secret key
- `JWT_EXPIRATION` - Hours (e.g., `24`)
- `CORS_ALLOWED_ORIGINS` - Frontend URL

### Frontend .env.production (Local)

**Location:** `/Users/shafikshaon/workplace/development/projects/daybook/frontend/.env.production`

```bash
# Edit
nano frontend/.env.production

# After editing, MUST rebuild
./scripts/deploy-local.sh --force
```

**Key variables:**
- `VITE_API_URL` - Backend API URL (e.g., `https://api.daybook.shafik.xyz/api/v1`)

---

## Logs

### Application Logs

```bash
# Backend application logs
sudo tail -f /var/log/daybook-backend.log
sudo tail -f /var/log/daybook-backend-error.log

# Nginx access logs
sudo tail -f /var/log/nginx/daybook-frontend-access.log
sudo tail -f /var/log/nginx/daybook-backend-access.log

# Nginx error logs
sudo tail -f /var/log/nginx/daybook-frontend-error.log
sudo tail -f /var/log/nginx/daybook-backend-error.log

# PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-15-main.log
```

### View All Recent Errors

```bash
# Backend errors (last 100 lines)
sudo tail -100 /var/log/daybook-backend-error.log

# Nginx errors (last 100 lines)
sudo tail -100 /var/log/nginx/daybook-backend-error.log

# Systemd journal (last 50 entries)
sudo journalctl -u daybook-backend -n 50 --no-pager
```

---

## Verification

### Check Everything is Running

```bash
# All services
sudo systemctl status daybook-backend postgresql redis-server nginx

# Port 8082 (backend)
sudo netstat -tulpn | grep 8082

# Test backend locally
curl http://localhost:8082

# Test backend through Nginx
curl https://api.daybook.shafik.xyz

# Test frontend
curl https://daybook.shafik.xyz
```

### Check Resource Usage

```bash
# Memory
free -h

# Disk
df -h

# Processes (install first: sudo apt install htop)
htop

# All ports
sudo netstat -tulpn
```

---

## Troubleshooting

### Backend Won't Start

```bash
# Stop the service
sudo systemctl stop daybook-backend

# Try running manually to see error
cd /var/www/daybook/backend
./daybook-backend

# Common issues:
# - Database password wrong → Edit .env, fix DB_PASSWORD
# - Redis password wrong → Edit .env, fix REDIS_PASSWORD
# - JWT_EXPIRATION wrong → Should be just number (24), not "24h"
# - Binary not executable → sudo chmod +x daybook-backend
```

### Frontend Shows Wrong API URL

```bash
# Frontend env vars are BAKED IN at build time
# You MUST rebuild on local machine:

# 1. Edit .env.production on local
nano frontend/.env.production

# 2. Force rebuild
./scripts/deploy-local.sh --force

# 3. Deploy to EC2
cd /home/ubuntu/daybook
./scripts/deploy-ec2.sh

# 4. Hard refresh browser (Ctrl+Shift+R)
```

### Database Connection Failed

```bash
# Test connection
psql -h localhost -U daybook_user -d daybook

# If password wrong, reset it
sudo -u postgres psql
ALTER USER daybook_user WITH PASSWORD 'your_password';
\q

# Update .env
nano /var/www/daybook/backend/.env
# Change DB_PASSWORD

# Restart backend
sudo systemctl restart daybook-backend
```

### Out of Memory

```bash
# Check memory
free -h

# Add 1GB swap
sudo fallocate -l 1G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

---

## File Locations Reference

| What | Where |
|------|-------|
| **Git Repo** | `/home/ubuntu/daybook` |
| **Deployed Backend** | `/var/www/daybook/backend/` |
| **Deployed Frontend** | `/var/www/daybook/frontend/` |
| **Backend Binary** | `/var/www/daybook/backend/daybook-backend` |
| **Backend .env** | `/var/www/daybook/backend/.env` |
| **Backend Logs** | `/var/log/daybook-backend*.log` |
| **Nginx Config (Frontend)** | `/etc/nginx/sites-available/daybook-frontend` |
| **Nginx Config (Backend)** | `/etc/nginx/sites-available/daybook-backend` |
| **Systemd Service** | `/etc/systemd/system/daybook-backend.service` |

---

## Commit to Git

Now that the scripts are updated, commit them:

```bash
# On local machine
cd /Users/shafikshaon/workplace/development/projects/daybook

git add scripts/deploy-ec2.sh scripts/DEPLOYMENT_README.md DEPLOY_COMMANDS.md
git commit -m "Update deployment scripts with correct EC2 paths"
git push origin master
```

Then on EC2:

```bash
cd /home/ubuntu/daybook
git pull origin master
```

Now your deployment scripts will work perfectly with your setup! 🚀
