import { defineStore } from 'pinia'
import axios from 'axios'

// Ambil cache settings dari localStorage jika ada, agar saat reload tidak berkedip
const cachedSettings = JSON.parse(localStorage.getItem('app_settings') || '{}')

function updateFavicon(url) {
  if (!url) return
  let link = document.querySelector("link[rel~='icon']")
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.getElementsByTagName('head')[0].appendChild(link)
  }
  link.href = url
}

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    app_name: cachedSettings.app_name || 'SMARTSIS',
    school_name: cachedSettings.school_name || 'Pondok Pesantren & Madrasah Terpadu',
    tagline: cachedSettings.tagline || 'Sistem Informasi Sekolah & Pesantren Terpadu',
    phone: cachedSettings.phone || '',
    email: cachedSettings.email || '',
    address: cachedSettings.address || '',
    logo_url: cachedSettings.logo_url || '',
    loading: false
  }),

  actions: {
    async fetchSettings() {
      try {
        const res = await axios.get('/api/settings')
        if (res.data) {
          this.app_name = res.data.app_name || 'SMARTSIS'
          this.school_name = res.data.school_name || 'Sistem Informasi Sekolah'
          this.tagline = res.data.tagline || ''
          this.phone = res.data.phone || ''
          this.email = res.data.email || ''
          this.address = res.data.address || ''
          this.logo_url = res.data.logo_url || ''

          // Simpan ke cache lokal agar reload mulus tanpa berkedip
          localStorage.setItem('app_settings', JSON.stringify({
            app_name: this.app_name,
            school_name: this.school_name,
            tagline: this.tagline,
            phone: this.phone,
            email: this.email,
            address: this.address,
            logo_url: this.logo_url
          }))

          // Update judul tab & Favicon
          document.title = `${this.app_name} - ${this.school_name}`
          if (this.logo_url) {
            updateFavicon(this.logo_url)
          }
        }
      } catch (err) {
        console.error('Gagal mengambil pengaturan:', err)
      }
    },

    async saveSettings(data, token) {
      this.loading = true
      try {
        await axios.post('/api/settings', data, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Object.assign(this, data)
        localStorage.setItem('app_settings', JSON.stringify(data))

        document.title = `${this.app_name} - ${this.school_name}`
        if (this.logo_url) {
          updateFavicon(this.logo_url)
        }

        return { success: true }
      } catch (err) {
        return {
          success: false,
          message: err.response?.data?.error || 'Gagal menyimpan pengaturan lembaga'
        }
      } finally {
        this.loading = false
      }
    }
  }
})
