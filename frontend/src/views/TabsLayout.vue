<template>
  <ion-page>
    <ion-split-pane content-id="main-content" when="md">
      
      <!-- SIDEBAR UNTUK TABLET & DESKTOP (MD/LG SCREEN) -->
      <ion-menu content-id="main-content" type="overlay">
        <ion-content class="bg-white">
          <div class="p-6 border-b border-slate-100 flex items-center gap-3">
            <div class="w-10 h-10 flex items-center justify-center shrink-0">
              <img v-if="settingsStore.logo_url" :src="settingsStore.logo_url" alt="Logo" class="w-full h-full object-contain" />
              <div v-else class="w-10 h-10 rounded-xl bg-brand-500 text-white flex items-center justify-center font-bold text-lg shadow-sm">
                {{ settingsStore.app_name.charAt(0) }}
              </div>
            </div>
            <div class="min-w-0">
              <h2 class="font-bold text-slate-800 text-sm leading-tight truncate">{{ settingsStore.app_name }}</h2>
              <p class="text-[10px] text-slate-400 truncate">{{ settingsStore.school_name }}</p>
              <span class="inline-block text-[10px] text-brand-600 font-semibold mt-0.5">{{ authStore.roleLabel }}</span>
            </div>
          </div>

          <!-- Switcher Siswa di Sidebar untuk Akun Wali -->
          <div v-if="authStore.isWali && authStore.daftarSiswa.length > 0" class="p-4 mx-4 mt-4 bg-emerald-50 rounded-xl border border-emerald-100">
            <label class="text-[10px] font-semibold text-emerald-800 uppercase tracking-wider block mb-1">
              Data Siswa Aktif:
            </label>
            <select
              :value="authStore.activeSiswa?.id"
              @change="onSwitchSiswa($event)"
              class="w-full text-xs font-medium text-slate-800 bg-white border border-emerald-200 rounded-lg p-2 focus:outline-none"
            >
              <option v-for="s in authStore.daftarSiswa" :key="s.id" :value="s.id">
                {{ s.nama_lengkap }} ({{ s.kelas }})
              </option>
            </select>
          </div>

          <!-- Daftar Menu Dinamis di Sidebar -->
          <ion-list lines="none" class="px-3 py-4 space-y-1">
            <ion-item
              v-for="menu in authStore.menus"
              :key="menu.code"
              :router-link="menu.path"
              router-direction="root"
              class="rounded-xl text-slate-700 hover:text-brand-600 cursor-pointer font-medium text-xs"
            >
              <ion-icon slot="start" :icon="getIcon(menu.icon)" class="text-lg text-slate-500"></ion-icon>
              <ion-label>{{ menu.title }}</ion-label>
            </ion-item>

            <ion-item @click="handleLogout" class="rounded-xl text-rose-600 cursor-pointer font-medium text-xs pt-4">
              <ion-icon slot="start" :icon="logOutOutline" class="text-lg text-rose-500"></ion-icon>
              <ion-label>Keluar Akun</ion-label>
            </ion-item>
          </ion-list>
        </ion-content>
      </ion-menu>

      <!-- KONTEN UTAMA & BOTTOM TABS (UNTUK SMARTPHONE) -->
      <div id="main-content" class="ion-page">

        <!-- Mobile Top App Bar (Hanya tampil di smartphone dan BUKAN di Beranda) -->
        <ion-header v-if="!isDashboard" class="ion-no-border shadow-sm bg-white md:hidden">
          <ion-toolbar class="bg-white px-2">
            <ion-buttons slot="start">
              <ion-back-button default-href="/tabs/dashboard" text="" color="dark"></ion-back-button>
            </ion-buttons>
            <ion-title class="text-sm font-bold text-slate-800">
              {{ currentTitle }}
            </ion-title>
          </ion-toolbar>
        </ion-header>

        <!-- Ionic Tabs Shell -->
        <ion-tabs>
          <ion-router-outlet></ion-router-outlet>

          <!-- Bottom Navigation Bar Mobile: Hanya 3 Menu Ringkas & Rapi -->
          <ion-tab-bar slot="bottom" class="border-t border-slate-100 bg-white/95 backdrop-blur md:hidden">
            <ion-tab-button tab="dashboard" href="/tabs/dashboard">
              <ion-icon :icon="homeOutline"></ion-icon>
              <ion-label class="text-[11px] font-medium">Beranda</ion-label>
            </ion-tab-button>

            <ion-tab-button tab="profil" href="/tabs/profil">
              <ion-icon :icon="personOutline"></ion-icon>
              <ion-label class="text-[11px] font-medium">Profil</ion-label>
            </ion-tab-button>

            <ion-tab-button tab="logout" @click="handleLogout">
              <ion-icon :icon="logOutOutline" class="text-rose-500"></ion-icon>
              <ion-label class="text-[11px] font-medium text-rose-500">Keluar</ion-label>
            </ion-tab-button>
          </ion-tab-bar>
        </ion-tabs>

      </div>

    </ion-split-pane>
  </ion-page>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  IonPage, IonSplitPane, IonMenu, IonContent, IonHeader, IonToolbar, IonTitle,
  IonList, IonItem, IonIcon, IonLabel, IonTabs, IonTabBar, IonTabButton, IonRouterOutlet,
  IonButtons, IonMenuButton, IonBackButton
} from '@ionic/vue'
import {
  homeOutline, calendarOutline, ribbonOutline, walletOutline,
  peopleOutline, briefcaseOutline, shieldCheckmarkOutline, personOutline, logOutOutline, settingsOutline
} from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'

const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const router = useRouter()
const route = useRoute()

// Cek apakah halaman saat ini adalah Beranda
const isDashboard = computed(() => {
  return route.path === '/tabs/dashboard' || route.path === '/tabs/' || route.path === '/'
})

// Judul halaman aktif untuk top bar mobile
const currentTitle = computed(() => {
  return route.meta?.title || 'SMARTSIS'
})

// Mapping icon string ke Ionicons
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

// Untuk tampilan HP, ambil maksimal 4-5 menu pertama di Tab Bar
const mobileTabs = computed(() => {
  return authStore.menus.slice(0, 5)
})

const onSwitchSiswa = (e) => {
  const selectedId = parseInt(e.target.value)
  const target = authStore.daftarSiswa.find(s => s.id === selectedId)
  if (target) {
    authStore.setActiveSiswa(target)
  }
}

const handleLogout = () => {
  authStore.logout()
  router.replace('/login')
}
</script>
