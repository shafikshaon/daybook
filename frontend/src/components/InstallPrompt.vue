<template>
  <div v-if="showInstallPrompt" class="install-prompt">
    <div class="install-prompt-content">
      <button type="button" class="btn-close" @click="dismissPrompt" aria-label="Close"></button>

      <div class="install-prompt-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" fill="currentColor" viewBox="0 0 16 16">
          <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5z"/>
          <path d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708l3 3z"/>
        </svg>
      </div>

      <div class="install-prompt-text">
        <h5 class="mb-1">Install Daybook</h5>
        <p class="mb-2 text-muted" v-if="!isIOS">
          Install Daybook on your {{ deviceType }} for quick access and offline support.
        </p>
        <p class="mb-2 text-muted" v-else>
          Add Daybook to your home screen for quick access and offline support.
        </p>
      </div>

      <!-- iOS specific instructions -->
      <div v-if="isIOS && !isStandalone" class="install-prompt-ios mb-3">
        <div class="alert alert-info mb-0">
          <strong>📱 How to Install on iPhone/iPad:</strong>
          <ol class="mb-0 mt-2">
            <li>Tap the <strong>Share</strong> button
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16" style="vertical-align: text-bottom;">
                <path d="M11 2.5a2.5 2.5 0 1 1 .603 1.628l-6.718 3.12a2.499 2.499 0 0 1 0 1.504l6.718 3.12a2.5 2.5 0 1 1-.488.876l-6.718-3.12a2.5 2.5 0 1 1 0-3.256l6.718-3.12A2.5 2.5 0 0 1 11 2.5z"/>
              </svg>
              in Safari (at the bottom)
            </li>
            <li>Scroll down and select <strong>"Add to Home Screen"</strong></li>
            <li>Tap <strong>"Add"</strong> to install the app</li>
          </ol>
        </div>
      </div>

      <!-- Install button (hidden on iOS since it doesn't work) -->
      <div v-if="!isIOS" class="install-prompt-actions">
        <button type="button" class="btn btn-primary btn-sm" @click="installApp">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="me-1" viewBox="0 0 16 16">
            <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5z"/>
            <path d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708l3 3z"/>
          </svg>
          Install
        </button>
        <button type="button" class="btn btn-link btn-sm text-muted" @click="dismissPrompt">
          Not now
        </button>
      </div>

      <!-- Dismiss button for iOS -->
      <div v-else class="install-prompt-actions">
        <button type="button" class="btn btn-secondary btn-sm w-100" @click="dismissPrompt">
          Got it, thanks!
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { promptInstall, isAppInstalled, getDeviceInfo } from '@/utils/pwa';

export default {
  name: 'InstallPrompt',
  setup() {
    const showInstallPrompt = ref(false);
    const deviceInfo = ref({});

    const isIOS = computed(() => deviceInfo.value.isIOS);
    const isStandalone = computed(() => deviceInfo.value.isStandalone);
    const deviceType = computed(() => {
      if (deviceInfo.value.isIOS) return 'iPhone/iPad';
      if (deviceInfo.value.isAndroid) return 'Android device';
      return 'device';
    });

    const installApp = async () => {
      const installed = await promptInstall();
      if (installed) {
        showInstallPrompt.value = false;
        localStorage.setItem('pwa-install-prompted', 'installed');
      }
    };

    const dismissPrompt = () => {
      showInstallPrompt.value = false;
      // Don't show again for 7 days
      const dismissedUntil = Date.now() + (7 * 24 * 60 * 60 * 1000);
      localStorage.setItem('pwa-install-dismissed', dismissedUntil.toString());
    };

    const checkIfShouldShow = () => {
      // Don't show if already installed
      if (isAppInstalled()) {
        return false;
      }

      // Don't show if dismissed recently
      const dismissedUntil = localStorage.getItem('pwa-install-dismissed');
      if (dismissedUntil && Date.now() < parseInt(dismissedUntil)) {
        return false;
      }

      // Don't show if already installed
      if (localStorage.getItem('pwa-install-prompted') === 'installed') {
        return false;
      }

      return true;
    };

    const handleInstallable = () => {
      if (checkIfShouldShow()) {
        // Show after a delay to not interrupt user
        setTimeout(() => {
          showInstallPrompt.value = true;
        }, 3000);
      }
    };

    const handleInstalled = () => {
      showInstallPrompt.value = false;
      localStorage.setItem('pwa-install-prompted', 'installed');
    };

    onMounted(() => {
      deviceInfo.value = getDeviceInfo();

      // Listen for installable event
      window.addEventListener('pwa-installable', handleInstallable);
      window.addEventListener('pwa-installed', handleInstalled);

      // For iOS, show manual instructions if not installed
      if (deviceInfo.value.isIOS && !deviceInfo.value.isStandalone && checkIfShouldShow()) {
        setTimeout(() => {
          showInstallPrompt.value = true;
        }, 5000);
      }
    });

    onUnmounted(() => {
      window.removeEventListener('pwa-installable', handleInstallable);
      window.removeEventListener('pwa-installed', handleInstalled);
    });

    return {
      showInstallPrompt,
      isIOS,
      isStandalone,
      deviceType,
      installApp,
      dismissPrompt
    };
  }
};
</script>

<style scoped>
.install-prompt {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1050;
  max-width: 90%;
  width: 400px;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateX(-50%) translateY(100px);
    opacity: 0;
  }
  to {
    transform: translateX(-50%) translateY(0);
    opacity: 1;
  }
}

.install-prompt-content {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  padding: 20px;
  position: relative;
}

.btn-close {
  position: absolute;
  top: 10px;
  right: 10px;
}

.install-prompt-icon {
  text-align: center;
  color: #0d6efd;
  margin-bottom: 12px;
}

.install-prompt-text {
  text-align: center;
}

.install-prompt-text h5 {
  font-weight: 600;
  color: #212529;
}

.install-prompt-text p {
  font-size: 0.9rem;
}

.install-prompt-actions {
  display: flex;
  gap: 10px;
  justify-content: center;
  align-items: center;
}

.install-prompt-ios ol {
  font-size: 0.875rem;
  padding-left: 20px;
}

.install-prompt-ios li {
  margin-bottom: 4px;
}

/* Mobile responsiveness */
@media (max-width: 576px) {
  .install-prompt {
    bottom: 10px;
    max-width: 95%;
  }

  .install-prompt-content {
    padding: 16px;
  }
}

/* Dark mode support */
@media (prefers-color-scheme: dark) {
  .install-prompt-content {
    background: #212529;
    color: #f8f9fa;
  }

  .install-prompt-text h5 {
    color: #f8f9fa;
  }
}
</style>
