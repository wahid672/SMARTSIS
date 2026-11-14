package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func MigrateSchema() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			label TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS menus (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL UNIQUE,
			title TEXT NOT NULL,
			icon TEXT NOT NULL,
			path TEXT NOT NULL,
			parent_id INTEGER DEFAULT 0,
			category TEXT DEFAULT '',
			sort_order INTEGER DEFAULT 0,
			is_active INTEGER DEFAULT 1
		);`,

		`CREATE TABLE IF NOT EXISTS role_menus (
			role_id INTEGER NOT NULL,
			menu_id INTEGER NOT NULL,
			can_view INTEGER DEFAULT 1,
			can_create INTEGER DEFAULT 0,
			can_edit INTEGER DEFAULT 0,
			can_delete INTEGER DEFAULT 0,
			PRIMARY KEY (role_id, menu_id),
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE,
			phone_number TEXT UNIQUE,
			password_hash TEXT NOT NULL,
			role_id INTEGER NOT NULL,
			user_type TEXT NOT NULL CHECK(user_type IN ('pegawai', 'wali')),
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (role_id) REFERENCES roles(id)
		);`,

		`CREATE TABLE IF NOT EXISTS kelas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nama_kelas TEXT NOT NULL UNIQUE,
			tingkat TEXT NOT NULL,
			wali_kelas TEXT,
			keterangan TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS siswa (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nis TEXT NOT NULL UNIQUE,
			nisn TEXT UNIQUE,
			nama_lengkap TEXT NOT NULL,
			jenis_kelamin TEXT CHECK(jenis_kelamin IN ('L', 'P')),
			kelas TEXT NOT NULL,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS wali_siswa (
			wali_user_id INTEGER NOT NULL,
			siswa_id INTEGER NOT NULL,
			hubungan TEXT NOT NULL DEFAULT 'Orang Tua / Wali',
			PRIMARY KEY (wali_user_id, siswa_id),
			FOREIGN KEY (wali_user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (siswa_id) REFERENCES siswa(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			description TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS tahun_ajaran (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nama_tahun TEXT NOT NULL,
			semester TEXT NOT NULL CHECK(semester IN ('Ganjil', 'Genap')),
			is_active INTEGER DEFAULT 0
		);`,

		`CREATE TABLE IF NOT EXISTS pos_keuangan (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nama_pos TEXT NOT NULL,
			tipe TEXT NOT NULL CHECK(tipe IN ('masuk', 'keluar')),
			keterangan TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS transaksi_kas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pos_id INTEGER NOT NULL,
			kode_transaksi TEXT NOT NULL,
			tanggal DATE NOT NULL,
			tipe TEXT NOT NULL CHECK(tipe IN ('masuk', 'keluar')),
			nominal REAL NOT NULL,
			keterangan TEXT,
			created_by INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (pos_id) REFERENCES pos_keuangan(id)
		);`,

		`CREATE TABLE IF NOT EXISTS pos_tagihan (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nama_tagihan TEXT NOT NULL,
			tipe_periode TEXT NOT NULL CHECK(tipe_periode IN ('bulanan', 'semester', 'tahunan')),
			nominal REAL NOT NULL,
			keterangan TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS tagihan_siswa (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			siswa_id INTEGER NOT NULL,
			pos_tagihan_id INTEGER NOT NULL,
			tahun_ajaran_id INTEGER NOT NULL,
			periode TEXT NOT NULL,
			nominal REAL NOT NULL,
			terbayar REAL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'belum_lunas' CHECK(status IN ('belum_lunas', 'lunas')),
			tanggal_bayar DATETIME,
			FOREIGN KEY (siswa_id) REFERENCES siswa(id) ON DELETE CASCADE,
			FOREIGN KEY (pos_tagihan_id) REFERENCES pos_tagihan(id),
			FOREIGN KEY (tahun_ajaran_id) REFERENCES tahun_ajaran(id)
		);`,

		`CREATE TABLE IF NOT EXISTS presensi_siswa (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			siswa_id INTEGER NOT NULL,
			tanggal DATE NOT NULL,
			status TEXT NOT NULL CHECK(status IN ('Hadir', 'Izin', 'Sakit', 'Alpa')),
			keterangan TEXT,
			recorded_by INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(siswa_id, tanggal),
			FOREIGN KEY (siswa_id) REFERENCES siswa(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS mata_pelajaran (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			kode_mapel TEXT NOT NULL UNIQUE,
			nama_mapel TEXT NOT NULL,
			kategori TEXT NOT NULL CHECK(kategori IN ('Diniyyah / Pesantren', 'Umum / Nasional')),
			kkm REAL DEFAULT 75
		);`,

		`CREATE TABLE IF NOT EXISTS nilai_siswa (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			siswa_id INTEGER NOT NULL,
			mapel_id INTEGER NOT NULL,
			tahun_ajaran_id INTEGER NOT NULL,
			nilai_tugas REAL DEFAULT 0,
			nilai_pts REAL DEFAULT 0,
			nilai_pas REAL DEFAULT 0,
			nilai_akhir REAL DEFAULT 0,
			predikat TEXT DEFAULT 'C',
			catatan TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(siswa_id, mapel_id, tahun_ajaran_id),
			FOREIGN KEY (siswa_id) REFERENCES siswa(id) ON DELETE CASCADE,
			FOREIGN KEY (mapel_id) REFERENCES mata_pelajaran(id) ON DELETE CASCADE,
			FOREIGN KEY (tahun_ajaran_id) REFERENCES tahun_ajaran(id)
		);`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatalf("Error migrasi schema: %v\nQuery: %s", err, q)
		}
	}

	// Pembersihan dan migrasi kolom baru
	DB.Exec(`DELETE FROM pos_keuangan WHERE LOWER(TRIM(nama_pos)) = 'spp santri'`)
	DB.Exec(`DELETE FROM pos_keuangan WHERE id NOT IN (
		SELECT MIN(id) FROM pos_keuangan GROUP BY LOWER(TRIM(nama_pos)), tipe
	)`)
	DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_pos_keuangan_nama_tipe ON pos_keuangan(nama_pos, tipe)`)

	DB.Exec(`ALTER TABLE pos_tagihan ADD COLUMN pos_id INTEGER DEFAULT 1`)
	DB.Exec(`ALTER TABLE pos_tagihan ADD COLUMN tahun_ajaran_id INTEGER DEFAULT 1`)
	DB.Exec(`ALTER TABLE pos_tagihan ADD COLUMN tipe_bayar TEXT DEFAULT 'bulanan'`)
	DB.Exec(`ALTER TABLE tagihan_siswa ADD COLUMN bulan TEXT DEFAULT ''`)

	log.Println("Migrasi tabel SQLite berhasil dijalankan.")
}

func SeedInitialData() {
	// 1. Roles Default
	roles := []struct {
		name, label, desc string
	}{
		{"admin", "Super Admin / Tata Usaha", "Akses penuh sistem administrasi"},
		{"guru", "Dewan Guru", "Akses presensi, nilai, dan jadwal mengajar"},
		{"wali", "Wali Murid", "Akses data perkembangan siswa, presensi, dan SPP"},
	}

	for _, r := range roles {
		DB.Exec(`INSERT OR IGNORE INTO roles (name, label, description) VALUES (?, ?, ?)`, r.name, r.label, r.desc)
	}

	// 2. Menus Default Berkelompok (Menu Utama, Akademik, Keuangan, Master Data, Pengaturan)
	DB.Exec(`DELETE FROM menus WHERE code = 'spp'`)

	menus := []struct {
		code, title, icon, path, category string
		sort                              int
	}{
		{"dashboard", "Beranda", "home-outline", "/tabs/dashboard", "UTAMA", 1},
		{"presensi", "Presensi & Absensi", "calendar-outline", "/tabs/presensi", "AKADEMIK", 2},
		{"nilai", "Nilai & Rapor", "ribbon-outline", "/tabs/nilai", "AKADEMIK", 3},
		{"kas_pos", "Buku Kas & Pos", "wallet-outline", "/keuangan/kas", "KEUANGAN", 4},
		{"tagihan_siswa", "Tagihan Siswa", "cash-outline", "/keuangan/tagihan", "KEUANGAN", 5},
		{"siswa_mgmt", "Data Siswa", "people-outline", "/admin/siswa", "DATA MASTER", 6},
		{"kelas_mgmt", "Data Kelas", "school-outline", "/admin/kelas", "DATA MASTER", 7},
		{"mapel_mgmt", "Mata Pelajaran", "book-outline", "/admin/mapel", "DATA MASTER", 8},
		{"tahun_ajaran_mgmt", "Tahun Ajaran", "calendar-number-outline", "/admin/tahun-ajaran", "DATA MASTER", 9},
		{"pegawai_mgmt", "Data Pegawai", "briefcase-outline", "/admin/pegawai", "DATA MASTER", 10},
		{"role_mgmt", "Manajemen Hak Akses", "shield-checkmark-outline", "/admin/roles", "PENGATURAN", 11},
		{"pengaturan", "Identitas Lembaga", "settings-outline", "/admin/pengaturan", "PENGATURAN", 12},
		{"profil", "Profil Akun", "person-outline", "/tabs/profil", "PENGATURAN", 13},
	}

	// Pastikan kolom category ada jika tabel sudah ada sebelumnya
	DB.Exec(`ALTER TABLE menus ADD COLUMN category TEXT DEFAULT ''`)

	for _, m := range menus {
		DB.Exec(`INSERT INTO menus (code, title, icon, path, category, sort_order) 
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(code) DO UPDATE SET title = excluded.title, icon = excluded.icon, path = excluded.path, category = excluded.category, sort_order = excluded.sort_order`,
			m.code, m.title, m.icon, m.path, m.category, m.sort)
	}

	// 3. Hak Akses Menu Role Admin (Bisa semua menu)
	DB.Exec(`INSERT OR IGNORE INTO role_menus (role_id, menu_id, can_view, can_create, can_edit, can_delete)
		SELECT 1, id, 1, 1, 1, 1 FROM menus;`)

	// 4. Hak Akses Menu Role Guru
	DB.Exec(`INSERT OR IGNORE INTO role_menus (role_id, menu_id, can_view, can_create, can_edit, can_delete)
		SELECT 2, id, 1, 1, 1, 0 FROM menus WHERE code IN ('dashboard', 'presensi', 'nilai', 'profil');`)

	// 5. Hak Akses Menu Role Wali Murid
	DB.Exec(`INSERT OR IGNORE INTO role_menus (role_id, menu_id, can_view, can_create, can_edit, can_delete)
		SELECT 3, id, 1, 0, 0, 0 FROM menus WHERE code IN ('dashboard', 'presensi', 'nilai', 'tagihan_siswa', 'profil');`)

	// 6. Akun Default Admin & Wali Contoh
	hashPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	DB.Exec(`INSERT OR IGNORE INTO users (id, name, email, phone_number, password_hash, role_id, user_type)
		VALUES (1, 'Administrator', 'admin@sekolah.sch.id', '081200000001', ?, 1, 'pegawai')`, string(hashPass))

	hashWali, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	DB.Exec(`INSERT OR IGNORE INTO users (id, name, email, phone_number, password_hash, role_id, user_type)
		VALUES (2, 'Bapak H. Ahmad (Wali Murid)', NULL, '081234567890', ?, 3, 'wali')`, string(hashWali))

	// 7. Data Siswa Contoh (2 Siswa untuk Wali Pak Ahmad)
	DB.Exec(`INSERT OR IGNORE INTO siswa (id, nis, nisn, nama_lengkap, jenis_kelamin, kelas)
		VALUES (1, '2026001', '0089123456', 'Muhammad Farhan', 'L', 'VII-A')`)
	DB.Exec(`INSERT OR IGNORE INTO siswa (id, nis, nisn, nama_lengkap, jenis_kelamin, kelas)
		VALUES (2, '2026002', '0089123457', 'Aisyah Zahra', 'P', 'IX-B')`)

	// Relasikan 2 Siswa ke Wali Pak Ahmad (ID: 2)
	DB.Exec(`INSERT OR IGNORE INTO wali_siswa (wali_user_id, siswa_id, hubungan) VALUES (2, 1, 'Ayah Kandung')`)
	DB.Exec(`INSERT OR IGNORE INTO wali_siswa (wali_user_id, siswa_id, hubungan) VALUES (2, 2, 'Ayah Kandung')`)

	// 7b. Data Kelas Contoh
	kelasDefault := []struct {
		nama, tingkat, ket string
	}{
		{"VII-A", "Kelas 7", "Putra - Reguler"},
		{"VII-B", "Kelas 7", "Putri - Reguler"},
		{"VIII-A", "Kelas 8", "Putra - Reguler"},
		{"IX-B", "Kelas 9", "Putri - Tahfidz"},
	}
	for _, k := range kelasDefault {
		DB.Exec(`INSERT OR IGNORE INTO kelas (nama_kelas, tingkat, keterangan) VALUES (?, ?, ?)`, k.nama, k.tingkat, k.ket)
	}

	// 7c. Tahun Ajaran Awal
	DB.Exec(`INSERT OR IGNORE INTO tahun_ajaran (id, nama_tahun, semester, is_active) VALUES (1, '2026/2027', 'Ganjil', 1)`)

	// 7d. Pos Keuangan & Pos Tagihan Awal
	posList := []struct {
		nama, tipe, ket string
	}{
		{"SPP Siswa", "masuk", "Penerimaan rutin iuran bulanan siswa"},
		{"Infaq Pembangunan", "masuk", "Penerimaan dana pembangunan/gedung"},
		{"Dana Bantuan Operasional", "masuk", "Dana bantuan operasional sekolah/pesantren"},
		{"Gaji Staf & Dewan Guru", "keluar", "Pengeluaran rutin honor dewan guru dan staf"},
		{"Operasional Siswa & Kegiatan", "keluar", "Pengeluaran kegiatan siswa & sarana harian"},
		{"Listrik, Air & Internet", "keluar", "Biaya utilitas sarana dan prasarana"},
	}
	for _, p := range posList {
		var exists int
		DB.QueryRow(`SELECT COUNT(*) FROM pos_keuangan WHERE LOWER(TRIM(nama_pos)) = LOWER(TRIM(?)) AND tipe = ?`, p.nama, p.tipe).Scan(&exists)
		if exists == 0 {
			DB.Exec(`INSERT INTO pos_keuangan (nama_pos, tipe, keterangan) VALUES (?, ?, ?)`, p.nama, p.tipe, p.ket)
		}
	}

	DB.Exec(`INSERT OR IGNORE INTO pos_tagihan (id, pos_id, tahun_ajaran_id, nama_tagihan, tipe_periode, tipe_bayar, nominal, keterangan)
		VALUES (1, 1, 1, 'SPP SMP - T.A 2026/2027', 'bulanan', 'bulanan', 350000, 'SPP reguler bulanan')`)
	DB.Exec(`INSERT OR IGNORE INTO pos_tagihan (id, pos_id, tahun_ajaran_id, nama_tagihan, tipe_periode, tipe_bayar, nominal, keterangan)
		VALUES (2, 2, 1, 'Daftar Ulang Siswa Baru', 'tahunan', 'bebas', 1500000, 'Iuran tahunan fasilitas dan seragam siswa')`)

	// Generate contoh tagihan siswa untuk 2 siswa demo
	DB.Exec(`INSERT OR IGNORE INTO tagihan_siswa (id, siswa_id, pos_tagihan_id, tahun_ajaran_id, periode, nominal, terbayar, status, tanggal_bayar)
		VALUES (1, 1, 1, 1, 'November 2026', 350000, 350000, 'lunas', '2026-11-05 08:30:00')`)
	DB.Exec(`INSERT OR IGNORE INTO tagihan_siswa (id, siswa_id, pos_tagihan_id, tahun_ajaran_id, periode, nominal, terbayar, status, tanggal_bayar)
		VALUES (2, 2, 1, 1, 'November 2026', 350000, 0, 'belum_lunas', NULL)`)

	// 7e. Data Presensi Contoh
	today := "2026-11-14"
	DB.Exec(`INSERT OR IGNORE INTO presensi_siswa (siswa_id, tanggal, status, keterangan, recorded_by)
		VALUES (1, ?, 'Hadir', 'Tepat waktu', 1)`, today)
	DB.Exec(`INSERT OR IGNORE INTO presensi_siswa (siswa_id, tanggal, status, keterangan, recorded_by)
		VALUES (2, ?, 'Hadir', 'Tepat waktu', 1)`, today)

	// 7f. Data Mata Pelajaran Contoh
	mapelDefault := []struct {
		kode, nama, kategori string
		kkm                  float64
	}{
		{"TFQ", "Tahfidz Al-Qur'an & Tajwid", "Diniyyah / Pesantren", 75},
		{"FQH", "Fiqih & Ibadah Praktis", "Diniyyah / Pesantren", 75},
		{"AQD", "Aqidah Akhlak", "Diniyyah / Pesantren", 75},
		{"ARB", "Bahasa Arab", "Diniyyah / Pesantren", 70},
		{"IND", "Bahasa Indonesia", "Umum / Nasional", 75},
		{"ENG", "Bahasa Inggris", "Umum / Nasional", 70},
		{"MTK", "Matematika", "Umum / Nasional", 70},
		{"IPA", "Ilmu Pengetahuan Alam (IPA)", "Umum / Nasional", 72},
	}
	for _, m := range mapelDefault {
		DB.Exec(`INSERT OR IGNORE INTO mata_pelajaran (kode_mapel, nama_mapel, kategori, kkm)
			VALUES (?, ?, ?, ?)`, m.kode, m.nama, m.kategori, m.kkm)
	}

	// 7g. Nilai Siswa Contoh untuk Demo Rapor
	contohNilai := []struct {
		siswaID, mapelID, taID           int
		tugas, pts, pas, akhir           float64
		predikat, catatan                string
	}{
		{1, 1, 1, 90, 88, 92, 90.2, "A", "Makharijul huruf dan hafalan sangat baik serta lancar"},
		{1, 2, 1, 85, 82, 85, 84.1, "B", "Memahami rukun dan syarat ibadah sholat dan thaharah"},
		{1, 3, 1, 88, 86, 90, 88.2, "A", "Berakhlak terpuji dan santun kepada sesama"},
		{1, 4, 1, 80, 78, 82, 80.2, "B", "Kosakata mufradat semakin bertambah"},
		{1, 5, 1, 85, 80, 84, 83.1, "B", "Keterampilan menulis narasi berkembang baik"},
		{1, 7, 1, 78, 75, 80, 77.9, "B", "Ditingkatkan lagi latihan soal operasi aljabar"},

		{2, 1, 1, 95, 92, 96, 94.5, "A", "Hafalan juz 30 mutqin dan tartil sangat indah"},
		{2, 2, 1, 88, 85, 88, 87.1, "B", "Aktif dalam kajian fiqih wanita dan ibadah harian"},
		{2, 6, 1, 85, 88, 90, 87.9, "B", "Grammar and conversation skills are great"},
		{2, 8, 1, 86, 84, 88, 86.2, "B", "Pemahaman konsep biologi dan ekosistem sangat baik"},
	}
	for _, n := range contohNilai {
		DB.Exec(`INSERT OR IGNORE INTO nilai_siswa (siswa_id, mapel_id, tahun_ajaran_id, nilai_tugas, nilai_pts, nilai_pas, nilai_akhir, predikat, catatan)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			n.siswaID, n.mapelID, n.taID, n.tugas, n.pts, n.pas, n.akhir, n.predikat, n.catatan)
	}

	// 8. Pengaturan Identitas Lembaga Default
	settingsDefault := []struct {
		key, val, desc string
	}{
		{"app_name", "SMARTSIS", "Nama Singkat Aplikasi"},
		{"school_name", "Pondok Pesantren & Madrasah Terpadu", "Nama Resmi Lembaga / Sekolah"},
		{"tagline", "Sistem Informasi Sekolah & Pesantren Terpadu", "Motto atau Deskripsi Singkat"},
		{"phone", "08123456789", "Nomor Kontak / WhatsApp Resmi Lembaga"},
		{"email", "info@sekolah.sch.id", "Email Resmi Lembaga"},
		{"address", "Jl. Pesantren No. 01, Kompleks Pendidikan", "Alamat Lembaga"},
		{"logo_url", "", "URL / Base64 Logo Lembaga"},
	}

	for _, s := range settingsDefault {
		DB.Exec(`INSERT OR IGNORE INTO settings (key, value, description) VALUES (?, ?, ?)`, s.key, s.val, s.desc)
	}

	log.Println("Seeding initial data selesai.")
}
