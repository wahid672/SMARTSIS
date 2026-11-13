<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-slate-50">
      <div class="max-w-4xl mx-auto py-4">

        <!-- Banner Selamat Datang & Status Siswa (Untuk Wali) -->
        <div v-if="authStore.isWali && authStore.activeSiswa" class="bg-gradient-to-r from-emerald-600 to-teal-600 rounded-2xl p-5 text-white shadow-md shadow-emerald-600/10 mb-6">
          <div class="flex items-center justify-between gap-3">
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <span class="inline-block text-[11px] bg-white/20 backdrop-blur px-2.5 py-0.5 rounded-full font-medium">
                  Santri / Siswa Aktif
                </span>
                <!-- Switcher Cepat jika wali punya > 1 anak -->
                <select
                  v-if="authStore.daftarSiswa.length > 1"
                  :value="authStore.activeSiswa.id"
                  @change="onSwitchSiswa($event)"
                  class="text-[11px] font-semibold text-emerald-900 bg-white/90 rounded-full px-2.5 py-0.5 focus:outline-none cursor-pointer"
                >
                  <option v-for="s in authStore.daftarSiswa" :key="s.id" :value="s.id">
                    Pilih: {{ s.nama_lengkap }} ({{ s.kelas }})
                  </option>
                </select>
              </div>
              <h1 class="text-xl font-bold">{{ authStore.activeSiswa.nama_lengkap }}</h1>
              <p class="text-xs text-emerald-100 mt-0.5">
                NIS: {{ authStore.activeSiswa.nis }} • Kelas: {{ authStore.activeSiswa.kelas }} • {{ authStore.activeSiswa.hubungan }}
              </p>
            </div>
            <div class="w-12 h-12 rounded-full bg-white/20 backdrop-blur flex items-center justify-center text-xl font-bold shrink-0">
              {{ authStore.activeSiswa.nama_lengkap.charAt(0) }}
            </div>
          </div>
        </div>

        <!-- Banner Pegawai / Guru / Admin -->
        <div v-else class="bg-white rounded-2xl p-5 border border-slate-100 shadow-sm mb-6">
          <h1 class="text-xl font-bold text-slate-800">Halo, {{ authStore.user?.name }} 👋</h1>
          <p class="text-xs text-slate-500 mt-1">
            Masuk sebagai <span class="font-semibold text-brand-600">{{ authStore.roleLabel }}</span>
          </p>
        </div>

        <!-- Statistik Ringkas / Dashboard Card -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3 sm:gap-4 mb-6">
          
          <div class="bg-white p-4 rounded-xl border border-slate-100 shadow-sm">
            <span class="text-[11px] text-slate-400 font-medium">Presensi Bulan Ini</span>
            <p class="text-xl font-bold text-slate-800 mt-1">98%</p>
            <span class="text-[10px] text-emerald-600 font-medium">Sangat Baik</span>
          </div>

          <div class="bg-white p-4 rounded-xl border border-slate-100 shadow-sm">
            <span class="text-[11px] text-slate-400 font-medium">Status Tagihan SPP</span>
            <p class="text-xl font-bold text-slate-800 mt-1">Lunas</p>
            <span class="text-[10px] text-emerald-600 font-medium">November 2026</span>
          </div>

          <div class="bg-white p-4 rounded-xl border border-slate-100 shadow-sm">
            <span class="text-[11px] text-slate-400 font-medium">Rata-rata Nilai</span>
            <p class="text-xl font-bold text-slate-800 mt-1">87.5</p>
            <span class="text-[10px] text-blue-600 font-medium">Semester Ganjil</span>
          </div>

          <div class="bg-white p-4 rounded-xl border border-slate-100 shadow-sm">
            <span class="text-[11px] text-slate-400 font-medium">Agenda Terdekat</span>
            <p class="text-xl font-bold text-slate-800 mt-1">Ujian PAS</p>
            <span class="text-[10px] text-amber-600 font-medium">2 Minggu Lagi</span>
          </div>

        </div>

        <!-- Menu Akses Cepat / Fitur Sekolah -->
        <h2 class="text-sm font-bold text-slate-800 mb-3">Layanan Informasi</h2>
        <div class="grid grid-cols-3 sm:grid-cols-4 gap-3">
          
          <div
            v-for="menu in authStore.menus"
            :key="menu.code"
            @click="navigateTo(menu.path)"
            class="bg-white p-4 rounded-2xl border border-slate-100 text-center flex flex-col items-center justify-center cursor-pointer hover:border-brand-500 hover:shadow-sm transition-all"
          >
            <div class="w-11 h-11 rounded-xl bg-emerald-50 text-brand-600 flex items-center justify-center mb-2">
              <ion-icon :icon="getIcon(menu.icon)" class="text-2xl"></ion-icon>
            </div>
            <span class="text-xs font-semibold text-slate-700">{{ menu.title }}</span>
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
  peopleOutline, briefcaseOutline, shieldCheckmarkOutline, personOutline
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
  'person-outline': personOutline
}

const getIcon = (name) => iconMap[name] || homeOutline

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
