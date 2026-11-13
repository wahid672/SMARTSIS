package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"smartsis-backend/database"
)

type SiswaDetail struct {
	ID           int    `json:"id"`
	NIS          string `json:"nis"`
	NISN         string `json:"nisn"`
	NamaLengkap  string `json:"nama_lengkap"`
	JenisKelamin string `json:"jenis_kelamin"`
	Kelas        string `json:"kelas"`
	IsActive     bool   `json:"is_active"`
	WaliUserID   int    `json:"wali_user_id"`
	NamaWali     string `json:"nama_wali"`
	PhoneWali    string `json:"phone_wali"`
	Hubungan     string `json:"hubungan"`
}

// GetAllSiswa mengambil semua siswa beserta informasi wali muridnya
func GetAllSiswa(c *gin.Context) {
	query := `
		SELECT s.id, s.nis, COALESCE(s.nisn, ''), s.nama_lengkap, s.jenis_kelamin, s.kelas, s.is_active,
		       COALESCE(u.id, 0) as wali_user_id,
		       COALESCE(u.name, '-') as nama_wali,
		       COALESCE(u.phone_number, '-') as phone_wali,
		       COALESCE(ws.hubungan, '-') as hubungan
		FROM siswa s
		LEFT JOIN wali_siswa ws ON ws.siswa_id = s.id
		LEFT JOIN users u ON u.id = ws.wali_user_id
		ORDER BY s.kelas ASC, s.nama_lengkap ASC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}
	defer rows.Close()

	var list []SiswaDetail
	for rows.Next() {
		var s SiswaDetail
		var active int
		if err := rows.Scan(&s.ID, &s.NIS, &s.NISN, &s.NamaLengkap, &s.JenisKelamin, &s.Kelas, &active,
			&s.WaliUserID, &s.NamaWali, &s.PhoneWali, &s.Hubungan); err == nil {
			s.IsActive = active == 1
			list = append(list, s)
		}
	}

	c.JSON(http.StatusOK, list)
}

// CreateSiswa menambah siswa baru dan otomatis mengaitkan / membuat akun login wali murid
func CreateSiswa(c *gin.Context) {
	var body struct {
		NIS          string `json:"nis" binding:"required"`
		NISN         string `json:"nisn"`
		NamaLengkap  string `json:"nama_lengkap" binding:"required"`
		JenisKelamin string `json:"jenis_kelamin" binding:"required"`
		Kelas        string `json:"kelas" binding:"required"`
		// Data Wali Murid
		NamaWali  string `json:"nama_wali" binding:"required"`
		PhoneWali string `json:"phone_wali" binding:"required"`
		Hubungan  string `json:"hubungan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lengkapi data siswa dan kontak wali murid"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	// 1. Simpan data Siswa
	resSiswa, err := tx.Exec(`
		INSERT INTO siswa (nis, nisn, nama_lengkap, jenis_kelamin, kelas, is_active)
		VALUES (?, ?, ?, ?, ?, 1)
	`, strings.TrimSpace(body.NIS), strings.TrimSpace(body.NISN), strings.TrimSpace(body.NamaLengkap), body.JenisKelamin, body.Kelas)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIS sudah terdaftar pada siswa lain"})
		return
	}

	siswaID, _ := resSiswa.LastInsertId()

	// 2. Cek apakah nomor HP wali sudah terdaftar di users
	cleanPhone := strings.ReplaceAll(body.PhoneWali, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	if strings.HasPrefix(cleanPhone, "+62") {
		cleanPhone = "0" + cleanPhone[3:]
	} else if strings.HasPrefix(cleanPhone, "62") {
		cleanPhone = "0" + cleanPhone[2:]
	}

	var waliUserID int
	err = tx.QueryRow(`SELECT id FROM users WHERE phone_number = ? AND user_type = 'wali'`, cleanPhone).Scan(&waliUserID)

	if err == sql.ErrNoRows {
		// Buat akun baru wali murid dengan PIN default 123456
		hashPass, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		resWali, err := tx.Exec(`
			INSERT INTO users (name, phone_number, password_hash, role_id, user_type, is_active)
			VALUES (?, ?, ?, 3, 'wali', 1)
		`, strings.TrimSpace(body.NamaWali), cleanPhone, string(hashPass))

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun wali murid"})
			return
		}
		newWaliID, _ := resWali.LastInsertId()
		waliUserID = int(newWaliID)
	}

	// 3. Hubungkan Siswa dengan Wali di wali_siswa
	hubungan := body.Hubungan
	if hubungan == "" {
		hubungan = "Orang Tua / Wali"
	}
	_, err = tx.Exec(`
		INSERT INTO wali_siswa (wali_user_id, siswa_id, hubungan)
		VALUES (?, ?, ?)
		ON CONFLICT(wali_user_id, siswa_id) DO UPDATE SET hubungan = excluded.hubungan
	`, waliUserID, siswaID, hubungan)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengaitkan wali dengan siswa"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Data siswa dan akun wali murid berhasil disimpan",
		"siswa_id":     siswaID,
		"wali_user_id": waliUserID,
	})
}

// UpdateSiswa mengubah data siswa dan informasi wali
func UpdateSiswa(c *gin.Context) {
	idStr := c.Param("id")
	siswaID, _ := strconv.Atoi(idStr)

	var body struct {
		NIS          string `json:"nis" binding:"required"`
		NISN         string `json:"nisn"`
		NamaLengkap  string `json:"nama_lengkap" binding:"required"`
		JenisKelamin string `json:"jenis_kelamin" binding:"required"`
		Kelas        string `json:"kelas" binding:"required"`
		IsActive     bool   `json:"is_active"`
		NamaWali     string `json:"nama_wali"`
		PhoneWali    string `json:"phone_wali"`
		Hubungan     string `json:"hubungan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	activeInt := 0
	if body.IsActive {
		activeInt = 1
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	_, err = tx.Exec(`
		UPDATE siswa SET nis = ?, nisn = ?, nama_lengkap = ?, jenis_kelamin = ?, kelas = ?, is_active = ?
		WHERE id = ?
	`, strings.TrimSpace(body.NIS), strings.TrimSpace(body.NISN), strings.TrimSpace(body.NamaLengkap),
		body.JenisKelamin, body.Kelas, activeInt, siswaID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal update: NIS mungkin sudah terpakai"})
		return
	}

	// Update data wali jika nomor telepon diisi
	if body.PhoneWali != "" {
		cleanPhone := strings.ReplaceAll(body.PhoneWali, " ", "")
		cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
		if strings.HasPrefix(cleanPhone, "+62") {
			cleanPhone = "0" + cleanPhone[3:]
		} else if strings.HasPrefix(cleanPhone, "62") {
			cleanPhone = "0" + cleanPhone[2:]
		}

		var waliID int
		err := tx.QueryRow(`SELECT id FROM users WHERE phone_number = ? AND user_type = 'wali'`, cleanPhone).Scan(&waliID)
		if err == sql.ErrNoRows {
			hashPass, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
			res, _ := tx.Exec(`
				INSERT INTO users (name, phone_number, password_hash, role_id, user_type, is_active)
				VALUES (?, ?, ?, 3, 'wali', 1)
			`, body.NamaWali, cleanPhone, string(hashPass))
			newWaliID, _ := res.LastInsertId()
			waliID = int(newWaliID)
		} else if body.NamaWali != "" {
			tx.Exec(`UPDATE users SET name = ? WHERE id = ?`, body.NamaWali, waliID)
		}

		// Update / pasang relasi
		tx.Exec(`DELETE FROM wali_siswa WHERE siswa_id = ?`, siswaID)
		tx.Exec(`INSERT INTO wali_siswa (wali_user_id, siswa_id, hubungan) VALUES (?, ?, ?)`,
			waliID, siswaID, body.Hubungan)
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data siswa berhasil diperbarui"})
}

// DeleteSiswa menghapus data siswa
func DeleteSiswa(c *gin.Context) {
	idStr := c.Param("id")
	siswaID, _ := strconv.Atoi(idStr)

	database.DB.Exec(`DELETE FROM wali_siswa WHERE siswa_id = ?`, siswaID)
	database.DB.Exec(`DELETE FROM siswa WHERE id = ?`, siswaID)

	c.JSON(http.StatusOK, gin.H{"message": "Siswa berhasil dihapus"})
}
