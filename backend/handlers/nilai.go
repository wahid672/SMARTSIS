package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"smartsis-backend/database"
)

type NilaiInputItem struct {
	SiswaID    int     `json:"siswa_id" binding:"required"`
	NilaiTugas float64 `json:"nilai_tugas"`
	NilaiPTS   float64 `json:"nilai_pts"`
	NilaiPAS   float64 `json:"nilai_pas"`
	Catatan    string  `json:"catatan"`
}

type BatchNilaiRequest struct {
	MapelID       int              `json:"mapel_id" binding:"required"`
	TahunAjaranID int              `json:"tahun_ajaran_id" binding:"required"`
	Items         []NilaiInputItem `json:"items" binding:"required"`
}

type SiswaNilaiRow struct {
	SiswaID     int     `json:"siswa_id"`
	NIS         string  `json:"nis"`
	NamaLengkap string  `json:"nama_lengkap"`
	Kelas       string  `json:"kelas"`
	NilaiTugas  float64 `json:"nilai_tugas"`
	NilaiPTS    float64 `json:"nilai_pts"`
	NilaiPAS    float64 `json:"nilai_pas"`
	NilaiAkhir  float64 `json:"nilai_akhir"`
	Predikat    string  `json:"predikat"`
	Catatan     string  `json:"catatan"`
}

// Helper menghitung predikat dari nilai akhir
func hitungPredikat(nilai float64) string {
	if nilai >= 88 {
		return "A"
	} else if nilai >= 78 {
		return "B"
	} else if nilai >= 65 {
		return "C"
	}
	return "D"
}

// GetNilaiKelas - Mengambil daftar siswa sekelas beserta nilainya pada mapel dan tahun ajaran tertentu
func GetNilaiKelas(c *gin.Context) {
	kelas := c.Query("kelas")
	mapelIDStr := c.Query("mapel_id")
	taIDStr := c.Query("tahun_ajaran_id")

	if kelas == "" || mapelIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter kelas dan mapel_id wajib diisi"})
		return
	}

	mapelID, _ := strconv.Atoi(mapelIDStr)
	var taID int
	if taIDStr != "" {
		taID, _ = strconv.Atoi(taIDStr)
	} else {
		// Ambil tahun ajaran aktif
		_ = database.DB.QueryRow(`SELECT id FROM tahun_ajaran WHERE is_active = 1 LIMIT 1`).Scan(&taID)
	}

	// Ambil data mapel
	var mapel Mapel
	err := database.DB.QueryRow(`SELECT id, kode_mapel, nama_mapel, kategori, kkm FROM mata_pelajaran WHERE id = ?`, mapelID).
		Scan(&mapel.ID, &mapel.KodeMapel, &mapel.NamaMapel, &mapel.Kategori, &mapel.KKM)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata pelajaran tidak ditemukan"})
		return
	}

	// Ambil data siswa di kelas ini beserta nilai yang sudah tercatat
	query := `
		SELECT 
			s.id, s.nis, s.nama_lengkap, s.kelas,
			COALESCE(n.nilai_tugas, 0),
			COALESCE(n.nilai_pts, 0),
			COALESCE(n.nilai_pas, 0),
			COALESCE(n.nilai_akhir, 0),
			COALESCE(n.predikat, '-'),
			COALESCE(n.catatan, '')
		FROM siswa s
		LEFT JOIN nilai_siswa n ON s.id = n.siswa_id AND n.mapel_id = ? AND n.tahun_ajaran_id = ?
		WHERE s.kelas = ?
		ORDER BY s.nama_lengkap ASC
	`

	rows, err := database.DB.Query(query, mapelID, taID, kelas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data nilai: " + err.Error()})
		return
	}
	defer rows.Close()

	var list []SiswaNilaiRow
	for rows.Next() {
		var r SiswaNilaiRow
		if err := rows.Scan(&r.SiswaID, &r.NIS, &r.NamaLengkap, &r.Kelas, &r.NilaiTugas, &r.NilaiPTS, &r.NilaiPAS, &r.NilaiAkhir, &r.Predikat, &r.Catatan); err == nil {
			list = append(list, r)
		}
	}

	if list == nil {
		list = []SiswaNilaiRow{}
	}

	c.JSON(http.StatusOK, gin.H{
		"mapel":           mapel,
		"tahun_ajaran_id": taID,
		"kelas":           kelas,
		"students":        list,
	})
}

// SaveNilaiBatch - Menyimpan atau memperbarui nilai banyak siswa sekaligus
func SaveNilaiBatch(c *gin.Context) {
	var req BatchNilaiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Daftar nilai tidak boleh kosong"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO nilai_siswa (siswa_id, mapel_id, tahun_ajaran_id, nilai_tugas, nilai_pts, nilai_pas, nilai_akhir, predikat, catatan, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(siswa_id, mapel_id, tahun_ajaran_id) DO UPDATE SET
			nilai_tugas = excluded.nilai_tugas,
			nilai_pts = excluded.nilai_pts,
			nilai_pas = excluded.nilai_pas,
			nilai_akhir = excluded.nilai_akhir,
			predikat = excluded.predikat,
			catatan = excluded.catatan,
			updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyiapkan query: " + err.Error()})
		return
	}
	defer stmt.Close()

	for _, item := range req.Items {
		// Hitung nilai akhir: 30% tugas + 30% PTS + 40% PAS
		akhir := (0.30 * item.NilaiTugas) + (0.30 * item.NilaiPTS) + (0.40 * item.NilaiPAS)
		akhir = math.Round(akhir*10) / 10
		predikat := hitungPredikat(akhir)

		_, err := stmt.Exec(
			item.SiswaID,
			req.MapelID,
			req.TahunAjaranID,
			item.NilaiTugas,
			item.NilaiPTS,
			item.NilaiPAS,
			akhir,
			predikat,
			item.Catatan,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan nilai siswa: " + err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menyimpan nilai untuk " + strconv.Itoa(len(req.Items)) + " siswa",
	})
}

// GetRaporSiswa - Menampilkan buku rapor digital siswa (lengkap dengan kehadiran & data lembaga)
func GetRaporSiswa(c *gin.Context) {
	siswaIDStr := c.Query("siswa_id")
	taIDStr := c.Query("tahun_ajaran_id")

	userType := c.GetString("user_type")
	userIDVal, _ := c.Get("user_id")

	var siswaID int
	if siswaIDStr != "" {
		siswaID, _ = strconv.Atoi(siswaIDStr)
	}

	// Jika akun adalah wali murid dan tidak spesifik id, ambil siswa pertama yang terafiliasi
	if userType == "wali" {
		var allowedSiswaID int
		if siswaID > 0 {
			err := database.DB.QueryRow(`SELECT siswa_id FROM wali_siswa WHERE wali_user_id = ? AND siswa_id = ?`, userIDVal, siswaID).Scan(&allowedSiswaID)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke data rapor santri ini"})
				return
			}
		} else {
			err := database.DB.QueryRow(`SELECT siswa_id FROM wali_siswa WHERE wali_user_id = ? ORDER BY id ASC LIMIT 1`, userIDVal).Scan(&siswaID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Belum ada santri yang terhubung dengan akun Anda"})
				return
			}
		}
	}

	if siswaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter siswa_id wajib diisi"})
		return
	}

	// 1. Ambil data Siswa
	var siswa struct {
		ID           int    `json:"id"`
		NIS          string `json:"nis"`
		NISN         string `json:"nisn"`
		NamaLengkap  string `json:"nama_lengkap"`
		JenisKelamin string `json:"jenis_kelamin"`
		Kelas        string `json:"kelas"`
	}
	err := database.DB.QueryRow(`SELECT id, nis, nisn, nama_lengkap, jenis_kelamin, kelas FROM siswa WHERE id = ?`, siswaID).
		Scan(&siswa.ID, &siswa.NIS, &siswa.NISN, &siswa.NamaLengkap, &siswa.JenisKelamin, &siswa.Kelas)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data santri tidak ditemukan"})
		return
	}

	// 2. Ambil data Tahun Ajaran
	var ta struct {
		ID        int    `json:"id"`
		NamaTahun string `json:"nama_tahun"`
		Semester  string `json:"semester"`
		IsActive  bool   `json:"is_active"`
	}

	if taIDStr != "" {
		taID, _ := strconv.Atoi(taIDStr)
		_ = database.DB.QueryRow(`SELECT id, nama_tahun, semester, is_active FROM tahun_ajaran WHERE id = ?`, taID).
			Scan(&ta.ID, &ta.NamaTahun, &ta.Semester, &ta.IsActive)
	} else {
		_ = database.DB.QueryRow(`SELECT id, nama_tahun, semester, is_active FROM tahun_ajaran WHERE is_active = 1 LIMIT 1`).
			Scan(&ta.ID, &ta.NamaTahun, &ta.Semester, &ta.IsActive)
	}

	// 3. Ambil Identitas Lembaga untuk Kop Rapor
	settings := make(map[string]string)
	sRows, err := database.DB.Query(`SELECT key, value FROM settings`)
	if err == nil {
		defer sRows.Close()
		for sRows.Next() {
			var k, v string
			if sRows.Scan(&k, &v) == nil {
				settings[k] = v
			}
		}
	}

	// 4. Ambil Nilai Semua Mata Pelajaran untuk Siswa ini
	type RaporMapelItem struct {
		MapelID    int     `json:"mapel_id"`
		KodeMapel  string  `json:"kode_mapel"`
		NamaMapel  string  `json:"nama_mapel"`
		Kategori   string  `json:"kategori"`
		KKM        float64 `json:"kkm"`
		NilaiTugas float64 `json:"nilai_tugas"`
		NilaiPTS   float64 `json:"nilai_pts"`
		NilaiPAS   float64 `json:"nilai_pas"`
		NilaiAkhir float64 `json:"nilai_akhir"`
		Predikat   string  `json:"predikat"`
		Catatan    string  `json:"catatan"`
		Tuntas     bool    `json:"tuntas"`
	}

	mapelRows, err := database.DB.Query(`
		SELECT 
			m.id, m.kode_mapel, m.nama_mapel, m.kategori, m.kkm,
			COALESCE(n.nilai_tugas, 0),
			COALESCE(n.nilai_pts, 0),
			COALESCE(n.nilai_pas, 0),
			COALESCE(n.nilai_akhir, 0),
			COALESCE(n.predikat, '-'),
			COALESCE(n.catatan, '')
		FROM mata_pelajaran m
		LEFT JOIN nilai_siswa n ON m.id = n.mapel_id AND n.siswa_id = ? AND n.tahun_ajaran_id = ?
		ORDER BY m.kategori ASC, m.nama_mapel ASC
	`, siswaID, ta.ID)

	var listMapel []RaporMapelItem
	var totalAkhir float64
	var countGraded int

	if err == nil {
		defer mapelRows.Close()
		for mapelRows.Next() {
			var it RaporMapelItem
			if err := mapelRows.Scan(
				&it.MapelID, &it.KodeMapel, &it.NamaMapel, &it.Kategori, &it.KKM,
				&it.NilaiTugas, &it.NilaiPTS, &it.NilaiPAS, &it.NilaiAkhir, &it.Predikat, &it.Catatan,
			); err == nil {
				it.Tuntas = it.NilaiAkhir >= it.KKM
				if it.NilaiAkhir > 0 {
					totalAkhir += it.NilaiAkhir
					countGraded++
				}
				listMapel = append(listMapel, it)
			}
		}
	}

	var rataRata float64
	if countGraded > 0 {
		rataRata = math.Round((totalAkhir/float64(countGraded))*10) / 10
	}

	// 5. Rekap Presensi Siswa
	presensiSummary := map[string]int{
		"hadir": 0,
		"sakit": 0,
		"izin":  0,
		"alpa":  0,
	}

	pRows, err := database.DB.Query(`
		SELECT status, COUNT(*) 
		FROM presensi_siswa 
		WHERE siswa_id = ? 
		GROUP BY status
	`, siswaID)
	if err == nil {
		defer pRows.Close()
		for pRows.Next() {
			var st string
			var count int
			if pRows.Scan(&st, &count) == nil {
				switch st {
				case "Hadir":
					presensiSummary["hadir"] = count
				case "Sakit":
					presensiSummary["sakit"] = count
				case "Izin":
					presensiSummary["izin"] = count
				case "Alpa":
					presensiSummary["alpa"] = count
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"siswa":        siswa,
		"tahun_ajaran": ta,
		"settings":     settings,
		"nilai":        listMapel,
		"presensi":     presensiSummary,
		"stats": gin.H{
			"rata_rata":    rataRata,
			"total_mapel":  len(listMapel),
			"count_graded": countGraded,
		},
	})
}
