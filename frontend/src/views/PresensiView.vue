<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc] pt-14 md:pt-4">
      <div class="max-w-6xl mx-auto py-2">

        <!-- ==================================================== -->
        <!-- 1. TAMPILAN KHUSUS WALI SISWA (REKAP PRESENSI ANAK)  -->
        <!-- ==================================================== -->
        <div v-if="authStore.isWali">
          <div class="mb-6">
            <h1 class="text-xl font-bold text-slate-800">Laporan Kehadiran Siswa</h1>
            <p class="text-xs text-slate-500 mt-1">
              Rekapitulasi kehadiran harian an. <span class="font-bold text-emerald-700">{{ authStore.activeSiswa?.nama_lengkap }}</span> ({{ authStore.activeSiswa?.kelas }})
            </p>
          </div>

          <!-- 4 Ringkasan Angka Kehadiran -->
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
            <div class="bg-gradient-to-br from-emerald-50 to-teal-100/40 p-5 rounded-2xl border border-emerald-200/60 shadow-sm text-center">
              <span class="text-[11px] font-bold text-emerald-800 uppercase tracking-wider block">Hadir</span>
              <p class="text-3xl font-black text-emerald-700 mt-1">{{ rekapWali.hadir }} Hari</p>
              <span class="text-[10px] text-emerald-600 font-semibold mt-1 inline-block bg-white px-2 py-0.5 rounded-full border border-emerald-200">
                {{ rekapWali.persentase ? rekapWali.persentase.toFixed(1) : 0 }}% Kehadiran
              </span>
            </div>

            <div class="bg-gradient-to-br from-blue-50 to-sky-100/40 p-5 rounded-2xl border border-blue-200/60 shadow-sm text-center">
              <span class="text-[11px] font-bold text-blue-800 uppercase tracking-wider block">Izin</span>
              <p class="text-3xl font-black text-blue-700 mt-1">{{ rekapWali.izin }} Hari</p>
              <span class="text-[10px] text-slate-400 mt-1 block">Dengan Surat/Pemberitahuan</span>
            </div>

            <div class="bg-gradient-to-br from-amber-50 to-orange-100/40 p-5 rounded-2xl border border-amber-200/60 shadow-sm text-center">
              <span class="text-[11px] font-bold text-amber-800 uppercase tracking-wider block">Sakit</span>
              <p class="text-3xl font-black text-amber-700 mt-1">{{ rekapWali.sakit }} Hari</p>
              <span class="text-[10px] text-slate-400 mt-1 block">Kondisi Kesehatan</span>
            </div>

            <div class="bg-gradient-to-br from-rose-50 to-red-100/40 p-5 rounded-2xl border border-rose-200/60 shadow-sm text-center">
              <span class="text-[11px] font-bold text-rose-800 uppercase tracking-wider block">Alpa</span>
              <p class="text-3xl font-black text-rose-700 mt-1">{{ rekapWali.alpa }} Hari</p>
              <span class="text-[10px] text-slate-400 mt-1 block">Tanpa Keterangan</span>
            </div>
          </div>

          <!-- Riwayat Absensi Terakhir -->
          <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
            <div class="p-4 border-b border-slate-100 bg-slate-50/50">
              <h2 class="text-xs font-bold uppercase tracking-wider text-slate-600">Catatan Harian Kehadiran</h2>
            </div>
            <table class="w-full text-left text-xs text-slate-600">
              <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                <tr>
                  <th class="py-3 px-4">Tanggal</th>
                  <th class="py-3 px-4">Status Kehadiran</th>
                  <th class="py-3 px-4">Keterangan</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="(r, idx) in rekapWali.riwayat" :key="idx" class="hover:bg-slate-50/60">
                  <td class="py-3 px-4 font-bold text-slate-800">{{ r.tanggal }}</td>
                  <td class="py-3 px-4">
                    <span
                      class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                      :class="{
                        'bg-emerald-100 text-emerald-800': r.status === 'Hadir',
                        'bg-blue-100 text-blue-800': r.status === 'Izin',
                        'bg-amber-100 text-amber-800': r.status === 'Sakit',
                        'bg-rose-100 text-rose-800': r.status === 'Alpa'
                      }"
                    >
                      {{ r.status }}
                    </span>
                  </td>
                  <td class="py-3 px-4 text-slate-500">{{ r.keterangan || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- ==================================================== -->
        <!-- 2. TAMPILAN GURU / WALI KELAS (INPUT ABSENSI KELAS) -->
        <!-- ==================================================== -->
        <div v-else>
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
            <div>
              <h1 class="text-xl font-bold text-slate-800">Presensi & Absensi Harian</h1>
              <p class="text-xs text-slate-500 mt-1">
                Catat kehadiran siswa per rombel/kelas dengan tombol cepat.
              </p>
            </div>
            <button
              @click="handleSavePresensi"
              :disabled="saving"
              class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all self-start sm:self-auto disabled:opacity-60"
            >
              <ion-icon :icon="saveOutline" class="text-base"></ion-icon>
              <span>{{ saving ? 'Menyimpan...' : 'Simpan Presensi' }}</span>
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

          <!-- Bar Filter Kelas, Tanggal & Tombol Massal -->
          <div class="bg-white p-4 rounded-2xl border border-slate-200/80 shadow-sm mb-6 flex flex-col sm:flex-row items-center justify-between gap-4">
            <div class="flex flex-col sm:flex-row items-center gap-3 w-full sm:w-auto">
              <div>
                <label class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1">Pilih Kelas</label>
                <select
                  v-model="filterKelas"
                  @change="fetchPresensiKelas"
                  class="px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-bold text-slate-700 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                >
                  <option v-for="k in kelasList" :key="k.id" :value="k.nama_kelas">
                    {{ k.nama_kelas }} ({{ k.tingkat }})
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1">Tanggal Absensi</label>
                <input
                  v-model="filterTanggal"
                  type="date"
                  @change="fetchPresensiKelas"
                  class="px-3 py-2 text-xs rounded-xl border border-slate-200 font-bold text-slate-700 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>
            </div>

            <div class="flex items-center gap-2 self-stretch sm:self-end">
              <button
                @click="setAllStatus('Hadir')"
                class="flex-1 sm:flex-initial px-3 py-2 rounded-xl bg-emerald-50 hover:bg-emerald-100 text-emerald-700 font-bold text-xs border border-emerald-200 transition-colors"
              >
                ✓ Semua Hadir
              </button>
            </div>
          </div>

          <!-- Tabel Presensi Siswa Satu Kelas -->
          <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
            <table class="w-full text-left text-xs text-slate-600">
              <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                <tr>
                  <th class="py-3.5 px-4">Nama Siswa</th>
                  <th class="py-3.5 px-4">NIS</th>
                  <th class="py-3.5 px-4 text-center">Status Kehadiran</th>
                  <th class="py-3.5 px-4">Catatan / Keterangan</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="item in presensiList" :key="item.siswa_id" class="hover:bg-slate-50/60 transition-colors">
                  <td class="py-3.5 px-4">
                    <span class="font-bold text-slate-800 text-xs block">{{ item.nama_lengkap }}</span>
                    <span class="text-[10px] text-slate-400">{{ item.kelas }}</span>
                  </td>

                  <td class="py-3.5 px-4 font-mono font-medium text-slate-700">
                    {{ item.nis }}
                  </td>

                  <!-- Tombol Pilihan Status (Hadir, Izin, Sakit, Alpa) -->
                  <td class="py-3.5 px-4 text-center">
                    <div class="inline-flex rounded-xl p-1 bg-slate-100 border border-slate-200/60 gap-1">
                      <button
                        type="button"
                        @click="item.status = 'Hadir'"
                        class="px-2.5 py-1 text-[11px] font-bold rounded-lg transition-all"
                        :class="item.status === 'Hadir' ? 'bg-emerald-500 text-white shadow-sm' : 'text-slate-600 hover:text-slate-900'"
                      >
                        Hadir
                      </button>
                      <button
                        type="button"
                        @click="item.status = 'Izin'"
                        class="px-2.5 py-1 text-[11px] font-bold rounded-lg transition-all"
                        :class="item.status === 'Izin' ? 'bg-blue-500 text-white shadow-sm' : 'text-slate-600 hover:text-slate-900'"
                      >
                        Izin
                      </button>
                      <button
                        type="button"
                        @click="item.status = 'Sakit'"
                        class="px-2.5 py-1 text-[11px] font-bold rounded-lg transition-all"
                        :class="item.status === 'Sakit' ? 'bg-amber-500 text-white shadow-sm' : 'text-slate-600 hover:text-slate-900'"
                      >
                        Sakit
                      </button>
                      <button
                        type="button"
                        @click="item.status = 'Alpa'"
                        class="px-2.5 py-1 text-[11px] font-bold rounded-lg transition-all"
                        :class="item.status === 'Alpa' ? 'bg-rose-500 text-white shadow-sm' : 'text-slate-600 hover:text-slate-900'"
                      >
                        Alpa
                      </button>
                    </div>
                  </td>

                  <td class="py-3.5 px-4">
                    <input
                      v-model="item.keterangan"
                      type="text"
                      placeholder="Contoh: Mengikuti lomba / Demam"
                      class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500"
                    />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

      </div>
    </ion-content>
  </ion-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import axios from 'axios'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import { saveOutline, checkmarkCircleOutline, alertCircleOutline } from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// State Guru / Staf
const kelasList = ref([])
const filterKelas = ref('')
const filterTanggal = ref(new Date().toISOString().slice(0, 10))
const presensiList = ref([])
const saving = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

// State Wali Murid
const rekapWali = ref({
  total_hari: 0,
  hadir: 0,
  izin: 0,
  sakit: 0,
  alpa: 0,
  persentase: 0,
  riwayat: []
})

const fetchKelas = async () => {
  try {
    const res = await axios.get('/api/admin/kelas')
    kelasList.value = res.data || []
    if (kelasList.value.length > 0 && !filterKelas.value) {
      filterKelas.value = kelasList.value[0].nama_kelas
      await fetchPresensiKelas()
    }
  } catch (err) {
    console.error(err)
  }
}

const fetchPresensiKelas = async () => {
  if (!filterKelas.value) return
  try {
    const res = await axios.get(`/api/presensi/kelas?kelas=${filterKelas.value}&tanggal=${filterTanggal.value}`)
    presensiList.value = res.data || []
  } catch (err) {
    errorMsg.value = 'Gagal memuat presensi kelas'
  }
}

const setAllStatus = (st) => {
  presensiList.value.forEach(item => {
    item.status = st
  })
}

const handleSavePresensi = async () => {
  saving.value = true
  successMsg.value = ''
  errorMsg.value = ''

  try {
    const payload = {
      tanggal: filterTanggal.value,
      items: presensiList.value.map(p => ({
        siswa_id: p.siswa_id,
        status: p.status,
        keterangan: p.keterangan
      }))
    }
    await axios.post('/api/presensi/batch', payload)
    successMsg.value = `Presensi kelas ${filterKelas.value} tanggal ${filterTanggal.value} berhasil disimpan!`
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan presensi'
  } finally {
    saving.value = false
  }
}

const fetchRekapWali = async () => {
  if (!authStore.activeSiswa) return
  try {
    const res = await axios.get(`/api/presensi/wali?siswa_id=${authStore.activeSiswa.id}`)
    rekapWali.value = res.data || { total_hari: 0, hadir: 0, izin: 0, sakit: 0, alpa: 0, persentase: 0, riwayat: [] }
  } catch (err) {
    console.error(err)
  }
}

watch(() => authStore.activeSiswa, () => {
  if (authStore.isWali) {
    fetchRekapWali()
  }
})

onMounted(() => {
  if (authStore.isWali) {
    fetchRekapWali()
  } else {
    fetchKelas()
  }
})
</script>
