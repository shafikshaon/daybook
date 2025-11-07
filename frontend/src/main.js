import { createApp } from 'vue'
import { pinia } from '@/stores'
import router from './router'
import App from './App.vue'

// Import Bootstrap and custom styles
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import './assets/styles/custom.scss'

// Create app instance
const app = createApp(App)

// Install Pinia FIRST - this must happen before router
app.use(pinia)

// Install router
app.use(router)

// Mount the app
app.mount('#app')
