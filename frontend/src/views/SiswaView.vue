<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc] pt-14 md:pt-4">
      <div class="max-w-6xl mx-auto py-2">

        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h1 class="text-xl font-bold text-slate-800">Manajemen Data Siswa</h1>
            <p class="text-xs text-slate-500 mt-1">
              Data siswa terhubung langsung dengan akun login wali murid via nomor WhatsApp/HP.
            </p>
          </div>
          <button
            @click="openAddModal"
            class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all self-start sm:self-auto"
          >
            <ion-icon :icon="personAddOutline" class="text-base"></ion-icon>
            Tambah Siswa Baru
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

        <!-- Filter Pencarian & Kelas -->
        <div class="bg-white p-3.5 rounded-2xl border border-slate-200/80 shadow-sm mb-4 flex flex-col sm:flex-row gap-3">
          <div class="flex-1 relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari berdasarkan nama siswa atau NIS..."
              class="w-full pl-9 pr-3.5 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 text-slate-800"
            />
            <ion-icon :icon="searchOutline" class="absolute left-3 top-2.5 text-slate-400 text-sm"></ion-icon>
          </div>
          <div class="w-full sm:w-48">
            <select
              v-model="selectedKelas"
              class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 text-slate-700 font-medium"
            >
              <option value="">Semua Kelas</option>
              <option v-for="k in kelasList" :key="k.id" :value="k.nama_kelas">
                {{ k.nama_kelas }} ({{ k.tingkat }})
              </option>
            </select>
          </div>
        </div>

        <!-- Tabel Data Siswa -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs text-slate-600">
              <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                <tr>
                  <th class="py-3.5 px-4">Nama Siswa</th>
                  <th class="py-3.5 px-4">NIS / NISN</th>
                  <th class="py-3.5 px-4">Kelas</th>
                  <th class="py-3.5 px-4">Akun Login Wali (No. WA)</th>
                  <th class="py-3.5 px-4 text-center">Status</th>
                  <th class="py-3.5 px-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="s in filteredSiswa" :key="s.id" class="hover:bg-slate-50/60 transition-colors">
                  <td class="py-3.5 px-4">
                    <div class="flex items-center gap-3">
                      <div class="w-9 h-9 rounded-xl flex items-center justify-center font-bold text-sm shrink-0 border"
                           :class="s.jenis_kelamin === 'L' ? 'bg-sky-50 text-sky-600 border-sky-100' : 'bg-rose-50 text-rose-600 border-rose-100'">
                        {{ s.nama_lengkap.charAt(0) }}
                      </div>
                      <div>
                        <span class="font-bold text-slate-800 text-xs block">{{ s.nama_lengkap }}</span>
                        <span class="text-[10px] text-slate-400 font-medium">
                          {{ s.jenis_kelamin === 'L' ? 'Laki-laki' : 'Perempuan' }}
                        </span>
                      </div>
                    </div>
                  </td>

                  <td class="py-3.5 px-4 font-mono">
                    <span class="text-slate-800 font-bold block">{{ s.nis }}</span>
                    <span class="text-[10px] text-slate-400">{{ s.nisn || '-' }}</span>
                  </td>

                  <td class="py-3.5 px-4">
                    <span class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-sky-100 text-sky-700">
                      {{ s.kelas }}
                    </span>
                  </td>

                  <td class="py-3.5 px-4">
                    <span class="font-semibold text-slate-800 block">{{ s.nama_wali }}</span>
                    <span class="text-[11px] font-mono text-emerald-600 flex items-center gap-1">
                      <ion-icon :icon="logoWhatsapp" class="text-xs"></ion-icon>
                      {{ s.phone_wali }}
                    </span>
                  </td>

                  <td class="py-3.5 px-4 text-center">
                    <span
                      class="inline-block px-2 py-0.5 rounded-full text-[10px] font-medium"
                      :class="s.is_active ? 'bg-green-100 text-green-700' : 'bg-rose-100 text-rose-700'"
                    >
                      {{ s.is_active ? 'Aktif' : 'Nonaktif' }}
                    </span>
                  </td>

                  <td class="py-3.5 px-4 text-right">
                    <div class="inline-flex items-center gap-1">
                      <button
                        @click="openEditModal(s)"
                        class="p-1.5 rounded-lg text-slate-500 hover:text-brand-600 hover:bg-emerald-50 transition-colors"
                        title="Edit Siswa"
                      >
                        <ion-icon :icon="createOutline" class="text-base"></ion-icon>
                      </button>
                      <button
                        @click="handleDelete(s)"
                        class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-colors"
                        title="Hapus Siswa"
                      >
                        <ion-icon :icon="trashOutline" class="text-base"></ion-icon>
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Modal Tambah / Edit Siswa -->
        <div v-if="showModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-lg w-full shadow-xl border border-slate-100 max-h-[90vh] overflow-y-auto">
            <h2 class="text-sm font-bold text-slate-800 mb-1">
              {{ isEditing ? 'Edit Data Siswa' : 'Tambah Siswa Baru' }}
            </h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Lengkapi data siswa dan nomor WhatsApp wali murid untuk akses login otomatis.
            </p>

            <form @submit.prevent="handleSubmit" class="space-y-4">
              
              <!-- Bagian Data Siswa -->
              <div class="bg-slate-50/70 p-3.5 rounded-xl border border-slate-200/60 space-y-3">
                <span class="text-[10px] font-bold uppercase tracking-wider text-slate-400 block">Identitas Siswa</span>
                
                <div>
                  <label class="block text-[11px] font-semibold text-slate-700 mb-1">Nama Lengkap Siswa</label>
                  <input
                    v-model="form.nama_lengkap"
                    type="text"
                    required
                    placeholder="Contoh: Muhammad Farhan"
                    class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                  />
                </div>

                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <label class="block text-[11px] font-semibold text-slate-700 mb-1">NIS (Nomor Induk)</label>
                    <input
                      v-model="form.nis"
                      type="text"
                      required
                      placeholder="2026001"
                      class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                    />
                  </div>
                  <div>
                    <label class="block text-[11px] font-semibold text-slate-700 mb-1">NISN (Nasional)</label>
                    <input
                      v-model="form.nisn"
                      type="text"
                      placeholder="0089123456"
                      class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                    />
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <label class="block text-[11px] font-semibold text-slate-700 mb-1">Jenis Kelamin</label>
                    <select
                      v-model="form.jenis_kelamin"
                      required
                      class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                    >
                      <option value="L">Laki-laki (Putra)</option>
                      <option value="P">Perempuan (Putri)</option>
                    </select>
                  </div>

                  <div>
                    <label class="block text-[11px] font-semibold text-slate-700 mb-1">Kelas / Rombel</label>
                    <select
                      v-model="form.kelas"
                      required
                      class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                    >
                      <option v-for="k in kelasList" :key="k.id" :value="k.nama_kelas">
                        {{ k.nama_kelas }} ({{ k.tingkat }})
                      </option>
                    </select>
                  </div>
                </div>
              </div>

              <!-- Bagian Akun Wali Murid -->
              <div class="bg-emerald-50/50 p-3.5 rounded-xl border border-emerald-200/60 space-y-3">
                <div class="flex items-center justify-between">
                  <span class="text-[10px] font-bold uppercase tracking-wider text-emerald-800">Akun Login Wali Murid</span>
                  <span class="text-[10px] text-emerald-600 bg-white px-2 py-0.5 rounded-full border border-emerald-200 font-semibold">
                    PIN Awal: 123456
                  </span>
                </div>

                <div>
                  <label class="block text-[11px] font-semibold text-slate-700 mb-1">Nama Orang Tua / Wali</label>
                  <input
                    v-model="form.nama_wali"
                    type="text"
                    required
                    placeholder="Contoh: Bapak H. Ahmad"
                    class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                  />
                </div>

                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <label class="block text-[11px] font-semibold text-slate-700 mb-1">No. WhatsApp / HP</label>
                    <input
                      v-model="form.phone_wali"
                      type="text"
                      required
                      placeholder="081234567890"
                      class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                    />
                  </div>
                  <div>
                    <label class="block text-[11px] font-semibold text-slate-700 mb-1">Status Hubungan</label>
                    <select
                      v-model="form.hubungan"
                      class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                    >
                      <option value="Ayah Kandung">Ayah Kandung</option>
                      <option value="Ibu Kandung">Ibu Kandung</option>
                      <option value="Wali Siswa">Wali Siswa</option>
                    </select>
                  </div>
                </div>
              </div>

              <div class="pt-2 flex justify-end gap-2">
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
                  {{ isEditing ? 'Simpan Perubahan' : 'Tambah Siswa' }}
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
import axios from 'axios'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import {
  personAddOutline, createOutline, trashOutline, searchOutline, logoWhatsapp,
  checkmarkCircleOutline, alertCircleOutline
} from 'ionicons/icons'

const siswaList = ref([])
const kelasList = ref([])
const searchQuery = ref('')
const selectedKelas = ref('')

const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const successMsg = ref('')
const errorMsg = ref('')

const form = ref({
  nis: '',
  nisn: '',
  nama_lengkap: '',
  jenis_kelamin: 'L',
  kelas: '',
  nama_wali: '',
  phone_wali: '',
  hubungan: 'Ayah Kandung',
  is_active: true
})

const fetchSiswa = async () => {
  try {
    const res = await axios.get('/api/admin/siswa')
    siswaList.value = res.data || []
  } catch (err) {
    errorMsg.value = 'Gagal memuat data siswa'
  }
}

const fetchKelas = async () => {
  try {
    const res = await axios.get('/api/admin/kelas')
    kelasList.value = res.data || []
    if (kelasList.value.length > 0 && !form.value.kelas) {
      form.value.kelas = kelasList.value[0].nama_kelas
    }
  } catch (err) {
    console.error(err)
  }
}

const filteredSiswa = computed(() => {
  return siswaList.value.filter(s => {
    const matchQuery = !searchQuery.value ||
      s.nama_lengkap.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      s.nis.includes(searchQuery.value)
    const matchKelas = !selectedKelas.value || s.kelas === selectedKelas.value
    return matchQuery && matchKelas
  })
})

const openAddModal = () => {
  isEditing.value = false
  editingId.value = null
  form.value = {
    nis: '',
    nisn: '',
    nama_lengkap: '',
    jenis_kelamin: 'L',
    kelas: kelasList.value[0]?.nama_kelas || '',
    nama_wali: '',
    phone_wali: '',
    hubungan: 'Ayah Kandung',
    is_active: true
  }
  showModal.value = true
}

const openEditModal = (s) => {
  isEditing.value = true
  editingId.value = s.id
  form.value = {
    nis: s.nis,
    nisn: s.nisn,
    nama_lengkap: s.nama_lengkap,
    jenis_kelamin: s.jenis_kelamin,
    kelas: s.kelas,
    nama_wali: s.nama_wali !== '-' ? s.nama_wali : '',
    phone_wali: s.phone_wali !== '-' ? s.phone_wali : '',
    hubungan: s.hubungan !== '-' ? s.hubungan : 'Ayah Kandung',
    is_active: s.is_active
  }
  showModal.value = true
}

const handleSubmit = async () => {
  successMsg.value = ''
  errorMsg.value = ''

  try {
    if (isEditing.value) {
      await axios.put(`/api/admin/siswa/${editingId.value}`, form.value)
      successMsg.value = 'Data siswa berhasil diperbarui!'
    } else {
      await axios.post('/api/admin/siswa', form.value)
      successMsg.value = 'Siswa baru dan akun login wali berhasil disimpan!'
    }
    showModal.value = false
    await fetchSiswa()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan data siswa'
  }
}

const handleDelete = async (s) => {
  if (!confirm(`Hapus data siswa "${s.nama_lengkap}"?`)) return
  try {
    await axios.delete(`/api/admin/siswa/${s.id}`)
    successMsg.value = 'Data siswa berhasil dihapus'
    await fetchSiswa()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menghapus siswa'
  }
}

onMounted(() => {
  fetchSiswa()
  fetchKelas()
})
</script>
