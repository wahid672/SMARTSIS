import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { IonicVue } from '@ionic/vue'
import axios from 'axios'

import App from './App.vue'
import router from './router'

/* Theme CSS */
import './theme/variables.css'

// Pasang Axios Interceptor agar setiap request API otomatis menyertakan Token JWT
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
}, (error) => {
  return Promise.reject(error)
})

// Tangani jika sesi expired (401 Unauthorized) otomatis redirect ke login
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      if (window.location.pathname !== '/login') {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(IonicVue, {
  mode: 'md' // Menggunakan Material Design modern yang konsisten
})
app.use(router)

router.isReady().then(() => {
  app.mount('#app')
})
