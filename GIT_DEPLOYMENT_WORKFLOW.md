# Git-Based Deployment Workflow

This document explains the git-based deployment workflow for the Daybook application on EC2.

## Overview

Unlike traditional deployment workflows that use `scp`, `rsync`, or CI/CD pipelines, this project uses **git as the deployment mechanism**. Build artifacts are committed to the repository and deployed via `git pull`.

## Why Git-Based Deployment?

### Advantages

1. **Simple** - No need for SSH file transfers or deployment tools
2. **Consistent** - Same process across all environments
3. **Version Controlled** - Build artifacts versioned with source code
4. **Atomic** - All files deploy together in a single commit
5. **Efficient** - Git's delta compression minimizes transfer size
6. **Low RAM Friendly** - Build locally, deploy pre-built binary (EC2 has only 0.5GB RAM)

### Trade-offs

- **Repository Size** - Binary files increase repository size
- **Not Typical** - Most projects don't commit build artifacts
- **Rebuild Required** - Must rebuild before every deployment

## Workflow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    LOCAL MACHINE                             │
├─────────────────────────────────────────────────────────────┤
│  1. Make code changes                                        │
│  2. ./scripts/local-build.sh                                 │
│     ├─ Build Go backend (Linux binary)                       │
│     └─ Build Vue.js frontend (static files)                  │
│  3. git add build/                                           │
│  4. git commit -m "Update deployment build"                  │
│  5. git push origin master                                   │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ GitHub
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    EC2 INSTANCE                              │
├─────────────────────────────────────────────────────────────┤
│  1. cd ~/daybook                                             │
│  2. git pull origin master                                   │
│  3. ./scripts/ec2-setup.sh (first time)                      │
│     OR                                                        │
│     Manual deployment (updates)                              │
│     ├─ Copy binary from build/                               │
│     ├─ Extract frontend tarball                              │
│     └─ Restart services                                      │
└─────────────────────────────────────────────────────────────┘
```

## Step-by-Step Guide

### Initial Deployment

**Step 1: Build Locally**

```bash
cd /path/to/daybook

# Build backend and frontend
./scripts/local-build.sh
```

This creates:
- `build/daybook-backend` - Go binary for Linux
- `build/frontend-dist.tar.gz` - Frontend tarball
- `build/*.env.example` - Environment templates
- `build/DEPLOYMENT_INFO.txt` - Build information

**Step 2: Commit Build Artifacts**

```bash
# Add build directory
git add build/

# Commit
git commit -m "Add deployment build artifacts"

# Push to GitHub
git push origin master
```

**Step 3: Deploy on EC2**

```bash
# SSH to EC2
ssh -i your-key.pem ubuntu@your-ec2-ip

# Clone repository (first time only)
cd ~
git clone https://github.com/shafikshaon/daybook.git
cd daybook

# Run automated setup
./scripts/ec2-setup.sh
```

The setup script will:
- Create directory structure
- Deploy backend binary
- Deploy frontend files
- Configure environment
- Create systemd service
- Configure Nginx

**Step 4: Start Services**

```bash
sudo systemctl daemon-reload
sudo systemctl enable daybook-backend
sudo systemctl start daybook-backend
sudo systemctl reload nginx
```

### Updating Deployment

When you make changes and need to redeploy:

**On Local Machine:**

```bash
# Make your code changes
# ...

# Rebuild
./scripts/local-build.sh

# Commit source code changes
git add backend/ frontend/
git commit -m "Update feature X"

# Commit build artifacts
git add build/
git commit -m "Update deployment build"

# Push everything
git push origin master
```

**On EC2:**

```bash
# Pull latest changes
cd ~/daybook
git pull origin master

# Stop backend
sudo systemctl stop daybook-backend

# Deploy backend
cp ~/daybook/build/daybook-backend /var/www/daybook/backend/daybook-backend
chmod +x /var/www/daybook/backend/daybook-backend

# Deploy frontend
cd /var/www/daybook/frontend
rm -rf *
tar -xzf ~/daybook/build/frontend-dist.tar.gz

# Restart services
sudo systemctl start daybook-backend
sudo systemctl reload nginx

# Verify
sudo systemctl status daybook-backend
```

## Directory Structure

```
daybook/
├── backend/               # Go source code
├── frontend/              # Vue.js source code
├── build/                 # Build artifacts (committed to git)
│   ├── daybook-backend   # Linux binary (~26MB)
│   ├── frontend-dist.tar.gz  # Frontend build (~1MB)
│   ├── README.md
│   └── .gitkeep
├── scripts/
│   ├── local-build.sh    # Build script
│   └── ec2-setup.sh      # Deployment script
├── .gitignore            # Does NOT ignore build/
├── EC2_DEPLOYMENT_GUIDE.md
├── DEPLOYMENT_QUICKSTART.md
└── DEPLOYMENT_CHECKLIST.md
```

## Git Configuration

### Root .gitignore

The root `.gitignore` explicitly **does NOT** ignore the `build/` directory:

```gitignore
# NOTE: Do NOT ignore /build/ directory!
# The build/ directory contains deployment artifacts that must be committed to git
# for EC2 deployment. The deployment workflow requires build artifacts in the repository.
```

### Backend .gitignore

Ignores local build artifacts but not the root `build/` directory:

```gitignore
/bin/
/build/      # This ignores backend/build/, not the root build/
/dist/
daybook-backend
```

### Frontend .gitignore

Ignores the local `dist/` directory:

```gitignore
dist         # Frontend build output (packaged in tarball instead)
node_modules
```

## Build Script Details

The `local-build.sh` script:

1. **Builds Backend** - Compiles Go code for Linux with optimizations
   ```bash
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/daybook-backend main.go
   ```

2. **Builds Frontend** - Creates production Vue.js build and tarballs it
   ```bash
   npm run build
   cd dist
   tar -czf ../../build/frontend-dist.tar.gz .
   ```

3. **Copies Templates** - Includes .env.example files

4. **Creates Info File** - Generates deployment information

## EC2 Setup Script Details

The `ec2-setup.sh` script:

1. **Validates** - Checks that build artifacts exist in `~/daybook/build/`
2. **Creates Directories** - Sets up `/var/www/daybook` structure
3. **Deploys Backend** - Copies binary from `build/` to deployment location
4. **Deploys Frontend** - Extracts tarball to deployment location
5. **Configures** - Creates `.env` file with user input
6. **Systemd** - Creates backend service
7. **Nginx** - Configures reverse proxy

## Security Considerations

### What's Committed to Git

✅ **Safe to commit:**
- Linux backend binary (compiled code)
- Frontend static files (HTML, JS, CSS)
- `.env.example` files (no secrets)
- Build metadata

❌ **Never commit:**
- `.env` files with actual credentials
- Database passwords
- API keys
- JWT secrets

### Binary Safety

The Linux binary is:
- Built from source on your local machine
- Not executable on macOS/Windows (platform-specific)
- Verifiable via `file` command on EC2
- Stripped of debug symbols (smaller, less info)

## Troubleshooting

### Build artifacts not in git

```bash
# Check if build/ is ignored
git check-ignore build/

# Should output nothing (not ignored)

# Check git status
git status

# build/ should appear in untracked or modified files
```

### Build artifacts not on EC2

```bash
# Verify you've pulled latest
cd ~/daybook
git pull origin master

# Check if build directory exists
ls -la ~/daybook/build/

# If empty, check git
git log --oneline | head -5
# Should see your build commit

# Check if artifacts are in repo
git ls-files build/
```

### Binary won't execute on EC2

```bash
# Check file type
file ~/daybook/build/daybook-backend
# Should be: ELF 64-bit LSB executable, x86-64

# Make executable
chmod +x ~/daybook/build/daybook-backend

# Test run
~/daybook/build/daybook-backend --help
```

## Alternative Workflows

If you don't want to commit build artifacts to git:

### Option 1: GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Build and Deploy
on:
  push:
    branches: [master]
jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build backend
        run: cd backend && go build
      - name: Deploy to EC2
        run: scp binary user@ec2:/path/
```

### Option 2: AWS S3

```bash
# Local machine
./scripts/local-build.sh
aws s3 cp build/ s3://my-bucket/builds/latest/ --recursive

# EC2
aws s3 cp s3://my-bucket/builds/latest/ ~/build/ --recursive
```

### Option 3: Separate Branch

```bash
# Local machine
./scripts/local-build.sh
git checkout -b deployment
git add -f build/
git commit -m "Deployment build"
git push origin deployment

# EC2
git pull origin deployment
```

### Option 4: Build on EC2 (Requires More RAM)

```bash
# Upgrade to t3.small (2GB RAM) or add swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Build on EC2
cd ~/daybook/backend
go build -o daybook-backend main.go

cd ~/daybook/frontend
npm install
npm run build
```

## Performance Metrics

### Repository Size Impact

- **Source Code**: ~2MB
- **Backend Binary**: ~26MB
- **Frontend Tarball**: ~1MB
- **Total Build Artifacts**: ~27MB per deployment

With git delta compression, incremental commits are much smaller.

### Deployment Speed

- **git pull**: ~5-10 seconds (depending on changes)
- **Extract + Copy**: ~2-3 seconds
- **Service Restart**: ~1-2 seconds
- **Total**: ~10-15 seconds for update deployment

## Best Practices

1. **Rebuild Before Deploy** - Always run `./scripts/local-build.sh` before committing
2. **Commit Message** - Use clear messages like "Update deployment build for feature X"
3. **Test Locally** - Verify backend and frontend work before building
4. **Verify Pull** - Always check `git pull` output on EC2
5. **Backup First** - Create database backup before deploying major changes
6. **Monitor Logs** - Watch logs after deployment to catch issues early

## Related Documentation

- **EC2_DEPLOYMENT_GUIDE.md** - Complete deployment guide (14 parts)
- **DEPLOYMENT_QUICKSTART.md** - Fast-track deployment (5 minutes)
- **DEPLOYMENT_CHECKLIST.md** - Deployment checklist with checkboxes
- **build/README.md** - Build directory documentation
- **scripts/local-build.sh** - Build automation script
- **scripts/ec2-setup.sh** - EC2 deployment automation

## Summary

This git-based deployment workflow is:

✅ Simple and consistent
✅ Optimized for low-RAM EC2 instances
✅ Version controlled and auditable
✅ No external dependencies (no scp, no CI/CD required)
✅ Atomic deployments

While committing binaries to git is unconventional, it's a pragmatic solution for this project's constraints and provides a straightforward deployment process.
