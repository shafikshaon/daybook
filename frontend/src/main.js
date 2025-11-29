import { createApp } from 'vue'
import { pinia } from '@/stores'
import router, { setupRouterGuard } from './router'
import App from './App.vue'
import { initVersionChecker } from '@/utils/versionChecker'

// Import Bootstrap and custom styles
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import './assets/styles/custom.scss'

// Unregister any service workers (we don't use PWA)
if ('serviceWorker' in navigator) {
  navigator.serviceWorker.getRegistrations().then(registrations => {
    registrations.forEach(registration => {
      registration.unregister()
      console.log('Service worker unregistered')
    })
  })
}

// Create app instance
const app = createApp(App)

// Install Pinia FIRST - this must happen before router
app.use(pinia)

// Install router
app.use(router)

// Setup router guard after both Pinia and router are installed
setupRouterGuard()

// Mount the app
app.mount('#app')

// Initialize version checker to detect updates
initVersionChecker()
