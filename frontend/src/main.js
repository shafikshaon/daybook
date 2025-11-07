import { createApp } from 'vue'
import { setActivePinia } from 'pinia'
import { pinia } from '@/stores'
import router from './router'
import App from './App.vue'

// Import Bootstrap and custom styles
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import './assets/styles/custom.scss'

// CRITICAL: Set this pinia instance as the active one globally
// This makes it available to all stores, even when used outside components
setActivePinia(pinia)

// Create app instance
const app = createApp(App)

// Install Pinia on the app
app.use(pinia)

// Install router
app.use(router)

// Mount the app
app.mount('#app')
