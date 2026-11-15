<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding-none bg-slate-100">
      <div class="min-h-full flex items-center justify-center p-0 sm:p-4 lg:p-10 font-sans text-slate-800">
        
        <!-- ======================================================== -->
        <!-- MAIN AUTH CONTAINER (SPLIT SCREEN DESKTOP & STACK MOBILE)-->
        <!-- ======================================================== -->
        <main class="w-full max-w-6xl bg-white rounded-none sm:rounded-3xl shadow-2xl overflow-hidden border border-slate-100 flex flex-col lg:flex-row min-h-[680px]">
          
          <!-- ==================================================== -->
          <!-- HERO BRANDING PANEL (LEFT DESKTOP / TOP MOBILE)      -->
          <!-- ==================================================== -->
          <section class="lg:w-5/12 islamic-pattern pattern-overlay relative p-6 sm:p-10 lg:p-14 text-white flex flex-col justify-between overflow-hidden">
            <!-- Glow Aesthetic Backdrops -->
            <div class="absolute -top-24 -left-24 w-80 h-80 bg-emerald-400/20 rounded-full blur-3xl pointer-events-none"></div>
            <div class="absolute -bottom-24 -right-24 w-80 h-80 bg-emerald-500/20 rounded-full blur-3xl pointer-events-none"></div>

            <div class="relative z-10">
              <!-- Logo Emblem & App Title -->
              <div class="flex items-center space-x-3.5 mb-6 sm:mb-8">
                <div class="w-14 h-14 sm:w-16 sm:h-16 rounded-2xl bg-white p-2 shadow-xl flex items-center justify-center ring-4 ring-emerald-500/30 shrink-0">
                  <img
                    v-if="settingsStore.logo_url"
                    :src="settingsStore.logo_url"
                    :alt="settingsStore.app_name"
                    class="w-full h-full object-contain"
                  />
                  <ion-icon v-else :icon="schoolOutline" class="text-3xl sm:text-4xl text-emerald-600"></ion-icon>
                </div>
                <div>
                  <span class="inline-block text-[10px] uppercase tracking-widest font-bold text-emerald-200 bg-emerald-900/70 px-2.5 py-0.5 rounded-full border border-emerald-500/30">
                    Sistem Terpadu
                  </span>
                  <h1 class="text-xl sm:text-2xl font-black tracking-tight text-white mt-0.5">
                    {{ settingsStore.app_name }}
                  </h1>
                </div>
              </div>

              <!-- Institution Title & Motto -->
              <div class="space-y-3">
                <div class="inline-flex items-center gap-2 text-emerald-300 font-semibold tracking-wide text-xs sm:text-sm uppercase">
                  <svg class="w-4 h-4 text-emerald-300 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path>
                  </svg>
                  <span>{{ settingsStore.school_name }}</span>
                </div>

                <h2 class="text-xl sm:text-2xl lg:text-3xl font-extrabold leading-snug text-white">
                  {{ settingsStore.tagline || 'Membentuk Generasi Unggul, Berakhlak Mulia, dan Berprestasi.' }}
                </h2>

                <p class="text-emerald-100/80 text-xs sm:text-sm leading-relaxed pt-1 hidden sm:block">
                  Platform komprehensif penghubung civitas lembaga, dewan guru, siswa, dan wali murid dalam satu ekosistem digital terintegrasi.
                </p>
              </div>
            </div>

            <!-- Feature Badges -->
            <div class="relative z-10 pt-6 sm:pt-10 mt-auto">
              <p class="text-[11px] sm:text-xs uppercase tracking-wider text-emerald-200/80 font-semibold mb-2.5">
                Layanan Terpadu Lembaga:
              </p>
              <div class="flex flex-wrap gap-2 overflow-x-auto no-scrollbar pb-1">
                <span class="glass-pill text-[11px] sm:text-xs font-medium px-3 py-1 rounded-full flex items-center gap-1.5 text-white shadow-sm whitespace-nowrap">
                  <svg class="w-3.5 h-3.5 text-emerald-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                  Portal Akademik
                </span>
                <span class="glass-pill text-[11px] sm:text-xs font-medium px-3 py-1 rounded-full flex items-center gap-1.5 text-white shadow-sm whitespace-nowrap">
                  <svg class="w-3.5 h-3.5 text-emerald-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path></svg>
                  Tahfidz & Rapor
                </span>
                <span class="glass-pill text-[11px] sm:text-xs font-medium px-3 py-1 rounded-full flex items-center gap-1.5 text-white shadow-sm whitespace-nowrap">
                  <svg class="w-3.5 h-3.5 text-emerald-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                  E-Presensi & Izin
                </span>
              </div>
            </div>
          </section>

          <!-- ==================================================== -->
          <!-- LOGIN FORM PANEL (RIGHT DESKTOP / BOTTOM MOBILE)     -->
          <!-- ==================================================== -->
          <section class="lg:w-7/12 p-6 sm:p-10 lg:p-14 flex flex-col justify-between bg-white rounded-t-[32px] lg:rounded-none -mt-4 lg:mt-0 relative z-20 shadow-sheet lg:shadow-none">
            
            <!-- Mobile Pull Bar Indicator -->
            <div aria-hidden="true" class="w-10 h-1 bg-slate-200 rounded-full mx-auto mb-4 lg:hidden"></div>

            <div class="max-w-md mx-auto w-full">
              
              <!-- Greeting & Subtitle -->
              <div class="mb-5 sm:mb-6">
                <div class="flex items-center gap-2 mb-1.5">
                  <h2 class="text-xl sm:text-2xl font-extrabold text-slate-900 tracking-tight">
                    Selamat Datang Kembali
                  </h2>
                  <span class="text-xl sm:text-2xl">👋</span>
                </div>
                <p class="text-xs sm:text-sm text-slate-500 font-normal">
                  Silakan masuk ke portal sistem informasi akademik & lembaga.
                </p>
              </div>

              <!-- Role Info Callout Banner -->
              <div class="mb-5 p-3.5 rounded-2xl bg-emerald-50/90 border border-emerald-100 flex items-start gap-2.5 shadow-sm">
                <div class="text-emerald-700 bg-emerald-100 rounded-full p-1 shrink-0 mt-0.5">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </div>
                <p class="text-[11px] sm:text-xs text-emerald-900 leading-relaxed font-medium">
                  Gunakan <strong class="font-bold text-emerald-950">Email</strong> untuk Pegawai/Guru atau <strong class="font-bold text-emerald-950">No. WhatsApp/HP</strong> untuk Wali Murid & Siswa.
                </p>
              </div>

              <!-- Alert Error Notification -->
              <div
                v-if="errorMessage"
                class="mb-4 p-3.5 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 text-xs font-semibold flex items-center gap-2.5 animate-shake"
              >
                <ion-icon :icon="alertCircleOutline" class="text-lg shrink-0 text-rose-500"></ion-icon>
                <span>{{ errorMessage }}</span>
              </div>

              <!-- Form Inputs -->
              <form @submit.prevent="handleLogin" class="space-y-4">
                
                <!-- Input 1: Identifier (Email / No. HP) -->
                <div>
                  <div class="flex items-center justify-between mb-1.5">
                    <label class="block text-[11px] sm:text-xs font-bold uppercase tracking-wider text-slate-700">
                      Email atau No. Handphone/WA
                    </label>
                    <span
                      v-if="detectedRole"
                      class="text-[10px] font-bold px-2 py-0.5 rounded-full"
                      :class="detectedRole === 'wali' ? 'bg-emerald-100 text-emerald-800' : 'bg-sky-100 text-sky-800'"
                    >
                      {{ detectedRole === 'wali' ? 'Wali Murid / Siswa' : 'Pegawai / Guru' }}
                    </span>
                  </div>

                  <div class="relative rounded-xl shadow-sm">
                    <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.206"></path>
                      </svg>
                    </div>
                    <input
                      v-model="identifier"
                      type="text"
                      required
                      placeholder="admin@sekolah.sch.id atau 08123456789"
                      class="block w-full pl-11 pr-4 py-3 text-xs sm:text-sm text-slate-900 placeholder:text-slate-400 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition duration-150"
                    />
                  </div>
                </div>

                <!-- Input 2: Kata Sandi (Password) -->
                <div>
                  <label class="block text-[11px] sm:text-xs font-bold uppercase tracking-wider text-slate-700 mb-1.5">
                    Kata Sandi
                  </label>

                  <div class="relative rounded-xl shadow-sm">
                    <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
                      </svg>
                    </div>

                    <input
                      v-model="password"
                      :type="showPassword ? 'text' : 'password'"
                      required
                      placeholder="Masukkan kata sandi"
                      class="block w-full pl-11 pr-12 py-3 text-xs sm:text-sm text-slate-900 placeholder:text-slate-400 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition duration-150"
                    />

                    <!-- Toggle Show/Hide Password Button -->
                    <button
                      type="button"
                      @click="showPassword = !showPassword"
                      class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none transition-colors"
                      title="Lihat / Sembunyikan Kata Sandi"
                    >
                      <svg v-if="!showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                      </svg>
                      <svg v-else class="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l18 18"></path>
                      </svg>
                    </button>
                  </div>
                </div>

                <!-- Remember Me & Forgot Password -->
                <div class="flex items-center justify-between text-xs pt-0.5">
                  <label class="flex items-center gap-2 cursor-pointer select-none">
                    <input
                      v-model="rememberMe"
                      type="checkbox"
                      class="w-4 h-4 text-emerald-600 border-slate-300 rounded focus:ring-emerald-500 focus:ring-offset-0 transition cursor-pointer"
                    />
                    <span class="text-slate-600 font-medium">Ingat Saya</span>
                  </label>

                  <button
                    type="button"
                    @click="handleForgotPassword"
                    class="font-bold text-emerald-700 hover:text-emerald-800 hover:underline transition"
                  >
                    Lupa Kata Sandi?
                  </button>
                </div>

                <!-- Primary Submit CTA Button -->
                <div class="pt-2">
                  <button
                    type="submit"
                    :disabled="loading"
                    class="w-full py-3.5 px-6 rounded-xl font-bold text-xs sm:text-sm text-white bg-gradient-to-r from-emerald-600 to-emerald-500 hover:from-emerald-700 hover:to-emerald-600 active:scale-[0.99] shadow-lg shadow-emerald-600/25 hover:shadow-emerald-600/35 transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-60 disabled:cursor-not-allowed"
                  >
                    <span v-if="loading" class="flex items-center gap-2">
                      <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      Memeriksa akun...
                    </span>
                    <span v-else class="flex items-center gap-2">
                      <span>Masuk ke Akun</span>
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path>
                      </svg>
                    </span>
                  </button>
                </div>
              </form>

              <!-- Divider -->
              <div class="relative my-5 sm:my-6">
                <div class="absolute inset-0 flex items-center">
                  <div class="w-full border-t border-slate-200"></div>
                </div>
                <div class="relative flex justify-center text-xs">
                  <span class="bg-white px-3 text-slate-400 font-medium">Atau butuh bantuan?</span>
                </div>
              </div>

              <!-- Helpdesk WhatsApp Action Button -->
              <div>
                <a
                  :href="helpdeskWhatsappUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="w-full flex items-center justify-center gap-2.5 py-3 px-4 rounded-xl border border-slate-200 hover:border-emerald-200 hover:bg-emerald-50/50 text-slate-700 hover:text-emerald-800 text-xs font-semibold transition duration-150 shadow-sm"
                >
                  <!-- WhatsApp Vector Icon -->
                  <svg class="w-4 h-4 text-emerald-600 shrink-0" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12.031 2c-5.508 0-9.986 4.475-9.986 9.982 0 1.761.46 3.479 1.334 4.996L2 22l5.176-1.357a9.988 9.988 0 004.855 1.258h.004c5.508 0 9.987-4.475 9.987-9.982 0-2.666-1.038-5.172-2.926-7.06A9.91 9.91 0 0012.031 2zm0 18.272a8.27 8.27 0 01-4.223-1.156l-.303-.18-3.14.823.839-3.061-.198-.314a8.26 8.26 0 01-1.272-4.402c0-4.57 3.719-8.289 8.297-8.289 2.215 0 4.297.863 5.863 2.43 1.566 1.567 2.428 3.65 2.428 5.867 0 4.571-3.719 8.282-8.291 8.282zm4.545-6.208c-.249-.125-1.474-.727-1.703-.81-.229-.083-.396-.125-.562.125-.166.25-.645.81-.79 1-.146.188-.292.209-.541.084-.25-.125-1.055-.389-2.01-1.24a7.514 7.514 0 01-1.391-1.729c-.146-.25-.016-.385.109-.51.112-.112.25-.292.375-.438.125-.146.166-.25.25-.417.083-.166.042-.312-.021-.437-.063-.125-.562-1.354-.77-1.854-.203-.487-.41-.421-.563-.429-.145-.007-.312-.008-.479-.008-.166 0-.437.063-.666.312-.229.25-.875.854-.875 2.083 0 1.229.896 2.417 1.021 2.584.125.166 1.763 2.692 4.27 3.776.597.258 1.063.412 1.427.527.6.19 1.146.163 1.577.099.481-.072 1.474-.603 1.682-1.186.208-.583.208-1.083.146-1.187-.063-.104-.229-.167-.479-.292z"></path>
                  </svg>
                  <span>Hubungi Helpdesk / Admin Lembaga</span>
                </a>
              </div>

            </div>

            <!-- Footer Copyright -->
            <footer class="mt-6 pt-4 border-t border-slate-100 text-center">
              <p class="text-[11px] sm:text-xs text-slate-400">
                © {{ new Date().getFullYear() }} {{ settingsStore.school_name || 'SMAS AL-QURAN RQ4' }}. Seluruh hak cipta dilindungi.
              </p>
            </footer>

          </section>

        </main>

      </div>
    </ion-content>
  </ion-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import { schoolOutline, alertCircleOutline } from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { showAlert } from '@/utils/alert'

const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const router = useRouter()

const identifier = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = ref(false)
const showPassword = ref(false)
const rememberMe = ref(false)

// Smart Role Detection
const detectedRole = computed(() => {
  const clean = identifier.value.trim()
  if (!clean) return null
  if (/^[0-9+]+$/.test(clean)) return 'wali'
  if (clean.includes('@')) return 'pegawai'
  return null
})

// Helpdesk WhatsApp URL
const helpdeskWhatsappUrl = computed(() => {
  let p = (settingsStore.phone || '').replace(/[^0-9]/g, '')
  if (p.startsWith('0')) {
    p = '62' + p.slice(1)
  }
  if (!p) p = '6281234567890'
  const text = encodeURIComponent(
    `Halo Admin / Tata Usaha ${settingsStore.school_name || 'Lembaga'}, saya butuh bantuan terkait akun login di ${settingsStore.app_name || 'SMARTSIS'}.`
  )
  return `https://wa.me/${p}?text=${text}`
})

const handleForgotPassword = () => {
  showAlert({
    title: 'Lupa Kata Sandi?',
    text: `Untuk menjaga keamanan data, pengaturan ulang (*reset*) kata sandi akun dilakukan oleh Administrator / Tata Usaha lembaga.<br/><br/>Silakan hubungi pihak Tata Usaha atau gunakan tombol bantuan WhatsApp di bawah.`,
    confirmButtonText: 'Tutup'
  })
}

const handleLogin = async () => {
  errorMessage.value = ''
  loading.value = true

  const res = await authStore.login(identifier.value, password.value)
  loading.value = false

  if (res.success) {
    if (rememberMe.value) {
      localStorage.setItem('remembered_identifier', identifier.value)
    } else {
      localStorage.removeItem('remembered_identifier')
    }
    router.replace('/tabs/dashboard')
  } else {
    errorMessage.value = res.message || 'Login gagal, periksa email/no. HP dan kata sandi Anda.'
  }
}

onMounted(() => {
  settingsStore.fetchSettings()
  const savedId = localStorage.getItem('remembered_identifier')
  if (savedId) {
    identifier.value = savedId
    rememberMe.value = true
  }
})
</script>

<style scoped>
/* Islamic & Dot Pattern matching the Stitch design */
.islamic-pattern {
  background-color: #064e3b;
  background-image: radial-gradient(circle at 100% 150%, #047857 24%, #064e3b 50%, #022c22 85%);
}

.pattern-overlay {
  background-image: radial-gradient(rgba(255, 255, 255, 0.12) 1.5px, transparent 0);
  background-size: 24px 24px;
}

.glass-pill {
  background: rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.shadow-sheet {
  box-shadow: 0 -10px 30px rgba(0, 0, 0, 0.08);
}

.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
