// Version Checker - Detects when a new version is deployed
// and prompts user to reload

let currentVersion = null
let checkInterval = null

export async function initVersionChecker() {
  try {
    // Fetch current version
    const response = await fetch('/version.json?t=' + Date.now())
    const data = await response.json()
    currentVersion = data.version

    console.log('Current app version:', currentVersion)

    // Check for updates every 5 minutes
    checkInterval = setInterval(checkForUpdate, 5 * 60 * 1000)
  } catch (error) {
    console.warn('Version checker initialization failed:', error)
  }
}

export async function checkForUpdate() {
  try {
    const response = await fetch('/version.json?t=' + Date.now())
    const data = await response.json()

    if (currentVersion && data.version !== currentVersion) {
      console.log('New version detected:', data.version, '(current:', currentVersion + ')')

      // Show update notification
      const shouldReload = confirm(
        'A new version of Daybook is available!\n\n' +
        'Click OK to reload and get the latest updates.'
      )

      if (shouldReload) {
        // Clear all caches and reload with cache busting
        if ('caches' in window) {
          const cacheNames = await caches.keys()
          await Promise.all(cacheNames.map(name => caches.delete(name)))
        }

        // Force a hard reload with cache bypass by adding version parameter
        const url = new URL(window.location.href)
        url.searchParams.set('v', data.version)
        window.location.href = url.toString()
      } else {
        // Update current version so we don't keep prompting
        currentVersion = data.version
      }
    }
  } catch (error) {
    console.warn('Version check failed:', error)
  }
}

export function stopVersionChecker() {
  if (checkInterval) {
    clearInterval(checkInterval)
    checkInterval = null
  }
}
