# PWA App Icons

This directory should contain the app icons for the Progressive Web App.

## Required Icons

You need to generate icons in the following sizes:

```
icons/
├── icon-72x72.png
├── icon-96x96.png
├── icon-128x128.png
├── icon-144x144.png
├── icon-152x152.png
├── icon-192x192.png
├── icon-384x384.png
└── icon-512x512.png
```

## Quick Generate Icons

### Option 1: Online Tool (Easiest)

1. Go to: https://www.pwabuilder.com/imageGenerator
2. Upload your logo/icon (512x512px or larger)
3. Download the generated package
4. Extract and copy PNG files here

### Option 2: Using ImageMagick

```bash
# Install ImageMagick first (macOS):
brew install imagemagick

# Or Ubuntu:
sudo apt install imagemagick

# Create your base icon (icon-base.png) in public/ folder
# Then run:
cd frontend/public

convert icon-base.png -resize 72x72 icons/icon-72x72.png
convert icon-base.png -resize 96x96 icons/icon-96x96.png
convert icon-base.png -resize 128x128 icons/icon-128x128.png
convert icon-base.png -resize 144x144 icons/icon-144x144.png
convert icon-base.png -resize 152x152 icons/icon-152x152.png
convert icon-base.png -resize 192x192 icons/icon-192x192.png
convert icon-base.png -resize 384x384 icons/icon-384x384.png
convert icon-base.png -resize 512x512 icons/icon-512x512.png
```

### Option 3: Using Node.js

```bash
# Install sharp
npm install --save-dev sharp

# Run the generator script
npm run generate-icons
```

## Icon Design Tips

1. **Simple Design**: Icons should be recognizable at small sizes
2. **Transparent Background**: Use PNG with transparency
3. **Safe Zone**: Keep important elements in the center 80%
4. **Square Canvas**: Design on a square canvas (512x512 or 1024x1024)
5. **Consistent Style**: Match your app's branding

## Example Design Ideas for Daybook

- 💰 Money bag icon
- 📊 Bar chart with dollar sign
- 💵 Stylized dollar bill
- 📱 Phone with money symbol
- 🏦 Bank/building icon
- 📈 Growth chart with currency

## Temporary Icons

Until you create custom icons, you can use:

1. **Placeholder service**: https://via.placeholder.com/512/0d6efd/fff?text=D
2. **Emoji to PNG**: Convert 💰 emoji to PNG using: https://emoji.to
3. **Free icon packs**: https://www.flaticon.com/ or https://icons8.com/

## Verify Icons

After generating, check that all files exist:

```bash
ls -la icons/
```

You should see all 8 icon files listed above.

## Testing

1. Build the app: `npm run build`
2. Check `dist/icons/` folder has all icons
3. Deploy and visit your site
4. Open DevTools → Application → Manifest
5. Check "Icons" section shows all icons loaded

---

**Important**: Don't forget to generate icons before deploying! The PWA install prompt won't work without them.
