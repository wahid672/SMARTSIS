<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc] pt-14 md:pt-4">
      <div class="max-w-6xl mx-auto py-2">

        <!-- ==================================================== -->
        <!-- 1. TAMPILAN KHUSUS WALI SISWA (STATUS TAGIHAN ANAK)  -->
        <!-- ==================================================== -->
        <div v-if="authStore.isWali">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
            <div>
              <h1 class="text-xl font-bold text-slate-800">Tagihan & Pembayaran Siswa</h1>
              <p class="text-xs text-slate-500 mt-1">
                Informasi status pembayaran SPP dan iuran an. 
                <span class="font-bold text-brand-600">{{ authStore.activeSiswa?.nama_lengkap }}</span> 
                ({{ authStore.activeSiswa?.kelas }})
              </p>
            </div>
          </div>

          <!-- 3 Ringkasan Angka Tagihan Wali -->
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
            <div class="bg-gradient-to-br from-emerald-50 to-teal-100/40 p-5 rounded-2xl border border-emerald-200/60 shadow-sm text-center">
              <span class="text-[11px] font-bold text-emerald-800 uppercase tracking-wider block">Sudah Terbayar</span>
              <p class="text-2xl font-black text-emerald-700 mt-1">
                Rp {{ formatRupiah(waliSummary.terbayar) }}
              </p>
              <span class="text-[10px] text-emerald-600 font-semibold mt-1 inline-block">Lunas</span>
            </div>

            <div class="bg-gradient-to-br from-amber-50 to-orange-100/40 p-5 rounded-2xl border border-amber-200/60 shadow-sm text-center">
              <span class="text-[11px] font-bold text-amber-800 uppercase tracking-wider block">Sisa Tagihan</span>
              <p class="text-2xl font-black text-amber-700 mt-1">
                Rp {{ formatRupiah(waliSummary.tunggakan) }}
              </p>
              <span class="text-[10px] text-amber-600 font-semibold mt-1 inline-block">Belum Lunas</span>
            </div>

            <div class="bg-gradient-to-br from-blue-50 to-sky-100/40 p-5 rounded-2xl border border-blue-200/60 shadow-sm text-center">
              <span class="text-[11px] font-bold text-blue-800 uppercase tracking-wider block">Total Tagihan Siswa</span>
              <p class="text-2xl font-black text-blue-700 mt-1">
                Rp {{ formatRupiah(waliSummary.total) }}
              </p>
              <span class="text-[10px] text-slate-400 mt-1 block">Tahun Ajaran Aktif</span>
            </div>
          </div>

          <!-- Tabel Riwayat Tagihan Siswa untuk Wali -->
          <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
            <div class="p-4 border-b border-slate-100">
              <h2 class="text-xs font-bold uppercase tracking-wider text-slate-700">Daftar Tagihan & Status Pembayaran</h2>
            </div>
            <div class="overflow-x-auto">
              <table class="w-full text-left text-xs text-slate-600 min-w-[600px]">
                <thead class="bg-slate-50 text-[10px] uppercase font-bold text-slate-400 border-b border-slate-100">
                  <tr>
                    <th class="py-3 px-4">Jenis Tagihan</th>
                    <th class="py-3 px-4">Bulan / Periode</th>
                    <th class="py-3 px-4">Nominal</th>
                    <th class="py-3 px-4 text-center">Status</th>
                    <th class="py-3 px-4">Tanggal Pembayaran</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr v-for="t in waliTagihanList" :key="t.id" class="hover:bg-slate-50">
                    <td class="py-3 px-4 font-bold text-slate-800">{{ t.nama_tagihan }}</td>
                    <td class="py-3 px-4 font-semibold text-slate-700">{{ t.periode }}</td>
                    <td class="py-3 px-4 font-mono font-bold text-slate-800">Rp {{ formatRupiah(t.nominal) }}</td>
                    <td class="py-3 px-4 text-center">
                      <span
                        class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                        :class="t.status === 'lunas' ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'"
                      >
                        {{ t.status === 'lunas' ? 'Lunas' : 'Belum Lunas' }}
                      </span>
                    </td>
                    <td class="py-3 px-4 text-slate-500 font-mono text-[11px]">
                      {{ t.tanggal_bayar ? t.tanggal_bayar.slice(0, 10) : '-' }}
                    </td>
                  </tr>
                  <tr v-if="waliTagihanList.length === 0">
                    <td colspan="5" class="p-8 text-center text-slate-400 text-xs">Belum ada data tagihan yang diterbitkan.</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- ==================================================== -->
        <!-- 2. TAMPILAN ADMIN / BENDAHARA (SETTING & KASIR)     -->
        <!-- ==================================================== -->
        <div v-else>

          <!-- Header Admin -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
            <div>
              <h1 class="text-xl font-bold text-slate-800">Manajemen Tagihan Siswa</h1>
              <p class="text-xs text-slate-500 mt-1">
                Atur jenis pos tagihan (Bulanan / Bebas), tetapkan tarif per kelas, dan proses pembayaran kasir.
              </p>
            </div>

            <!-- Switcher Tab Admin -->
            <div class="flex items-center gap-1.5 p-1 bg-slate-200/80 rounded-xl self-start sm:self-auto">
              <button
                @click="adminMode = 'jenis'"
                class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all"
                :class="adminMode === 'jenis' || adminMode === 'tarif' ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
              >
                Pos Tagihan & Tarif
              </button>
              <button
                @click="adminMode = 'kasir'"
                class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all"
                :class="adminMode === 'kasir' ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
              >
                Kasir & Rekap Pembayaran
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

          <!-- ==================================================== -->
          <!-- VIEW 1: DAFTAR JENIS TAGIHAN (SESUAI SCREENSHOT 2)   -->
          <!-- ==================================================== -->
          <div v-if="adminMode === 'jenis'" class="space-y-4">
            <div class="flex items-center justify-between">
              <h2 class="text-sm font-bold text-slate-800">Daftar Jenis Pembayaran / Tagihan</h2>
              <button
                @click="openAddJenisModal"
                class="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all"
              >
                <ion-icon :icon="addOutline" class="text-base"></ion-icon>
                Tambah Jenis Pembayaran
              </button>
            </div>

            <!-- Tabel Jenis Tagihan (Sesuai Screenshot 2) -->
            <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
              <div class="overflow-x-auto">
                <table class="w-full text-left text-xs text-slate-600 min-w-[700px]">
                  <thead class="bg-slate-50 text-[10px] uppercase font-bold text-slate-400 border-b border-slate-100">
                    <tr>
                      <th class="py-3 px-4 w-12 text-center">No</th>
                      <th class="py-3 px-4">POS</th>
                      <th class="py-3 px-4">Nama Pembayaran</th>
                      <th class="py-3 px-4 text-center">Tipe</th>
                      <th class="py-3 px-4 text-center">Tahun</th>
                      <th class="py-3 px-4 text-center">Tarif Pembayaran</th>
                      <th class="py-3 px-4 text-right">Aksi</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr v-for="(p, idx) in posTagihanList" :key="p.id" class="hover:bg-slate-50/60">
                      <td class="py-3.5 px-4 text-center font-mono text-slate-400">{{ idx + 1 }}</td>
                      <td class="py-3.5 px-4 font-bold text-slate-800">{{ p.nama_pos }}</td>
                      <td class="py-3.5 px-4 font-semibold text-slate-700">{{ p.nama_tagihan }}</td>
                      <td class="py-3.5 px-4 text-center">
                        <span
                          class="px-2.5 py-0.5 rounded-full text-[10px] font-bold capitalize"
                          :class="p.tipe_bayar === 'bulanan' ? 'bg-sky-100 text-sky-800' : 'bg-purple-100 text-purple-800'"
                        >
                          {{ p.tipe_bayar }}
                        </span>
                      </td>
                      <td class="py-3.5 px-4 text-center font-semibold text-slate-600">{{ p.nama_tahun }}</td>
                      <td class="py-3.5 px-4 text-center">
                        <button
                          @click="openSettingTarif(p)"
                          class="px-3.5 py-1.5 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs shadow-sm transition-all inline-flex items-center gap-1.5"
                        >
                          <ion-icon :icon="settingsOutline" class="text-sm"></ion-icon>
                          Setting Tarif Pembayaran
                        </button>
                      </td>
                      <td class="py-3.5 px-4 text-right">
                        <div class="inline-flex items-center gap-1.5">
                          <button
                            @click="openEditJenisModal(p)"
                            class="p-1.5 rounded-lg text-slate-400 hover:text-amber-600 hover:bg-amber-50"
                            title="Edit"
                          >
                            <ion-icon :icon="createOutline" class="text-base text-amber-500"></ion-icon>
                          </button>
                          <button
                            @click="handleDeleteJenis(p)"
                            class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50"
                            title="Hapus"
                          >
                            <ion-icon :icon="trashOutline" class="text-base"></ion-icon>
                          </button>
                        </div>
                      </td>
                    </tr>
                    <tr v-if="posTagihanList.length === 0">
                      <td colspan="7" class="p-8 text-center text-slate-400 text-xs">
                        Belum ada pos tagihan pembayaran. Klik tombol "Tambah Jenis Pembayaran" untuk memulai.
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <!-- ==================================================== -->
          <!-- VIEW 2: SETTING TARIF KELAS (SESUAI SCREENSHOT 3 & 5) -->
          <!-- ==================================================== -->
          <div v-else-if="adminMode === 'tarif'" class="space-y-4">
            
            <!-- Top Card Navigasi Tarif (Screenshot 3) -->
            <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm space-y-4">
              <div class="flex items-center justify-between border-b border-slate-100 pb-3">
                <div>
                  <h2 class="text-sm font-bold text-slate-800">
                    Tarif - {{ currentPosTagihan?.nama_tagihan }}
                  </h2>
                  <p class="text-[11px] text-slate-400">
                    POS: {{ currentPosTagihan?.nama_pos }} • Tipe: {{ currentPosTagihan?.tipe_bayar }}
                  </p>
                </div>
                <button
                  @click="adminMode = 'jenis'"
                  class="px-3.5 py-1.5 rounded-xl border border-slate-200 hover:bg-slate-50 text-xs font-bold text-slate-600 flex items-center gap-1.5"
                >
                  <ion-icon :icon="arrowBackOutline" class="text-sm"></ion-icon>
                  Kembali
                </button>
              </div>

              <!-- Baris Filter Kelas & Tahun Ajaran -->
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 items-end">
                <div>
                  <label class="block text-[11px] font-bold text-slate-600 mb-1">Tahun Ajaran</label>
                  <input
                    type="text"
                    disabled
                    :value="currentPosTagihan?.nama_tahun"
                    class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-100 text-slate-700 font-bold"
                  />
                </div>

                <div>
                  <label class="block text-[11px] font-bold text-slate-600 mb-1">Kelas</label>
                  <select
                    v-model="selectedTarifKelas"
                    class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-bold text-slate-700 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                  >
                    <option v-for="k in kelasList" :key="k.id" :value="k.nama_kelas">
                      {{ k.nama_kelas }} ({{ k.tingkat }})
                    </option>
                  </select>
                </div>

                <div>
                  <button
                    @click="loadTarifKelas"
                    class="w-full px-4 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold shadow-sm flex items-center justify-center gap-1.5"
                  >
                    <ion-icon :icon="searchOutline" class="text-sm"></ion-icon>
                    Cari / Tampilkan
                  </button>
                </div>
              </div>

              <!-- Tombol Setting Tarif (Screenshot 3 & 5) -->
              <div class="flex items-center gap-2 pt-2 border-t border-slate-100 flex-wrap">
                <span class="text-xs font-bold text-slate-600 mr-2">Setting Tarif:</span>
                <button
                  @click="openModalSetKelas"
                  class="px-3.5 py-1.5 rounded-xl bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold shadow-sm flex items-center gap-1.5"
                >
                  <span>+ Berdasarkan Kelas</span>
                </button>
                <button
                  @click="handleHapusSemuaTarif"
                  class="px-3.5 py-1.5 rounded-xl bg-slate-700 hover:bg-slate-800 text-white text-xs font-bold shadow-sm flex items-center gap-1.5 ml-auto"
                >
                  <ion-icon :icon="trashOutline" class="text-sm"></ion-icon>
                  Hapus Semua Tagihan Kelas Ini
                </button>
              </div>
            </div>

            <!-- Tabel Siswa & Tarif Tagihan (Screenshot 5) -->
            <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
              <div v-if="isLoadingTarif" class="p-8 text-center text-xs text-slate-400">
                Memuat data tagihan siswa kelas {{ selectedTarifKelas }}...
              </div>
              <div v-else class="overflow-x-auto">
                <table class="w-full text-left text-xs text-slate-600 min-w-[750px]">
                  <thead class="bg-slate-50 text-[10px] uppercase font-bold text-slate-400 border-b border-slate-100">
                    <tr>
                      <th class="py-3 px-4 w-12 text-center">No</th>
                      <th class="py-3 px-4">NIS</th>
                      <th class="py-3 px-4">Nama Siswa</th>
                      <th class="py-3 px-4">Kelas</th>
                      <th class="py-3 px-4">Total Tagihan</th>
                      <th class="py-3 px-4">Rincian Bulan / Periode</th>
                      <th class="py-3 px-4 text-center">Aksi</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr v-for="(s, idx) in tarifStudents" :key="s.siswa_id" class="hover:bg-slate-50/60">
                      <td class="py-3 px-4 text-center font-mono text-slate-400">{{ idx + 1 }}</td>
                      <td class="py-3 px-4 font-mono font-semibold">{{ s.nis }}</td>
                      <td class="py-3 px-4 font-bold text-slate-800">{{ s.nama_lengkap }}</td>
                      <td class="py-3 px-4 font-semibold text-slate-600">{{ s.kelas }}</td>
                      <td class="py-3 px-4">
                        <span v-if="s.total_tagihan > 0" class="font-mono font-black text-slate-800">
                          Rp {{ formatRupiah(s.total_tagihan) }}
                        </span>
                        <span v-else class="px-2 py-0.5 rounded text-[10px] font-semibold bg-slate-100 text-slate-400">
                          Belum Diset (Rp 0)
                        </span>
                      </td>
                      <td class="py-3 px-4">
                        <div v-if="s.bulan_list && s.bulan_list.length > 0" class="flex flex-wrap gap-1 max-w-md">
                          <span
                            v-for="b in s.bulan_list"
                            :key="b.id"
                            class="px-2 py-0.5 rounded text-[10px] font-mono border"
                            :class="b.status === 'lunas' ? 'bg-emerald-50 text-emerald-800 border-emerald-200 font-bold' : 'bg-slate-50 text-slate-700 border-slate-200'"
                            :title="`${b.bulan}: Rp ${formatRupiah(b.nominal)} (${b.status})`"
                          >
                            {{ b.bulan.slice(0, 3) }}: {{ b.status === 'lunas' ? '✓' : formatRupiah(b.nominal) }}
                          </span>
                        </div>
                        <span v-else class="text-slate-400 italic text-[11px]">- Belum ada tagihan -</span>
                      </td>
                      <td class="py-3 px-4 text-center">
                        <button
                          v-if="s.total_tagihan > 0"
                          @click="handleHapusTarifSiswa(s)"
                          class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-rose-50 hover:bg-rose-100 text-rose-700 font-bold text-[11px] border border-rose-200 transition-colors"
                          title="Hapus seluruh tagihan siswa ini pada pos ini"
                        >
                          <ion-icon :icon="trashOutline" class="text-xs"></ion-icon>
                          Hapus
                        </button>
                        <span v-else class="text-slate-300 text-xs">-</span>
                      </td>
                    </tr>
                    <tr v-if="tarifStudents.length === 0">
                      <td colspan="7" class="p-8 text-center text-slate-400 text-xs">
                        Tidak ada siswa di kelas {{ selectedTarifKelas }} atau tagihan belum diset. Klik "+ Berdasarkan Kelas" untuk mengisi tarif.
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <!-- ==================================================== -->
          <!-- VIEW 3: KASIR & REKAP PEMBAYARAN TAGIHAN SISWA       -->
          <!-- ==================================================== -->
          <div v-else-if="adminMode === 'kasir'" class="space-y-4">
            
            <!-- Filter Bar Kasir -->
            <div class="bg-white p-4 rounded-2xl border border-slate-200/80 shadow-sm flex flex-col sm:flex-row gap-3 items-center justify-between">
              <div class="w-full sm:w-72">
                <input
                  v-model="searchKasir"
                  type="text"
                  placeholder="Cari nama siswa / NIS..."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div class="flex items-center gap-2 w-full sm:w-auto">
                <select
                  v-model="filterKasirStatus"
                  class="px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-bold text-slate-700 focus:outline-none"
                >
                  <option value="">Semua Status</option>
                  <option value="belum_lunas">Belum Lunas</option>
                  <option value="lunas">Lunas</option>
                </select>

                <select
                  v-model="filterKasirKelas"
                  class="px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white font-bold text-slate-700 focus:outline-none"
                >
                  <option value="">Semua Kelas</option>
                  <option v-for="k in kelasList" :key="k.id" :value="k.nama_kelas">
                    Kelas {{ k.nama_kelas }}
                  </option>
                </select>
              </div>
            </div>

            <!-- Tabel Kasir Tagihan Siswa -->
            <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
              <div class="overflow-x-auto">
                <table class="w-full text-left text-xs text-slate-600 min-w-[700px]">
                  <thead class="bg-slate-50 text-[10px] uppercase font-bold text-slate-400 border-b border-slate-100">
                    <tr>
                      <th class="py-3 px-4">Nama Siswa</th>
                      <th class="py-3 px-4">Jenis Tagihan</th>
                      <th class="py-3 px-4">Bulan / Periode</th>
                      <th class="py-3 px-4">Nominal</th>
                      <th class="py-3 px-4 text-center">Status</th>
                      <th class="py-3 px-4 text-right">Aksi Kasir</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr v-for="t in filteredAdminTagihan" :key="t.id" class="hover:bg-slate-50">
                      <td class="py-3 px-4">
                        <span class="font-bold text-slate-800 block">{{ t.nama_siswa }}</span>
                        <span class="text-[10px] text-slate-400 font-mono">NIS: {{ t.nis }} • {{ t.kelas }}</span>
                      </td>
                      <td class="py-3 px-4 font-semibold text-slate-700">{{ t.nama_tagihan }}</td>
                      <td class="py-3 px-4 font-medium text-slate-600">{{ t.periode }}</td>
                      <td class="py-3 px-4 font-mono font-bold text-slate-800">Rp {{ formatRupiah(t.nominal) }}</td>
                      <td class="py-3 px-4 text-center">
                        <span
                          class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                          :class="t.status === 'lunas' ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'"
                        >
                          {{ t.status === 'lunas' ? 'Lunas' : 'Belum Lunas' }}
                        </span>
                      </td>
                      <td class="py-3 px-4 text-right">
                        <div class="inline-flex items-center gap-2 justify-end">
                          <button
                            v-if="t.status === 'belum_lunas'"
                            @click="handleBayarKasir(t)"
                            class="px-3 py-1.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white font-bold text-xs shadow-sm transition-all"
                          >
                            Terima Bayar
                          </button>
                          <span v-else class="text-[10px] text-slate-400 font-mono">
                            {{ t.tanggal_bayar.slice(0, 10) }}
                          </span>
                          <button
                            @click="handleDeleteTagihan(t)"
                            class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-colors"
                            title="Hapus Data Tagihan"
                          >
                            <ion-icon :icon="trashOutline" class="text-base"></ion-icon>
                          </button>
                        </div>
                      </td>
                    </tr>
                    <tr v-if="filteredAdminTagihan.length === 0">
                      <td colspan="6" class="p-8 text-center text-slate-400 text-xs">
                        Tidak ada tagihan yang cocok dengan filter.
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

          </div>

        </div>

        <!-- ==================================================== -->
        <!-- MODAL: TAMBAH / EDIT POS TAGIHAN (SCREENSHOT 2 FORM) -->
        <!-- ==================================================== -->
        <div v-if="showModalJenis" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-md w-full shadow-xl border border-slate-100">
            <h2 class="text-sm font-bold text-slate-800 mb-1">
              {{ isEditJenis ? 'Edit Jenis Pembayaran' : 'Tambah Jenis Pembayaran Baru' }}
            </h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Tentukan pos keuangan rujukan dan tipe pembayaran (Bulanan / Bebas).
            </p>

            <form @submit.prevent="handleSubmitJenis" class="space-y-3.5">
              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Nama Pembayaran</label>
                <input
                  v-model="formJenis.nama_tagihan"
                  type="text"
                  required
                  placeholder="Contoh: SPP SMP - T.A 2026/2027"
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500"
                />
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Pos Keuangan (Buku Kas)</label>
                <select
                  v-model="formJenis.pos_id"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none"
                >
                  <option v-for="p in posKeuanganList" :key="p.id" :value="p.id">
                    {{ p.nama_pos }} ({{ p.tipe }})
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Tahun Ajaran</label>
                <select
                  v-model="formJenis.tahun_ajaran_id"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none"
                >
                  <option v-for="ta in tahunList" :key="ta.id" :value="ta.id">
                    {{ ta.nama_tahun }} ({{ ta.semester }}) {{ ta.is_active ? '★ Aktif' : '' }}
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Tipe Pembayaran</label>
                <select
                  v-model="formJenis.tipe_bayar"
                  required
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 bg-white focus:outline-none"
                >
                  <option value="bulanan">Bulanan (SPP 12 Bulan: Juli s.d. Juni)</option>
                  <option value="bebas">Bebas (Non-Bulanan: Uang Pangkal / Daftar Ulang)</option>
                </select>
              </div>

              <div>
                <label class="block text-[11px] font-semibold text-slate-700 mb-1">Keterangan Singkat</label>
                <input
                  v-model="formJenis.keterangan"
                  type="text"
                  placeholder="Catatan tambahan (opsional)..."
                  class="w-full px-3 py-2 text-xs rounded-xl border border-slate-200 focus:outline-none"
                />
              </div>

              <div class="pt-3 border-t border-slate-100 flex justify-end gap-2">
                <button
                  type="button"
                  @click="showModalJenis = false"
                  class="px-3.5 py-2 rounded-xl text-slate-500 hover:bg-slate-100 text-xs font-medium"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-sm"
                >
                  {{ isEditJenis ? 'Perbarui' : 'Simpan Pembayaran' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- ==================================================== -->
        <!-- MODAL: SETTING TARIF KELAS (SESUAI SCREENSHOT 4)     -->
        <!-- ==================================================== -->
        <div v-if="showModalSetKelas" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-2xl p-6 max-w-3xl w-full shadow-xl border border-slate-100 max-h-[90vh] overflow-y-auto">
            <h2 class="text-sm font-bold text-slate-800 mb-1">
              Tambah Berdasarkan Kelas Tarif Pembayaran {{ currentPosTagihan?.tipe_bayar === 'bulanan' ? 'Bulanan' : 'Bebas' }}
            </h2>
            <p class="text-[11px] text-slate-500 mb-4">
              Isi tarif untuk seluruh siswa di kelas {{ selectedTarifKelas }}.
            </p>

            <form @submit.prevent="handleSubmitTarifKelas">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-5">
                
                <!-- Kolom Kiri: Info Readonly & Quick Fill (Screenshot 4) -->
                <div class="space-y-3">
                  <div>
                    <label class="block text-[10px] font-bold text-slate-400 uppercase">Jenis Bayar</label>
                    <input
                      type="text"
                      disabled
                      :value="currentPosTagihan?.nama_tagihan"
                      class="w-full mt-1 px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-100 font-bold text-slate-700"
                    />
                  </div>

                  <div>
                    <label class="block text-[10px] font-bold text-slate-400 uppercase">Tahun Ajaran</label>
                    <input
                      type="text"
                      disabled
                      :value="currentPosTagihan?.nama_tahun"
                      class="w-full mt-1 px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-100 font-bold text-slate-700"
                    />
                  </div>

                  <div>
                    <label class="block text-[10px] font-bold text-slate-400 uppercase">Tipe Bayar</label>
                    <input
                      type="text"
                      disabled
                      :value="currentPosTagihan?.tipe_bayar"
                      class="w-full mt-1 px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-100 font-bold text-slate-700 capitalize"
                    />
                  </div>

                  <div>
                    <label class="block text-[10px] font-bold text-slate-400 uppercase">Kelas Dituju</label>
                    <input
                      type="text"
                      disabled
                      :value="selectedTarifKelas"
                      class="w-full mt-1 px-3 py-2 text-xs rounded-xl border border-slate-200 bg-slate-100 font-bold text-slate-700"
                    />
                  </div>

                  <!-- Quick Fill Tarif Sama Rata untuk 12 Bulan (Screenshot 4) -->
                  <div v-if="currentPosTagihan?.tipe_bayar === 'bulanan'" class="p-3.5 rounded-xl bg-amber-50/70 border border-amber-200/80 mt-4">
                    <label class="block text-[11px] font-bold text-amber-900 mb-1">Tarif Setiap Bulan Sama</label>
                    <div class="flex gap-2">
                      <input
                        v-model.number="quickTarifBulanan"
                        type="number"
                        placeholder="Masukkan nilai (Rp)..."
                        class="flex-1 px-3 py-2 text-xs rounded-lg border border-amber-300 font-mono font-bold focus:outline-none"
                      />
                      <button
                        type="button"
                        @click="terapkanSemuaBulan"
                        class="px-3 py-2 rounded-lg bg-amber-400 hover:bg-amber-500 text-amber-950 font-bold text-xs shrink-0 transition-all shadow-sm"
                      >
                        Terapkan 12 Bulan
                      </button>
                    </div>
                  </div>

                  <!-- Tarif Bebas Nominal -->
                  <div v-else class="p-3.5 rounded-xl bg-purple-50/70 border border-purple-200/80 mt-4">
                    <label class="block text-[11px] font-bold text-purple-900 mb-1">Nominal Tagihan Bebas (Rp)</label>
                    <input
                      v-model.number="nominalBebas"
                      type="number"
                      required
                      placeholder="Contoh: 1500000"
                      class="w-full px-3 py-2 text-xs rounded-lg border border-purple-300 font-mono font-bold focus:outline-none"
                    />
                  </div>
                </div>

                <!-- Kolom Kanan: Rincian 12 Bulan Juli s.d. Juni (Screenshot 4) -->
                <div v-if="currentPosTagihan?.tipe_bayar === 'bulanan'" class="border border-slate-200/80 rounded-xl p-4 bg-slate-50/50">
                  <h3 class="text-xs font-bold text-slate-700 mb-3">Tarif Setiap Bulan (Juli s.d. Juni)</h3>
                  <div class="grid grid-cols-2 gap-2 max-h-72 overflow-y-auto pr-1">
                    <div v-for="bln in daftarBulan" :key="bln" class="space-y-1">
                      <label class="block text-[10px] font-bold text-slate-600">{{ bln }}</label>
                      <input
                        v-model.number="formBulanTarif[bln]"
                        type="number"
                        min="0"
                        class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 bg-white font-mono focus:outline-none focus:ring-1 focus:ring-brand-500"
                      />
                    </div>
                  </div>
                </div>

              </div>

              <div class="pt-3 border-t border-slate-100 flex justify-end gap-2">
                <button
                  type="button"
                  @click="showModalSetKelas = false"
                  class="px-4 py-2 rounded-xl text-slate-500 hover:bg-slate-100 text-xs font-medium"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  :disabled="isSubmittingTarif"
                  class="px-5 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 disabled:opacity-50"
                >
                  {{ isSubmittingTarif ? 'Menerapkan...' : 'Simpan & Terapkan Tarif' }}
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
import { ref, computed, onMounted, watch } from 'vue'
import axios from 'axios'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import {
  addOutline, settingsOutline, searchOutline, arrowBackOutline,
  createOutline, trashOutline, checkmarkCircleOutline, alertCircleOutline
} from 'ionicons/icons'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// State Admin Mode
const adminMode = ref('jenis') // 'jenis' | 'tarif' | 'kasir'
const successMsg = ref('')
const errorMsg = ref('')

// Dropdowns
const posTagihanList = ref([])
const posKeuanganList = ref([])
const tahunList = ref([])
const kelasList = ref([])

// Current Active Setting Pos Tagihan
const currentPosTagihan = ref(null)
const selectedTarifKelas = ref('')
const tarifStudents = ref([])
const isLoadingTarif = ref(false)
const isSubmittingTarif = ref(false)

// Modals
const showModalJenis = ref(false)
const isEditJenis = ref(false)
const editJenisId = ref(null)

const formJenis = ref({
  nama_tagihan: '',
  pos_id: null,
  tahun_ajaran_id: null,
  tipe_bayar: 'bulanan',
  keterangan: ''
})

const showModalSetKelas = ref(false)
const quickTarifBulanan = ref(350000)
const nominalBebas = ref(1500000)
const daftarBulan = ['Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember', 'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni']
const formBulanTarif = ref({})

// Admin Kasir State
const allAdminTagihan = ref([])
const searchKasir = ref('')
const filterKasirStatus = ref('')
const filterKasirKelas = ref('')

// Wali State
const waliTagihanList = ref([])

const formatRupiah = (val) => {
  if (!val) return '0'
  return Number(val).toLocaleString('id-ID')
}

// Summary Wali Murid
const waliSummary = computed(() => {
  let terbayar = 0
  let total = 0
  waliTagihanList.value.forEach(t => {
    total += t.nominal
    terbayar += t.terbayar
  })
  return {
    total,
    terbayar,
    tunggakan: total - terbayar
  }
})

const filteredAdminTagihan = computed(() => {
  return allAdminTagihan.value.filter(t => {
    const matchStatus = !filterKasirStatus.value || t.status === filterKasirStatus.value
    const matchKelas = !filterKasirKelas.value || t.kelas === filterKasirKelas.value
    const q = searchKasir.value.toLowerCase().trim()
    const matchSearch = !q || t.nama_siswa.toLowerCase().includes(q) || t.nis.toLowerCase().includes(q)
    return matchStatus && matchKelas && matchSearch
  })
})

const fetchDropdowns = async () => {
  try {
    const [resPos, resTA, resKelas] = await Promise.all([
      axios.get('/api/keuangan/pos'),
      axios.get('/api/admin/tahun-ajaran'),
      axios.get('/api/admin/kelas')
    ])
    posKeuanganList.value = (resPos.data || []).filter(p => p.tipe === 'masuk')
    tahunList.value = resTA.data || []
    kelasList.value = resKelas.data || []

    if (kelasList.value.length > 0) {
      selectedTarifKelas.value = kelasList.value[0].nama_kelas
    }
    const activeTA = tahunList.value.find(t => t.is_active)
    if (activeTA) {
      formJenis.value.tahun_ajaran_id = activeTA.id
    } else if (tahunList.value.length > 0) {
      formJenis.value.tahun_ajaran_id = tahunList.value[0].id
    }
    if (posKeuanganList.value.length > 0) {
      formJenis.value.pos_id = posKeuanganList.value[0].id
    }
  } catch (err) {
    console.error(err)
  }
}

const fetchPosTagihan = async () => {
  try {
    const res = await axios.get('/api/keuangan/pos-tagihan')
    posTagihanList.value = res.data || []
  } catch (err) {
    console.error(err)
  }
}

const fetchAllAdminTagihan = async () => {
  try {
    const res = await axios.get('/api/keuangan/tagihan')
    allAdminTagihan.value = res.data || []
  } catch (err) {
    console.error(err)
  }
}

const fetchWaliTagihan = async () => {
  if (!authStore.activeSiswa) return
  try {
    const res = await axios.get('/api/keuangan/tagihan', {
      params: { siswa_id: authStore.activeSiswa.id }
    })
    waliTagihanList.value = res.data || []
  } catch (err) {
    console.error(err)
  }
}

// Buka Layar Setting Tarif untuk Pos Tagihan tertentu (Screenshot 3 & 5)
const openSettingTarif = async (pos) => {
  currentPosTagihan.value = pos
  adminMode.value = 'tarif'
  await loadTarifKelas()
}

// Load Siswa & Tarif Tagihan Kelas
const loadTarifKelas = async () => {
  if (!currentPosTagihan.value || !selectedTarifKelas.value) return
  isLoadingTarif.value = true
  try {
    const res = await axios.get('/api/keuangan/tagihan/tarif', {
      params: {
        pos_tagihan_id: currentPosTagihan.value.id,
        kelas: selectedTarifKelas.value
      }
    })
    tarifStudents.value = res.data.students || []
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal memuat tarif kelas'
  } finally {
    isLoadingTarif.value = false
  }
}

// Modal Set Kelas
const openModalSetKelas = () => {
  daftarBulan.forEach(b => {
    formBulanTarif.value[b] = quickTarifBulanan.value
  })
  showModalSetKelas.value = true
}

const terapkanSemuaBulan = () => {
  daftarBulan.forEach(b => {
    formBulanTarif.value[b] = quickTarifBulanan.value
  })
}

const handleSubmitTarifKelas = async () => {
  isSubmittingTarif.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    const payload = {
      pos_tagihan_id: currentPosTagihan.value.id,
      tahun_ajaran_id: currentPosTagihan.value.tahun_ajaran_id,
      kelas: selectedTarifKelas.value,
      tipe_bayar: currentPosTagihan.value.tipe_bayar,
      nominal_bebas: nominalBebas.value,
      bulan_tarif: formBulanTarif.value
    }
    const res = await axios.post('/api/keuangan/tagihan/tarif/kelas', payload)
    successMsg.value = res.data.message || 'Tarif berhasil diterapkan!'
    showModalSetKelas.value = false
    await loadTarifKelas()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menerapkan tarif kelas'
  } finally {
    isSubmittingTarif.value = false
    setTimeout(() => { successMsg.value = ''; errorMsg.value = '' }, 4000)
  }
}

const handleHapusSemuaTarif = async () => {
  const adaTagihan = tarifStudents.value.some(s => s.total_tagihan > 0)
  if (!adaTagihan) {
    alert(`Belum ada tagihan yang diset pada kelas ${selectedTarifKelas.value} untuk pos ${currentPosTagihan.value.nama_tagihan}.`)
    return
  }

  if (!confirm(`Hapus semua tagihan untuk kelas ${selectedTarifKelas.value} pada ${currentPosTagihan.value.nama_tagihan}?\n\nCatatan: Seluruh riwayat mutasi kas terkait juga akan dibatalkan/dihapus agar saldo kas tetap sinkron.`)) return
  try {
    const res = await axios.delete('/api/keuangan/tagihan/tarif/kelas', {
      params: {
        pos_tagihan_id: currentPosTagihan.value.id,
        kelas: selectedTarifKelas.value
      }
    })
    successMsg.value = res.data.message || 'Tagihan kelas berhasil dihapus'
    await loadTarifKelas()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menghapus tagihan kelas'
  } finally {
    setTimeout(() => { successMsg.value = '' }, 4000)
  }
}

const handleHapusTarifSiswa = async (s) => {
  if (!confirm(`Hapus semua data tagihan ${currentPosTagihan.value.nama_tagihan} untuk siswa "${s.nama_lengkap}" (${s.kelas})?\n\nJika ada tagihan yang sudah dibayar, mutasi kas masuk terkait juga akan otomatis dibatalkan/dihapus.`)) return
  try {
    const res = await axios.delete('/api/keuangan/tagihan/tarif/siswa', {
      params: {
        pos_tagihan_id: currentPosTagihan.value.id,
        siswa_id: s.siswa_id
      }
    })
    successMsg.value = res.data.message || `Tagihan an. ${s.nama_lengkap} berhasil dihapus`
    await loadTarifKelas()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menghapus tagihan siswa'
  } finally {
    setTimeout(() => { successMsg.value = '' }, 4000)
  }
}

// Modal Tambah / Edit Jenis
const openAddJenisModal = () => {
  isEditJenis.value = false
  editJenisId.value = null
  formJenis.value.nama_tagihan = ''
  formJenis.value.keterangan = ''
  showModalJenis.value = true
}

const openEditJenisModal = (item) => {
  isEditJenis.value = true
  editJenisId.value = item.id
  formJenis.value = {
    nama_tagihan: item.nama_tagihan,
    pos_id: item.pos_id,
    tahun_ajaran_id: item.tahun_ajaran_id,
    tipe_bayar: item.tipe_bayar,
    keterangan: item.keterangan
  }
  showModalJenis.value = true
}

const handleSubmitJenis = async () => {
  try {
    if (isEditJenis.value) {
      await axios.put(`/api/keuangan/pos-tagihan/${editJenisId.value}`, formJenis.value)
      successMsg.value = 'Jenis tagihan pembayaran berhasil diperbarui'
    } else {
      await axios.post('/api/keuangan/pos-tagihan', formJenis.value)
      successMsg.value = 'Jenis tagihan pembayaran baru berhasil disimpan'
    }
    showModalJenis.value = false
    await fetchPosTagihan()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan jenis tagihan'
  } finally {
    setTimeout(() => { successMsg.value = '' }, 4000)
  }
}

const handleDeleteJenis = async (item) => {
  if (!confirm(`Hapus jenis tagihan "${item.nama_tagihan}"?`)) return
  try {
    await axios.delete(`/api/keuangan/pos-tagihan/${item.id}`)
    successMsg.value = 'Jenis tagihan berhasil dihapus'
    await fetchPosTagihan()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menghapus jenis tagihan'
  } finally {
    setTimeout(() => { successMsg.value = '' }, 4000)
  }
}

const handleBayarKasir = async (t) => {
  if (!confirm(`Terima pembayaran ${t.nama_tagihan} (${t.periode}) an. ${t.nama_siswa} sebesar Rp ${formatRupiah(t.nominal)}?`)) return
  try {
    await axios.post(`/api/keuangan/tagihan/${t.id}/bayar`)
    successMsg.value = `Pembayaran an. ${t.nama_siswa} berhasil diproses!`
    await fetchAllAdminTagihan()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal memproses pembayaran kasir'
  } finally {
    setTimeout(() => { successMsg.value = '' }, 4000)
  }
}

const handleDeleteTagihan = async (t) => {
  if (!confirm(`Hapus data tagihan ${t.nama_tagihan} (${t.periode}) an. ${t.nama_siswa}?\n\nCatatan: Jika tagihan ini telah lunas/dibayar, riwayat mutasi kas masuk di Buku Kas Umum juga akan otomatis dibatalkan/dihapus agar saldo kas tetap sinkron.`)) return
  try {
    await axios.delete(`/api/keuangan/tagihan/${t.id}`)
    successMsg.value = `Data tagihan an. ${t.nama_siswa} dan mutasi kas terkait berhasil dihapus!`
    await fetchAllAdminTagihan()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menghapus tagihan'
  } finally {
    setTimeout(() => { successMsg.value = '' }, 4000)
  }
}

watch(() => authStore.activeSiswa?.id, () => {
  if (authStore.isWali) {
    fetchWaliTagihan()
  }
})

onMounted(async () => {
  await fetchDropdowns()
  if (authStore.isWali) {
    await fetchWaliTagihan()
  } else {
    await fetchPosTagihan()
    await fetchAllAdminTagihan()
  }
})
</script>
