// PWA Utilities

let deferredPrompt = null;

/**
 * Register service worker
 */
export function registerServiceWorker() {
  if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
      navigator.serviceWorker
        .register('/sw.js')
        .then((registration) => {
          console.log('Service Worker registered successfully:', registration.scope);

          // Check for updates every hour
          setInterval(() => {
            registration.update();
          }, 60 * 60 * 1000);

          // Listen for updates
          registration.addEventListener('updatefound', () => {
            const newWorker = registration.installing;
            newWorker.addEventListener('statechange', () => {
              if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                // New service worker available
                showUpdateNotification();
              }
            });
          });
        })
        .catch((error) => {
          console.error('Service Worker registration failed:', error);
        });
    });
  }
}

/**
 * Show update notification when new version is available
 */
function showUpdateNotification() {
  if (confirm('A new version of Daybook is available. Reload to update?')) {
    window.location.reload();
  }
}

/**
 * Setup install prompt
 */
export function setupInstallPrompt() {
  window.addEventListener('beforeinstallprompt', (e) => {
    // Prevent the mini-infobar from appearing on mobile
    e.preventDefault();
    // Stash the event so it can be triggered later
    deferredPrompt = e;

    // Show install button
    const event = new CustomEvent('pwa-installable', { detail: { canInstall: true } });
    window.dispatchEvent(event);
  });

  window.addEventListener('appinstalled', () => {
    console.log('PWA was installed');
    deferredPrompt = null;

    // Hide install button
    const event = new CustomEvent('pwa-installed', { detail: { installed: true } });
    window.dispatchEvent(event);
  });
}

/**
 * Trigger install prompt
 */
export async function promptInstall() {
  if (!deferredPrompt) {
    return false;
  }

  // Show the install prompt
  deferredPrompt.prompt();

  // Wait for the user to respond to the prompt
  const { outcome } = await deferredPrompt.userChoice;
  console.log(`User response to install prompt: ${outcome}`);

  // Clear the deferredPrompt for reuse
  deferredPrompt = null;

  return outcome === 'accepted';
}

/**
 * Check if app is installed
 */
export function isAppInstalled() {
  // Check if running in standalone mode
  if (window.matchMedia('(display-mode: standalone)').matches) {
    return true;
  }

  // Check if running as iOS PWA
  if (window.navigator.standalone === true) {
    return true;
  }

  return false;
}

/**
 * Check if app can be installed
 */
export function canInstall() {
  return deferredPrompt !== null;
}

/**
 * Request notification permission
 */
export async function requestNotificationPermission() {
  if (!('Notification' in window)) {
    console.log('This browser does not support notifications');
    return false;
  }

  if (Notification.permission === 'granted') {
    return true;
  }

  if (Notification.permission !== 'denied') {
    const permission = await Notification.requestPermission();
    return permission === 'granted';
  }

  return false;
}

/**
 * Show notification
 */
export async function showNotification(title, options = {}) {
  if (!('Notification' in window)) {
    return false;
  }

  const hasPermission = await requestNotificationPermission();
  if (!hasPermission) {
    return false;
  }

  const defaultOptions = {
    icon: '/icons/icon-192x192.png',
    badge: '/icons/icon-72x72.png',
    vibrate: [200, 100, 200],
    ...options
  };

  if ('serviceWorker' in navigator && navigator.serviceWorker.controller) {
    // Use service worker to show notification
    navigator.serviceWorker.ready.then((registration) => {
      registration.showNotification(title, defaultOptions);
    });
  } else {
    // Fallback to regular notification
    new Notification(title, defaultOptions);
  }

  return true;
}

/**
 * Get device info
 */
export function getDeviceInfo() {
  const ua = navigator.userAgent;

  return {
    isIOS: /iPhone|iPad|iPod/.test(ua),
    isAndroid: /Android/.test(ua),
    isChrome: /Chrome/.test(ua) && !/Edge/.test(ua),
    isFirefox: /Firefox/.test(ua),
    isSafari: /Safari/.test(ua) && !/Chrome/.test(ua),
    isMobile: /Mobile|Android|iPhone|iPad|iPod/.test(ua),
    isStandalone: isAppInstalled()
  };
}

/**
 * Initialize PWA
 */
export function initPWA() {
  registerServiceWorker();
  setupInstallPrompt();

  const deviceInfo = getDeviceInfo();
  console.log('Device Info:', deviceInfo);

  // Log install status
  if (isAppInstalled()) {
    console.log('App is running as installed PWA');
  } else {
    console.log('App is running in browser');
  }
}
