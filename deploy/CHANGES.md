# Deployment System - Changes Summary

## Overview

The deployment system has been completely refactored to fix build issues and ensure your git repository stays clean.

## Problems Fixed

### 1. **Git Repository Contamination**
- **Before**: Scripts built directly in `/opt/daybook`, creating build artifacts tracked by git
- **After**: Scripts copy code to `/opt/daybook-build`, build there, then deploy to `/opt/daybook-app`
- **Result**: Your git repository at `/opt/daybook` remains pristine

### 2. **Path Mismatches**
- **Before**: Config said `/home/ubuntu/projects/daybook` but actual location was `/opt/daybook`
- **After**: Config updated to use correct paths (`/opt/daybook`)

### 3. **Missing Package Lock**
- **Before**: `npm ci` failed when `package-lock.json` wasn't present on server
- **After**: Script uses `npm install` which handles missing lockfile gracefully

### 4. **Missing Migrations Script**
- **Before**: `deploy.sh` referenced `08_run_migrations.sh` which didn't exist
- **After**: Created `08_run_migrations.sh` (migrations run automatically with backend)

## Files Changed

### 1. `config/deploy.conf`
**Changes:**
- Added `SOURCE_DIR="/opt/daybook"` - Your git repository (read-only)
- Added `BUILD_DIR="/opt/daybook-build"` - Temporary build directory
- Changed `APP_DIR="/opt/daybook-app"` - Final deployment location
- Added `BUILD_BACKEND_DIR` and `BUILD_FRONTEND_DIR` variables

```diff
# Before
- APP_DIR="/home/ubuntu/projects/daybook"
- BACKEND_DIR="${APP_DIR}/backend"
- FRONTEND_DIR="${APP_DIR}/frontend"

# After
+ SOURCE_DIR="/opt/daybook"                  # Git repo (read-only)
+ BUILD_DIR="/opt/daybook-build"             # Temporary build
+ BUILD_BACKEND_DIR="${BUILD_DIR}/backend"
+ BUILD_FRONTEND_DIR="${BUILD_DIR}/frontend"
+ APP_DIR="/opt/daybook-app"                 # Deployment
+ BACKEND_DIR="${APP_DIR}/backend"
+ FRONTEND_DIR="${APP_DIR}/frontend"
```

### 2. `scripts/06_deploy_backend.sh`
**Changes:**
- Now copies source from `SOURCE_DIR` to `BUILD_BACKEND_DIR`
- Builds in `BUILD_BACKEND_DIR`
- Deploys binary to `BACKEND_DIR`
- Cleans up build directory after deployment
- Never modifies source code

**Process:**
```
/opt/daybook/backend → [copy] → /opt/daybook-build/backend → [build] → /opt/daybook-app/backend
```

### 3. `scripts/07_deploy_frontend.sh`
**Changes:**
- Now copies source from `SOURCE_DIR` to `BUILD_FRONTEND_DIR`
- Builds in `BUILD_FRONTEND_DIR`
- Uses `npm install` instead of `npm ci` (handles missing package-lock.json)
- Deploys `dist` folder to `FRONTEND_DIR`
- Cleans up build directory after deployment
- Never modifies source code

**Process:**
```
/opt/daybook/frontend → [copy] → /opt/daybook-build/frontend → [build] → /opt/daybook-app/frontend
```

### 4. `scripts/08_run_migrations.sh` (NEW)
**Created:**
- Verifies database connection
- Documents that migrations run automatically via GORM
- Shows current database tables
- Provides helpful logging for troubleshooting

### 5. `deploy.sh`
**Changes:**
- Added verification that `SOURCE_DIR` exists
- Updated configuration display to show:
  - Source Directory (git repository)
  - Build Directory (temporary)
  - App Directory (deployment)

## New Files Created

1. **`DEPLOYMENT_GUIDE.md`** - Comprehensive deployment documentation
2. **`QUICK_START.md`** - Quick reference for common tasks
3. **`CHANGES.md`** - This file, documenting all changes
4. **`scripts/08_run_migrations.sh`** - Database migration script

## Directory Structure

### Before
```
/opt/daybook/
├── backend/
│   ├── main.go
│   ├── daybook-backend  ← Build artifacts in git repo!
│   └── .env
└── frontend/
    ├── src/
    ├── dist/            ← Build artifacts in git repo!
    └── node_modules/    ← Dependencies in git repo!
```

### After
```
/opt/daybook/              ← Git repository (CLEAN)
├── backend/
│   └── main.go
└── frontend/
    └── src/

/opt/daybook-build/        ← Temporary (deleted after build)
├── backend/
└── frontend/

/opt/daybook-app/          ← Production deployment
├── backend/
│   ├── daybook-backend   ← Binary
│   ├── .env
│   └── uploads/
└── frontend/
    └── dist/             ← Static files
```

## Benefits

### 1. **Clean Git Repository**
- No build artifacts in source
- No `node_modules` clutter
- No compiled binaries
- Easy to track actual code changes

### 2. **Isolated Builds**
- Build in separate directory
- No interference with source
- Can delete build dir safely
- Easy to debug build issues

### 3. **Safe Deployments**
- Source code never modified
- Can always rebuild from git
- Easy rollbacks (just checkout and redeploy)
- No risk of losing source changes

### 4. **Flexible Workflow**
```bash
# On server
cd /opt/daybook
git pull origin master    # Get new code
cd deploy
./deploy.sh              # Deploy - source stays clean!
```

## Testing the Changes

### 1. Verify Source is Clean
```bash
cd /opt/daybook
git status
# Should show: "working tree clean"
```

### 2. Deploy
```bash
cd /opt/daybook/deploy
./deploy.sh
```

### 3. Verify Deployment
```bash
# Check backend
sudo systemctl status daybook-backend
curl http://localhost:8080/health

# Check frontend
ls -la /opt/daybook-app/frontend/dist

# Check source is still clean
cd /opt/daybook
git status  # Still clean!
```

## Backward Compatibility

⚠️ **Important**: This is a breaking change in directory structure.

**Migration Path:**
1. Pull these changes to your server
2. Ensure code is at `/opt/daybook`
3. Run `./deploy.sh --fresh`
4. The new structure will be created automatically

**Old deployments** at `/home/ubuntu/projects/daybook` will need to be migrated or you can update `SOURCE_DIR` in `deploy.conf` to point to your location.

## Configuration Variables

| Variable | Path | Purpose |
|----------|------|---------|
| `SOURCE_DIR` | `/opt/daybook` | Git repository (read-only) |
| `BUILD_DIR` | `/opt/daybook-build` | Temporary build location |
| `APP_DIR` | `/opt/daybook-app` | Production deployment |
| `BACKEND_DIR` | `${APP_DIR}/backend` | Backend deployment |
| `FRONTEND_DIR` | `${APP_DIR}/frontend` | Frontend deployment |

## Upgrade Instructions

If you have an existing deployment:

```bash
# 1. SSH to server
ssh ubuntu@your-server

# 2. Pull latest changes
cd /opt/daybook  # or wherever your repo is
git pull

# 3. Update config if needed
nano deploy/config/deploy.conf
# Make sure SOURCE_DIR points to your git repo location

# 4. Run fresh deployment
cd deploy
./deploy.sh --fresh
```

## Rollback (if needed)

If something goes wrong:

```bash
# 1. Check out previous commit
cd /opt/daybook
git log --oneline
git checkout <previous-commit-hash>

# 2. Redeploy
cd deploy
./deploy.sh
```

## Support

- Read `DEPLOYMENT_GUIDE.md` for detailed documentation
- Read `QUICK_START.md` for quick reference
- Check logs: `sudo journalctl -u daybook-backend -f`
- Check Nginx: `sudo nginx -t`

## Summary

✅ **Git repository stays clean** - No build artifacts
✅ **Separate build directory** - Isolated compilation
✅ **Fixed all path issues** - Correct directories
✅ **Fixed package-lock issue** - Graceful handling
✅ **Created missing scripts** - Complete deployment system
✅ **Comprehensive documentation** - Easy to use
✅ **Safe and reliable** - No risk to source code
