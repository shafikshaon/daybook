#!/usr/bin/env node

/**
 * Updates version.json with current build information
 * Run this script before building to ensure version tracking works
 */

import { execSync } from 'child_process'
import { writeFileSync } from 'fs'
import { fileURLToPath } from 'url'
import { dirname, join } from 'path'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

try {
  // Get git commit hash (short version)
  let commit = 'unknown'
  try {
    commit = execSync('git rev-parse --short HEAD', { encoding: 'utf-8' }).trim()
  } catch (e) {
    console.warn('Warning: Could not get git commit hash:', e.message)
  }

  // Generate version number (timestamp-based for uniqueness)
  const now = new Date()
  const version = `${now.getFullYear()}.${String(now.getMonth() + 1).padStart(2, '0')}.${String(now.getDate()).padStart(2, '0')}.${String(now.getHours()).padStart(2, '0')}${String(now.getMinutes()).padStart(2, '0')}`

  const versionData = {
    version: version,
    buildTime: now.toISOString(),
    commit: commit
  }

  // Write to public/version.json
  const versionPath = join(__dirname, '..', 'public', 'version.json')
  writeFileSync(versionPath, JSON.stringify(versionData, null, 2) + '\n')

  console.log('✓ Updated version.json:')
  console.log(`  Version: ${version}`)
  console.log(`  Build Time: ${versionData.buildTime}`)
  console.log(`  Commit: ${commit}`)
} catch (error) {
  console.error('Error updating version.json:', error)
  process.exit(1)
}
