<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc]">
      <div class="max-w-6xl mx-auto py-2">

        <!-- Banner Selamat Datang Bergaya Modern (Sesuai Gambar Desain) -->
        <div class="relative bg-white rounded-2xl p-6 sm:p-7 border border-slate-100 shadow-sm overflow-hidden mb-6 flex flex-col md:flex-row md:items-center justify-between gap-4">
          <!-- Hiasan Gradien Langit Biru di Belakang -->
          <div class="absolute inset-0 bg-gradient-to-r from-sky-50/80 via-blue-50/50 to-sky-100/60 pointer-events-none"></div>

          <div class="relative z-10 max-w-xl">
            <!-- Jika Akun Wali -->
            <div v-if="authStore.isWali && authStore.activeSiswa">
              <div class="flex items-center gap-2 mb-2">
                <span class="inline-block text-[11px] font-bold text-amber-900 bg-amber-400 px-3 py-1 rounded-full shadow-sm">
                  Siswa Aktif
                </span>
                <select
                  v-if="authStore.daftarSiswa.length > 1"
                  :value="authStore.activeSiswa.id"
                  @change="onSwitchSiswa($event)"
                  class="text-[11px] font-semibold text-slate-800 bg-white border border-slate-200 rounded-full px-3 py-1 focus:outline-none shadow-sm cursor-pointer"
                >
                  <option v-for="s in authStore.daftarSiswa" :key="s.id" :value="s.id">
                    Pilih: {{ s.nama_lengkap }} ({{ s.kelas }})
                  </option>
                </select>
              </div>
              <h1 class="text-2xl font-black text-slate-800 tracking-tight">
                {{ authStore.activeSiswa.nama_lengkap }}
              </h1>
              <p class="text-xs text-slate-600 mt-1 font-medium">
                NIS: {{ authStore.activeSiswa.nis }} • Kelas: {{ authStore.activeSiswa.kelas }} • {{ authStore.activeSiswa.hubungan }}
              </p>
            </div>

            <!-- Jika Akun Pegawai / Admin -->
            <div v-else>
              <h1 class="text-2xl font-black text-slate-800 tracking-tight flex items-center gap-2">
                Halo, {{ authStore.user?.name }} <span>👋</span>
              </h1>
              <div class="flex items-center gap-2 mt-2">
                <span class="text-xs text-slate-500 font-medium">Masuk sebagai</span>
                <span class="inline-block px-3 py-0.5 rounded-full text-xs font-bold bg-amber-400 text-amber-950 shadow-sm">
                  {{ authStore.roleLabel }}
                </span>
              </div>
            </div>
          </div>

          <!-- Ilustrasi Gedung / Banner Sekolah di Kanan -->
          <div class="relative z-10 hidden sm:flex items-center justify-end shrink-0">
            <div class="w-40 h-24 rounded-2xl bg-gradient-to-tr from-sky-400 to-indigo-500 text-white flex flex-col items-center justify-center p-3 text-center shadow-md shadow-sky-400/20">
              <ion-icon :icon="schoolOutline" class="text-4xl mb-1"></ion-icon>
              <span class="text-[10px] font-bold tracking-wider uppercase opacity-90">SIAKAD TERPADU</span>
            </div>
          </div>
        </div>

        <!-- 4 Statistik Card Berwarna-Warni Lembut (Soft Pastel Cards Sesuai Desain) -->
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          
          <!-- Card 1: Presensi (Biru Muda) -->
          <div class="bg-gradient-to-br from-blue-50/70 to-sky-100/40 p-4 sm:p-5 rounded-2xl border border-blue-200/60 shadow-sm flex items-center justify-between">
            <div>
              <span class="text-[11px] font-semibold text-slate-500 block">Presensi Bulan Ini</span>
              <p class="text-2xl font-black text-slate-800 mt-0.5">98%</p>
              <span class="inline-block mt-1 text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-700">
                Sangat Baik
              </span>
            </div>
            <div class="w-12 h-12 rounded-2xl bg-sky-200/80 text-sky-600 flex items-center justify-center text-2xl shadow-inner shrink-0">
              <ion-icon :icon="calendarOutline"></ion-icon>
            </div>
          </div>

          <!-- Card 2: Tagihan SPP (Hijau Muda) -->
          <div class="bg-gradient-to-br from-emerald-50/70 to-teal-100/40 p-4 sm:p-5 rounded-2xl border border-emerald-200/60 shadow-sm flex items-center justify-between">
            <div>
              <span class="text-[11px] font-semibold text-slate-500 block">Status Tagihan SPP</span>
              <p class="text-xl font-black text-slate-800 mt-0.5">Lunas</p>
              <span class="inline-block mt-1 text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-700">
                Bulan Ini
              </span>
            </div>
            <div class="w-12 h-12 rounded-2xl bg-emerald-200/80 text-emerald-600 flex items-center justify-center text-2xl shadow-inner shrink-0">
              <ion-icon :icon="walletOutline"></ion-icon>
            </div>
          </div>

          <!-- Card 3: Nilai Rata-rata (Ungu Muda) -->
          <div class="bg-gradient-to-br from-purple-50/70 to-fuchsia-100/40 p-4 sm:p-5 rounded-2xl border border-purple-200/60 shadow-sm flex items-center justify-between">
            <div>
              <span class="text-[11px] font-semibold text-slate-500 block">Rata-rata Nilai</span>
              <p class="text-2xl font-black text-slate-800 mt-0.5">87.5</p>
              <span class="inline-block mt-1 text-[10px] font-bold px-2 py-0.5 rounded-full bg-purple-100 text-purple-700">
                Semester Ganjil
              </span>
            </div>
            <div class="w-12 h-12 rounded-2xl bg-purple-200/80 text-purple-600 flex items-center justify-center text-2xl shadow-inner shrink-0">
              <ion-icon :icon="ribbonOutline"></ion-icon>
            </div>
          </div>

          <!-- Card 4: Agenda Terdekat (Kuning / Oranye Muda) -->
          <div class="bg-gradient-to-br from-amber-50/70 to-orange-100/40 p-4 sm:p-5 rounded-2xl border border-amber-200/60 shadow-sm flex items-center justify-between">
            <div>
              <span class="text-[11px] font-semibold text-slate-500 block">Agenda Terdekat</span>
              <p class="text-lg font-black text-slate-800 mt-0.5">Ujian PAS</p>
              <span class="inline-block mt-1 text-[10px] font-bold px-2 py-0.5 rounded-full bg-amber-100 text-amber-800">
                2 Minggu Lagi
              </span>
            </div>
            <div class="w-12 h-12 rounded-2xl bg-amber-200/80 text-amber-600 flex items-center justify-center text-2xl shadow-inner shrink-0">
              <ion-icon :icon="timeOutline"></ion-icon>
            </div>
          </div>

        </div>

        <!-- Bagian Layanan Informasi (Tombol Card Berwarna dengan Panah Kanan Sesuai Gambar) -->
        <div class="mb-4 flex items-center gap-2">
          <div class="w-5 h-5 rounded-lg bg-amber-400 text-amber-950 flex items-center justify-center text-xs font-bold shadow-sm">
            <ion-icon :icon="gridOutline"></ion-icon>
          </div>
          <h2 class="text-base font-bold text-slate-800 tracking-tight">Layanan Informasi</h2>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div
            v-for="menu in authStore.menus"
            :key="menu.code"
            @click="navigateTo(menu.path)"
            class="bg-white p-4 rounded-2xl border border-slate-200/70 shadow-sm hover:shadow-md hover:border-slate-300 transition-all cursor-pointer flex items-center justify-between group"
          >
            <div class="flex items-center gap-3.5">
              <!-- Kotak Icon Berwarna Pastel Unik per Menu -->
              <div
                class="w-12 h-12 rounded-2xl flex items-center justify-center text-xl shrink-0 transition-transform group-hover:scale-105"
                :class="getMenuTheme(menu.code)"
              >
                <ion-icon :icon="getIcon(menu.icon)"></ion-icon>
              </div>
              <span class="text-xs font-bold text-slate-800 group-hover:text-blue-600 transition-colors">
                {{ menu.title }}
              </span>
            </div>

            <!-- Ikon Panah Kanan (Chevron) Sesuai Gambar -->
            <ion-icon :icon="chevronForwardOutline" class="text-slate-400 group-hover:text-slate-700 text-sm group-hover:translate-x-0.5 transition-all"></ion-icon>
          </div>
        </div>

      </div>
    </ion-content>
  </ion-page>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import {
  homeOutline, calendarOutline, ribbonOutline, walletOutline,
  peopleOutline, briefcaseOutline, shieldCheckmarkOutline, personOutline,
  settingsOutline, timeOutline, gridOutline, chevronForwardOutline, schoolOutline
} from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const iconMap = {
  'home-outline': homeOutline,
  'calendar-outline': calendarOutline,
  'ribbon-outline': ribbonOutline,
  'wallet-outline': walletOutline,
  'people-outline': peopleOutline,
  'briefcase-outline': briefcaseOutline,
  'shield-checkmark-outline': shieldCheckmarkOutline,
  'settings-outline': settingsOutline,
  'person-outline': personOutline
}

const getIcon = (name) => iconMap[name] || homeOutline

// Warna pastel khas tiap menu sesuai desain di gambar
const getMenuTheme = (code) => {
  switch (code) {
    case 'dashboard':
      return 'bg-blue-100 text-blue-600'
    case 'presensi':
      return 'bg-emerald-100 text-emerald-600'
    case 'nilai':
      return 'bg-purple-100 text-purple-600'
    case 'spp':
      return 'bg-amber-100 text-amber-600'
    case 'siswa_mgmt':
      return 'bg-sky-100 text-sky-600'
    case 'pegawai_mgmt':
      return 'bg-rose-100 text-rose-600'
    case 'role_mgmt':
      return 'bg-teal-100 text-teal-600'
    case 'pengaturan':
      return 'bg-indigo-100 text-indigo-600'
    case 'profil':
      return 'bg-orange-100 text-orange-600'
    default:
      return 'bg-slate-100 text-slate-600'
  }
}

const navigateTo = (path) => {
  router.push(path)
}

const onSwitchSiswa = (e) => {
  const selectedId = parseInt(e.target.value)
  const target = authStore.daftarSiswa.find(s => s.id === selectedId)
  if (target) {
    authStore.setActiveSiswa(target)
  }
}
</script>
