package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"smartsis-backend/database"
)

type PresensiItem struct {
	SiswaID      int    `json:"siswa_id"`
	NIS          string `json:"nis"`
	NamaLengkap  string `json:"nama_lengkap"`
	JenisKelamin string `json:"jenis_kelamin"`
	Kelas        string `json:"kelas"`
	Status       string `json:"status"` // 'Hadir', 'Izin', 'Sakit', 'Alpa', atau ''
	Keterangan   string `json:"keterangan"`
	Tanggal      string `json:"tanggal"`
}

// GetPresensiByKelasDanTanggal mengambil daftar siswa di kelas tertentu beserta status presensinya pada tanggal yang dipilih
func GetPresensiByKelasDanTanggal(c *gin.Context) {
	kelas := c.Query("kelas")
	tanggal := c.Query("tanggal")

	if tanggal == "" {
		tanggal = time.Now().Format("2006-01-02")
	}

	query := `
		SELECT s.id, s.nis, s.nama_lengkap, s.jenis_kelamin, s.kelas,
		       COALESCE(p.status, 'Hadir') as status,
		       COALESCE(p.keterangan, '') as keterangan,
		       ? as tanggal
		FROM siswa s
		LEFT JOIN presensi_siswa p ON p.siswa_id = s.id AND p.tanggal = ?
		WHERE s.is_active = 1
	`
	var rows *sql.Rows
	var err error

	if kelas != "" {
		query += ` AND s.kelas = ? ORDER BY s.nama_lengkap ASC`
		rows, err = database.DB.Query(query, tanggal, tanggal, kelas)
	} else {
		query += ` ORDER BY s.kelas ASC, s.nama_lengkap ASC`
		rows, err = database.DB.Query(query, tanggal, tanggal)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data presensi"})
		return
	}
	defer rows.Close()

	var list []PresensiItem
	for rows.Next() {
		var item PresensiItem
		if err := rows.Scan(&item.SiswaID, &item.NIS, &item.NamaLengkap, &item.JenisKelamin, &item.Kelas,
			&item.Status, &item.Keterangan, &item.Tanggal); err == nil {
			list = append(list, item)
		}
	}

	c.JSON(http.StatusOK, list)
}

// SavePresensiBatch menyimpan / memperbarui absensi siswa satu kelas secara massal
func SavePresensiBatch(c *gin.Context) {
	var body struct {
		Tanggal string `json:"tanggal" binding:"required"`
		Items   []struct {
			SiswaID    int    `json:"siswa_id" binding:"required"`
			Status     string `json:"status" binding:"required"`
			Keterangan string `json:"keterangan"`
		} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data presensi tidak valid"})
		return
	}

	userID := c.MustGet("user_id")

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	stmt, err := tx.Prepare(`
		INSERT INTO presensi_siswa (siswa_id, tanggal, status, keterangan, recorded_by)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(siswa_id, tanggal) DO UPDATE SET
			status = excluded.status,
			keterangan = excluded.keterangan,
			recorded_by = excluded.recorded_by
	`)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyiapkan query"})
		return
	}
	defer stmt.Close()

	for _, item := range body.Items {
		_, err := stmt.Exec(item.SiswaID, body.Tanggal, item.Status, item.Keterangan, userID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan presensi siswa ID: " + strconv.Itoa(item.SiswaID)})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan presensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data presensi berhasil disimpan"})
}

// GetRekapPresensiWali mengambil rekap persentase dan riwayat absensi siswa untuk Wali Murid
func GetRekapPresensiWali(c *gin.Context) {
	siswaIDStr := c.Query("siswa_id")
	siswaID, _ := strconv.Atoi(siswaIDStr)

	if siswaID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Siswa ID wajib disertakan"})
		return
	}

	// Hitung rekap Hadir, Izin, Sakit, Alpa
	var hadir, izin, sakit, alpa int
	database.DB.QueryRow(`SELECT COUNT(*) FROM presensi_siswa WHERE siswa_id = ? AND status = 'Hadir'`, siswaID).Scan(&hadir)
	database.DB.QueryRow(`SELECT COUNT(*) FROM presensi_siswa WHERE siswa_id = ? AND status = 'Izin'`, siswaID).Scan(&izin)
	database.DB.QueryRow(`SELECT COUNT(*) FROM presensi_siswa WHERE siswa_id = ? AND status = 'Sakit'`, siswaID).Scan(&sakit)
	database.DB.QueryRow(`SELECT COUNT(*) FROM presensi_siswa WHERE siswa_id = ? AND status = 'Alpa'`, siswaID).Scan(&alpa)

	totalHari := hadir + izin + sakit + alpa
	persentase := 0.0
	if totalHari > 0 {
		persentase = (float64(hadir) / float64(totalHari)) * 100.0
	}

	// Ambil 30 riwayat absensi terakhir
	rows, err := database.DB.Query(`
		SELECT tanggal, status, COALESCE(keterangan, '')
		FROM presensi_siswa
		WHERE siswa_id = ?
		ORDER BY tanggal DESC
		LIMIT 30
	`, siswaID)

	type RiwayatItem struct {
		Tanggal    string `json:"tanggal"`
		Status     string `json:"status"`
		Keterangan string `json:"keterangan"`
	}

	var riwayat []RiwayatItem
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var r RiwayatItem
			if err := rows.Scan(&r.Tanggal, &r.Status, &r.Keterangan); err == nil {
				riwayat = append(riwayat, r)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_hari": totalHari,
		"hadir":      hadir,
		"izin":       izin,
		"sakit":      sakit,
		"alpa":       alpa,
		"persentase": persentase,
		"riwayat":    riwayat,
	})
}
