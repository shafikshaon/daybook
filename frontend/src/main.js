import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// Import Bootstrap and custom styles
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import './assets/styles/custom.scss'

// Create Pinia instance FIRST (before router)
const pinia = createPinia()

// Create app instance
const app = createApp(App)

// Use Pinia BEFORE router to ensure stores are available in navigation guards
app.use(pinia)
app.use(router)

// Mount app
app.mount('#app')
