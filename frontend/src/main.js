import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// Import Bootstrap and custom styles
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import './assets/styles/custom.scss'

// Import PWA utilities
import { initPWA } from './utils/pwa'

// Create app instance
const app = createApp(App)

// Use plugins
app.use(createPinia())
app.use(router)

// Mount app
app.mount('#app')

// Initialize PWA (service worker, install prompt, etc.)
initPWA()
