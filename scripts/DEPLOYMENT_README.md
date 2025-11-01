# Automated Deployment Scripts

These scripts automate the build and deployment process for Daybook.

## Available Scripts

### 1. `deploy-local.sh` (Run on Local Machine)

**Purpose**: Detects changes in backend/frontend, builds them, and pushes to git.

**Features**:
- ✅ Auto-detects changes in backend/frontend
- ✅ Builds only what changed
- ✅ Prompts for confirmation before building
- ✅ Commits and pushes to git
- ✅ Shows deployment info

**Usage**:

```bash
# Auto-detect and build changes
./scripts/deploy-local.sh

# Force build both backend and frontend
./scripts/deploy-local.sh --force
./scripts/deploy-local.sh -f
```

**Workflow**:
1. Detects git changes in `backend/` and `frontend/` directories
2. Prompts you to build what changed
3. Builds backend (Go binary for Linux)
4. Builds frontend (Vue.js static files)
5. Creates tarball with proper formatting
6. Asks if you want to commit and push
7. Commits with automatic/custom message
8. Pushes to GitHub

**Example Output**:
```
==========================================
Daybook Auto-Deploy (Local)
==========================================

✓ Backend changes detected
✓ Frontend changes detected

Build Plan:
  ✓ Backend will be rebuilt
  ✓ Frontend will be rebuilt

[1/2] Building Backend...
✓ Backend build successful

[2/2] Building Frontend...
✓ Frontend build successful

Do you want to commit and push to git? (y/n)
```

---

### 2. `deploy-ec2.sh` (Run on EC2)

**Purpose**: Pulls latest changes from git and deploys to EC2.

**Features**:
- ✅ Pulls latest code from GitHub
- ✅ Auto-detects what changed
- ✅ Deploys only what changed
- ✅ Restarts services automatically
- ✅ Verifies deployment

**Usage**:

```bash
# Auto-deploy what changed
./scripts/deploy-ec2.sh

# Force deploy everything
./scripts/deploy-ec2.sh --force
./scripts/deploy-ec2.sh -f
```

**Workflow**:
1. Pulls latest changes from GitHub
2. Detects if backend or frontend changed
3. Stops backend service (if deploying backend)
4. Copies new backend binary
5. Extracts new frontend files
6. Restarts backend service
7. Reloads Nginx
8. Verifies everything is running

**Example Output**:
```
==========================================
Daybook Auto-Deploy (EC2)
==========================================

[1/5] Pulling latest changes from git...
✓ Successfully pulled latest changes

Deployment Plan:
  ✓ Backend will be deployed
  ✓ Frontend will be deployed

[2/5] Deploying Backend...
✓ Backend deployed

[3/5] Deploying Frontend...
✓ Frontend deployed

[4/5] Restarting Services...
✓ Backend service started
✓ Services restarted

[5/5] Verifying Deployment...
Backend service: ✓ Running
  Port 8082: ✓ Listening
Frontend files: ✓ Deployed
Nginx: ✓ Running

==========================================
Deployment Complete!
==========================================

Your application is available at:
  Frontend: https://daybook.shafik.xyz
  Backend:  https://api.daybook.shafik.xyz
```

---

## Complete Deployment Workflow

### Scenario: You changed frontend code

**On Local Machine:**
```bash
# 1. Make your changes
nano frontend/src/views/Dashboard.vue

# 2. Test locally
cd frontend
npm run dev

# 3. Deploy
cd ..
./scripts/deploy-local.sh

# Output:
# ✓ Frontend changes detected
# [2/2] Building Frontend...
# ✓ Frontend build successful
# Do you want to commit and push to git? (y)
# ✓ Successfully pushed to remote
```

**On EC2:**
```bash
# 1. Deploy
./scripts/deploy-ec2.sh

# Output:
# [1/5] Pulling latest changes from git...
# [3/5] Deploying Frontend...
# ✓ Frontend deployed
# Deployment Complete!
```

**Done!** Frontend updated in ~30 seconds.

---

### Scenario: You changed backend code

**On Local Machine:**
```bash
# 1. Make your changes
nano backend/handlers/auth_handler.go

# 2. Test locally
cd backend
go run main.go

# 3. Deploy
cd ..
./scripts/deploy-local.sh

# Output:
# ✓ Backend changes detected
# [1/2] Building Backend...
# ✓ Backend build successful
```

**On EC2:**
```bash
./scripts/deploy-ec2.sh

# Output:
# [2/5] Deploying Backend...
# Stopping backend service...
# ✓ Backend deployed
# ✓ Backend service started
```

**Done!** Backend updated and restarted.

---

### Scenario: You changed both

**On Local Machine:**
```bash
./scripts/deploy-local.sh

# Builds both backend and frontend
```

**On EC2:**
```bash
./scripts/deploy-ec2.sh

# Deploys both
```

---

## Force Deployment

Sometimes you want to rebuild/redeploy even if there are no detected changes (e.g., environment variable changes).

**Local:**
```bash
./scripts/deploy-local.sh --force
```

**EC2:**
```bash
./scripts/deploy-ec2.sh --force
```

---

## Safety Features

### Local Script:
- ✅ Checks for uncommitted changes and warns you
- ✅ Prompts before building
- ✅ Prompts before committing
- ✅ Prompts for commit message
- ✅ Verifies build success before committing
- ✅ Excludes macOS metadata from tarballs

### EC2 Script:
- ✅ Checks repository exists
- ✅ Shows current and new commit info
- ✅ Verifies files exist before deploying
- ✅ Gracefully stops backend before updating
- ✅ Verifies backend started after deployment
- ✅ Shows useful troubleshooting commands if something fails

---

## Troubleshooting

### "No changes detected"

If you know you made changes but the script doesn't detect them:

1. **Check git status:**
   ```bash
   git status
   ```

2. **Use force mode:**
   ```bash
   ./scripts/deploy-local.sh --force
   ```

### Backend fails to start on EC2

```bash
# Check logs
sudo journalctl -u daybook-backend -n 50

# Try manual start
cd /var/www/daybook/backend
./daybook-backend
```

### Frontend not updating

1. **Clear browser cache** (Ctrl+Shift+R)
2. **Check if files deployed:**
   ```bash
   ls -la /var/www/daybook/frontend/index.html
   ```
3. **Check Nginx logs:**
   ```bash
   sudo tail -f /var/log/nginx/daybook-frontend-error.log
   ```

---

## Environment Variables

### Frontend (.env.production)

Created automatically by `deploy-local.sh` if missing:

```env
VITE_API_URL=https://api.daybook.shafik.xyz/api/v1
VITE_APP_NAME=Daybook
VITE_APP_VERSION=1.0.0
```

**Important**: Frontend env vars are baked in at build time. To change them:
1. Edit `frontend/.env.production`
2. Run `./scripts/deploy-local.sh --force`
3. Deploy to EC2

### Backend (.env on EC2)

Located at `/var/www/daybook/backend/.env`

To update:
1. SSH to EC2
2. Edit: `nano /var/www/daybook/backend/.env`
3. Restart: `sudo systemctl restart daybook-backend`

---

## Quick Reference

### Local Machine Commands

```bash
# Build and deploy changes
./scripts/deploy-local.sh

# Force rebuild everything
./scripts/deploy-local.sh --force

# Check what changed
git status
git diff
```

### EC2 Commands

```bash
# Deploy latest changes
./scripts/deploy-ec2.sh

# Force redeploy everything
./scripts/deploy-ec2.sh --force

# Check service status
sudo systemctl status daybook-backend

# View logs
sudo journalctl -u daybook-backend -f

# Restart backend
sudo systemctl restart daybook-backend

# Reload Nginx
sudo systemctl reload nginx
```

---

## Related Documentation

- **EC2_DEPLOYMENT_GUIDE.md** - Complete deployment guide
- **DEPLOYMENT_QUICKSTART.md** - Quick start guide
- **TROUBLESHOOTING.md** - Troubleshooting guide
- **GIT_DEPLOYMENT_WORKFLOW.md** - Git workflow explanation

---

## Tips

1. **Always test locally first** before deploying to EC2
2. **Use meaningful commit messages** when deploying
3. **Check backend logs** after deploying backend changes
4. **Clear browser cache** after deploying frontend changes
5. **Monitor resource usage** on EC2 (run `htop` or `free -h`)
6. **Create database backups** before deploying major changes

---

**Happy Deploying!** 🚀
