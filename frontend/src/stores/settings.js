import { defineStore } from 'pinia'
import axios from 'axios'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    app_name: 'SMARTSIS',
    school_name: 'Pondok Pesantren & Madrasah Terpadu',
    tagline: 'Sistem Informasi Sekolah & Pesantren Terpadu',
    phone: '',
    email: '',
    address: '',
    logo_url: '',
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
