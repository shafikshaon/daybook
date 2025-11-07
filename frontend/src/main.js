import { createApp } from 'vue'
import { createPinia, setActivePinia } from 'pinia'
import router from './router'
import App from './App.vue'

// Import Bootstrap and custom styles
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import './assets/styles/custom.scss'

// Create Pinia instance FIRST
const pinia = createPinia()

// CRITICAL: Set Pinia as active BEFORE router setup
// This ensures router guards can access stores
setActivePinia(pinia)

// Create app instance
const app = createApp(App)

// Install Pinia on the app
app.use(pinia)

// Install router AFTER Pinia is active and installed
app.use(router)

// Wait for router to be ready before mounting the app
// This ensures router guards complete before mounting
router.isReady().then(() => {
  app.mount('#app')
})
