<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc] pt-14 md:pt-4">
      <div class="max-w-6xl mx-auto py-2">

        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h1 class="text-xl font-bold text-slate-800">Manajemen Keuangan</h1>
            <p class="text-xs text-slate-500 mt-1">
              Buku kas umum (kas masuk/keluar), pos anggaran, dan penagihan SPP siswa.
            </p>
          </div>

          <!-- Tab Navigasi Sub-Modul (Ringkasan, Kas Masuk/Keluar, Tagihan SPP) -->
          <div class="inline-flex bg-white p-1 rounded-2xl border border-slate-200/80 shadow-sm self-start sm:self-auto">
            <button
              @click="activeSubTab = 'tagihan'"
              class="px-4 py-1.5 rounded-xl text-xs font-bold transition-all"
              :class="activeSubTab === 'tagihan' ? 'bg-amber-400 text-amber-950 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
            >
              Tagihan SPP
            </button>
            <button
              @click="activeSubTab = 'kas'"
              class="px-4 py-1.5 rounded-xl text-xs font-bold transition-all"
              :class="activeSubTab === 'kas' ? 'bg-amber-400 text-amber-950 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
            >
              Buku Kas Masuk / Keluar
            </button>
          </div>
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
            <span class="text-[10px] text-slate-400 mt-0.5 block">Gaji, operasional, & utilitas</span>
          </div>

          <div class="bg-gradient-to-br from-blue-50 to-sky-100/40 p-5 rounded-2xl border border-blue-200/60 shadow-sm">
            <span class="text-[11px] font-bold text-blue-800 uppercase tracking-wider block">Saldo Kas Saat Ini</span>
            <p class="text-2xl font-black text-blue-700 mt-1">
              Rp {{ formatRupiah(summary.saldo_kas) }}
            </p>
            <span class="text-[10px] text-slate-400 mt-0.5 block">Selisih kas masuk & keluar</span>
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

        <!-- ==================================================== -->
        <!-- SUB-TAB 1: TAGIHAN SISWA / SPP                       -->
        <!-- ==================================================== -->
        <div v-if="activeSubTab === 'tagihan'" class="space-y-4">
          
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <h2 class="text-sm font-bold text-slate-800">Daftar Tagihan SPP & Iuran Siswa</h2>
            <!-- Filter jika akun wali vs admin -->
            <span v-if="authStore.isWali" class="text-xs font-semibold text-emerald-700 bg-emerald-50 px-3 py-1 rounded-full border border-emerald-200">
              Menampilkan tagihan siswa: {{ authStore.activeSiswa?.nama_lengkap }}
            </span>
          </div>

          <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
            <table class="w-full text-left text-xs text-slate-600">
              <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                <tr>
                  <th class="py-3.5 px-4">Nama Siswa</th>
                  <th class="py-3.5 px-4">Jenis Tagihan</th>
                  <th class="py-3.5 px-4">Periode</th>
                  <th class="py-3.5 px-4">Nominal</th>
                  <th class="py-3.5 px-4 text-center">Status</th>
                  <th class="py-3.5 px-4 text-right" v-if="!authStore.isWali">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="tag in filteredTagihan" :key="tag.id" class="hover:bg-slate-50/60 transition-colors">
                  <td class="py-3.5 px-4">
                    <span class="font-bold text-slate-800 block">{{ tag.nama_siswa }}</span>
                    <span class="text-[10px] text-slate-400 font-mono">NIS: {{ tag.nis }} • {{ tag.kelas }}</span>
                  </td>

                  <td class="py-3.5 px-4">
                    <span class="font-semibold text-slate-700 block">{{ tag.nama_tagihan }}</span>
                    <span class="text-[10px] px-2 py-0.5 rounded-full bg-slate-100 text-slate-500 font-medium">
                      {{ tag.tipe_periode }}
                    </span>
                  </td>

                  <td class="py-3.5 px-4 font-medium text-slate-700">
                    {{ tag.periode }} ({{ tag.nama_tahun }})
                  </td>

                  <td class="py-3.5 px-4 font-mono font-bold text-slate-800">
                    Rp {{ formatRupiah(tag.nominal) }}
                  </td>

                  <td class="py-3.5 px-4 text-center">
                    <span
                      v-if="tag.status === 'lunas'"
                      class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800 border border-emerald-200"
                    >
                      Lunas
                    </span>
                    <span
                      v-else
                      class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-amber-100 text-amber-800 border border-amber-200"
                    >
                      Belum Lunas
                    </span>
                  </td>

                  <td class="py-3.5 px-4 text-right" v-if="!authStore.isWali">
                    <button
                      v-if="tag.status === 'belum_lunas'"
                      @click="handleBayar(tag)"
                      class="px-3 py-1.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white font-bold text-[11px] shadow-sm transition-all"
                    >
                      Terima Bayar
                    </button>
                    <span v-else class="text-[10px] text-slate-400 font-medium">
                      {{ tag.tanggal_bayar.slice(0, 10) }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- ==================================================== -->
        <!-- SUB-TAB 2: BUKU KAS MASUK / KELUAR                   -->
        <!-- ==================================================== -->
        <div v-if="activeSubTab === 'kas'" class="space-y-4">
          
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <h2 class="text-sm font-bold text-slate-800">Buku Kas Umum</h2>
            <button
              @click="showModalKas = true"
              class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all self-start sm:self-auto"
            >
              <ion-icon :icon="addOutline" class="text-base"></ion-icon>
              Catat Transaksi Kas
            </button>
          </div>

          <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
            <table class="w-full text-left text-xs text-slate-600">
              <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                <tr>
                  <th class="py-3.5 px-4">Tanggal & Kode</th>
                  <th class="py-3.5 px-4">Pos Anggaran</th>
                  <th class="py-3.5 px-4">Keterangan</th>
                  <th class="py-3.5 px-4 text-right">Nominal</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="k in kasList" :key="k.id" class="hover:bg-slate-50/60 transition-colors">
                  <td class="py-3.5 px-4">
                    <span class="font-bold text-slate-800 block">{{ k.tanggal }}</span>
                    <span class="text-[10px] text-slate-400 font-mono">{{ k.kode_transaksi }}</span>
                  </td>

                  <td class="py-3.5 px-4">
                    <span class="font-semibold text-slate-700 block">{{ k.nama_pos }}</span>
                    <span
                      class="inline-block px-2 py-0.5 rounded-full text-[9px] font-bold uppercase"
                      :class="k.tipe === 'masuk' ? 'bg-emerald-100 text-emerald-800' : 'bg-rose-100 text-rose-800'"
                    >
                      Kas {{ k.tipe }}
                    </span>
                  </td>

                  <td class="py-3.5 px-4 text-slate-600 max-w-xs truncate">
                    {{ k.keterangan || '-' }}
                  </td>

                  <td class="py-3.5 px-4 text-right font-mono font-bold"
                      :class="k.tipe === 'masuk' ? 'text-emerald-600' : 'text-rose-600'">
                    {{ k.tipe === 'masuk' ? '+' : '-' }} Rp {{ formatRupiah(k.nominal) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Modal Tambah Transaksi Kas -->
        <div v-if="showModalKas" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-md w-full shadow-xl border border-slate-100">
            <h2 class="text-sm font-bold text-slate-800 mb-1">Catat Kas Masuk / Keluar</h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Masukkan transaksi penerimaan atau pengeluaran dana operasional.
            </p>

            <form @submit.prevent="handleSubmitKas" class="space-y-3.5">
              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Jenis Transaksi</label>
                <div class="grid grid-cols-2 gap-2">
                  <button
                    type="button"
                    @click="formKas.tipe = 'masuk'"
                    class="py-2 text-xs font-bold rounded-xl border transition-all"
                    :class="formKas.tipe === 'masuk' ? 'bg-emerald-500 text-white border-emerald-500' : 'border-slate-200 text-slate-600'"
                  >
                    📥 Kas Masuk
                  </button>
                  <button
                    type="button"
                    @click="formKas.tipe = 'keluar'"
                    class="py-2 text-xs font-bold rounded-xl border transition-all"
                    :class="formKas.tipe === 'keluar' ? 'bg-rose-500 text-white border-rose-500' : 'border-slate-200 text-slate-600'"
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
                  <option v-for="p in filteredPos" :key="p.id" :value="p.id">
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
                  v-model="formKas.nominal"
                  type="number"
                  required
                  placeholder="Contoh: 500000"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 font-mono"
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
                  class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-sm"
                >
                  Simpan Transaksi
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
import { addOutline, checkmarkCircleOutline, alertCircleOutline } from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const activeSubTab = ref('tagihan')
const summary = ref({ total_masuk: 0, total_keluar: 0, saldo_kas: 0 })
const tagihanList = ref([])
const kasList = ref([])
const posList = ref([])

const showModalKas = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const formKas = ref({
  tipe: 'masuk',
  pos_id: 1,
  tanggal: new Date().toISOString().slice(0, 10),
  nominal: '',
  keterangan: ''
})

const formatRupiah = (val) => {
  if (!val) return '0'
  return Number(val).toLocaleString('id-ID')
}

const fetchSummary = async () => {
  try {
    const res = await axios.get('/api/keuangan/summary')
    summary.value = res.data || { total_masuk: 0, total_keluar: 0, saldo_kas: 0 }
  } catch (err) {
    console.error(err)
  }
}

const fetchTagihan = async () => {
  try {
    const url = authStore.isWali && authStore.activeSiswa
      ? `/api/keuangan/tagihan?siswa_id=${authStore.activeSiswa.id}`
      : '/api/keuangan/tagihan'
    const res = await axios.get(url)
    tagihanList.value = res.data || []
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
    if (posList.value.length > 0) {
      formKas.value.pos_id = posList.value[0].id
    }
  } catch (err) {
    console.error(err)
  }
}

const filteredPos = computed(() => {
  return posList.value.filter(p => p.tipe === formKas.value.tipe)
})

const filteredTagihan = computed(() => {
  if (authStore.isWali && authStore.activeSiswa) {
    return tagihanList.value.filter(t => t.siswa_id === authStore.activeSiswa.id)
  }
  return tagihanList.value
})

const handleBayar = async (tag) => {
  if (!confirm(`Konfirmasi pembayaran tagihan ${tag.nama_tagihan} (${tag.periode}) an. ${tag.nama_siswa} sebesar Rp ${formatRupiah(tag.nominal)}?`)) return
  try {
    await axios.post(`/api/keuangan/tagihan/${tag.id}/bayar`)
    successMsg.value = `Pembayaran ${tag.nama_tagihan} berhasil dicatat & masuk kas!`
    await fetchTagihan()
    await fetchSummary()
    await fetchKas()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal memproses pembayaran'
  }
}

const handleSubmitKas = async () => {
  try {
    await axios.post('/api/keuangan/kas', formKas.value)
    successMsg.value = 'Transaksi kas berhasil dicatat!'
    showModalKas.value = false
    formKas.value.nominal = ''
    formKas.value.keterangan = ''
    await fetchSummary()
    await fetchKas()
    setTimeout(() => { successMsg.value = '' }, 4000)
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan transaksi'
  }
}

onMounted(() => {
  fetchSummary()
  fetchTagihan()
  fetchKas()
  fetchPos()
})
</script>
