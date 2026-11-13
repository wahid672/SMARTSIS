<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-slate-50 pt-14 md:pt-4">
      <div class="max-w-5xl mx-auto py-2">

        <!-- Header Halaman -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h1 class="text-xl font-bold text-slate-800">Manajemen Pegawai & Guru</h1>
            <p class="text-xs text-slate-500 mt-1">
              Kelola data staf sekolah, dewan guru, dan tentukan level hak akses (role) masing-masing.
            </p>
          </div>
          <button
            @click="openAddModal"
            class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-medium shadow-md shadow-emerald-500/20 transition-all self-start sm:self-auto"
          >
            <ion-icon :icon="personAddOutline" class="text-base"></ion-icon>
            Tambah Pegawai Baru
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

        <!-- Tabel Data Pegawai (Gaya Admin Dashboard) -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs text-slate-600">
              <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                <tr>
                  <th class="py-3.5 px-4">Nama Pegawai / Guru</th>
                  <th class="py-3.5 px-4">Kontak (Email & WhatsApp)</th>
                  <th class="py-3.5 px-4">Hak Akses / Role</th>
                  <th class="py-3.5 px-4 text-center">Status</th>
                  <th class="py-3.5 px-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="p in pegawaiList" :key="p.id" class="hover:bg-slate-50/60 transition-colors">
                  <td class="py-3.5 px-4">
                    <div class="flex items-center gap-3">
                      <div class="w-9 h-9 rounded-xl bg-emerald-50 text-brand-600 flex items-center justify-center font-bold text-sm shrink-0 border border-emerald-100">
                        {{ p.name.charAt(0) }}
                      </div>
                      <div>
                        <span class="font-bold text-slate-800 text-xs block">{{ p.name }}</span>
                        <span class="text-[10px] text-slate-400">ID: #{{ p.id }}</span>
                      </div>
                    </div>
                  </td>

                  <td class="py-3.5 px-4">
                    <span class="text-slate-700 font-medium block">{{ p.email }}</span>
                    <span class="text-[11px] text-slate-400">{{ p.phone_number || '-' }}</span>
                  </td>

                  <td class="py-3.5 px-4">
                    <span class="inline-block px-2.5 py-1 rounded-full text-[10px] font-semibold bg-emerald-50 text-brand-700 border border-emerald-200">
                      {{ p.role_label }}
                    </span>
                  </td>

                  <td class="py-3.5 px-4 text-center">
                    <span
                      class="inline-block px-2 py-0.5 rounded-full text-[10px] font-medium"
                      :class="p.is_active ? 'bg-green-100 text-green-700' : 'bg-rose-100 text-rose-700'"
                    >
                      {{ p.is_active ? 'Aktif' : 'Nonaktif' }}
                    </span>
                  </td>

                  <td class="py-3.5 px-4 text-right">
                    <div class="inline-flex items-center gap-1">
                      <button
                        @click="openEditModal(p)"
                        class="p-1.5 rounded-lg text-slate-500 hover:text-brand-600 hover:bg-emerald-50 transition-colors"
                        title="Edit Pegawai"
                      >
                        <ion-icon :icon="createOutline" class="text-base"></ion-icon>
                      </button>
                      <button
                        v-if="p.id !== 1"
                        @click="handleDelete(p)"
                        class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-colors"
                        title="Hapus Pegawai"
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

        <!-- Modal Tambah / Edit Pegawai -->
        <div v-if="showModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-md w-full shadow-xl border border-slate-100">
            <h2 class="text-sm font-bold text-slate-800 mb-1">
              {{ isEditing ? 'Edit Akun Pegawai' : 'Tambah Pegawai Baru' }}
            </h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Isi data lengkap pegawai dan tentukan wewenang hak aksesnya.
            </p>

            <form @submit.prevent="handleSubmit" class="space-y-3.5">
              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Nama Lengkap & Gelar</label>
                <input
                  v-model="form.name"
                  type="text"
                  required
                  placeholder="Contoh: Ustadz Ahmad Fauzi, S.Pd."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Email Login</label>
                <input
                  v-model="form.email"
                  type="email"
                  required
                  placeholder="fauzi@sekolah.sch.id"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Nomor WhatsApp / HP</label>
                <input
                  v-model="form.phone_number"
                  type="text"
                  placeholder="08123456789"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">
                  Kata Sandi {{ isEditing ? '(Kosongkan jika tidak diubah)' : '' }}
                </label>
                <input
                  v-model="form.password"
                  type="password"
                  :required="!isEditing"
                  placeholder="Minimal 6 karakter"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Level Hak Akses (Role)</label>
                <select
                  v-model="form.role_id"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                >
                  <option v-for="r in roles" :key="r.id" :value="r.id">
                    {{ r.label }}
                  </option>
                </select>
              </div>

              <div v-if="isEditing" class="flex items-center gap-2 pt-1">
                <input
                  type="checkbox"
                  id="activeCheck"
                  v-model="form.is_active"
                  class="w-4 h-4 rounded text-brand-600 focus:ring-brand-500 border-slate-300 cursor-pointer"
                />
                <label for="activeCheck" class="text-xs font-medium text-slate-700 cursor-pointer">
                  Akun Aktif (Bisa Login)
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
                  class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-medium shadow-sm"
                >
                  {{ isEditing ? 'Simpan Perubahan' : 'Tambah Pegawai' }}
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
  personAddOutline, createOutline, trashOutline,
  checkmarkCircleOutline, alertCircleOutline
} from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const pegawaiList = ref([])
const roles = ref([])
const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const successMsg = ref('')
const errorMsg = ref('')

const form = ref({
  name: '',
  email: '',
  phone_number: '',
  password: '',
  role_id: 2,
  is_active: true
})

const fetchPegawai = async () => {
  try {
    const res = await axios.get('/api/admin/pegawai', {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    pegawaiList.value = res.data || []
  } catch (err) {
    errorMsg.value = 'Gagal memuat daftar pegawai'
  }
}

const fetchRoles = async () => {
  try {
    const res = await axios.get('/api/admin/roles', {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    // Filter agar role wali tidak muncul untuk akun pegawai
    roles.value = (res.data || []).filter(r => r.name !== 'wali')
  } catch (err) {
    console.error(err)
  }
}

const openAddModal = () => {
  isEditing.value = false
  editingId.value = null
  form.value = {
    name: '',
    email: '',
    phone_number: '',
    password: '',
    role_id: roles.value[0]?.id || 2,
    is_active: true
  }
  showModal.value = true
}

const openEditModal = (p) => {
  isEditing.value = true
  editingId.value = p.id
  form.value = {
    name: p.name,
    email: p.email,
    phone_number: p.phone_number,
    password: '',
    role_id: p.role_id,
    is_active: p.is_active
  }
  showModal.value = true
}

const handleSubmit = async () => {
  successMsg.value = ''
  errorMsg.value = ''

  try {
    if (isEditing.value) {
      await axios.put(`/api/admin/pegawai/${editingId.value}`, form.value, {
        headers: { Authorization: `Bearer ${authStore.token}` }
      })
      successMsg.value = 'Data pegawai berhasil diperbarui!'
    } else {
      await axios.post('/api/admin/pegawai', form.value, {
        headers: { Authorization: `Bearer ${authStore.token}` }
      })
      successMsg.value = 'Pegawai baru berhasil ditambahkan!'
    }
    showModal.value = false
    await fetchPegawai()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan data pegawai'
  }
}

const handleDelete = async (p) => {
  if (!confirm(`Yakin ingin menghapus akun "${p.name}"?`)) return
  try {
    await axios.delete(`/api/admin/pegawai/${p.id}`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    successMsg.value = 'Akun pegawai berhasil dihapus'
    await fetchPegawai()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menghapus pegawai'
  }
}

onMounted(() => {
  fetchPegawai()
  fetchRoles()
})
</script>
