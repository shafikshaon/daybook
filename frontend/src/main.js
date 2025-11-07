import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// Import Bootstrap and custom styles
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import './assets/styles/custom.scss'

// Create Pinia instance FIRST
const pinia = createPinia()

// Create app instance
const app = createApp(App)

// CRITICAL: Install Pinia FIRST, synchronously, before anything else
app.use(pinia)

// Install router AFTER Pinia is installed
app.use(router)

// Wait for router to be ready before mounting the app
// This ensures router guards complete and Pinia is fully active before component setup
router.isReady().then(() => {
  app.mount('#app')
})
