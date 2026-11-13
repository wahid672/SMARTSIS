<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-slate-50">
      <div class="max-w-4xl mx-auto py-6">

        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h1 class="text-xl font-bold text-slate-800">Manajemen Hak Akses & Role</h1>
            <p class="text-xs text-slate-500 mt-1">
              Atur level pegawai dan tentukan menu apa saja yang boleh diakses oleh setiap role.
            </p>
          </div>
          <button
            @click="showAddModal = true"
            class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-medium shadow-md shadow-emerald-500/20 transition-all self-start sm:self-auto"
          >
            <ion-icon :icon="addOutline" class="text-base"></ion-icon>
            Tambah Role Baru
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

        <!-- Daftar Role Cards -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          <div
            v-for="r in roles"
            :key="r.id"
            @click="selectRole(r)"
            class="p-5 rounded-2xl border transition-all cursor-pointer relative"
            :class="selectedRole?.id === r.id ? 'bg-emerald-50/70 border-brand-500 shadow-sm' : 'bg-white border-slate-200 hover:border-slate-300'"
          >
            <div class="flex items-start justify-between">
              <div>
                <span class="inline-block text-[10px] font-bold uppercase tracking-wider text-slate-400 mb-1">
                  Role ID: #{{ r.id }}
                </span>
                <h3 class="text-sm font-bold text-slate-800">{{ r.label }}</h3>
                <code class="text-[11px] text-emerald-600 font-mono">{{ r.name }}</code>
              </div>
              <span class="text-[10px] px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 font-medium">
                {{ r.user_count }} Akun
              </span>
            </div>

            <p class="text-xs text-slate-500 mt-2 line-clamp-2">
              {{ r.description || 'Tidak ada keterangan tambahan.' }}
            </p>

            <div class="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-xs">
              <span class="font-semibold text-brand-600 flex items-center gap-1">
                <ion-icon :icon="shieldCheckmarkOutline"></ion-icon>
                Kelola Izin Menu
              </span>
              <button
                v-if="r.id > 3 && r.user_count === 0"
                @click.stop="handleDeleteRole(r)"
                class="text-rose-500 hover:text-rose-700 p-1 rounded-lg hover:bg-rose-50"
                title="Hapus Role"
              >
                <ion-icon :icon="trashOutline" class="text-sm"></ion-icon>
              </button>
            </div>
          </div>
        </div>

        <!-- Matriks Centang Hak Akses Menu untuk Role yang Dipilih -->
        <div v-if="selectedRole" class="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
          
          <div class="p-5 border-b border-slate-100 bg-slate-50/50 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
            <div>
              <h2 class="text-sm font-bold text-slate-800">
                Pengaturan Menu: <span class="text-brand-600">{{ selectedRole.label }}</span>
              </h2>
              <p class="text-[11px] text-slate-500">
                Centang menu dan hak aksi apa saja yang boleh dibuka oleh akun dengan role ini.
              </p>
            </div>
            <button
              @click="savePermissions"
              :disabled="savingPerms"
              class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-medium shadow-sm transition-all flex items-center gap-1.5 self-start sm:self-auto disabled:opacity-60"
            >
              <ion-icon :icon="saveOutline" class="text-sm"></ion-icon>
              <span>{{ savingPerms ? 'Menyimpan...' : 'Simpan Hak Akses' }}</span>
            </button>
          </div>

          <!-- Tabel Matriks Menu -->
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs text-slate-600">
              <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                <tr>
                  <th class="py-3 px-4">Menu Aplikasi</th>
                  <th class="py-3 px-3 text-center">Tampilkan (Lihat)</th>
                  <th class="py-3 px-3 text-center">Tambah (Create)</th>
                  <th class="py-3 px-3 text-center">Edit (Update)</th>
                  <th class="py-3 px-3 text-center">Hapus (Delete)</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="m in permissions" :key="m.menu_id" class="hover:bg-slate-50/50">
                  <td class="py-3 px-4">
                    <div class="flex items-center gap-2">
                      <div class="w-7 h-7 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center text-sm shrink-0">
                        <ion-icon :icon="sparklesOutline"></ion-icon>
                      </div>
                      <div>
                        <span class="font-semibold text-slate-800 block">{{ m.title }}</span>
                        <code class="text-[10px] text-slate-400">{{ m.path }}</code>
                      </div>
                    </div>
                  </td>

                  <td class="py-3 px-3 text-center">
                    <input
                      type="checkbox"
                      v-model="m.can_view"
                      class="w-4 h-4 rounded text-brand-600 focus:ring-brand-500 border-slate-300 cursor-pointer"
                    />
                  </td>
                  <td class="py-3 px-3 text-center">
                    <input
                      type="checkbox"
                      v-model="m.can_create"
                      :disabled="!m.can_view"
                      class="w-4 h-4 rounded text-brand-600 focus:ring-brand-500 border-slate-300 cursor-pointer disabled:opacity-30"
                    />
                  </td>
                  <td class="py-3 px-3 text-center">
                    <input
                      type="checkbox"
                      v-model="m.can_edit"
                      :disabled="!m.can_view"
                      class="w-4 h-4 rounded text-brand-600 focus:ring-brand-500 border-slate-300 cursor-pointer disabled:opacity-30"
                    />
                  </td>
                  <td class="py-3 px-3 text-center">
                    <input
                      type="checkbox"
                      v-model="m.can_delete"
                      :disabled="!m.can_view"
                      class="w-4 h-4 rounded text-brand-600 focus:ring-brand-500 border-slate-300 cursor-pointer disabled:opacity-30"
                    />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

        </div>

        <!-- Modal Tambah Role Baru -->
        <div v-if="showAddModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-md w-full shadow-xl border border-slate-100">
            <h2 class="text-sm font-bold text-slate-800 mb-1">Tambah Role / Level Pegawai</h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Buat level akses baru untuk staf atau pengurus sekolah.
            </p>

            <form @submit.prevent="handleAddRole" class="space-y-3.5">
              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Nama Jabatan / Label</label>
                <input
                  v-model="newRole.label"
                  type="text"
                  required
                  placeholder="Contoh: Bendahara SPP / Guru Piket"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Kode Role (Huruf kecil tanpa spasi)</label>
                <input
                  v-model="newRole.name"
                  type="text"
                  required
                  placeholder="Contoh: bendahara_spp"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 font-mono"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Keterangan / Tugas Role</label>
                <textarea
                  v-model="newRole.description"
                  rows="2"
                  placeholder="Penjelasan wewenang role ini"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                ></textarea>
              </div>

              <div class="pt-3 border-t border-slate-100 flex justify-end gap-2">
                <button
                  type="button"
                  @click="showAddModal = false"
                  class="px-3.5 py-2 rounded-xl text-slate-500 hover:bg-slate-100 text-xs font-medium"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-medium shadow-sm"
                >
                  Simpan Role
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
  addOutline, shieldCheckmarkOutline, trashOutline, saveOutline,
  checkmarkCircleOutline, alertCircleOutline, sparklesOutline
} from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const roles = ref([])
const selectedRole = ref(null)
const permissions = ref([])
const showAddModal = ref(false)
const savingPerms = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const newRole = ref({
  name: '',
  label: '',
  description: ''
})

const fetchRoles = async () => {
  try {
    const res = await axios.get('/api/admin/roles', {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    roles.value = res.data || []
    if (roles.value.length > 0 && !selectedRole.value) {
      selectRole(roles.value[0])
    }
  } catch (err) {
    errorMsg.value = 'Gagal memuat daftar role'
  }
}

const selectRole = async (r) => {
  selectedRole.value = r
  try {
    const res = await axios.get(`/api/admin/roles/${r.id}/permissions`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    permissions.value = res.data || []
  } catch (err) {
    errorMsg.value = 'Gagal memuat izin menu role'
  }
}

const savePermissions = async () => {
  if (!selectedRole.value) return
  savingPerms.value = true
  successMsg.value = ''
  errorMsg.value = ''

  try {
    await axios.post(`/api/admin/roles/${selectedRole.value.id}/permissions`, permissions.value, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    successMsg.value = `Hak akses untuk role "${selectedRole.value.label}" berhasil disimpan!`
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan hak akses'
  } finally {
    savingPerms.value = false
  }
}

const handleAddRole = async () => {
  try {
    await axios.post('/api/admin/roles', newRole.value, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    showAddModal.value = false
    newRole.value = { name: '', label: '', description: '' }
    successMsg.value = 'Role baru berhasil dibuat!'
    await fetchRoles()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal membuat role'
  }
}

const handleDeleteRole = async (r) => {
  if (!confirm(`Yakin ingin menghapus role "${r.label}"?`)) return
  try {
    await axios.delete(`/api/admin/roles/${r.id}`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    successMsg.value = 'Role berhasil dihapus'
    selectedRole.value = null
    await fetchRoles()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menghapus role'
  }
}

onMounted(() => {
  fetchRoles()
})
</script>
