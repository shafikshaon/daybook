# Console Errors - Diagnosis and Fixes

## Issues Found

### 1. ✅ Service Worker Error (FIXED)
**Error**: `Failed to execute 'put' on 'Cache': Request scheme 'chrome-extension' is unsupported`

**Cause**: Service worker was trying to cache chrome extension requests.

**Fix Applied**: Updated `frontend/public/sw.js` to:
- Skip non-http(s) requests (chrome-extension, etc.)
- Only cache GET requests from same origin
- Add proper error handling

**Status**: ✅ Fixed - Deploy to apply

---

### 2. ⚠️ Missing PWA Icons (ACTION REQUIRED)
**Error**: `GET https://daybook.shafik.xyz/icons/icon-144x144.png 404 (Not Found)`

**Cause**: PWA icons haven't been generated yet.

**Fix**: Generate icons using the automated script:

```bash
cd frontend

# Install dependencies (includes sharp for icon generation)
npm install

# Generate all PWA icons automatically
npm run generate-icons
```

This will create 8 icon sizes (72x72 to 512x512) from the placeholder SVG.

**Status**: ⚠️ Pending - Run the commands above

---

### 3. 🔴 Missing Build Files (CRITICAL)
**Errors**:
- `GET .../assets/BudgetsView-CRMrZmmf.js 404`
- `GET .../assets/ReportsView-DJ9kE3uQ.js 404`
- `GET .../assets/SettingsView-Dn72ayZt.js 404`

**Cause**: Build files weren't deployed correctly to EC2. This happens when:
1. Frontend wasn't rebuilt after recent changes
2. Build files weren't committed to git
3. Deployment script didn't copy all files

**Fix**: Complete rebuild and redeploy:

```bash
# 1. Clean and rebuild frontend
cd frontend
rm -rf dist node_modules/.vite
npm run build

# 2. Verify build succeeded
ls -la dist/assets/*.js
# Should see all view files

# 3. Deploy using scripts
cd ..
./scripts/deploy-local.sh
# When prompted, commit the changes

# 4. On EC2, pull and deploy
./scripts/deploy-ec2.sh
```

**Status**: 🔴 Critical - Must fix immediately

---

### 4. ⚠️ Backend API Missing (BACKEND ISSUE)
**Error**: `GET https://api.daybook.shafik.xyz/api/v1/fixed-deposits 404`

**Cause**: Backend doesn't have the fixed-deposits endpoint.

**Impact**: Non-critical - TransactionsView expects this endpoint but handles the error gracefully.

**Fix Options**:
1. **Remove the call** (quick fix):
   - Remove `fixedDepositsStore.fetchFixedDeposits()` from TransactionsView

2. **Add backend endpoint** (proper fix):
   - Implement `/api/v1/fixed-deposits` in the backend Go code

**For now**: The app works despite this error. Fix later if needed.

**Status**: ⚠️ Low priority - App still functional

---

## Complete Fix Procedure

### Step 1: Fix Service Worker & Generate Icons

```bash
cd frontend

# Install dependencies
npm install

# Generate PWA icons
npm run generate-icons

# Verify icons created
ls -la public/icons/*.png
# Should see 8 icon files
```

### Step 2: Clean Rebuild

```bash
# Still in frontend directory

# Clean previous build
rm -rf dist node_modules/.vite

# Fresh build
npm run build

# Verify build output
ls -la dist/
ls -la dist/assets/*.js
# Should see all view files (BudgetsView, ReportsView, etc.)
```

### Step 3: Deploy

```bash
cd ..

# Deploy from local machine
./scripts/deploy-local.sh
# Answer 'y' to commit changes

# On EC2 (or SSH to EC2 and run there)
./scripts/deploy-ec2.sh
```

### Step 4: Verify

Visit `https://daybook.shafik.xyz` and check:

1. **Console should be clean** (no service worker errors)
2. **PWA icons load** (check Application tab → Manifest)
3. **All pages work** (Dashboard, Budgets, Reports, Settings)
4. **Install prompt appears** on mobile/desktop

---

## Quick Verification Commands

### Check Build Files Locally
```bash
cd frontend/dist/assets
ls -la *.js | grep -E "(BudgetsView|ReportsView|SettingsView)"
```

### Check Icons Generated
```bash
ls frontend/public/icons/*.png | wc -l
# Should output: 10 (8 app icons + 2 shortcut icons)
```

### Check Deployment on EC2
```bash
# On EC2
ls -la /var/www/daybook/frontend/assets/*.js | grep -E "(BudgetsView|ReportsView|SettingsView)"
ls -la /var/www/daybook/frontend/icons/*.png
```

---

## Expected Console After Fix

After deploying all fixes, console should show:

✅ `Service Worker registered successfully: https://daybook.shafik.xyz/`
✅ No chrome-extension errors
✅ All icon URLs return 200
✅ All JS files load successfully
⚠️ `fixed-deposits 404` (expected - low priority backend issue)

---

## Files Modified

1. ✅ `frontend/public/sw.js` - Fixed service worker caching logic
2. ⏳ `frontend/public/icons/*.png` - Need to be generated
3. ⏳ `frontend/dist/` - Need to be rebuilt

---

## Deployment Order

**IMPORTANT**: Follow this exact order:

1. ✅ Generate icons (`npm run generate-icons`)
2. ✅ Build frontend (`npm run build`)
3. ✅ Verify build locally (check dist/assets)
4. ✅ Deploy local (`./scripts/deploy-local.sh`)
5. ✅ Deploy to EC2 (`./scripts/deploy-ec2.sh`)
6. ✅ Test in browser
7. ✅ Clear browser cache if needed (Ctrl+Shift+R)

---

## If Issues Persist

### Clear Service Worker Cache
```javascript
// In browser console:
navigator.serviceWorker.getRegistrations()
  .then(regs => regs.forEach(reg => reg.unregister()))
  .then(() => location.reload());
```

### Clear Browser Cache
- Chrome: Ctrl+Shift+Delete → Clear cached images and files
- Safari: Cmd+Option+E
- Mobile: Settings → Clear browser data

### Check Nginx on EC2
```bash
# On EC2
sudo nginx -t
sudo systemctl status nginx
sudo tail -f /var/log/nginx/error.log
```

---

## Summary

**Immediate Actions Required**:

1. 🔴 **CRITICAL**: Rebuild and redeploy frontend (missing JS files)
2. ⚠️ **IMPORTANT**: Generate PWA icons
3. ✅ **DONE**: Service worker fixed (deploy to apply)
4. ⚠️ **OPTIONAL**: Fix fixed-deposits backend endpoint

**Estimated Time**: 10-15 minutes total

**Commands to Run**:
```bash
cd frontend
npm install
npm run generate-icons
npm run build
cd ..
./scripts/deploy-local.sh
./scripts/deploy-ec2.sh
```

After these steps, your app should work perfectly with no console errors (except the optional fixed-deposits 404).
