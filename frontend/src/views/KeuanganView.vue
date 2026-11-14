<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc] pt-14 md:pt-4">
      <div class="max-w-6xl mx-auto py-2">

        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h1 class="text-xl font-bold text-slate-800">Buku Kas & Pos Keuangan</h1>
            <p class="text-xs text-slate-500 mt-1">
              Catat mutasi arus kas umum (masuk/keluar) serta pantau saldo per pos anggaran operasional lembaga.
            </p>
          </div>

          <!-- Tombol Aksi Header -->
          <div class="flex items-center gap-2 self-start sm:self-auto flex-wrap">
            <button
              @click="openModalMasterPos"
              class="inline-flex items-center gap-2 px-3.5 py-2.5 rounded-xl bg-white hover:bg-slate-50 text-slate-700 border border-slate-200 text-xs font-bold shadow-sm transition-all"
            >
              <ion-icon :icon="folderOpenOutline" class="text-base text-amber-500"></ion-icon>
              Master Pos Keuangan
            </button>

            <button
              @click="openModalCatatKas"
              class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all"
            >
              <ion-icon :icon="addOutline" class="text-base"></ion-icon>
              Catat Kas Masuk / Keluar
            </button>
          </div>
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

        <!-- 3 Kartu Ringkasan Cashflow Utama -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
          <div class="bg-gradient-to-br from-emerald-50 to-teal-100/40 p-5 rounded-2xl border border-emerald-200/60 shadow-sm">
            <span class="text-[11px] font-bold text-emerald-800 uppercase tracking-wider block">Total Kas Masuk</span>
            <p class="text-2xl font-black text-emerald-700 mt-1">
              Rp {{ formatRupiah(summary.total_masuk) }}
            </p>
            <span class="text-[10px] text-slate-400 mt-0.5 block">Akumulasi seluruh pos penerimaan</span>
          </div>

          <div class="bg-gradient-to-br from-rose-50 to-red-100/40 p-5 rounded-2xl border border-rose-200/60 shadow-sm">
            <span class="text-[11px] font-bold text-rose-800 uppercase tracking-wider block">Total Kas Keluar</span>
            <p class="text-2xl font-black text-rose-700 mt-1">
              Rp {{ formatRupiah(summary.total_keluar) }}
            </p>
            <span class="text-[10px] text-slate-400 mt-0.5 block">Gaji, operasional, & sarana</span>
          </div>

          <div class="bg-gradient-to-br from-blue-50 to-sky-100/40 p-5 rounded-2xl border border-blue-200/60 shadow-sm">
            <span class="text-[11px] font-bold text-blue-800 uppercase tracking-wider block">Saldo Kas Saat Ini</span>
            <p class="text-2xl font-black text-blue-700 mt-1">
              Rp {{ formatRupiah(summary.saldo_kas) }}
            </p>
            <span class="text-[10px] text-slate-400 mt-0.5 block">Selisih kas masuk & keluar</span>
          </div>
        </div>

        <!-- SECTION: SALDO MASING-MASING POS KEUANGAN -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-5 mb-6">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h2 class="text-sm font-bold text-slate-800">Saldo Masing-masing Pos Keuangan</h2>
              <p class="text-[11px] text-slate-400">Rincian mutasi kas masuk & keluar per pos anggaran</p>
            </div>
            <button
              @click="openModalMasterPos"
              class="text-xs text-brand-600 hover:text-brand-700 font-bold flex items-center gap-1"
            >
              <span>+ Atur Pos</span>
            </button>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            <div
              v-for="pos in posList"
              :key="pos.id"
              class="p-3.5 rounded-xl border transition-all"
              :class="pos.tipe === 'masuk' ? 'bg-emerald-50/40 border-emerald-200/70' : 'bg-rose-50/40 border-rose-200/70'"
            >
              <div class="flex items-start justify-between gap-2 mb-1">
                <span class="font-bold text-xs text-slate-800 truncate" :title="pos.nama_pos">
                  {{ pos.nama_pos }}
                </span>
                <span
                  class="px-2 py-0.5 rounded-full text-[9px] font-black uppercase shrink-0"
                  :class="pos.tipe === 'masuk' ? 'bg-emerald-100 text-emerald-800' : 'bg-rose-100 text-rose-800'"
                >
                  {{ pos.tipe }}
                </span>
              </div>
              <p class="text-base font-black font-mono mt-2" :class="pos.tipe === 'masuk' ? 'text-emerald-700' : 'text-rose-700'">
                {{ pos.tipe === 'masuk' ? '+' : '-' }} Rp {{ formatRupiah(pos.tipe === 'masuk' ? pos.total_masuk : pos.total_keluar) }}
              </p>
              <span class="text-[10px] text-slate-400 block mt-0.5 truncate">{{ pos.keterangan || 'Pos anggaran umum' }}</span>
            </div>
          </div>
        </div>

        <!-- SECTION: TABEL BUKU KAS UMUM -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden mb-6">
          <div class="p-4 border-b border-slate-100 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div>
              <h2 class="text-sm font-bold text-slate-800">Riwayat Mutasi Buku Kas</h2>
              <p class="text-[11px] text-slate-400">100 transaksi operasional kas terbaru</p>
            </div>

            <!-- Filter Tipe Kas -->
            <div class="flex items-center gap-1.5 p-1 bg-slate-100 rounded-xl self-start sm:self-auto">
              <button
                @click="filterTipeKas = ''"
                class="px-3 py-1 rounded-lg text-xs font-bold transition-all"
                :class="filterTipeKas === '' ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
              >
                Semua
              </button>
              <button
                @click="filterTipeKas = 'masuk'"
                class="px-3 py-1 rounded-lg text-xs font-bold transition-all"
                :class="filterTipeKas === 'masuk' ? 'bg-white text-emerald-700 shadow-sm' : 'text-slate-500 hover:text-emerald-700'"
              >
                Kas Masuk
              </button>
              <button
                @click="filterTipeKas = 'keluar'"
                class="px-3 py-1 rounded-lg text-xs font-bold transition-all"
                :class="filterTipeKas === 'keluar' ? 'bg-white text-rose-700 shadow-sm' : 'text-slate-500 hover:text-rose-700'"
              >
                Kas Keluar
              </button>
            </div>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs text-slate-600 min-w-[650px]">
              <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                <tr>
                  <th class="py-3.5 px-4 w-36">Tanggal & Kode</th>
                  <th class="py-3.5 px-4">Pos Anggaran</th>
                  <th class="py-3.5 px-4">Rincian / Keterangan</th>
                  <th class="py-3.5 px-4 text-right">Nominal</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="k in filteredKasList" :key="k.id" class="hover:bg-slate-50/60 transition-colors">
                  <td class="py-3.5 px-4">
                    <span class="font-bold text-slate-800 block">{{ k.tanggal }}</span>
                    <span class="text-[10px] text-slate-400 font-mono">{{ k.kode_transaksi }}</span>
                  </td>

                  <td class="py-3.5 px-4">
                    <span class="font-semibold text-slate-800 block">{{ k.nama_pos }}</span>
                    <span
                      class="inline-block px-2 py-0.5 rounded-full text-[9px] font-bold uppercase mt-0.5"
                      :class="k.tipe === 'masuk' ? 'bg-emerald-100 text-emerald-800' : 'bg-rose-100 text-rose-800'"
                    >
                      Kas {{ k.tipe }}
                    </span>
                  </td>

                  <td class="py-3.5 px-4 text-slate-600 max-w-xs truncate">
                    {{ k.keterangan || '-' }}
                  </td>

                  <td class="py-3.5 px-4 text-right font-mono font-bold text-sm"
                      :class="k.tipe === 'masuk' ? 'text-emerald-600' : 'text-rose-600'">
                    {{ k.tipe === 'masuk' ? '+' : '-' }} Rp {{ formatRupiah(k.nominal) }}
                  </td>
                </tr>
                <tr v-if="filteredKasList.length === 0">
                  <td colspan="4" class="p-8 text-center text-slate-400 text-xs">
                    Belum ada transaksi kas yang dicatat.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- ==================================================== -->
        <!-- MODAL 1: CATAT KAS MASUK / KELUAR                    -->
        <!-- ==================================================== -->
        <div v-if="showModalKas" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-md w-full shadow-xl border border-slate-100">
            <h2 class="text-sm font-bold text-slate-800 mb-1">Catat Kas Masuk / Keluar</h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Masukkan transaksi penerimaan atau pengeluaran dana operasional lembaga.
            </p>

            <form @submit.prevent="handleSubmitKas" class="space-y-3.5">
              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Jenis Transaksi</label>
                <div class="grid grid-cols-2 gap-2">
                  <button
                    type="button"
                    @click="setTipeKasForm('masuk')"
                    class="py-2 text-xs font-bold rounded-xl border transition-all flex items-center justify-center gap-1.5"
                    :class="formKas.tipe === 'masuk' ? 'bg-emerald-500 text-white border-emerald-500 shadow-sm' : 'border-slate-200 text-slate-600'"
                  >
                    📥 Kas Masuk
                  </button>
                  <button
                    type="button"
                    @click="setTipeKasForm('keluar')"
                    class="py-2 text-xs font-bold rounded-xl border transition-all flex items-center justify-center gap-1.5"
                    :class="formKas.tipe === 'keluar' ? 'bg-rose-500 text-white border-rose-500 shadow-sm' : 'border-slate-200 text-slate-600'"
                  >
                    📤 Kas Keluar
                  </button>
                </div>
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Pos Anggaran</label>
                <select
                  v-model="formKas.pos_id"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                >
                  <option v-for="p in dropdownPosOptions" :key="p.id" :value="p.id">
                    {{ p.nama_pos }} ({{ p.tipe }})
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Tanggal Transaksi</label>
                <input
                  v-model="formKas.tanggal"
                  type="date"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Nominal (Rp)</label>
                <input
                  v-model.number="formKas.nominal"
                  type="number"
                  min="1"
                  required
                  placeholder="Contoh: 500000"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 font-mono font-bold"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Keterangan / Rincian</label>
                <textarea
                  v-model="formKas.keterangan"
                  rows="2"
                  placeholder="Catatan transaksi..."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                ></textarea>
              </div>

              <div class="pt-3 border-t border-slate-100 flex justify-end gap-2">
                <button
                  type="button"
                  @click="showModalKas = false"
                  class="px-3.5 py-2 rounded-xl text-slate-500 hover:bg-slate-100 text-xs font-medium"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  :disabled="isSubmittingKas"
                  class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-sm disabled:opacity-50"
                >
                  {{ isSubmittingKas ? 'Menyimpan...' : 'Simpan Transaksi' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- ==================================================== -->
        <!-- MODAL 2: MASTER POS KEUANGAN & SALDO                -->
        <!-- ==================================================== -->
        <div v-if="showModalPos" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-2xl w-full shadow-xl border border-slate-100 max-h-[90vh] overflow-y-auto">
            <div class="flex items-center justify-between pb-3 border-b border-slate-100 mb-4">
              <div>
                <h2 class="text-sm font-bold text-slate-800">Master Data Pos Keuangan</h2>
                <p class="text-[11px] text-slate-400">Kelola pos penerimaan kas masuk dan pos pengeluaran</p>
              </div>
              <button @click="showModalPos = false" class="text-slate-400 hover:text-slate-600 text-lg">✕</button>
            </div>

            <!-- Form Tambah Pos Baru -->
            <form @submit.prevent="handleSubmitPos" class="bg-slate-50 p-4 rounded-xl border border-slate-200/80 mb-5">
              <h3 class="text-xs font-bold text-slate-700 mb-2">
                {{ isEditPos ? 'Edit Pos Keuangan' : '+ Tambah Pos Keuangan Baru' }}
              </h3>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-2.5 mb-2.5">
                <div>
                  <label class="block text-[10px] font-bold text-slate-500 mb-1">Nama Pos</label>
                  <input
                    v-model="formPos.nama_pos"
                    type="text"
                    required
                    placeholder="Contoh: Infaq Gedung"
                    class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 bg-white focus:outline-none focus:ring-1 focus:ring-brand-500"
                  />
                </div>
                <div>
                  <label class="block text-[10px] font-bold text-slate-500 mb-1">Tipe Arus Kas</label>
                  <select
                    v-model="formPos.tipe"
                    class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 bg-white focus:outline-none focus:ring-1 focus:ring-brand-500"
                  >
                    <option value="masuk">Kas Masuk (Penerimaan)</option>
                    <option value="keluar">Kas Keluar (Pengeluaran)</option>
                  </select>
                </div>
                <div>
                  <label class="block text-[10px] font-bold text-slate-500 mb-1">Keterangan Singkat</label>
                  <input
                    v-model="formPos.keterangan"
                    type="text"
                    placeholder="Opsional..."
                    class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 bg-white focus:outline-none focus:ring-1 focus:ring-brand-500"
                  />
                </div>
              </div>
              <div class="flex items-center justify-end gap-2">
                <button
                  v-if="isEditPos"
                  type="button"
                  @click="cancelEditPos"
                  class="px-3 py-1 rounded-lg text-xs font-medium text-slate-500 hover:bg-slate-200"
                >
                  Batal Edit
                </button>
                <button
                  type="submit"
                  class="px-4 py-1.5 rounded-lg bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-sm"
                >
                  {{ isEditPos ? 'Perbarui Pos' : 'Simpan Pos' }}
                </button>
              </div>
            </form>

            <!-- Tabel Daftar Pos & Saldo -->
            <div class="border border-slate-200/80 rounded-xl overflow-hidden">
              <table class="w-full text-left text-xs text-slate-600">
                <thead class="bg-slate-50 text-[10px] uppercase font-bold text-slate-400 border-b border-slate-100">
                  <tr>
                    <th class="py-2.5 px-3">Nama Pos</th>
                    <th class="py-2.5 px-3">Tipe</th>
                    <th class="py-2.5 px-3 text-right">Akumulasi Saldo</th>
                    <th class="py-2.5 px-3 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr v-for="p in posList" :key="p.id" class="hover:bg-slate-50">
                    <td class="py-2.5 px-3 font-bold text-slate-800">
                      {{ p.nama_pos }}
                      <span class="text-[10px] text-slate-400 block font-normal">{{ p.keterangan }}</span>
                    </td>
                    <td class="py-2.5 px-3">
                      <span
                        class="px-2 py-0.5 rounded-full text-[9px] font-black uppercase"
                        :class="p.tipe === 'masuk' ? 'bg-emerald-100 text-emerald-800' : 'bg-rose-100 text-rose-800'"
                      >
                        {{ p.tipe }}
                      </span>
                    </td>
                    <td class="py-2.5 px-3 text-right font-mono font-bold" :class="p.tipe === 'masuk' ? 'text-emerald-700' : 'text-rose-700'">
                      Rp {{ formatRupiah(p.tipe === 'masuk' ? p.total_masuk : p.total_keluar) }}
                    </td>
                    <td class="py-2.5 px-3 text-right">
                      <div class="inline-flex items-center gap-1">
                        <button
                          @click="openEditPos(p)"
                          class="p-1 rounded text-slate-400 hover:text-brand-600 hover:bg-brand-50"
                          title="Edit"
                        >
                          <ion-icon :icon="createOutline" class="text-sm"></ion-icon>
                        </button>
                        <button
                          @click="handleDeletePos(p)"
                          class="p-1 rounded text-slate-400 hover:text-rose-600 hover:bg-rose-50"
                          title="Hapus"
                        >
                          <ion-icon :icon="trashOutline" class="text-sm"></ion-icon>
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
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
  addOutline, folderOpenOutline, createOutline, trashOutline,
  checkmarkCircleOutline, alertCircleOutline
} from 'ionicons/icons'
import { showConfirm, showAlert, showToast } from '@/utils/alert'

const summary = ref({ total_masuk: 0, total_keluar: 0, saldo_kas: 0 })
const kasList = ref([])
const posList = ref([])

const filterTipeKas = ref('')

const showModalKas = ref(false)
const showModalPos = ref(false)
const isSubmittingKas = ref(false)

const successMsg = ref('')
const errorMsg = ref('')

const formKas = ref({
  tipe: 'masuk',
  pos_id: null,
  tanggal: new Date().toISOString().slice(0, 10),
  nominal: '',
  keterangan: ''
})

const formPos = ref({
  nama_pos: '',
  tipe: 'masuk',
  keterangan: ''
})
const isEditPos = ref(false)
const editPosId = ref(null)

const formatRupiah = (val) => {
  if (!val) return '0'
  return Number(val).toLocaleString('id-ID')
}

const dropdownPosOptions = computed(() => {
  return posList.value.filter(p => p.tipe === formKas.value.tipe)
})

const filteredKasList = computed(() => {
  if (!filterTipeKas.value) return kasList.value
  return kasList.value.filter(k => k.tipe === filterTipeKas.value)
})

const fetchSummary = async () => {
  try {
    const res = await axios.get('/api/keuangan/summary')
    summary.value = res.data || { total_masuk: 0, total_keluar: 0, saldo_kas: 0 }
  } catch (err) {
    console.error(err)
  }
}

const fetchKas = async () => {
  try {
    const res = await axios.get('/api/keuangan/kas')
    kasList.value = res.data || []
  } catch (err) {
    console.error(err)
  }
}

const fetchPos = async () => {
  try {
    const res = await axios.get('/api/keuangan/pos')
    posList.value = res.data || []
    if (posList.value.length > 0 && !formKas.value.pos_id) {
      const firstMatch = posList.value.find(p => p.tipe === formKas.value.tipe)
      if (firstMatch) formKas.value.pos_id = firstMatch.id
    }
  } catch (err) {
    console.error(err)
  }
}

const setTipeKasForm = (tipe) => {
  formKas.value.tipe = tipe
  const firstMatch = posList.value.find(p => p.tipe === tipe)
  if (firstMatch) {
    formKas.value.pos_id = firstMatch.id
  }
}

const openModalCatatKas = () => {
  formKas.value.tanggal = new Date().toISOString().slice(0, 10)
  formKas.value.nominal = ''
  formKas.value.keterangan = ''
  setTipeKasForm('masuk')
  showModalKas.value = true
}

const openModalMasterPos = () => {
  isEditPos.value = false
  editPosId.value = null
  formPos.value = { nama_pos: '', tipe: 'masuk', keterangan: '' }
  showModalPos.value = true
}

const openEditPos = (pos) => {
  isEditPos.value = true
  editPosId.value = pos.id
  formPos.value = { nama_pos: pos.nama_pos, tipe: pos.tipe, keterangan: pos.keterangan }
}

const cancelEditPos = () => {
  isEditPos.value = false
  editPosId.value = null
  formPos.value = { nama_pos: '', tipe: 'masuk', keterangan: '' }
}

const handleSubmitKas = async () => {
  isSubmittingKas.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    await axios.post('/api/keuangan/kas', formKas.value)
    successMsg.value = 'Transaksi kas berhasil dicatat ke buku kas!'
    showModalKas.value = false
    await fetchSummary()
    await fetchKas()
    await fetchPos()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan transaksi'
  } finally {
    isSubmittingKas.value = false
    setTimeout(() => { successMsg.value = ''; errorMsg.value = '' }, 4000)
  }
}

const handleSubmitPos = async () => {
  successMsg.value = ''
  errorMsg.value = ''
  try {
    if (isEditPos.value) {
      await axios.put(`/api/keuangan/pos/${editPosId.value}`, formPos.value)
      successMsg.value = 'Pos keuangan berhasil diperbarui'
    } else {
      await axios.post('/api/keuangan/pos', formPos.value)
      successMsg.value = 'Pos keuangan baru berhasil ditambahkan'
    }
    cancelEditPos()
    await fetchPos()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan pos keuangan'
  } finally {
    setTimeout(() => { successMsg.value = ''; errorMsg.value = '' }, 4000)
  }
}

const handleDeletePos = async (p) => {
  const confirmed = await showConfirm({
    title: 'Hapus Pos Keuangan?',
    text: `Yakin ingin menghapus pos anggaran <b>"${p.nama_pos}"</b>?<br/><br/><span class="text-xs text-slate-500">Pos yang memiliki riwayat transaksi buku kas tidak dapat dihapus.</span>`,
    icon: 'warning',
    confirmButtonText: 'Ya, Hapus Pos',
    cancelButtonText: 'Batal',
    isDanger: true
  })
  if (!confirmed) return

  try {
    await axios.delete(`/api/keuangan/pos/${p.id}`)
    showToast('Pos keuangan berhasil dihapus', 'success')
    await fetchPos()
  } catch (err) {
    showAlert({
      title: 'Gagal Menghapus Pos',
      text: err.response?.data?.error || 'Gagal menghapus pos keuangan',
      icon: 'error'
    })
  }
}

onMounted(() => {
  fetchSummary()
  fetchKas()
  fetchPos()
})
</script>
