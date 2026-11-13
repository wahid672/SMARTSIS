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

	// 2. Menus Default
	menus := []struct {
		code, title, icon, path string
		sort                    int
	}{
		{"dashboard", "Beranda", "home-outline", "/tabs/dashboard", 1},
		{"presensi", "Presensi & Absensi", "calendar-outline", "/tabs/presensi", 2},
		{"nilai", "Nilai & Rapor", "ribbon-outline", "/tabs/nilai", 3},
		{"spp", "Keuangan / SPP", "wallet-outline", "/tabs/spp", 4},
		{"siswa_mgmt", "Data Siswa", "people-outline", "/admin/siswa", 5},
		{"pegawai_mgmt", "Data Pegawai", "briefcase-outline", "/admin/pegawai", 6},
		{"role_mgmt", "Manajemen Hak Akses", "shield-checkmark-outline", "/admin/roles", 7},
		{"profil", "Profil Akun", "person-outline", "/tabs/profil", 8},
	}

	for _, m := range menus {
		DB.Exec(`INSERT OR IGNORE INTO menus (code, title, icon, path, sort_order) VALUES (?, ?, ?, ?, ?)`,
			m.code, m.title, m.icon, m.path, m.sort)
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

	log.Println("Seeding initial data selesai.")
}
