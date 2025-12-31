# Quick Deploy Guide: /api/v1/auth/me Fix

## What Was Fixed

The `/api/v1/auth/me` endpoint was returning "User not found" because it was trying to query the `users` table with a non-existent `user_id` column.

**Fix**: Added `FindByUserID` method that correctly queries users by their `id` column.

## Deploy to Production (5 minutes)

### Step 1: Build Fixed Backend (on your local machine)

```bash
cd /Users/shafikshaon/workplace/development/projects/daybook/backend
go build -o daybook-backend
```

**Expected output**: Binary file `daybook-backend` created

### Step 2: Copy to Production Server

```bash
# Replace with your actual server details
scp daybook-backend user@your-server:/path/to/daybook/backend/
```

**Example**:
```bash
scp daybook-backend ubuntu@api.daybook.shafik.xyz:/home/ubuntu/daybook/backend/
```

### Step 3: SSH into Production and Restart Backend

```bash
# SSH into your server
ssh user@your-server

# Navigate to backend directory
cd /path/to/daybook/backend

# Stop old backend
pkill daybook-backend
# OR if using systemd
sudo systemctl stop daybook-backend
# OR if using pm2
pm2 stop daybook-backend

# Start new backend
./daybook-backend
# OR if using systemd
sudo systemctl start daybook-backend
# OR if using pm2
pm2 restart daybook-backend
```

### Step 4: Verify the Fix

```bash
# Test health endpoint (should return 200)
curl https://api.daybook.shafik.xyz/health

# Test /auth/me without token (should return 401, NOT 404)
curl https://api.daybook.shafik.xyz/api/v1/auth/me
```

**Expected responses**:
```bash
# Health check
{"status":"ok","message":"Daybook API is running"}

# /auth/me without token
{"success":false,"error":"Authorization header required"}
# HTTP Status: 401 (NOT 404!)
```

### Step 5: Test from Frontend

1. Open https://daybook.shafik.xyz in browser
2. Login with your credentials
3. Check browser console (F12) - should be no errors
4. Navigate to Settings - should show your profile
5. Try export feature - should work

## One-Liner Deploy (if you have SSH keys set up)

```bash
cd /Users/shafikshaon/workplace/development/projects/daybook/backend && \
go build -o daybook-backend && \
scp daybook-backend user@your-server:/path/to/daybook/backend/ && \
ssh user@your-server "cd /path/to/daybook/backend && pkill daybook-backend; nohup ./daybook-backend > /tmp/daybook-backend.log 2>&1 &" && \
echo "Deployed! Verifying..." && \
sleep 3 && \
curl -s https://api.daybook.shafik.xyz/health
```

**Replace**:
- `user@your-server` with your actual SSH login
- `/path/to/daybook` with actual path on server

## Troubleshooting

### Issue: "Permission denied" when copying

**Solution**: Make sure you have write permissions
```bash
# On production server, change ownership
sudo chown -R $USER:$USER /path/to/daybook/backend
```

### Issue: "Port 8080 already in use"

**Solution**: Kill the old process
```bash
sudo lsof -i :8080
sudo kill -9 <PID>
```

### Issue: Backend starts but immediately stops

**Check logs**:
```bash
# If using systemd
sudo journalctl -u daybook-backend -f

# If using log file
tail -f /tmp/daybook-backend.log

# If running directly
./daybook-backend
# (will show errors in console)
```

**Common causes**:
- Database connection issue → Check .env DB settings
- Missing .env file → Copy .env to production
- Port 8080 blocked → Check firewall

### Issue: Still getting 404 on /auth/me

**Possible causes**:
1. Old binary still running → Kill all instances: `pkill -9 daybook-backend`
2. Nginx not routing correctly → Check nginx config and restart
3. Wrong binary deployed → Verify: `./daybook-backend --version`

## Verification Checklist

- [ ] Backend built successfully
- [ ] Binary copied to production server
- [ ] Old backend process stopped
- [ ] New backend started
- [ ] Health check returns 200
- [ ] /auth/me returns 401 (not 404) without token
- [ ] Can login from frontend
- [ ] Settings page loads user profile
- [ ] Export features work

## Rollback (if needed)

If something goes wrong, you can rollback:

```bash
# SSH into production
ssh user@your-server
cd /path/to/daybook/backend

# Stop new backend
pkill daybook-backend

# Restore old binary (if you backed it up)
mv daybook-backend.backup daybook-backend

# Start old backend
./daybook-backend
```

**Prevention**: Always backup before deploying
```bash
# On production server, before deploying
cp daybook-backend daybook-backend.backup
```

## Files Changed in This Fix

1. `repository/user_repository.go` - Added FindByUserID method
2. `services/auth_service.go` - Updated GetProfile to use new method

**Total changes**: ~15 lines across 2 files

## Testing Locally Before Deploy

Before deploying to production, test locally:

```bash
cd /Users/shafikshaon/workplace/development/projects/daybook
./test-auth-me-fix.sh
```

This ensures the fix works on your local environment first.

---

**Ready to Deploy?** Follow steps 1-5 above!

**Questions?** Check `AUTH_ME_FIX_SUMMARY.md` for technical details.
