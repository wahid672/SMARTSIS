<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc] pt-14 md:pt-4">
      <div class="max-w-6xl mx-auto py-2">

        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h1 class="text-xl font-bold text-slate-800">Master Mata Pelajaran</h1>
            <p class="text-xs text-slate-500 mt-1">
              Kelola daftar kurikulum pesantren (Diniyyah) dan kurikulum umum/nasional beserta standar KKM.
            </p>
          </div>
          <button
            @click="openAddModal"
            class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all self-start sm:self-auto"
          >
            <ion-icon :icon="addOutline" class="text-base"></ion-icon>
            Tambah Mata Pelajaran
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

        <!-- Filter & Search Bar -->
        <div class="bg-white p-4 rounded-2xl border border-slate-200/80 shadow-sm mb-5 flex flex-col md:flex-row md:items-center justify-between gap-3">
          <!-- Kategori Tabs -->
          <div class="flex items-center gap-1.5 p-1 bg-slate-100/80 rounded-xl self-start">
            <button
              v-for="cat in categoryTabs"
              :key="cat.value"
              @click="activeCategory = cat.value"
              class="px-3 py-1.5 rounded-lg text-xs font-bold transition-all"
              :class="activeCategory === cat.value ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
            >
              {{ cat.label }}
            </button>
          </div>

          <!-- Input Search -->
          <div class="relative w-full md:w-72">
            <ion-icon :icon="searchOutline" class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-sm"></ion-icon>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari kode atau nama mapel..."
              class="w-full pl-9 pr-3 py-2 rounded-xl border border-slate-200 text-xs text-slate-700 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
            />
          </div>
        </div>

        <!-- Tabel Mata Pelajaran -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
          <div v-if="isLoading" class="p-8 text-center text-xs text-slate-400">
            Memuat data mata pelajaran...
          </div>
          <div v-else-if="filteredList.length === 0" class="p-12 text-center">
            <div class="w-12 h-12 rounded-full bg-slate-100 flex items-center justify-center mx-auto text-slate-400 mb-3">
              <ion-icon :icon="bookOutline" class="text-2xl"></ion-icon>
            </div>
            <p class="text-xs font-medium text-slate-600">Tidak ada mata pelajaran yang ditemukan</p>
            <p class="text-[11px] text-slate-400 mt-0.5">Silakan tambahkan mata pelajaran baru untuk memulai.</p>
          </div>
          <table v-else class="w-full text-left text-xs text-slate-600">
            <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
              <tr>
                <th class="py-3.5 px-4 w-28">Kode</th>
                <th class="py-3.5 px-4">Nama Mata Pelajaran</th>
                <th class="py-3.5 px-4">Kategori Kurikulum</th>
                <th class="py-3.5 px-4 text-center">KKM</th>
                <th class="py-3.5 px-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="m in filteredList" :key="m.id" class="hover:bg-slate-50/60 transition-colors">
                <td class="py-3.5 px-4">
                  <span class="inline-block px-2.5 py-1 rounded-lg bg-slate-100 text-slate-800 font-mono font-bold text-[11px] tracking-wide">
                    {{ m.kode_mapel }}
                  </span>
                </td>
                <td class="py-3.5 px-4 font-bold text-slate-800 text-sm">
                  {{ m.nama_mapel }}
                </td>
                <td class="py-3.5 px-4">
                  <span
                    class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                    :class="m.kategori === 'Diniyyah / Pesantren' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-sky-50 text-sky-700 border border-sky-200'"
                  >
                    {{ m.kategori }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-center font-bold text-slate-700">
                  <span class="px-2 py-0.5 rounded bg-slate-100 text-slate-800 font-mono">
                    {{ m.kkm }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-right">
                  <div class="inline-flex items-center gap-1.5">
                    <button
                      @click="openEditModal(m)"
                      class="p-1.5 rounded-lg text-slate-400 hover:text-brand-600 hover:bg-brand-50 transition-colors"
                      title="Edit"
                    >
                      <ion-icon :icon="createOutline" class="text-base"></ion-icon>
                    </button>
                    <button
                      @click="handleDelete(m)"
                      class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-colors"
                      title="Hapus"
                    >
                      <ion-icon :icon="trashOutline" class="text-base"></ion-icon>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Modal Tambah / Edit -->
        <div v-if="showModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-md w-full shadow-xl border border-slate-100">
            <h2 class="text-sm font-bold text-slate-800 mb-1">
              {{ isEditing ? 'Edit Mata Pelajaran' : 'Tambah Mata Pelajaran Baru' }}
            </h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Masukkan informasi mata pelajaran dan standar kriteria ketuntasan minimal (KKM).
            </p>

            <form @submit.prevent="handleSubmit" class="space-y-3.5">
              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">Kode Mapel</label>
                <input
                  v-model="form.kode_mapel"
                  type="text"
                  placeholder="Contoh: TFQ, FQH, MTK"
                  required
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-xs uppercase text-slate-800 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">Nama Mata Pelajaran</label>
                <input
                  v-model="form.nama_mapel"
                  type="text"
                  placeholder="Contoh: Tahfidz Al-Qur'an & Tajwid"
                  required
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">Kategori Kurikulum</label>
                <select
                  v-model="form.kategori"
                  required
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 bg-white"
                >
                  <option value="Diniyyah / Pesantren">Diniyyah / Pesantren</option>
                  <option value="Umum / Nasional">Umum / Nasional</option>
                </select>
              </div>

              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">KKM (Kriteria Ketuntasan Minimal)</label>
                <input
                  v-model.number="form.kkm"
                  type="number"
                  min="0"
                  max="100"
                  step="0.5"
                  placeholder="75"
                  required
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-100">
                <button
                  type="button"
                  @click="showModal = false"
                  class="px-4 py-2 rounded-xl border border-slate-200 text-xs font-bold text-slate-600 hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  :disabled="isSubmitting"
                  class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all disabled:opacity-50"
                >
                  {{ isSubmitting ? 'Menyimpan...' : (isEditing ? 'Perbarui' : 'Simpan') }}
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
import { ref, computed, onMounted } from 'vue'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import {
  addOutline, bookOutline, searchOutline, createOutline, trashOutline,
  checkmarkCircleOutline, alertCircleOutline
} from 'ionicons/icons'
import axios from 'axios'

const mapelList = ref([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const showModal = ref(false)
const isEditing = ref(false)
const editId = ref(null)

const searchQuery = ref('')
const activeCategory = ref('all')

const successMsg = ref('')
const errorMsg = ref('')

const categoryTabs = [
  { label: 'Semua', value: 'all' },
  { label: 'Diniyyah / Pesantren', value: 'Diniyyah / Pesantren' },
  { label: 'Umum / Nasional', value: 'Umum / Nasional' }
]

const form = ref({
  kode_mapel: '',
  nama_mapel: '',
  kategori: 'Diniyyah / Pesantren',
  kkm: 75
})

const filteredList = computed(() => {
  return mapelList.value.filter(m => {
    const matchCat = activeCategory.value === 'all' || m.kategori === activeCategory.value
    const q = searchQuery.value.toLowerCase().trim()
    const matchQuery = !q || m.nama_mapel.toLowerCase().includes(q) || m.kode_mapel.toLowerCase().includes(q)
    return matchCat && matchQuery
  })
})

const fetchMapel = async () => {
  isLoading.value = true
  try {
    const res = await axios.get('/api/admin/mapel')
    mapelList.value = res.data || []
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal memuat daftar mata pelajaran'
  } finally {
    isLoading.value = false
  }
}

const openAddModal = () => {
  isEditing.value = false
  editId.value = null
  form.value = {
    kode_mapel: '',
    nama_mapel: '',
    kategori: 'Diniyyah / Pesantren',
    kkm: 75
  }
  showModal.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  editId.value = item.id
  form.value = {
    kode_mapel: item.kode_mapel,
    nama_mapel: item.nama_mapel,
    kategori: item.kategori,
    kkm: item.kkm
  }
  showModal.value = true
}

const handleSubmit = async () => {
  isSubmitting.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    if (isEditing.value) {
      await axios.put(`/api/admin/mapel/${editId.value}`, form.value)
      successMsg.value = 'Mata pelajaran berhasil diperbarui'
    } else {
      await axios.post('/api/admin/mapel', form.value)
      successMsg.value = 'Mata pelajaran baru berhasil ditambahkan'
    }
    showModal.value = false
    await fetchMapel()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan mata pelajaran'
  } finally {
    isSubmitting.value = false
    setTimeout(() => {
      successMsg.value = ''
      errorMsg.value = ''
    }, 4000)
  }
}

const handleDelete = async (item) => {
  if (!confirm(`Hapus mata pelajaran "${item.nama_mapel}"?`)) return
  successMsg.value = ''
  errorMsg.value = ''
  try {
    await axios.delete(`/api/admin/mapel/${item.id}`)
    successMsg.value = 'Mata pelajaran berhasil dihapus'
    await fetchMapel()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menghapus mata pelajaran'
  } finally {
    setTimeout(() => {
      successMsg.value = ''
      errorMsg.value = ''
    }, 4000)
  }
}

onMounted(() => {
  fetchMapel()
})
</script>
