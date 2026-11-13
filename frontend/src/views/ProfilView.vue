<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-slate-50">
      <div class="max-w-2xl mx-auto py-6">
        <div class="bg-white rounded-2xl p-6 border border-slate-100 shadow-sm text-center">
          
          <div class="w-16 h-16 rounded-2xl bg-brand-50 text-brand-600 mx-auto flex items-center justify-center text-3xl font-bold mb-4">
            <ion-icon :icon="personOutline"></ion-icon>
          </div>

          <h2 class="text-lg font-bold text-slate-800">{{ authStore.user?.name }}</h2>
          <p class="text-xs text-brand-600 font-medium">{{ authStore.roleLabel }}</p>

          <div class="mt-6 border-t border-slate-100 pt-4 text-left space-y-3 text-xs">
            <div class="flex justify-between py-1 border-b border-slate-50">
              <span class="text-slate-400">Tipe Akun</span>
              <span class="font-medium uppercase text-slate-700">{{ authStore.user?.user_type }}</span>
            </div>
            <div class="flex justify-between py-1 border-b border-slate-50">
              <span class="text-slate-400">Email</span>
              <span class="font-medium text-slate-700">{{ authStore.user?.email || '-' }}</span>
            </div>
            <div class="flex justify-between py-1 border-b border-slate-50">
              <span class="text-slate-400">No. WhatsApp / HP</span>
              <span class="font-medium text-slate-700">{{ authStore.user?.phone_number || '-' }}</span>
            </div>
          </div>

          <button
            @click="handleLogout"
            class="w-full mt-6 py-2.5 rounded-xl bg-rose-50 hover:bg-rose-100 text-rose-600 font-medium text-xs transition-colors flex items-center justify-center gap-2"
          >
            <ion-icon :icon="logOutOutline" class="text-base"></ion-icon>
            Keluar dari Akun
          </button>

        </div>
      </div>
    </ion-content>
  </ion-page>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import { personOutline, logOutOutline } from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
  router.replace('/login')
}
</script>
