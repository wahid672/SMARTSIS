<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc] pt-14 md:pt-4">
      <div class="max-w-5xl mx-auto py-2">

        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h1 class="text-xl font-bold text-slate-800">Master Tahun Ajaran</h1>
            <p class="text-xs text-slate-500 mt-1">
              Atur tahun ajaran aktif dan semester sebagai rujukan akademik serta tagihan siswa.
            </p>
          </div>
          <button
            @click="openAddModal"
            class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all self-start sm:self-auto"
          >
            <ion-icon :icon="addOutline" class="text-base"></ion-icon>
            Tambah Tahun Ajaran
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

        <!-- Tabel Tahun Ajaran -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
          <table class="w-full text-left text-xs text-slate-600">
            <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
              <tr>
                <th class="py-3.5 px-4">Tahun Ajaran</th>
                <th class="py-3.5 px-4">Semester</th>
                <th class="py-3.5 px-4 text-center">Status</th>
                <th class="py-3.5 px-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="t in tahunList" :key="t.id" class="hover:bg-slate-50/60 transition-colors">
                <td class="py-3.5 px-4">
                  <span class="font-bold text-slate-800 text-sm block">{{ t.nama_tahun }}</span>
                  <span class="text-[10px] text-slate-400">ID: #{{ t.id }}</span>
                </td>

                <td class="py-3.5 px-4 font-semibold text-slate-700">
                  Semester {{ t.semester }}
                </td>

                <td class="py-3.5 px-4 text-center">
                  <span
                    v-if="t.is_active"
                    class="inline-block px-3 py-1 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800 border border-emerald-200"
                  >
                    Aktif Berjalan
                  </span>
                  <span
                    v-else
                    class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-medium bg-slate-100 text-slate-500"
                  >
                    Tidak Aktif
                  </span>
                </td>

                <td class="py-3.5 px-4 text-right">
                  <div class="inline-flex items-center gap-2">
                    <button
                      v-if="!t.is_active"
                      @click="handleActivate(t)"
                      class="px-2.5 py-1 rounded-lg text-emerald-700 bg-emerald-50 hover:bg-emerald-100 text-[11px] font-bold transition-colors"
                    >
                      Aktifkan
                    </button>
                    <button
                      v-if="!t.is_active"
                      @click="handleDelete(t)"
                      class="p-1 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50"
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

        <!-- Modal Tambah -->
        <div v-if="showModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-md w-full shadow-xl border border-slate-100">
            <h2 class="text-sm font-bold text-slate-800 mb-1">Tambah Tahun Ajaran Baru</h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Tentukan periode tahun ajaran dan semester aktif.
            </p>

            <form @submit.prevent="handleSubmit" class="space-y-3.5">
              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Tahun Ajaran</label>
                <input
                  v-model="form.nama_tahun"
                  type="text"
                  required
                  placeholder="Contoh: 2026/2027"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Semester</label>
                <select
                  v-model="form.semester"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                >
                  <option value="Ganjil">Semester Ganjil</option>
                  <option value="Genap">Semester Genap</option>
                </select>
              </div>

              <div class="flex items-center gap-2 pt-1">
                <input
                  type="checkbox"
                  id="activeTahun"
                  v-model="form.is_active"
                  class="w-4 h-4 rounded text-brand-600 focus:ring-brand-500 border-slate-300 cursor-pointer"
                />
                <label for="activeTahun" class="text-xs font-medium text-slate-700 cursor-pointer">
                  Jadikan Tahun Ajaran Aktif Saat Ini
                </label>
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
                  Simpan Tahun Ajaran
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
import { addOutline, trashOutline, checkmarkCircleOutline, alertCircleOutline } from 'ionicons/icons'
import { showConfirm, showAlert, showToast } from '@/utils/alert'

const tahunList = ref([])
const showModal = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const form = ref({
  nama_tahun: '',
  semester: 'Ganjil',
  is_active: true
})

const fetchTahun = async () => {
  try {
    const res = await axios.get('/api/admin/tahun-ajaran')
    tahunList.value = res.data || []
  } catch (err) {
    errorMsg.value = 'Gagal memuat tahun ajaran'
  }
}

const openAddModal = () => {
  form.value = { nama_tahun: '', semester: 'Ganjil', is_active: true }
  showModal.value = true
}

const handleSubmit = async () => {
  successMsg.value = ''
  errorMsg.value = ''
  try {
    await axios.post('/api/admin/tahun-ajaran', form.value)
    successMsg.value = 'Tahun ajaran berhasil ditambahkan!'
    showModal.value = false
    await fetchTahun()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan tahun ajaran'
  }
}

const handleActivate = async (t) => {
  try {
    await axios.put(`/api/admin/tahun-ajaran/${t.id}/activate`)
    successMsg.value = `Tahun ajaran ${t.nama_tahun} (${t.semester}) sekarang aktif!`
    await fetchTahun()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = 'Gagal mengaktifkan tahun ajaran'
  }
}

const handleDelete = async (t) => {
  const confirmed = await showConfirm({
    title: 'Hapus Tahun Ajaran?',
    text: `Yakin ingin menghapus tahun ajaran <b>"${t.nama_tahun}"</b>?`,
    icon: 'warning',
    confirmButtonText: 'Ya, Hapus',
    cancelButtonText: 'Batal',
    isDanger: true
  })
  if (!confirmed) return

  try {
    await axios.delete(`/api/admin/tahun-ajaran/${t.id}`)
    showToast('Tahun ajaran berhasil dihapus', 'success')
    await fetchTahun()
  } catch (err) {
    showAlert({
      title: 'Gagal Menghapus',
      text: err.response?.data?.error || 'Gagal menghapus tahun ajaran',
      icon: 'error'
    })
  }
}

onMounted(() => {
  fetchTahun()
})
</script>
