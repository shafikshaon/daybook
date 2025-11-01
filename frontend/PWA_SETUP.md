# PWA (Progressive Web App) Setup Guide

Your Daybook frontend is now configured as a Progressive Web App! Users can install it on their devices and use it like a native app.

## What's Been Implemented

✅ **Web App Manifest** (`public/manifest.json`)
✅ **Service Worker** (`public/sw.js`) - Offline support & caching
✅ **PWA Utilities** (`src/utils/pwa.js`) - Install prompts, notifications
✅ **Install Prompt Component** (`src/components/InstallPrompt.vue`)
✅ **Meta Tags** (in `index.html`) - iOS, Android, Windows support
✅ **Auto-registration** (in `main.js`)

## Required: Generate App Icons

You need to create app icons in multiple sizes. Follow these steps:

### Option 1: Using Online Tools (Easiest)

1. **Create a base icon** (512x512px or 1024x1024px)
   - Use any design tool (Figma, Canva, Photoshop)
   - Simple, recognizable design
   - Works well at small sizes
   - Transparent background for maskable icons

2. **Generate all sizes using:**
   - **PWA Asset Generator**: https://www.pwabuilder.com/imageGenerator
   - **RealFaviconGenerator**: https://realfavicongenerator.net/
   - **App Icon Generator**: https://appicon.co/

3. **Download and place in** `frontend/public/icons/`

### Option 2: Using ImageMagick (Command Line)

```bash
cd frontend/public

# Create a base icon first (icon-512x512.png)
# Then generate all sizes:

convert icon-512x512.png -resize 72x72 icons/icon-72x72.png
convert icon-512x512.png -resize 96x96 icons/icon-96x96.png
convert icon-512x512.png -resize 128x128 icons/icon-128x128.png
convert icon-512x512.png -resize 144x144 icons/icon-144x144.png
convert icon-512x512.png -resize 152x152 icons/icon-152x152.png
convert icon-512x512.png -resize 192x192 icons/icon-192x192.png
convert icon-512x512.png -resize 384x384 icons/icon-384x384.png
cp icon-512x512.png icons/icon-512x512.png
```

### Option 3: Using Node.js Script

```bash
npm install --save-dev sharp

# Create generate-icons.js:
```

```javascript
const sharp = require('sharp');
const sizes = [72, 96, 128, 144, 152, 192, 384, 512];

async function generateIcons() {
  for (const size of sizes) {
    await sharp('public/icon-base.png')
      .resize(size, size)
      .toFile(`public/icons/icon-${size}x${size}.png`);
    console.log(`Generated icon-${size}x${size}.png`);
  }
}

generateIcons();
```

```bash
node generate-icons.js
```

### Required Icon Sizes

| Size | Purpose |
|------|---------|
| 72x72 | Android small icon, badge |
| 96x96 | Android medium icon |
| 128x128 | Chrome Web Store |
| 144x144 | Windows tiles |
| 152x152 | iOS (iPad) |
| 192x192 | Android home screen (standard) |
| 384x384 | Android splash screen |
| 512x512 | PWA install prompt, splash screen |

### Optional: Screenshots

Create screenshots for better install experience:

1. **Desktop screenshot**: 1280x720px → `public/screenshots/desktop.png`
2. **Mobile screenshot**: 750x1334px → `public/screenshots/mobile.png`

## Features

### 1. Install Prompt

The app automatically shows an install prompt after 3 seconds when:
- User visits the site
- App is installable (HTTPS, has manifest, has service worker)
- User hasn't dismissed it recently (7 days)
- User hasn't already installed it

**For iOS users**, it shows manual instructions since iOS doesn't support the automatic prompt.

### 2. Offline Support

The service worker caches:
- HTML pages
- JavaScript bundles
- CSS files
- Images and assets

API calls are NOT cached to ensure data freshness.

### 3. Push Notifications

Ready for push notifications (requires backend integration):

```javascript
import { requestNotificationPermission, showNotification } from '@/utils/pwa';

// Request permission
const granted = await requestNotificationPermission();

// Show notification
await showNotification('Transaction Added', {
  body: 'Your transaction has been saved',
  icon: '/icons/icon-192x192.png'
});
```

### 4. App Shortcuts

Users can right-click the installed app icon to access shortcuts:
- Add Transaction
- View Dashboard

### 5. Standalone Mode

When installed, the app runs in standalone mode (no browser UI).

## Testing PWA

### Local Testing (HTTPS Required)

PWAs require HTTPS. For local testing:

```bash
# Option 1: Use vite's dev server (automatically uses localhost)
npm run dev

# Option 2: Build and serve with HTTPS
npm run build
npx serve -s dist --ssl-cert cert.pem --ssl-key key.pem
```

### Production Testing

After deploying to EC2:

1. Visit `https://daybook.shafik.xyz`
2. Check for install prompt (Chrome/Edge)
3. Open DevTools → Application → Manifest
4. Check Service Worker registration
5. Test offline mode (DevTools → Network → Offline)

### Browser DevTools

**Chrome/Edge:**
- Open DevTools (F12)
- Go to **Application** tab
- Check:
  - **Manifest**: Should load without errors
  - **Service Workers**: Should show registered worker
  - **Storage**: Check cache storage

**Lighthouse Audit:**
```bash
# Run PWA audit
npm install -g lighthouse
lighthouse https://daybook.shafik.xyz --view
```

## Platform-Specific Behavior

### Android (Chrome/Edge)

✅ **Automatic install prompt**
✅ **Home screen icon**
✅ **Splash screen**
✅ **Notification support**
✅ **Background sync**

### iOS (Safari)

⚠️ **Manual installation** (no automatic prompt)
✅ **Home screen icon**
✅ **Splash screen**
⚠️ **Limited notification support**
❌ **No background sync**

**To install on iOS:**
1. Tap Share button
2. Tap "Add to Home Screen"
3. Tap "Add"

### Desktop (Chrome/Edge/Firefox)

✅ **Install prompt**
✅ **Desktop shortcut**
✅ **Standalone window**
✅ **Notification support**

## Customization

### Update App Name/Colors

Edit `frontend/public/manifest.json`:

```json
{
  "name": "Your App Name",
  "short_name": "AppName",
  "theme_color": "#your-color",
  "background_color": "#your-bg-color"
}
```

### Update Caching Strategy

Edit `frontend/public/sw.js`:

```javascript
// Change cache name to force update
const CACHE_NAME = 'daybook-v2';

// Add more URLs to cache
const urlsToCache = [
  '/',
  '/index.html',
  '/manifest.json',
  '/your-critical-file.css'
];
```

### Customize Install Prompt

Edit `frontend/src/components/InstallPrompt.vue`:

- Change appearance
- Modify delay (default: 3 seconds)
- Change dismiss duration (default: 7 days)

## Deployment

### Build with PWA

```bash
cd frontend

# Ensure icons exist in public/icons/
ls -la public/icons/

# Build
npm run build

# Icons and manifest will be in dist/
```

### Deploy to EC2

```bash
# On local machine
./scripts/deploy-local.sh

# On EC2
./scripts/deploy-ec2.sh
```

### Verify Deployment

1. Visit `https://daybook.shafik.xyz`
2. Check manifest: `https://daybook.shafik.xyz/manifest.json`
3. Check service worker: `https://daybook.shafik.xyz/sw.js`
4. Check icons: `https://daybook.shafik.xyz/icons/icon-192x192.png`

## Troubleshooting

### Install Prompt Not Showing

**Check:**
1. Site must be HTTPS (not HTTP)
2. Must have valid manifest.json
3. Must have registered service worker
4. Icons must exist and be accessible
5. User hasn't dismissed recently
6. Not already installed

**Debug:**
```javascript
// In browser console
navigator.serviceWorker.getRegistrations()
  .then(regs => console.log('Service Workers:', regs));
```

### Service Worker Not Registering

**Check console for errors:**
- Open DevTools (F12)
- Go to Console tab
- Look for service worker errors

**Common issues:**
- Service worker file not found (must be in `public/sw.js`)
- HTTPS required (except localhost)
- JavaScript errors in service worker

**Force update:**
```javascript
// In browser console
navigator.serviceWorker.getRegistrations()
  .then(regs => regs.forEach(reg => reg.unregister()))
  .then(() => location.reload());
```

### Icons Not Loading

**Check:**
1. Icons exist in `public/icons/` folder
2. Build includes icons (check `dist/icons/`)
3. No 404 errors in Network tab
4. Correct paths in manifest.json

### iOS Not Showing Install Prompt

**This is normal!** iOS doesn't support automatic install prompts. Users must:
1. Open Safari
2. Tap Share button
3. Tap "Add to Home Screen"

The app shows manual instructions for iOS users.

## Best Practices

1. **Keep service worker simple** - Complex logic can cause issues
2. **Version your caches** - Increment cache name when updating
3. **Test offline mode** - Ensure critical features work offline
4. **Optimize icons** - Use PNG with transparent backgrounds
5. **Update regularly** - Check for service worker updates
6. **Monitor errors** - Log service worker errors

## Advanced Features

### Background Sync

Enable background sync for offline transactions:

```javascript
// Register sync
navigator.serviceWorker.ready.then(reg => {
  reg.sync.register('sync-transactions');
});

// Handle in service worker (sw.js)
self.addEventListener('sync', event => {
  if (event.tag === 'sync-transactions') {
    event.waitUntil(syncPendingTransactions());
  }
});
```

### Push Notifications (Requires Backend)

```javascript
// Request permission
const permission = await Notification.requestPermission();

// Subscribe to push
const registration = await navigator.serviceWorker.ready;
const subscription = await registration.pushManager.subscribe({
  userVisibleOnly: true,
  applicationServerKey: 'your-vapid-public-key'
});

// Send subscription to backend
await fetch('/api/v1/subscribe', {
  method: 'POST',
  body: JSON.stringify(subscription)
});
```

## Resources

- [PWA Documentation](https://web.dev/progressive-web-apps/)
- [Service Worker API](https://developer.mozilla.org/en-US/docs/Web/API/Service_Worker_API)
- [Web App Manifest](https://developer.mozilla.org/en-US/docs/Web/Manifest)
- [Workbox](https://developers.google.com/web/tools/workbox) - Advanced service worker library
- [PWABuilder](https://www.pwabuilder.com/) - PWA tools and testing

## Quick Checklist

- [ ] Generated all required icon sizes
- [ ] Icons placed in `public/icons/` folder
- [ ] Tested manifest loads without errors
- [ ] Service worker registers successfully
- [ ] Install prompt appears (on supported browsers)
- [ ] App works offline (basic functionality)
- [ ] Tested on mobile device
- [ ] Tested on desktop browser
- [ ] Lighthouse PWA audit passes
- [ ] Deployed to production with HTTPS

---

**Your app is now installable!** 🎉

Users can install Daybook on their home screen and use it like a native app with offline support.
