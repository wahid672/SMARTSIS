<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-slate-50 pt-14 md:pt-4">
      <div class="max-w-2xl mx-auto py-8 text-center">
        <div class="w-16 h-16 rounded-2xl bg-emerald-50 text-emerald-600 mx-auto flex items-center justify-center text-3xl font-bold mb-4">
          <ion-icon :icon="sparklesOutline"></ion-icon>
        </div>
        <h2 class="text-lg font-bold text-slate-800">{{ pageTitle }}</h2>
        <p class="text-xs text-slate-500 mt-1 max-w-sm mx-auto">
          Modul ini terhubung dengan hak akses Anda. Tampilan dan data form/tabel siap disesuaikan pada fase pengembangan modul spesifik.
        </p>
      </div>
    </ion-content>
  </ion-page>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import { sparklesOutline } from 'ionicons/icons'

const route = useRoute()
const pageTitle = computed(() => {
  return route.meta?.title || 'Modul Sistem Sekolah'
})
</script>
