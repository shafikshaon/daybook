# Daybook Deployment

## Super Simple Deployment

### On Your Server:

```bash
# 1. Clone the repository
cd /opt
sudo git clone https://github.com/shafikshaon/daybook.git
sudo chown -R $USER:$USER daybook

# 2. Run deployment
cd daybook/deploy
./deploy.sh

# That's it! 🎉
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

## Updating Your App

```bash
cd /opt/daybook
git pull
cd deploy
./deploy.sh
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
