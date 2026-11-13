<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc] pt-14 md:pt-4">
      <div class="max-w-5xl mx-auto py-2">

        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h1 class="text-xl font-bold text-slate-800">Master Data Kelas</h1>
            <p class="text-xs text-slate-500 mt-1">
              Atur daftar rombel / kelas, tingkatan, serta wali kelas untuk santri/siswa.
            </p>
          </div>
          <button
            @click="openAddModal"
            class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all self-start sm:self-auto"
          >
            <ion-icon :icon="addOutline" class="text-base"></ion-icon>
            Tambah Kelas Baru
          </button>
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

        <!-- Daftar Kelas Card Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
          <div
            v-for="k in kelasList"
            :key="k.id"
            class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-all flex flex-col justify-between"
          >
            <div>
              <div class="flex items-center justify-between mb-3">
                <span class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-sky-100 text-sky-700">
                  {{ k.tingkat }}
                </span>
                <span class="text-xs font-bold text-emerald-600 bg-emerald-50 px-2.5 py-0.5 rounded-full border border-emerald-100">
                  {{ k.siswa_count }} Siswa
                </span>
              </div>

              <h3 class="text-lg font-black text-slate-800">{{ k.nama_kelas }}</h3>
              <p class="text-xs text-slate-500 mt-1">
                Wali Kelas: <span class="font-semibold text-slate-700">{{ k.wali_kelas || '-' }}</span>
              </p>
              <p class="text-[11px] text-slate-400 mt-0.5">
                {{ k.keterangan || 'Tidak ada keterangan' }}
              </p>
            </div>

            <div class="pt-4 mt-4 border-t border-slate-100 flex items-center justify-end gap-1">
              <button
                @click="openEditModal(k)"
                class="p-1.5 rounded-lg text-slate-500 hover:text-brand-600 hover:bg-emerald-50 transition-colors"
                title="Edit Kelas"
              >
                <ion-icon :icon="createOutline" class="text-base"></ion-icon>
              </button>
              <button
                v-if="k.siswa_count === 0"
                @click="handleDelete(k)"
                class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-colors"
                title="Hapus Kelas"
              >
                <ion-icon :icon="trashOutline" class="text-base"></ion-icon>
              </button>
            </div>
          </div>
        </div>

        <!-- Modal Tambah / Edit Kelas -->
        <div v-if="showModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-md w-full shadow-xl border border-slate-100">
            <h2 class="text-sm font-bold text-slate-800 mb-1">
              {{ isEditing ? 'Edit Data Kelas' : 'Tambah Kelas Baru' }}
            </h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Tentukan nama rombel, tingkatan, dan guru wali kelasnya.
            </p>

            <form @submit.prevent="handleSubmit" class="space-y-3.5">
              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Nama Rombel / Kelas</label>
                <input
                  v-model="form.nama_kelas"
                  type="text"
                  required
                  placeholder="Contoh: VII-A / X-IPA 1"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Tingkatan / Jenjang</label>
                <input
                  v-model="form.tingkat"
                  type="text"
                  required
                  placeholder="Contoh: Kelas 7 / Kelas 10"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Nama Wali Kelas</label>
                <input
                  v-model="form.wali_kelas"
                  type="text"
                  placeholder="Contoh: Ustadz Ahmad, S.Pd."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Keterangan Tambahan</label>
                <input
                  v-model="form.keterangan"
                  type="text"
                  placeholder="Contoh: Putra - Reguler"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div class="pt-3 border-t border-slate-100 flex justify-end gap-2">
                <button
                  type="button"
                  @click="showModal = false"
                  class="px-3.5 py-2 rounded-xl text-slate-500 hover:bg-slate-100 text-xs font-medium"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-sm"
                >
                  {{ isEditing ? 'Simpan Perubahan' : 'Buat Kelas' }}
                </button>
              </div>
            </form>
          </div>
        </div>

      </div>
    </ion-content>
  </ion-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import {
  addOutline, createOutline, trashOutline,
  checkmarkCircleOutline, alertCircleOutline
} from 'ionicons/icons'

const kelasList = ref([])
const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const successMsg = ref('')
const errorMsg = ref('')

const form = ref({
  nama_kelas: '',
  tingkat: '',
  wali_kelas: '',
  keterangan: ''
})

const fetchKelas = async () => {
  try {
    const res = await axios.get('/api/admin/kelas')
    kelasList.value = res.data || []
  } catch (err) {
    errorMsg.value = 'Gagal memuat data kelas'
  }
}

const openAddModal = () => {
  isEditing.value = false
  editingId.value = null
  form.value = {
    nama_kelas: '',
    tingkat: '',
    wali_kelas: '',
    keterangan: ''
  }
  showModal.value = true
}

const openEditModal = (k) => {
  isEditing.value = true
  editingId.value = k.id
  form.value = {
    nama_kelas: k.nama_kelas,
    tingkat: k.tingkat,
    wali_kelas: k.wali_kelas,
    keterangan: k.keterangan
  }
  showModal.value = true
}

const handleSubmit = async () => {
  successMsg.value = ''
  errorMsg.value = ''

  try {
    if (isEditing.value) {
      await axios.put(`/api/admin/kelas/${editingId.value}`, form.value)
      successMsg.value = 'Data kelas berhasil diperbarui!'
    } else {
      await axios.post('/api/admin/kelas', form.value)
      successMsg.value = 'Kelas baru berhasil ditambahkan!'
    }
    showModal.value = false
    await fetchKelas()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan data kelas'
  }
}

const handleDelete = async (k) => {
  if (!confirm(`Hapus kelas "${k.nama_kelas}"?`)) return
  try {
    await axios.delete(`/api/admin/kelas/${k.id}`)
    successMsg.value = 'Kelas berhasil dihapus'
    await fetchKelas()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menghapus kelas'
  }
}

onMounted(() => {
  fetchKelas()
})
</script>
