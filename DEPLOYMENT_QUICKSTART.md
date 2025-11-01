# Daybook EC2 Deployment - Quick Start

This guide provides a fast-track deployment process for your Daybook application on EC2.

## What You Have

Your deployment documentation includes:

1. **EC2_DEPLOYMENT_GUIDE.md** - Complete step-by-step deployment guide (14 parts)
2. **DEPLOYMENT_CHECKLIST.md** - Printable checklist with checkboxes
3. **scripts/local-build.sh** - Automated build script for local machine
4. **scripts/ec2-setup.sh** - Automated setup script for EC2 instance

## Target Architecture

```
┌─────────────────────────────────────────────────┐
│                  EC2 Instance                    │
│              (10GB storage, 0.5GB RAM)          │
├─────────────────────────────────────────────────┤
│                                                  │
│  ┌──────────────┐         ┌──────────────┐     │
│  │   Nginx      │         │   Nginx      │     │
│  │   :8081      │         │   :8081      │     │
│  │ (Frontend)   │         │  (Backend)   │     │
│  └──────┬───────┘         └──────┬───────┘     │
│         │                        │             │
│  ┌──────▼───────┐         ┌──────▼───────┐     │
│  │   Vue.js     │         │   Go API     │     │
│  │   Static     │         │   :8082      │     │
│  │   Files      │         │              │     │
│  └──────────────┘         └──────┬───────┘     │
│                                  │             │
│                           ┌──────▼───────┐     │
│                           │  PostgreSQL  │     │
│                           │    :5432     │     │
│                           └──────────────┘     │
│                                                 │
│                           ┌──────────────┐     │
│                           │    Redis     │     │
│                           │    :6379     │     │
│                           └──────────────┘     │
└─────────────────────────────────────────────────┘
         │                      │
         │                      │
    daybook.shafik.xyz   api.daybook.shafik.xyz
```

## Prerequisites Checklist

- [ ] EC2 instance running Ubuntu 22.04 LTS
- [ ] SSH key (.pem file) for EC2 access
- [ ] DNS access to configure A records
- [ ] Local machine with Go 1.25+ and Node.js 22+

## 5-Minute Quick Deploy

### Step 1: Build and Push (2 minutes)

On your **local machine**:

```bash
cd /path/to/daybook

# Build the application
./scripts/local-build.sh

# Commit and push build artifacts
git add build/
git commit -m "Add deployment build artifacts"
git push origin master
```

This creates and pushes:
- `build/daybook-backend` - Go binary for Linux
- `build/frontend-dist.tar.gz` - Frontend production build

### Step 2: Configure DNS (5-30 minutes for propagation)

Add these A records at your domain registrar:

```
daybook.shafik.xyz     → Your EC2 Public IP
api.daybook.shafik.xyz → Your EC2 Public IP
```

Verify propagation:
```bash
nslookup daybook.shafik.xyz
nslookup api.daybook.shafik.xyz
```

### Step 3: Prepare EC2 (5 minutes)

Connect to EC2:
```bash
ssh -i your-key.pem ubuntu@your-ec2-ip
```

Install dependencies:
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install PostgreSQL
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt update
sudo apt install -y postgresql-15 postgresql-contrib-15 nginx redis-server certbot python3-certbot-nginx

# Configure firewall
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8081/tcp
sudo ufw allow 8082/tcp
sudo ufw --force enable
```

Create database:
```bash
sudo -u postgres psql -c "CREATE DATABASE daybook;"
sudo -u postgres psql -c "CREATE USER daybook_user WITH ENCRYPTED PASSWORD 'YourStrongPassword123!';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE daybook TO daybook_user;"
```

Configure Redis password:
```bash
sudo sed -i 's/# requirepass foobared/requirepass YourRedisPassword123/' /etc/redis/redis.conf
sudo systemctl restart redis-server
```

### Step 4: Clone and Deploy (2 minutes)

On **EC2**, clone the repository and run setup:
```bash
# Clone repository
cd ~
git clone https://github.com/shafikshaon/daybook.git
cd daybook

# Run deployment script
./scripts/ec2-setup.sh
```

The script will prompt for:
- Database password
- Redis password
- Other configuration (with sensible defaults)

### Step 5: Start Services (1 minute)

```bash
# Start backend
sudo systemctl daemon-reload
sudo systemctl enable daybook-backend
sudo systemctl start daybook-backend

# Reload Nginx
sudo systemctl reload nginx

# Check status
sudo systemctl status daybook-backend
```

### Step 6: Enable HTTPS (2 minutes)

```bash
# Frontend
sudo certbot --nginx -d daybook.shafik.xyz

# Backend
sudo certbot --nginx -d api.daybook.shafik.xyz
```

Select option 2 to redirect HTTP to HTTPS.

### Step 7: Verify (1 minute)

```bash
# Check services
sudo systemctl status daybook-backend postgresql redis-server nginx

# Check logs
sudo journalctl -u daybook-backend -n 20

# Test endpoints
curl https://daybook.shafik.xyz
curl https://api.daybook.shafik.xyz
```

Open in browser:
- https://daybook.shafik.xyz

## Important Environment Variables

### Backend (.env location: /var/www/daybook/backend/.env)

Key variables to set:
```env
DB_PASSWORD=YourStrongPassword123!
REDIS_PASSWORD=YourRedisPassword123
JWT_SECRET=<generated-by-script>
CORS_ALLOWED_ORIGINS=https://daybook.shafik.xyz
```

### Frontend (built into static files)

```env
VITE_API_URL=https://api.daybook.shafik.xyz
```

## Common Commands

### Service Management
```bash
# Restart backend
sudo systemctl restart daybook-backend

# View backend logs
sudo journalctl -u daybook-backend -f

# Reload Nginx
sudo systemctl reload nginx
```

### Database
```bash
# Connect to database
psql -h localhost -U daybook_user -d daybook

# Backup
pg_dump -h localhost -U daybook_user -d daybook > backup.sql

# Restore
psql -h localhost -U daybook_user -d daybook < backup.sql
```

### Monitoring
```bash
# Resource usage
free -h        # Memory
df -h          # Disk
htop           # Processes

# Application logs
sudo tail -f /var/log/daybook-backend.log
sudo tail -f /var/log/nginx/daybook-frontend-access.log
sudo tail -f /var/log/nginx/daybook-backend-access.log
```

## Deployment Updates

To deploy code changes:

1. **Build and push (local machine):**
   ```bash
   # Build
   ./scripts/local-build.sh

   # Commit and push
   git add build/
   git commit -m "Update deployment build"
   git push origin master
   ```

2. **Pull and update (EC2):**
   ```bash
   # Pull latest changes
   cd ~/daybook
   git pull origin master

   # Stop backend
   sudo systemctl stop daybook-backend

   # Update backend
   cp ~/daybook/build/daybook-backend /var/www/daybook/backend/daybook-backend
   chmod +x /var/www/daybook/backend/daybook-backend

   # Update frontend
   cd /var/www/daybook/frontend
   rm -rf *
   tar -xzf ~/daybook/build/frontend-dist.tar.gz

   # Start backend
   sudo systemctl start daybook-backend

   # Reload Nginx
   sudo systemctl reload nginx
   ```

## Troubleshooting

### Backend won't start
```bash
# Check logs
sudo journalctl -u daybook-backend -n 50 --no-pager

# Check if port is in use
sudo netstat -tulpn | grep 8082

# Test database connection
psql -h localhost -U daybook_user -d daybook
```

### 502 Bad Gateway
```bash
# Check backend status
sudo systemctl status daybook-backend

# Check Nginx config
sudo nginx -t

# Check logs
sudo tail -50 /var/log/nginx/daybook-backend-error.log
```

### Out of Memory
```bash
# Add 1GB swap
sudo fallocate -l 1G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

### SSL Certificate Issues
```bash
# Check certificates
sudo certbot certificates

# Renew manually
sudo certbot renew --force-renewal

# Test auto-renewal
sudo certbot renew --dry-run
```

## Resource Optimization Tips

Given your limited resources (0.5GB RAM):

1. **PostgreSQL** is configured for low memory (64MB shared_buffers)
2. **Redis** is capped at 50MB
3. **Go binary** is statically compiled (no runtime overhead)
4. **Vue.js** is pre-built (no Node.js runtime needed)
5. **Swap** recommended (1GB swap file)

### Adding Swap
```bash
sudo fallocate -l 1G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

## Security Checklist

- [ ] Strong database password
- [ ] Strong Redis password
- [ ] Unique JWT secret (32+ characters)
- [ ] HTTPS enabled (Let's Encrypt)
- [ ] Firewall configured (UFW)
- [ ] .env file permissions (600)
- [ ] Regular database backups
- [ ] SSH key-based authentication

## Next Steps

After deployment:

1. **Set up monitoring** - CloudWatch, custom scripts, or third-party
2. **Configure backups** - Database, uploads directory
3. **Document credentials** - Store securely (1Password, AWS Secrets Manager)
4. **Test functionality** - Sign up, create accounts, add transactions
5. **Monitor resources** - Use `htop` and `free -h` regularly
6. **Plan for scaling** - Consider RDS/ElastiCache if traffic grows

## Additional Resources

- **Full Guide:** `EC2_DEPLOYMENT_GUIDE.md` - Detailed 14-part guide
- **Checklist:** `DEPLOYMENT_CHECKLIST.md` - Printable deployment checklist
- **Build Script:** `scripts/local-build.sh` - Automated local build
- **Setup Script:** `scripts/ec2-setup.sh` - Automated EC2 setup

## Support

If you encounter issues:

1. Check logs: `sudo journalctl -u daybook-backend -f`
2. Review the full deployment guide: `EC2_DEPLOYMENT_GUIDE.md`
3. Verify all services are running: `sudo systemctl status daybook-backend postgresql redis-server nginx`
4. Check resource usage: `free -h` and `df -h`

---

**Happy Deploying!** 🚀

Your application will be available at:
- Frontend: https://daybook.shafik.xyz
- Backend API: https://api.daybook.shafik.xyz
