import { defineStore } from 'pinia'
import axios from 'axios'

export const THEME_PALETTES = {
  emerald: {
    id: 'emerald',
    name: 'Emerald Green',
    subtitle: 'Pesantren, MI/MTs/MA & Nuansa Islami',
    primary: '#10b981',
    primaryHover: '#059669',
    base: '#064e3b',
    gradient: 'radial-gradient(circle at 100% 150%, #047857 24%, #064e3b 50%, #022c22 85%)',
    sidebarBg: '#064e3b',
    sidebarHeader: 'rgba(2, 44, 34, 0.5)',
    accent: '#34d399',
    buttonGradient: 'from-emerald-600 to-emerald-500 hover:from-emerald-700 hover:to-emerald-600',
    buttonShadow: 'shadow-emerald-600/25 hover:shadow-emerald-600/35',
    ring: 'ring-emerald-500/30',
    swatch: '#064e3b'
  },
  navy: {
    id: 'navy',
    name: 'Ocean Navy',
    subtitle: 'Sekolah Umum, SMP/SMA/SMK & Formal',
    primary: '#0284c7',
    primaryHover: '#0369a1',
    base: '#0b3b60',
    gradient: 'radial-gradient(circle at 100% 150%, #0f3057 24%, #0b3b60 50%, #071e33 85%)',
    sidebarBg: '#0b3b60',
    sidebarHeader: 'rgba(7, 30, 51, 0.5)',
    accent: '#38bdf8',
    buttonGradient: 'from-sky-600 to-blue-600 hover:from-sky-700 hover:to-blue-700',
    buttonShadow: 'shadow-blue-600/25 hover:shadow-blue-600/35',
    ring: 'ring-sky-500/30',
    swatch: '#0b3b60'
  },
  maroon: {
    id: 'maroon',
    name: 'Royal Maroon',
    subtitle: 'Wibawa, Tradisi & Kehormatan',
    primary: '#e11d48',
    primaryHover: '#be123c',
    base: '#4c0519',
    gradient: 'radial-gradient(circle at 100% 150%, #881337 24%, #4c0519 50%, #310410 85%)',
    sidebarBg: '#4c0519',
    sidebarHeader: 'rgba(49, 4, 16, 0.5)',
    accent: '#fb7185',
    buttonGradient: 'from-rose-600 to-pink-600 hover:from-rose-700 hover:to-pink-700',
    buttonShadow: 'shadow-rose-600/25 hover:shadow-rose-600/35',
    ring: 'ring-rose-500/30',
    swatch: '#4c0519'
  },
  teal: {
    id: 'teal',
    name: 'Midnight Teal',
    subtitle: 'Modern, Sejuk & Ramah Edukasi',
    primary: '#0d9488',
    primaryHover: '#0f766e',
    base: '#042f2e',
    gradient: 'radial-gradient(circle at 100% 150%, #115e59 24%, #042f2e 50%, #021a19 85%)',
    sidebarBg: '#042f2e',
    sidebarHeader: 'rgba(2, 26, 25, 0.5)',
    accent: '#2dd4bf',
    buttonGradient: 'from-teal-600 to-emerald-600 hover:from-teal-700 hover:to-emerald-700',
    buttonShadow: 'shadow-teal-600/25 hover:shadow-teal-600/35',
    ring: 'ring-teal-500/30',
    swatch: '#042f2e'
  },
  slate: {
    id: 'slate',
    name: 'Dark Slate',
    subtitle: 'Minimalis, Netral & Tech-Forward',
    primary: '#475569',
    primaryHover: '#334155',
    base: '#0f172a',
    gradient: 'radial-gradient(circle at 100% 150%, #1e293b 24%, #0f172a 50%, #020617 85%)',
    sidebarBg: '#0f172a',
    sidebarHeader: 'rgba(2, 6, 23, 0.5)',
    accent: '#94a3b8',
    buttonGradient: 'from-slate-700 to-slate-600 hover:from-slate-800 hover:to-slate-700',
    buttonShadow: 'shadow-slate-700/25 hover:shadow-slate-700/35',
    ring: 'ring-slate-500/30',
    swatch: '#0f172a'
  },
  indigo: {
    id: 'indigo',
    name: 'Deep Indigo',
    subtitle: 'Akademis, Visioner & Eksklusif',
    primary: '#6366f1',
    primaryHover: '#4f46e5',
    base: '#1e1b4b',
    gradient: 'radial-gradient(circle at 100% 150%, #312e81 24%, #1e1b4b 50%, #0f0d2e 85%)',
    sidebarBg: '#1e1b4b',
    sidebarHeader: 'rgba(15, 13, 46, 0.5)',
    accent: '#818cf8',
    buttonGradient: 'from-indigo-600 to-violet-600 hover:from-indigo-700 hover:to-violet-700',
    buttonShadow: 'shadow-indigo-600/25 hover:shadow-indigo-600/35',
    ring: 'ring-indigo-500/30',
    swatch: '#1e1b4b'
  }
}

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

export function applyThemeToDom(themeKey) {
  const theme = THEME_PALETTES[themeKey] || THEME_PALETTES.emerald
  if (typeof document !== 'undefined') {
    const root = document.documentElement
    root.style.setProperty('--theme-sidebar-bg', theme.sidebarBg)
    root.style.setProperty('--theme-base', theme.base)
    root.style.setProperty('--theme-gradient', theme.gradient)
    root.style.setProperty('--theme-primary', theme.primary)
    root.style.setProperty('--theme-accent', theme.accent)
  }
}

// Langsung terapkan theme cache saat file ini dimuat
applyThemeToDom(cachedSettings.theme_color || 'emerald')

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    app_name: cachedSettings.app_name || 'SMARTSIS',
    school_name: cachedSettings.school_name || 'Pondok Pesantren & Madrasah Terpadu',
    tagline: cachedSettings.tagline || 'Sistem Informasi Sekolah & Pesantren Terpadu',
    phone: cachedSettings.phone || '',
    email: cachedSettings.email || '',
    address: cachedSettings.address || '',
    logo_url: cachedSettings.logo_url || '',
    theme_color: cachedSettings.theme_color || 'emerald',
    loading: false
  }),

  getters: {
    currentTheme: (state) => THEME_PALETTES[state.theme_color] || THEME_PALETTES.emerald
  },

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
          this.theme_color = res.data.theme_color || 'emerald'

          // Terapkan variabel warna ke DOM
          applyThemeToDom(this.theme_color)

          // Simpan ke cache lokal agar reload mulus tanpa berkedip
          localStorage.setItem('app_settings', JSON.stringify({
            app_name: this.app_name,
            school_name: this.school_name,
            tagline: this.tagline,
            phone: this.phone,
            email: this.email,
            address: this.address,
            logo_url: this.logo_url,
            theme_color: this.theme_color
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

        // Terapkan tema yang baru disimpan
        if (data.theme_color) {
          this.theme_color = data.theme_color
          applyThemeToDom(data.theme_color)
        }

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
