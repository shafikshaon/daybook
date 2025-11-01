# Build Artifacts Directory

This directory contains production build artifacts for EC2 deployment.

## Contents

After running `./scripts/local-build.sh`, this directory will contain:

- **daybook-backend** - Go binary compiled for Linux (x86_64)
- **frontend-dist.tar.gz** - Gzipped tarball of Vue.js production build
- **backend.env.example** - Backend environment template
- **frontend.env.example** - Frontend environment template
- **DEPLOYMENT_INFO.txt** - Build information and deployment instructions

## Git Workflow

### Why are build artifacts committed to git?

This project uses a git-based deployment workflow instead of scp/rsync because:

1. **Simplicity** - Just `git pull` on EC2, no file transfers needed
2. **Version Control** - Build artifacts are versioned alongside source code
3. **Consistency** - Same deployment process for all environments
4. **Low RAM EC2** - Building on EC2 (0.5GB RAM) causes out-of-memory errors
5. **Atomic Deployments** - All files deploy together in a single git commit

### Deployment Workflow

**Local Machine:**
```bash
# Build the application
./scripts/local-build.sh

# Commit build artifacts
git add build/
git commit -m "Add deployment build artifacts"
git push origin master
```

**EC2 Instance:**
```bash
# Pull latest changes
cd ~/daybook
git pull origin master

# Deploy using script
./scripts/ec2-setup.sh
```

## Build Details

### Backend Binary
- **Platform**: Linux
- **Architecture**: x86_64 (amd64)
- **CGO**: Disabled (static binary)
- **Optimization**: Strip debug symbols (-ldflags="-s -w")
- **Size**: ~26MB

### Frontend Tarball
- **Framework**: Vue.js 3
- **Build Tool**: Vite
- **Contents**: Static files (HTML, JS, CSS, assets)
- **Compression**: gzip
- **Size**: ~500KB - 1MB

## Important Notes

1. **Do NOT ignore this directory** - The root `.gitignore` explicitly allows this directory
2. **Binary files in git** - Yes, this is intentional for this deployment workflow
3. **Rebuild before deploy** - Always run `./scripts/local-build.sh` before committing
4. **Platform-specific** - Backend binary is Linux-only (won't run on macOS/Windows)

## Alternative Workflows

If you prefer not to commit build artifacts to git, you can:

1. Use GitHub Actions / CI/CD to build and store artifacts
2. Use AWS S3 to store build artifacts
3. Build directly on EC2 (requires more RAM, consider upgrading instance)
4. Use a separate deployment branch for build artifacts

However, the current git-based workflow is optimized for the project's constraints (low RAM EC2).

## Security

- **Secrets**: Never commit `.env` files with secrets to this directory
- **Environment Templates**: Only `.env.example` files should be committed
- **Binary Verification**: Always verify the binary source before deploying

## Troubleshooting

### Build artifacts not found on EC2

```bash
# On EC2, check if build directory exists
ls -la ~/daybook/build/

# If empty, ensure you've pulled the latest changes
git pull origin master

# Verify build artifacts are in git history
git log --stat | grep build/
```

### Binary won't execute

```bash
# Check permissions
ls -l ~/daybook/build/daybook-backend

# Make executable
chmod +x ~/daybook/build/daybook-backend

# Verify it's a Linux binary
file ~/daybook/build/daybook-backend
# Should output: ELF 64-bit LSB executable, x86-64
```

## Related Documentation

- **EC2_DEPLOYMENT_GUIDE.md** - Complete deployment guide
- **DEPLOYMENT_QUICKSTART.md** - Quick start guide
- **DEPLOYMENT_CHECKLIST.md** - Deployment checklist
- **scripts/local-build.sh** - Build automation script
- **scripts/ec2-setup.sh** - EC2 deployment script
