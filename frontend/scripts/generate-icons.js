#!/usr/bin/env node

/**
 * PWA Icon Generator
 *
 * Generates all required PWA icons from a base image.
 * Usage: npm run generate-icons
 */

import sharp from 'sharp';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import { existsSync, mkdirSync } from 'fs';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Icon sizes required for PWA
const ICON_SIZES = [72, 96, 128, 144, 152, 192, 384, 512];

// Paths
const publicDir = join(__dirname, '..', 'public');
const iconsDir = join(publicDir, 'icons');
const baseIconPath = join(publicDir, 'icons', 'placeholder.svg');

// Ensure icons directory exists
if (!existsSync(iconsDir)) {
  mkdirSync(iconsDir, { recursive: true });
}

// Check if base icon exists
if (!existsSync(baseIconPath)) {
  console.error('❌ Error: Base icon not found at:', baseIconPath);
  console.log('\nPlease provide a base icon in one of these locations:');
  console.log('  - public/icons/placeholder.svg (SVG format)');
  console.log('  - public/icon-base.png (PNG format, 512x512 or larger)');
  process.exit(1);
}

console.log('🎨 Generating PWA icons...\n');
console.log('Base icon:', baseIconPath);
console.log('Output directory:', iconsDir);
console.log('');

// Generate icons
async function generateIcons() {
  let successCount = 0;
  let errorCount = 0;

  for (const size of ICON_SIZES) {
    const outputPath = join(iconsDir, `icon-${size}x${size}.png`);

    try {
      await sharp(baseIconPath)
        .resize(size, size, {
          fit: 'contain',
          background: { r: 13, g: 110, b: 253, alpha: 1 } // #0d6efd
        })
        .png()
        .toFile(outputPath);

      console.log(`✅ Generated: icon-${size}x${size}.png`);
      successCount++;
    } catch (error) {
      console.error(`❌ Failed to generate icon-${size}x${size}.png:`, error.message);
      errorCount++;
    }
  }

  console.log('\n📊 Summary:');
  console.log(`   ✅ Successfully generated: ${successCount} icons`);
  if (errorCount > 0) {
    console.log(`   ❌ Failed: ${errorCount} icons`);
  }
  console.log('');

  if (successCount === ICON_SIZES.length) {
    console.log('🎉 All PWA icons generated successfully!');
    console.log('');
    console.log('Next steps:');
    console.log('1. Run: npm run build');
    console.log('2. Deploy using: ../scripts/deploy-local.sh');
    console.log('3. Test PWA install prompt on mobile/desktop');
    console.log('');
    console.log('💡 Tip: Replace public/icons/placeholder.svg with your custom logo');
    console.log('   and re-run this script to regenerate icons.');
  } else {
    process.exit(1);
  }
}

// Also generate shortcut icons
async function generateShortcutIcons() {
  const shortcuts = [
    { name: 'shortcut-transaction.png', emoji: '💳' },
    { name: 'shortcut-dashboard.png', emoji: '📊' }
  ];

  console.log('\n🔗 Generating shortcut icons...\n');

  for (const shortcut of shortcuts) {
    const outputPath = join(iconsDir, shortcut.name);

    try {
      // Create a simple colored square with emoji-inspired design
      const svg = `
        <svg width="96" height="96" xmlns="http://www.w3.org/2000/svg">
          <rect width="96" height="96" fill="#0d6efd" rx="20"/>
          <text x="48" y="65" font-family="Arial, sans-serif" font-size="40" fill="#ffffff" text-anchor="middle">${shortcut.emoji}</text>
        </svg>
      `;

      await sharp(Buffer.from(svg))
        .resize(96, 96)
        .png()
        .toFile(outputPath);

      console.log(`✅ Generated: ${shortcut.name}`);
    } catch (error) {
      console.error(`❌ Failed to generate ${shortcut.name}:`, error.message);
    }
  }
}

// Run generation
(async () => {
  try {
    await generateIcons();
    await generateShortcutIcons();
  } catch (error) {
    console.error('\n❌ Error generating icons:', error);
    process.exit(1);
  }
})();
