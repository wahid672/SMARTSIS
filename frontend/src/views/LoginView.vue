<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-slate-50">
      <div class="min-h-full flex flex-col justify-center items-center py-8 px-4">
        
        <!-- Header / Logo -->
        <div class="text-center mb-8">
          <div class="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-brand-500 text-white shadow-lg shadow-emerald-200 mb-3">
            <ion-icon :icon="schoolOutline" class="text-4xl"></ion-icon>
          </div>
          <h1 class="text-2xl font-bold tracking-tight text-slate-800">SMARTSIS</h1>
          <p class="text-sm text-slate-500 mt-1">Sistem Informasi Sekolah Terpadu</p>
        </div>

        <!-- Kartu Login Pintar -->
        <div class="w-full max-w-md bg-white rounded-2xl shadow-sm border border-slate-100 p-6 sm:p-8">
          
          <div class="mb-6">
            <h2 class="text-lg font-semibold text-slate-800">Selamat Datang</h2>
            <p class="text-xs text-slate-500 mt-0.5">
              Gunakan <span class="font-medium text-slate-700">Email</span> (Pegawai/Guru) atau <span class="font-medium text-slate-700">No. WhatsApp/HP</span> (Wali Murid)
            </p>
          </div>

          <!-- Alert Error -->
          <div v-if="errorMessage" class="mb-4 p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 text-xs flex items-center gap-2">
            <ion-icon :icon="alertCircleOutline" class="text-base shrink-0"></ion-icon>
            <span>{{ errorMessage }}</span>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-4">
            
            <!-- Input Pintar (Email / No HP) -->
            <div>
              <label class="block text-xs font-medium text-slate-700 mb-1">
                Email atau No. Handphone/WA
              </label>
              <div class="relative">
                <input
                  v-model="identifier"
                  type="text"
                  required
                  placeholder="admin@sekolah.sch.id atau 08123456789"
                  class="w-full px-3.5 py-2.5 text-sm rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all text-slate-800 placeholder-slate-400"
                />
              </div>
            </div>

            <!-- Input Password -->
            <div>
              <label class="block text-xs font-medium text-slate-700 mb-1">
                Kata Sandi
              </label>
              <input
                v-model="password"
                type="password"
                required
                placeholder="Masukkan kata sandi"
                class="w-full px-3.5 py-2.5 text-sm rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all text-slate-800 placeholder-slate-400"
              />
            </div>

            <!-- Tombol Masuk -->
            <button
              type="submit"
              :disabled="loading"
              class="w-full mt-2 py-3 px-4 rounded-xl bg-brand-500 hover:bg-brand-600 active:scale-[0.99] text-white font-medium text-sm shadow-md shadow-emerald-500/20 transition-all flex items-center justify-center gap-2 disabled:opacity-60"
            >
              <span v-if="loading">Memeriksa akun...</span>
              <span v-else>Masuk ke Akun</span>
            </button>
          </form>

          <!-- Hint akun demo untuk kemudahan pengujian -->
          <div class="mt-6 pt-5 border-t border-slate-100 text-[11px] text-slate-400 space-y-1">
            <p class="font-medium text-slate-600">Demo Akun Cepat:</p>
            <p>• Pegawai/Admin: <code class="text-slate-600 bg-slate-100 px-1 py-0.5 rounded">admin@sekolah.sch.id</code> / <code class="text-slate-600 bg-slate-100 px-1 py-0.5 rounded">admin123</code></p>
            <p>• Wali Santri: <code class="text-slate-600 bg-slate-100 px-1 py-0.5 rounded">081234567890</code> / <code class="text-slate-600 bg-slate-100 px-1 py-0.5 rounded">123456</code></p>
          </div>

        </div>

      </div>
    </ion-content>
  </ion-page>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import { schoolOutline, alertCircleOutline } from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const identifier = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = ref(false)

const handleLogin = async () => {
  errorMessage.value = ''
  loading.value = true

  const res = await authStore.login(identifier.value, password.value)
  loading.value = false

  if (res.success) {
    router.replace('/tabs/dashboard')
  } else {
    errorMessage.value = res.message
  }
}
</script>
