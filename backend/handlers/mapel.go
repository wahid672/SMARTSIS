package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"smartsis-backend/database"
)

type Mapel struct {
	ID        int     `json:"id"`
	KodeMapel string  `json:"kode_mapel"`
	NamaMapel string  `json:"nama_mapel"`
	Kategori  string  `json:"kategori"`
	KKM       float64 `json:"kkm"`
}

// GetAllMapel - Mendapatkan semua data mata pelajaran
func GetAllMapel(c *gin.Context) {
	search := c.Query("search")
	kategori := c.Query("kategori")

	query := `SELECT id, kode_mapel, nama_mapel, kategori, kkm FROM mata_pelajaran WHERE 1=1`
	var args []interface{}

	if search != "" {
		query += ` AND (nama_mapel LIKE ? OR kode_mapel LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s)
	}

	if kategori != "" {
		query += ` AND kategori = ?`
		args = append(args, kategori)
	}

	query += ` ORDER BY kategori ASC, nama_mapel ASC`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data mata pelajaran: " + err.Error()})
		return
	}
	defer rows.Close()

	var list []Mapel
	for rows.Next() {
		var m Mapel
		if err := rows.Scan(&m.ID, &m.KodeMapel, &m.NamaMapel, &m.Kategori, &m.KKM); err == nil {
			list = append(list, m)
		}
	}

	if list == nil {
		list = []Mapel{}
	}

	c.JSON(http.StatusOK, list)
}

// CreateMapel - Menambah mata pelajaran baru
func CreateMapel(c *gin.Context) {
	var req struct {
		KodeMapel string  `json:"kode_mapel" binding:"required"`
		NamaMapel string  `json:"nama_mapel" binding:"required"`
		Kategori  string  `json:"kategori" binding:"required"`
		KKM       float64 `json:"kkm"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap: kode, nama, dan kategori wajib diisi"})
		return
	}

	req.KodeMapel = strings.ToUpper(strings.TrimSpace(req.KodeMapel))
	req.NamaMapel = strings.TrimSpace(req.NamaMapel)
	if req.KKM <= 0 {
		req.KKM = 75
	}

	// Validasi kategori
	if req.Kategori != "Diniyyah / Pesantren" && req.Kategori != "Umum / Nasional" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori harus 'Diniyyah / Pesantren' atau 'Umum / Nasional'"})
		return
	}

	// Cek duplicate kode
	var exists int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM mata_pelajaran WHERE UPPER(kode_mapel) = ?`, req.KodeMapel).Scan(&exists)
	if err == nil && exists > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode mata pelajaran sudah digunakan"})
		return
	}

	res, err := database.DB.Exec(`
		INSERT INTO mata_pelajaran (kode_mapel, nama_mapel, kategori, kkm)
		VALUES (?, ?, ?, ?)
	`, req.KodeMapel, req.NamaMapel, req.Kategori, req.KKM)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan mata pelajaran: " + err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"message": "Mata pelajaran berhasil ditambahkan",
		"id":      id,
	})
}

// UpdateMapel - Memperbarui mata pelajaran
func UpdateMapel(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		KodeMapel string  `json:"kode_mapel" binding:"required"`
		NamaMapel string  `json:"nama_mapel" binding:"required"`
		Kategori  string  `json:"kategori" binding:"required"`
		KKM       float64 `json:"kkm"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap"})
		return
	}

	req.KodeMapel = strings.ToUpper(strings.TrimSpace(req.KodeMapel))
	req.NamaMapel = strings.TrimSpace(req.NamaMapel)
	if req.KKM <= 0 {
		req.KKM = 75
	}

	// Cek duplicate kode untuk ID lain
	var exists int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM mata_pelajaran WHERE UPPER(kode_mapel) = ? AND id != ?`, req.KodeMapel, id).Scan(&exists)
	if err == nil && exists > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode mata pelajaran sudah digunakan pada mapel lain"})
		return
	}

	_, err = database.DB.Exec(`
		UPDATE mata_pelajaran
		SET kode_mapel = ?, nama_mapel = ?, kategori = ?, kkm = ?
		WHERE id = ?
	`, req.KodeMapel, req.NamaMapel, req.Kategori, req.KKM, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui mata pelajaran: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata pelajaran berhasil diperbarui"})
}

// DeleteMapel - Menghapus mata pelajaran
func DeleteMapel(c *gin.Context) {
	id := c.Param("id")

	var count int
	_ = database.DB.QueryRow(`SELECT COUNT(*) FROM nilai_siswa WHERE mapel_id = ?`, id).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak dapat menghapus mata pelajaran karena sudah memiliki riwayat nilai siswa."})
		return
	}

	res, err := database.DB.Exec(`DELETE FROM mata_pelajaran WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus mata pelajaran"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata pelajaran tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata pelajaran berhasil dihapus"})
}
