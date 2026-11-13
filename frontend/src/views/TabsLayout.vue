<template>
  <ion-page>
    <ion-split-pane content-id="main-content" when="md">
      
      <!-- SIDEBAR UNTUK TABLET & DESKTOP (MD/LG SCREEN) -->
      <ion-menu content-id="main-content" type="overlay" class="border-r border-sky-950/20 shadow-lg">
        <ion-content class="bg-[#0b3b60] text-slate-100">
          
          <!-- Header Sidebar dengan Brand Sekolah & Badge Role Kuning Emas -->
          <div class="p-5 border-b border-white/10 flex items-center gap-3 bg-[#082d4a]/70">
            <div class="w-11 h-11 flex items-center justify-center shrink-0">
              <img v-if="settingsStore.logo_url" :src="settingsStore.logo_url" alt="Logo" class="w-full h-full object-contain" />
              <div v-else class="w-11 h-11 rounded-xl bg-amber-400 text-amber-950 flex items-center justify-center font-bold text-lg shadow-sm">
                {{ settingsStore.app_name.charAt(0) }}
              </div>
            </div>
            <div class="min-w-0">
              <h2 class="font-extrabold text-white text-sm leading-tight truncate tracking-wide">{{ settingsStore.app_name }}</h2>
              <p class="text-[10px] text-sky-200 truncate">{{ settingsStore.school_name }}</p>
              <span class="inline-block text-[9px] font-bold px-2 py-0.5 rounded-full bg-amber-400 text-amber-950 shadow-sm mt-1">
                {{ authStore.roleLabel }}
              </span>
            </div>
          </div>

          <!-- Switcher Siswa di Sidebar untuk Akun Wali -->
          <div v-if="authStore.isWali && authStore.daftarSiswa.length > 0" class="p-3 mx-3 mt-3 bg-[#082d4a]/80 rounded-xl border border-white/10">
            <label class="text-[10px] font-bold text-amber-300 uppercase tracking-wider block mb-1">
              Santri / Siswa Aktif:
            </label>
            <select
              :value="authStore.activeSiswa?.id"
              @change="onSwitchSiswa($event)"
              class="w-full text-xs font-semibold text-white bg-[#0b3b60] border border-white/20 rounded-lg p-2 focus:outline-none"
            >
              <option v-for="s in authStore.daftarSiswa" :key="s.id" :value="s.id">
                {{ s.nama_lengkap }} ({{ s.kelas }})
              </option>
            </select>
          </div>

          <!-- Daftar Menu Berkelompok (Sub-kategori Sesuai Gambar Desain) -->
          <div class="px-3 py-4 space-y-4">
            <div v-for="(group, catName) in groupedMenus" :key="catName">
              <h3 class="text-[10px] font-extrabold uppercase tracking-wider text-sky-200/80 px-3 mb-1.5">
                {{ catName }}
              </h3>
              <div class="space-y-1">
                <router-link
                  v-for="menu in group"
                  :key="menu.code"
                  :to="menu.path"
                  class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-xs font-bold transition-all group"
                  :class="route.path === menu.path ? 'bg-amber-400 text-amber-950 shadow-md font-extrabold' : 'text-sky-100 hover:bg-white/10 hover:text-white'"
                >
                  <ion-icon :icon="getIcon(menu.icon)" class="text-base shrink-0 group-hover:scale-110 transition-transform"></ion-icon>
                  <span class="truncate">{{ menu.title }}</span>
                </router-link>
              </div>
            </div>

            <!-- Tombol Keluar di Sidebar -->
            <div class="pt-3 border-t border-white/10">
              <button
                @click="handleLogout"
                class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-xs font-bold text-rose-300 hover:bg-rose-500/20 hover:text-white transition-colors border border-rose-500/20"
              >
                <ion-icon :icon="logOutOutline" class="text-base shrink-0"></ion-icon>
                <span>Keluar Akun</span>
              </button>
            </div>
          </div>

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

// Mengelompokkan menu berdasarkan Kategori (Submenu/Grup Menu ala Admin Dashboard)
const groupedMenus = computed(() => {
  const groups = {}
  authStore.menus.forEach(m => {
    const cat = m.category || 'MENU LAINNYA'
    if (!groups[cat]) {
      groups[cat] = []
    }
    groups[cat].push(m)
  })
  return groups
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
