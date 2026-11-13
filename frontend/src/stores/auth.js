import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    menus: JSON.parse(localStorage.getItem('menus') || '[]'),
    daftarSiswa: JSON.parse(localStorage.getItem('daftarSiswa') || '[]'),
    activeSiswa: JSON.parse(localStorage.getItem('activeSiswa') || 'null')
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isWali: (state) => state.user?.user_type === 'wali',
    isPegawai: (state) => state.user?.user_type === 'pegawai',
    roleLabel: (state) => state.user?.role_label || ''
  },

  actions: {
    async login(identifier, password) {
      try {
        const response = await axios.post('/api/auth/login', {
          identifier,
          password
        })

        const data = response.data
        this.token = data.token
        this.user = data.user
        this.menus = data.menus || []
        this.daftarSiswa = data.daftar_siswa || []

        // Jika wali murid punya siswa, set default aktif ke anak pertama
        if (this.daftarSiswa.length > 0) {
          this.activeSiswa = this.daftarSiswa[0]
        } else {
          this.activeSiswa = null
        }

        // Simpan sesi ke localStorage
        localStorage.setItem('token', this.token)
        localStorage.setItem('user', JSON.stringify(this.user))
        localStorage.setItem('menus', JSON.stringify(this.menus))
        localStorage.setItem('daftarSiswa', JSON.stringify(this.daftarSiswa))
        localStorage.setItem('activeSiswa', JSON.stringify(this.activeSiswa))

        return { success: true }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.error || 'Gagal login, periksa koneksi Anda.'
        }
      }
    },

    setActiveSiswa(siswa) {
      this.activeSiswa = siswa
      localStorage.setItem('activeSiswa', JSON.stringify(siswa))
    },

    logout() {
      this.token = ''
      this.user = null
      this.menus = []
      this.daftarSiswa = []
      this.activeSiswa = null
      localStorage.clear()
    }
  }
})
