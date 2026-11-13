<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-slate-50">
      <div class="max-w-3xl mx-auto py-6">

        <!-- Header Halaman -->
        <div class="mb-6">
          <h1 class="text-xl font-bold text-slate-800">Identitas & Pengaturan Lembaga</h1>
          <p class="text-xs text-slate-500 mt-1">
            Ubah nama aplikasi, nama resmi sekolah/pesantren, kontak, dan logo yang tampil di sistem.
          </p>
        </div>

        <!-- Alert Notifikasi -->
        <div v-if="successMsg" class="mb-4 p-3.5 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs flex items-center gap-2">
          <ion-icon :icon="checkmarkCircleOutline" class="text-lg shrink-0"></ion-icon>
          <span>{{ successMsg }}</span>
        </div>
        <div v-if="errorMsg" class="mb-4 p-3.5 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 text-xs flex items-center gap-2">
          <ion-icon :icon="alertCircleOutline" class="text-lg shrink-0"></ion-icon>
          <span>{{ errorMsg }}</span>
        </div>

        <!-- Form Pengaturan -->
        <form @submit.prevent="handleSave" class="bg-white rounded-2xl p-6 border border-slate-100 shadow-sm space-y-5">
          
          <!-- Logo & Preview -->
          <div class="border-b border-slate-100 pb-5">
            <label class="block text-xs font-semibold text-slate-700 mb-2">Logo Lembaga</label>
            <div class="flex items-center gap-4">
              <div class="w-16 h-16 rounded-2xl bg-emerald-50 border border-emerald-100 flex items-center justify-center overflow-hidden shrink-0">
                <img v-if="form.logo_url" :src="form.logo_url" alt="Logo" class="w-full h-full object-contain p-1" />
                <ion-icon v-else :icon="schoolOutline" class="text-3xl text-emerald-600"></ion-icon>
              </div>
              <div class="flex-1">
                <input
                  type="file"
                  accept="image/*"
                  @change="handleFileUpload"
                  class="block w-full text-xs text-slate-500 file:mr-3 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-xs file:font-semibold file:bg-emerald-50 file:text-emerald-700 hover:file:bg-emerald-100 cursor-pointer"
                />
                <p class="text-[11px] text-slate-400 mt-1">Format PNG, JPG, atau SVG (Maks. 1MB)</p>
              </div>
            </div>
          </div>

          <!-- Nama Aplikasi & Nama Resmi Sekolah -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Nama Singkat Aplikasi</label>
              <input
                v-model="form.app_name"
                type="text"
                required
                placeholder="Contoh: SMARTSIS"
                class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all text-slate-800"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Nama Resmi Lembaga / Sekolah</label>
              <input
                v-model="form.school_name"
                type="text"
                required
                placeholder="Contoh: Pondok Pesantren & Madrasah Terpadu"
                class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all text-slate-800"
              />
            </div>
          </div>

          <!-- Tagline / Slogan -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Tagline / Deskripsi Singkat</label>
            <input
              v-model="form.tagline"
              type="text"
              placeholder="Contoh: Sistem Informasi Sekolah & Pesantren Terpadu"
              class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all text-slate-800"
            />
          </div>

          <!-- Kontak (No WhatsApp & Email) -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">No. WhatsApp / Telepon Resmi</label>
              <input
                v-model="form.phone"
                type="text"
                placeholder="08123456789"
                class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all text-slate-800"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Email Resmi</label>
              <input
                v-model="form.email"
                type="email"
                placeholder="info@sekolah.sch.id"
                class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all text-slate-800"
              />
            </div>
          </div>

          <!-- Alamat Lembaga -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Alamat Lembaga</label>
            <textarea
              v-model="form.address"
              rows="2"
              placeholder="Jl. Pesantren No. 01, Kompleks Pendidikan"
              class="w-full px-3.5 py-2.5 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all text-slate-800"
            ></textarea>
          </div>

          <!-- Tombol Simpan -->
          <div class="pt-3 border-t border-slate-100 flex justify-end">
            <button
              type="submit"
              :disabled="saving"
              class="px-5 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white font-medium text-xs shadow-md shadow-emerald-500/20 transition-all flex items-center gap-2 disabled:opacity-60"
            >
              <ion-icon :icon="saveOutline" class="text-base"></ion-icon>
              <span>{{ saving ? 'Menyimpan...' : 'Simpan Pengaturan' }}</span>
            </button>
          </div>

        </form>

      </div>
    </ion-content>
  </ion-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import { schoolOutline, saveOutline, checkmarkCircleOutline, alertCircleOutline } from 'ionicons/icons'
import { useSettingsStore } from '@/stores/settings'
import { useAuthStore } from '@/stores/auth'

const settingsStore = useSettingsStore()
const authStore = useAuthStore()

const saving = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const form = ref({
  app_name: '',
  school_name: '',
  tagline: '',
  phone: '',
  email: '',
  address: '',
  logo_url: ''
})

onMounted(async () => {
  await settingsStore.fetchSettings()
  form.value = {
    app_name: settingsStore.app_name,
    school_name: settingsStore.school_name,
    tagline: settingsStore.tagline,
    phone: settingsStore.phone,
    email: settingsStore.email,
    address: settingsStore.address,
    logo_url: settingsStore.logo_url
  }
})

const handleFileUpload = (e) => {
  const file = e.target.files[0]
  if (!file) return

  if (file.size > 1024 * 1024) {
    errorMsg.value = 'Ukuran file logo maksimal 1MB'
    return
  }

  const reader = new FileReader()
  reader.onload = (event) => {
    form.value.logo_url = event.target.result
  }
  reader.readAsDataURL(file)
}

const handleSave = async () => {
  saving.value = true
  successMsg.value = ''
  errorMsg.value = ''

  const res = await settingsStore.saveSettings(form.value, authStore.token)
  saving.value = false

  if (res.success) {
    successMsg.value = 'Identitas lembaga berhasil diperbarui!'
    setTimeout(() => { successMsg.value = '' }, 4000)
  } else {
    errorMsg.value = res.message
  }
}
</script>
