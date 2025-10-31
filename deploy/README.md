# Daybook Deployment

## Two Deployment Options

### Option 1: Pre-Built Binary (RECOMMENDED - Much Faster!)

Build locally, push binary to git, deploy in seconds:

**On your local machine:**
```bash
# Build the binary for Linux
./build-for-deployment.sh

# Commit and push
git add backend/daybook-backend
git commit -m "Add pre-built binary"
git push
```

**On your server:**
```bash
# 1. Clone the repository
cd /opt
sudo git clone https://github.com/shafikshaon/daybook.git
sudo chown -R $USER:$USER daybook

# 2. Run deployment (uses pre-built binary - takes ~2 minutes)
cd daybook/deploy
./deploy.sh

# That's it! 🎉
```

### Option 2: Build on Server (Slower)

If you don't have the pre-built binary, it will build on the server (takes ~10 minutes):

```bash
# Clone and deploy (will build from source)
cd /opt
sudo git clone https://github.com/shafikshaon/daybook.git
sudo chown -R $USER:$USER daybook
cd daybook/deploy
./deploy.sh
```

## What It Does

The script automatically:
- ✅ Installs all dependencies (Go, Node.js, PostgreSQL, Redis, Nginx)
- ✅ Sets up database with secure password
- ✅ Builds backend from source
- ✅ Builds frontend from source
- ✅ Configures and starts all services
- ✅ Sets up Nginx reverse proxy

## Configuration

Edit these values in `deploy.sh` before running:

```bash
DOMAIN="daybook.shafik.xyz"    # Your domain or server IP
```

Everything else is automatic!

## After Deployment

Your app will be available at:
- **Frontend**: http://your-domain.com
- **Backend API**: http://your-domain.com/api/v1
- **Health Check**: http://your-domain.com/health

## Important: .env Files

**The `.env` file is created on the server, not in your code!**

- ✅ Local `.env` - for your local development
- ✅ Server `.env` - auto-created during deployment with production settings
- ✅ The binary reads `.env` at runtime, not at compile time
- ✅ So you can build locally and it works on the server with different settings!

**Never commit `.env` to git!** It's in `.gitignore` for security.

## Updating Your App

**With pre-built binary:**
```bash
# On local machine
./build-for-deployment.sh
git add backend/daybook-backend
git commit -m "Update binary"
git push

# On server
cd /opt/daybook
git pull
cd deploy
./deploy.sh  # Fast! Uses pre-built binary
```

**Without pre-built binary:**
```bash
cd /opt/daybook
git pull
cd deploy
./deploy.sh  # Slower, builds from source
```

## Checking Status

```bash
# Backend service
sudo systemctl status daybook-backend

# View logs
sudo journalctl -u daybook-backend -f

# Nginx
sudo systemctl status nginx
```

## Troubleshooting

If deployment fails:

```bash
# Check backend logs
sudo journalctl -u daybook-backend -n 100

# Check nginx
sudo nginx -t

# Restart services
sudo systemctl restart daybook-backend
sudo systemctl restart nginx
```

## That's All!

No complex configuration files.
No multiple scripts.
Just one command: `./deploy.sh`
