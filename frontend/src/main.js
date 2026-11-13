import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { IonicVue } from '@ionic/vue'

import App from './App.vue'
import router from './router'

/* Theme CSS */
import './theme/variables.css'

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
