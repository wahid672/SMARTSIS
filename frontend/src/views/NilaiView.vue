<template>
  <ion-page>
    <ion-content :fullscreen="true" class="ion-padding bg-[#f4f7fc] pt-14 md:pt-4">
      <div class="max-w-6xl mx-auto py-2">

        <!-- ==================================================== -->
        <!-- 1. TAMPILAN KHUSUS WALI SANTRI (RAPOR DIGITAL ANAK) -->
        <!-- ==================================================== -->
        <div v-if="authStore.isWali">
          <!-- Header Rapor -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
            <div>
              <h1 class="text-xl font-bold text-slate-800">Buku Rapor & Capaian Akademik</h1>
              <p class="text-xs text-slate-500 mt-1">
                Laporan penilaian hasil belajar santri an. 
                <span class="font-bold text-brand-600">{{ raporData.siswa?.nama_lengkap || authStore.activeSiswa?.nama_lengkap }}</span> 
                ({{ raporData.siswa?.kelas || authStore.activeSiswa?.kelas }})
              </p>
            </div>
            <button
              @click="printRapor"
              class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-900 text-white text-xs font-bold shadow-md transition-all self-start sm:self-auto"
            >
              <ion-icon :icon="printOutline" class="text-base"></ion-icon>
              Cetak Buku Rapor
            </button>
          </div>

          <!-- Ringkasan Statistik Siswa -->
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
            <div class="bg-gradient-to-br from-emerald-50 to-teal-100/50 p-4 rounded-2xl border border-emerald-200/60 shadow-sm text-center">
              <span class="text-[10px] font-bold text-emerald-800 uppercase tracking-wider block">Rata-rata Nilai</span>
              <p class="text-3xl font-black text-emerald-700 mt-1">{{ raporData.stats?.rata_rata || 0 }}</p>
              <span class="text-[10px] text-emerald-600 font-semibold mt-1 inline-block">Skala 100</span>
            </div>

            <div class="bg-gradient-to-br from-blue-50 to-sky-100/50 p-4 rounded-2xl border border-blue-200/60 shadow-sm text-center">
              <span class="text-[10px] font-bold text-blue-800 uppercase tracking-wider block">Total Mata Pelajaran</span>
              <p class="text-3xl font-black text-blue-700 mt-1">{{ raporData.stats?.total_mapel || 0 }}</p>
              <span class="text-[10px] text-slate-400 mt-1 block">Kurikulum Terpadu</span>
            </div>

            <div class="bg-gradient-to-br from-amber-50 to-orange-100/50 p-4 rounded-2xl border border-amber-200/60 shadow-sm text-center">
              <span class="text-[10px] font-bold text-amber-800 uppercase tracking-wider block">Tuntas Belajar</span>
              <p class="text-3xl font-black text-amber-700 mt-1">{{ tuntasCount }} Mapel</p>
              <span class="text-[10px] text-slate-400 mt-1 block">Memenuhi Standar KKM</span>
            </div>

            <div class="bg-gradient-to-br from-indigo-50 to-purple-100/50 p-4 rounded-2xl border border-indigo-200/60 shadow-sm text-center">
              <span class="text-[10px] font-bold text-indigo-800 uppercase tracking-wider block">Presensi Santri</span>
              <p class="text-3xl font-black text-indigo-700 mt-1">{{ raporData.presensi?.hadir || 0 }} Hari</p>
              <span class="text-[10px] text-indigo-600 font-semibold mt-1 block">
                S: {{ raporData.presensi?.sakit || 0 }} | I: {{ raporData.presensi?.izin || 0 }} | A: {{ raporData.presensi?.alpa || 0 }}
              </span>
            </div>
          </div>

          <!-- Rapor Container (Ready for Print) -->
          <div id="printable-rapor" class="space-y-6">

            <!-- Tabel Kurikulum Diniyyah / Pesantren -->
            <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
              <div class="px-5 py-3.5 bg-emerald-50/70 border-b border-emerald-100 flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
                  <h3 class="text-xs font-bold uppercase tracking-wider text-emerald-900">
                    A. Kelompok Mata Pelajaran Diniyyah / Kepesantrenan
                  </h3>
                </div>
                <span class="text-[11px] font-semibold text-emerald-700">Muatan Kepesantrenan & Al-Qur'an</span>
              </div>

              <div class="overflow-x-auto">
                <table class="w-full text-left text-xs text-slate-600 min-w-[650px]">
                  <thead class="bg-slate-50 text-[10px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                    <tr>
                      <th class="py-3 px-4 w-12 text-center">No</th>
                      <th class="py-3 px-4">Mata Pelajaran</th>
                      <th class="py-3 px-3 text-center w-16">KKM</th>
                      <th class="py-3 px-3 text-center w-16">Tugas</th>
                      <th class="py-3 px-3 text-center w-16">PTS</th>
                      <th class="py-3 px-3 text-center w-16">PAS</th>
                      <th class="py-3 px-3 text-center w-20">Nilai Akhir</th>
                      <th class="py-3 px-3 text-center w-16">Predikat</th>
                      <th class="py-3 px-4">Deskripsi / Catatan Pembinaan</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr v-for="(n, idx) in listDiniyyah" :key="n.mapel_id" class="hover:bg-slate-50/50">
                      <td class="py-3 px-4 text-center text-slate-400 font-mono">{{ idx + 1 }}</td>
                      <td class="py-3 px-4 font-bold text-slate-800">
                        {{ n.nama_mapel }}
                        <span class="text-[10px] text-slate-400 block font-normal">{{ n.kode_mapel }}</span>
                      </td>
                      <td class="py-3 px-3 text-center font-mono text-slate-600">{{ n.kkm }}</td>
                      <td class="py-3 px-3 text-center font-mono">{{ n.nilai_tugas }}</td>
                      <td class="py-3 px-3 text-center font-mono">{{ n.nilai_pts }}</td>
                      <td class="py-3 px-3 text-center font-mono">{{ n.nilai_pas }}</td>
                      <td class="py-3 px-3 text-center font-mono font-bold text-sm" :class="n.nilai_akhir >= n.kkm ? 'text-emerald-700' : 'text-rose-600'">
                        {{ n.nilai_akhir }}
                      </td>
                      <td class="py-3 px-3 text-center">
                        <span class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-black" :class="getPredikatClass(n.predikat)">
                          {{ n.predikat }}
                        </span>
                      </td>
                      <td class="py-3 px-4 text-slate-500 text-[11px] leading-relaxed">
                        {{ n.catatan || 'Tercapai dengan sangat memuaskan sesuai KKM.' }}
                      </td>
                    </tr>
                    <tr v-if="listDiniyyah.length === 0">
                      <td colspan="9" class="p-6 text-center text-slate-400 text-xs">Belum ada mata pelajaran diniyyah yang dinilai.</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- Tabel Kurikulum Umum / Nasional -->
            <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
              <div class="px-5 py-3.5 bg-sky-50/70 border-b border-sky-100 flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="w-2.5 h-2.5 rounded-full bg-sky-500"></span>
                  <h3 class="text-xs font-bold uppercase tracking-wider text-sky-900">
                    B. Kelompok Mata Pelajaran Umum & Nasional
                  </h3>
                </div>
                <span class="text-[11px] font-semibold text-sky-700">Kurikulum Standar Nasional</span>
              </div>

              <div class="overflow-x-auto">
                <table class="w-full text-left text-xs text-slate-600 min-w-[650px]">
                  <thead class="bg-slate-50 text-[10px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                    <tr>
                      <th class="py-3 px-4 w-12 text-center">No</th>
                      <th class="py-3 px-4">Mata Pelajaran</th>
                      <th class="py-3 px-3 text-center w-16">KKM</th>
                      <th class="py-3 px-3 text-center w-16">Tugas</th>
                      <th class="py-3 px-3 text-center w-16">PTS</th>
                      <th class="py-3 px-3 text-center w-16">PAS</th>
                      <th class="py-3 px-3 text-center w-20">Nilai Akhir</th>
                      <th class="py-3 px-3 text-center w-16">Predikat</th>
                      <th class="py-3 px-4">Deskripsi / Catatan Pembinaan</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr v-for="(n, idx) in listUmum" :key="n.mapel_id" class="hover:bg-slate-50/50">
                      <td class="py-3 px-4 text-center text-slate-400 font-mono">{{ idx + 1 }}</td>
                      <td class="py-3 px-4 font-bold text-slate-800">
                        {{ n.nama_mapel }}
                        <span class="text-[10px] text-slate-400 block font-normal">{{ n.kode_mapel }}</span>
                      </td>
                      <td class="py-3 px-3 text-center font-mono text-slate-600">{{ n.kkm }}</td>
                      <td class="py-3 px-3 text-center font-mono">{{ n.nilai_tugas }}</td>
                      <td class="py-3 px-3 text-center font-mono">{{ n.nilai_pts }}</td>
                      <td class="py-3 px-3 text-center font-mono">{{ n.nilai_pas }}</td>
                      <td class="py-3 px-3 text-center font-mono font-bold text-sm" :class="n.nilai_akhir >= n.kkm ? 'text-emerald-700' : 'text-rose-600'">
                        {{ n.nilai_akhir }}
                      </td>
                      <td class="py-3 px-3 text-center">
                        <span class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-black" :class="getPredikatClass(n.predikat)">
                          {{ n.predikat }}
                        </span>
                      </td>
                      <td class="py-3 px-4 text-slate-500 text-[11px] leading-relaxed">
                        {{ n.catatan || 'Tercapai dengan sangat memuaskan sesuai KKM.' }}
                      </td>
                    </tr>
                    <tr v-if="listUmum.length === 0">
                      <td colspan="9" class="p-6 text-center text-slate-400 text-xs">Belum ada mata pelajaran umum yang dinilai.</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- Catatan Wali Kelas & Pengesahan Dokumen -->
            <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-sm flex flex-col sm:flex-row justify-between gap-6">
              <div class="max-w-md">
                <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider mb-2">Kriteria Predikat Kelulusan:</h4>
                <div class="grid grid-cols-2 gap-2 text-[11px] text-slate-600">
                  <div class="flex items-center gap-2"><span class="font-bold text-emerald-700">A (88 - 100):</span> Sangat Baik</div>
                  <div class="flex items-center gap-2"><span class="font-bold text-sky-700">B (78 - 87):</span> Baik</div>
                  <div class="flex items-center gap-2"><span class="font-bold text-amber-600">C (65 - 77):</span> Cukup</div>
                  <div class="flex items-center gap-2"><span class="font-bold text-rose-600">D (&lt; 65):</span> Perlu Bimbingan</div>
                </div>
              </div>

              <div class="text-right text-xs text-slate-500 flex flex-col justify-between">
                <div>
                  <p class="font-semibold text-slate-700">Mengetahui & Mengesahkan,</p>
                  <p class="text-[11px] text-slate-400">Kepala Sekolah / Pimpinan Pondok</p>
                </div>
                <div class="mt-12">
                  <p class="font-bold text-slate-800 underline decoration-slate-300 decoration-2">Ust. H. Ahmad Dahlan, M.Pd</p>
                  <p class="text-[10px] text-slate-400">NIP / NIY: 198504122010011002</p>
                </div>
              </div>
            </div>

          </div>
        </div>

        <!-- ==================================================== -->
        <!-- 2. TAMPILAN GURU / ADMIN (INPUT NILAI KELAS & MAPEL) -->
        <!-- ==================================================== -->
        <div v-else>
          <!-- Switcher Mode Guru: Input Nilai vs Lihat Buku Rapor -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
            <div>
              <h1 class="text-xl font-bold text-slate-800">Manajemen Nilai & Rapor</h1>
              <p class="text-xs text-slate-500 mt-1">
                Input evaluasi belajar siswa (Tugas, PTS, PAS) dan terbitkan rapor digital santri.
              </p>
            </div>

            <!-- Tab Mode -->
            <div class="flex items-center gap-1.5 p-1 bg-slate-200/80 rounded-xl self-start sm:self-auto">
              <button
                @click="adminTab = 'input'"
                class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all flex items-center gap-2"
                :class="adminTab === 'input' ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
              >
                <ion-icon :icon="createOutline" class="text-sm"></ion-icon>
                Input Nilai Kelas
              </button>
              <button
                @click="adminTab = 'rapor'"
                class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all flex items-center gap-2"
                :class="adminTab === 'rapor' ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
              >
                <ion-icon :icon="ribbonOutline" class="text-sm"></ion-icon>
                Lihat Buku Rapor
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

          <!-- TAB 1: INPUT NILAI BATCH KELAS -->
          <div v-if="adminTab === 'input'">
            <!-- Filter Bar: Pilih Tahun Ajaran, Kelas, Mapel -->
            <div class="bg-white p-4 rounded-2xl border border-slate-200/80 shadow-sm mb-6 grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">Tahun Ajaran</label>
                <select
                  v-model="filterInput.tahun_ajaran_id"
                  @change="loadNilaiKelas"
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 bg-white"
                >
                  <option v-for="t in tahunList" :key="t.id" :value="t.id">
                    {{ t.nama_tahun }} ({{ t.semester }}) {{ t.is_active ? '★ Aktif' : '' }}
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">Kelas</label>
                <select
                  v-model="filterInput.kelas"
                  @change="loadNilaiKelas"
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 bg-white"
                >
                  <option v-for="k in kelasList" :key="k.id" :value="k.nama_kelas">
                    {{ k.nama_kelas }} ({{ k.tingkat }})
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">Mata Pelajaran</label>
                <select
                  v-model="filterInput.mapel_id"
                  @change="loadNilaiKelas"
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 bg-white"
                >
                  <option v-for="m in mapelList" :key="m.id" :value="m.id">
                    {{ m.nama_mapel }} (KKM: {{ m.kkm }})
                  </option>
                </select>
              </div>
            </div>

            <!-- Info Mapel Aktif & Tombol Simpan -->
            <div v-if="selectedMapel" class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 bg-emerald-50/80 p-3.5 rounded-xl border border-emerald-200/80">
              <div class="flex items-center gap-2">
                <span class="p-2 rounded-lg bg-emerald-100 text-emerald-800 font-bold text-xs">{{ selectedMapel.kode_mapel }}</span>
                <div>
                  <h3 class="text-xs font-bold text-emerald-950">{{ selectedMapel.nama_mapel }}</h3>
                  <p class="text-[10px] text-emerald-700">
                    Kategori: {{ selectedMapel.kategori }} | Standar KKM: <span class="font-bold">{{ selectedMapel.kkm }}</span> | Rumus: (30% Tugas + 30% PTS + 40% PAS)
                  </p>
                </div>
              </div>
              <button
                @click="simpanSemuaNilai"
                :disabled="isSaving || inputStudents.length === 0"
                class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md shadow-emerald-500/20 transition-all disabled:opacity-50 flex items-center gap-2 self-start sm:self-auto"
              >
                <ion-icon :icon="checkmarkDoneOutline" class="text-base"></ion-icon>
                {{ isSaving ? 'Menyimpan...' : 'Simpan Semua Nilai' }}
              </button>
            </div>

            <!-- Tabel Input Nilai Siswa -->
            <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden mb-6">
              <div v-if="isLoadingInput" class="p-8 text-center text-xs text-slate-400">
                Memuat data siswa kelas {{ filterInput.kelas }}...
              </div>
              <div v-else-if="inputStudents.length === 0" class="p-12 text-center text-slate-400 text-xs">
                Tidak ada data santri di kelas ini atau silakan pilih kelas & mapel terlebih dahulu.
              </div>
              <div v-else class="overflow-x-auto">
                <table class="w-full text-left text-xs text-slate-600 min-w-[700px]">
                  <thead class="bg-slate-50 text-[10px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                    <tr>
                      <th class="py-3 px-4 w-12 text-center">No</th>
                      <th class="py-3 px-4 w-48">Nama Siswa</th>
                      <th class="py-3 px-3 text-center w-24">Tugas (30%)</th>
                      <th class="py-3 px-3 text-center w-24">PTS (30%)</th>
                      <th class="py-3 px-3 text-center w-24">PAS (40%)</th>
                      <th class="py-3 px-3 text-center w-24">Nilai Akhir</th>
                      <th class="py-3 px-3 text-center w-20">Predikat</th>
                      <th class="py-3 px-4">Catatan Guru</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr v-for="(s, idx) in inputStudents" :key="s.siswa_id" class="hover:bg-slate-50/50">
                      <td class="py-3 px-4 text-center text-slate-400 font-mono">{{ idx + 1 }}</td>
                      <td class="py-3 px-4">
                        <span class="font-bold text-slate-800 text-xs block">{{ s.nama_lengkap }}</span>
                        <span class="text-[10px] text-slate-400 font-mono">NIS: {{ s.nis }}</span>
                      </td>
                      <td class="py-2 px-2 text-center">
                        <input
                          v-model.number="s.nilai_tugas"
                          @input="recalculateRow(s)"
                          type="number"
                          min="0"
                          max="100"
                          step="0.5"
                          class="w-16 px-2 py-1.5 text-center font-mono font-bold text-xs rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500"
                        />
                      </td>
                      <td class="py-2 px-2 text-center">
                        <input
                          v-model.number="s.nilai_pts"
                          @input="recalculateRow(s)"
                          type="number"
                          min="0"
                          max="100"
                          step="0.5"
                          class="w-16 px-2 py-1.5 text-center font-mono font-bold text-xs rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500"
                        />
                      </td>
                      <td class="py-2 px-2 text-center">
                        <input
                          v-model.number="s.nilai_pas"
                          @input="recalculateRow(s)"
                          type="number"
                          min="0"
                          max="100"
                          step="0.5"
                          class="w-16 px-2 py-1.5 text-center font-mono font-bold text-xs rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500"
                        />
                      </td>
                      <td class="py-3 px-3 text-center font-mono font-bold text-sm" :class="s.nilai_akhir >= (selectedMapel?.kkm || 75) ? 'text-emerald-700' : 'text-rose-600'">
                        {{ s.nilai_akhir }}
                      </td>
                      <td class="py-3 px-3 text-center">
                        <span class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-black" :class="getPredikatClass(s.predikat)">
                          {{ s.predikat }}
                        </span>
                      </td>
                      <td class="py-2 px-3">
                        <input
                          v-model="s.catatan"
                          type="text"
                          placeholder="Catatan perkembangan..."
                          class="w-full px-2.5 py-1.5 text-xs rounded-lg border border-slate-200 focus:outline-none focus:ring-1 focus:ring-brand-500"
                        />
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <!-- TAB 2: LIHAT RAPOR SISWA PILIHAN GURU / ADMIN -->
          <div v-else-if="adminTab === 'rapor'">
            <!-- Selector Santri & Tahun Ajaran -->
            <div class="bg-white p-4 rounded-2xl border border-slate-200/80 shadow-sm mb-6 grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">Pilih Siswa / Santri</label>
                <select
                  v-model="selectedRaporSiswaID"
                  @change="loadAdminRapor"
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 bg-white"
                >
                  <option v-for="s in allStudentsList" :key="s.id" :value="s.id">
                    {{ s.nama_lengkap }} (Kelas: {{ s.kelas }} | NIS: {{ s.nis }})
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-[11px] font-bold text-slate-700 mb-1">Tahun Ajaran</label>
                <select
                  v-model="selectedRaporTaID"
                  @change="loadAdminRapor"
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 bg-white"
                >
                  <option v-for="t in tahunList" :key="t.id" :value="t.id">
                    {{ t.nama_tahun }} ({{ t.semester }}) {{ t.is_active ? '★ Aktif' : '' }}
                  </option>
                </select>
              </div>
            </div>

            <!-- Pratinjau Rapor Terpilih -->
            <div v-if="raporData.siswa" class="space-y-6">
              <div class="flex items-center justify-between">
                <div>
                  <h2 class="text-sm font-bold text-slate-800">
                    Pratinjau Rapor: {{ raporData.siswa.nama_lengkap }} ({{ raporData.siswa.kelas }})
                  </h2>
                  <p class="text-[11px] text-slate-500">Tahun Ajaran {{ raporData.tahun_ajaran?.nama_tahun }} - Semester {{ raporData.tahun_ajaran?.semester }}</p>
                </div>
                <button
                  @click="printRapor"
                  class="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-900 text-white text-xs font-bold shadow-md transition-all"
                >
                  <ion-icon :icon="printOutline" class="text-base"></ion-icon>
                  Cetak Rapor Santri
                </button>
              </div>

              <!-- Tabel Rapor Kurikulum Diniyyah -->
              <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
                <div class="px-5 py-3.5 bg-emerald-50/70 border-b border-emerald-100 flex items-center justify-between">
                  <h3 class="text-xs font-bold uppercase tracking-wider text-emerald-900">
                    A. Mata Pelajaran Diniyyah / Pesantren
                  </h3>
                </div>
                <div class="overflow-x-auto">
                  <table class="w-full text-left text-xs text-slate-600 min-w-[650px]">
                    <thead class="bg-slate-50 text-[10px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                      <tr>
                        <th class="py-3 px-4 w-12 text-center">No</th>
                        <th class="py-3 px-4">Mata Pelajaran</th>
                        <th class="py-3 px-3 text-center w-16">KKM</th>
                        <th class="py-3 px-3 text-center w-16">Tugas</th>
                        <th class="py-3 px-3 text-center w-16">PTS</th>
                        <th class="py-3 px-3 text-center w-16">PAS</th>
                        <th class="py-3 px-3 text-center w-20">Nilai Akhir</th>
                        <th class="py-3 px-3 text-center w-16">Predikat</th>
                        <th class="py-3 px-4">Catatan</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-100">
                      <tr v-for="(n, idx) in listDiniyyah" :key="n.mapel_id">
                        <td class="py-3 px-4 text-center font-mono">{{ idx + 1 }}</td>
                        <td class="py-3 px-4 font-bold text-slate-800">{{ n.nama_mapel }}</td>
                        <td class="py-3 px-3 text-center font-mono">{{ n.kkm }}</td>
                        <td class="py-3 px-3 text-center font-mono">{{ n.nilai_tugas }}</td>
                        <td class="py-3 px-3 text-center font-mono">{{ n.nilai_pts }}</td>
                        <td class="py-3 px-3 text-center font-mono">{{ n.nilai_pas }}</td>
                        <td class="py-3 px-3 text-center font-mono font-bold" :class="n.nilai_akhir >= n.kkm ? 'text-emerald-700' : 'text-rose-600'">
                          {{ n.nilai_akhir }}
                        </td>
                        <td class="py-3 px-3 text-center">
                          <span class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-black" :class="getPredikatClass(n.predikat)">
                            {{ n.predikat }}
                          </span>
                        </td>
                        <td class="py-3 px-4 text-slate-500 text-[11px]">{{ n.catatan || '-' }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- Tabel Rapor Kurikulum Umum -->
              <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
                <div class="px-5 py-3.5 bg-sky-50/70 border-b border-sky-100 flex items-center justify-between">
                  <h3 class="text-xs font-bold uppercase tracking-wider text-sky-900">
                    B. Mata Pelajaran Umum / Nasional
                  </h3>
                </div>
                <div class="overflow-x-auto">
                  <table class="w-full text-left text-xs text-slate-600 min-w-[650px]">
                    <thead class="bg-slate-50 text-[10px] uppercase tracking-wider text-slate-400 font-semibold border-b border-slate-100">
                      <tr>
                        <th class="py-3 px-4 w-12 text-center">No</th>
                        <th class="py-3 px-4">Mata Pelajaran</th>
                        <th class="py-3 px-3 text-center w-16">KKM</th>
                        <th class="py-3 px-3 text-center w-16">Tugas</th>
                        <th class="py-3 px-3 text-center w-16">PTS</th>
                        <th class="py-3 px-3 text-center w-16">PAS</th>
                        <th class="py-3 px-3 text-center w-20">Nilai Akhir</th>
                        <th class="py-3 px-3 text-center w-16">Predikat</th>
                        <th class="py-3 px-4">Catatan</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-100">
                      <tr v-for="(n, idx) in listUmum" :key="n.mapel_id">
                        <td class="py-3 px-4 text-center font-mono">{{ idx + 1 }}</td>
                        <td class="py-3 px-4 font-bold text-slate-800">{{ n.nama_mapel }}</td>
                        <td class="py-3 px-3 text-center font-mono">{{ n.kkm }}</td>
                        <td class="py-3 px-3 text-center font-mono">{{ n.nilai_tugas }}</td>
                        <td class="py-3 px-3 text-center font-mono">{{ n.nilai_pts }}</td>
                        <td class="py-3 px-3 text-center font-mono">{{ n.nilai_pas }}</td>
                        <td class="py-3 px-3 text-center font-mono font-bold" :class="n.nilai_akhir >= n.kkm ? 'text-emerald-700' : 'text-rose-600'">
                          {{ n.nilai_akhir }}
                        </td>
                        <td class="py-3 px-3 text-center">
                          <span class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-black" :class="getPredikatClass(n.predikat)">
                            {{ n.predikat }}
                          </span>
                        </td>
                        <td class="py-3 px-4 text-slate-500 text-[11px]">{{ n.catatan || '-' }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>

      </div>
    </ion-content>
  </ion-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { IonPage, IonContent, IonIcon } from '@ionic/vue'
import {
  ribbonOutline, createOutline, printOutline, checkmarkCircleOutline,
  alertCircleOutline, checkmarkDoneOutline
} from 'ionicons/icons'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// State Data Umum
const adminTab = ref('input')
const successMsg = ref('')
const errorMsg = ref('')

// Dropdowns
const tahunList = ref([])
const kelasList = ref([])
const mapelList = ref([])
const allStudentsList = ref([])

// State Guru: Input Batch
const filterInput = ref({
  tahun_ajaran_id: 1,
  kelas: '',
  mapel_id: 1
})
const selectedMapel = ref(null)
const inputStudents = ref([])
const isLoadingInput = ref(false)
const isSaving = ref(false)

// State Rapor (Wali & Guru Rapor View)
const raporData = ref({
  siswa: null,
  tahun_ajaran: null,
  settings: {},
  nilai: [],
  presensi: {},
  stats: {}
})
const selectedRaporSiswaID = ref(null)
const selectedRaporTaID = ref(1)

// Split nilai by Kategori
const listDiniyyah = computed(() => {
  return (raporData.value.nilai || []).filter(n => n.kategori === 'Diniyyah / Pesantren')
})
const listUmum = computed(() => {
  return (raporData.value.nilai || []).filter(n => n.kategori === 'Umum / Nasional')
})
const tuntasCount = computed(() => {
  return (raporData.value.nilai || []).filter(n => n.tuntas).length
})

const getPredikatClass = (p) => {
  switch (p) {
    case 'A': return 'bg-emerald-100 text-emerald-800'
    case 'B': return 'bg-sky-100 text-sky-800'
    case 'C': return 'bg-amber-100 text-amber-800'
    default: return 'bg-rose-100 text-rose-800'
  }
}

// Hitung nilai akhir real-time di baris tabel guru
const recalculateRow = (row) => {
  const tugas = Number(row.nilai_tugas) || 0
  const pts = Number(row.nilai_pts) || 0
  const pas = Number(row.nilai_pas) || 0
  const akhir = Math.round(((0.30 * tugas) + (0.30 * pts) + (0.40 * pas)) * 10) / 10
  row.nilai_akhir = akhir

  if (akhir >= 88) row.predikat = 'A'
  else if (akhir >= 78) row.predikat = 'B'
  else if (akhir >= 65) row.predikat = 'C'
  else row.predikat = 'D'
}

// Fetch master options
const fetchDropdowns = async () => {
  try {
    const [resTA, resKelas, resMapel] = await Promise.all([
      axios.get('/api/admin/tahun-ajaran'),
      axios.get('/api/admin/kelas'),
      axios.get('/api/admin/mapel')
    ])

    tahunList.value = resTA.data || []
    kelasList.value = resKelas.data || []
    mapelList.value = resMapel.data || []

    const activeTA = tahunList.value.find(t => t.is_active)
    if (activeTA) {
      filterInput.value.tahun_ajaran_id = activeTA.id
      selectedRaporTaID.value = activeTA.id
    } else if (tahunList.value.length > 0) {
      filterInput.value.tahun_ajaran_id = tahunList.value[0].id
      selectedRaporTaID.value = tahunList.value[0].id
    }

    if (kelasList.value.length > 0) {
      filterInput.value.kelas = kelasList.value[0].nama_kelas
    }
    if (mapelList.value.length > 0) {
      filterInput.value.mapel_id = mapelList.value[0].id
    }

    if (!authStore.isWali) {
      // Ambil seluruh siswa untuk opsi rapor guru
      const resSiswa = await axios.get('/api/admin/siswa')
      allStudentsList.value = resSiswa.data || []
      if (allStudentsList.value.length > 0) {
        selectedRaporSiswaID.value = allStudentsList.value[0].id
      }
    }
  } catch (err) {
    console.error('Error fetching dropdowns:', err)
  }
}

// Load nilai siswa di kelas
const loadNilaiKelas = async () => {
  if (!filterInput.value.kelas || !filterInput.value.mapel_id) return
  isLoadingInput.value = true
  try {
    const res = await axios.get('/api/nilai/kelas', {
      params: {
        kelas: filterInput.value.kelas,
        mapel_id: filterInput.value.mapel_id,
        tahun_ajaran_id: filterInput.value.tahun_ajaran_id
      }
    })
    selectedMapel.value = res.data.mapel
    inputStudents.value = (res.data.students || []).map(s => ({
      siswa_id: s.siswa_id,
      nis: s.nis,
      nama_lengkap: s.nama_lengkap,
      kelas: s.kelas,
      nilai_tugas: s.nilai_tugas,
      nilai_pts: s.nilai_pts,
      nilai_pas: s.nilai_pas,
      nilai_akhir: s.nilai_akhir,
      predikat: s.predikat,
      catatan: s.catatan
    }))
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal memuat data nilai kelas'
  } finally {
    isLoadingInput.value = false
  }
}

// Simpan batch nilai oleh guru
const simpanSemuaNilai = async () => {
  isSaving.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    const payload = {
      mapel_id: filterInput.value.mapel_id,
      tahun_ajaran_id: filterInput.value.tahun_ajaran_id,
      items: inputStudents.value.map(s => ({
        siswa_id: s.siswa_id,
        nilai_tugas: s.nilai_tugas || 0,
        nilai_pts: s.nilai_pts || 0,
        nilai_pas: s.nilai_pas || 0,
        catatan: s.catatan || ''
      }))
    }
    const res = await axios.post('/api/nilai/batch', payload)
    successMsg.value = res.data.message || 'Semua nilai berhasil disimpan!'
    await loadNilaiKelas()
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Gagal menyimpan nilai'
  } finally {
    isSaving.value = false
    setTimeout(() => {
      successMsg.value = ''
      errorMsg.value = ''
    }, 4000)
  }
}

// Load Rapor untuk Wali Santri
const loadWaliRapor = async () => {
  try {
    const siswaID = authStore.activeSiswa?.id
    const res = await axios.get('/api/nilai/rapor', {
      params: {
        siswa_id: siswaID
      }
    })
    raporData.value = res.data
  } catch (err) {
    console.error('Error load rapor wali:', err)
  }
}

// Load Rapor untuk Guru/Admin
const loadAdminRapor = async () => {
  if (!selectedRaporSiswaID.value) return
  try {
    const res = await axios.get('/api/nilai/rapor', {
      params: {
        siswa_id: selectedRaporSiswaID.value,
        tahun_ajaran_id: selectedRaporTaID.value
      }
    })
    raporData.value = res.data
  } catch (err) {
    console.error('Error load rapor admin:', err)
  }
}

// Cetak Rapor
const printRapor = () => {
  window.print()
}

// Watcher pergantian activeSiswa di Wali
watch(() => authStore.activeSiswa?.id, () => {
  if (authStore.isWali) {
    loadWaliRapor()
  }
})

onMounted(async () => {
  await fetchDropdowns()
  if (authStore.isWali) {
    await loadWaliRapor()
  } else {
    await loadNilaiKelas()
    await loadAdminRapor()
  }
})
</script>

<style>
@media print {
  /* Sembunyikan elemen navigasi saat cetak */
  ion-menu, ion-header, ion-tab-bar, button, .no-print {
    display: none !important;
  }
  body, ion-content {
    background: white !important;
    color: black !important;
    font-size: 10pt;
  }
  #printable-rapor {
    box-shadow: none !important;
    border: none !important;
    margin: 0 !important;
    padding: 0 !important;
  }
}
</style>
