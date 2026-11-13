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
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatalf("Error migrasi schema: %v\nQuery: %s", err, q)
		}
	}
	log.Println("Migrasi tabel SQLite berhasil dijalankan.")
}

func SeedInitialData() {
	// 1. Roles Default
	roles := []struct {
		name, label, desc string
	}{
		{"admin", "Super Admin / Tata Usaha", "Akses penuh sistem administrasi"},
		{"guru", "Dewan Guru", "Akses presensi, nilai, dan jadwal mengajar"},
		{"wali", "Wali Murid / Santri", "Akses data perkembangan siswa, presensi, dan SPP"},
	}

	for _, r := range roles {
		DB.Exec(`INSERT OR IGNORE INTO roles (name, label, description) VALUES (?, ?, ?)`, r.name, r.label, r.desc)
	}

	// 2. Menus Default Berkelompok (Menu Utama, Akademik, Master Data, Pengaturan)
	menus := []struct {
		code, title, icon, path, category string
		sort                              int
	}{
		{"dashboard", "Beranda", "home-outline", "/tabs/dashboard", "UTAMA", 1},
		{"presensi", "Presensi & Absensi", "calendar-outline", "/tabs/presensi", "AKADEMIK", 2},
		{"nilai", "Nilai & Rapor", "ribbon-outline", "/tabs/nilai", "AKADEMIK", 3},
		{"spp", "Keuangan / SPP", "wallet-outline", "/tabs/spp", "AKADEMIK", 4},
		{"siswa_mgmt", "Data Siswa", "people-outline", "/admin/siswa", "DATA MASTER", 5},
		{"kelas_mgmt", "Data Kelas", "school-outline", "/admin/kelas", "DATA MASTER", 6},
		{"pegawai_mgmt", "Data Pegawai", "briefcase-outline", "/admin/pegawai", "DATA MASTER", 7},
		{"role_mgmt", "Manajemen Hak Akses", "shield-checkmark-outline", "/admin/roles", "PENGATURAN", 8},
		{"pengaturan", "Identitas Lembaga", "settings-outline", "/admin/pengaturan", "PENGATURAN", 9},
		{"profil", "Profil Akun", "person-outline", "/tabs/profil", "PENGATURAN", 10},
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
		SELECT 3, id, 1, 0, 0, 0 FROM menus WHERE code IN ('dashboard', 'presensi', 'nilai', 'spp', 'profil');`)

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
