package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"smartsis-backend/database"
)

type Kelas struct {
	ID         int    `json:"id"`
	NamaKelas  string `json:"nama_kelas"`
	Tingkat    string `json:"tingkat"`
	WaliKelas  string `json:"wali_kelas"`
	Keterangan string `json:"keterangan"`
	SiswaCount int    `json:"siswa_count"`
}

// GetAllKelas mengambil semua kelas beserta jumlah siswa di dalamnya
func GetAllKelas(c *gin.Context) {
	query := `
		SELECT k.id, k.nama_kelas, k.tingkat, COALESCE(k.wali_kelas, ''), COALESCE(k.keterangan, ''),
		       (SELECT COUNT(*) FROM siswa s WHERE s.kelas = k.nama_kelas AND s.is_active = 1) as siswa_count
		FROM kelas k
		ORDER BY k.tingkat ASC, k.nama_kelas ASC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kelas"})
		return
	}
	defer rows.Close()

	var list []Kelas
	for rows.Next() {
		var k Kelas
		if err := rows.Scan(&k.ID, &k.NamaKelas, &k.Tingkat, &k.WaliKelas, &k.Keterangan, &k.SiswaCount); err == nil {
			list = append(list, k)
		}
	}

	c.JSON(http.StatusOK, list)
}

// CreateKelas menambah kelas baru
func CreateKelas(c *gin.Context) {
	var body struct {
		NamaKelas  string `json:"nama_kelas" binding:"required"`
		Tingkat    string `json:"tingkat" binding:"required"`
		WaliKelas  string `json:"wali_kelas"`
		Keterangan string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kelas dan tingkat wajib diisi"})
		return
	}

	namaClean := strings.TrimSpace(body.NamaKelas)
	res, err := database.DB.Exec(`
		INSERT INTO kelas (nama_kelas, tingkat, wali_kelas, keterangan)
		VALUES (?, ?, ?, ?)
	`, namaClean, body.Tingkat, body.WaliKelas, body.Keterangan)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kelas sudah terdaftar"})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Kelas baru berhasil dibuat", "id": id})
}

// UpdateKelas mengubah data kelas
func UpdateKelas(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var body struct {
		NamaKelas  string `json:"nama_kelas" binding:"required"`
		Tingkat    string `json:"tingkat" binding:"required"`
		WaliKelas  string `json:"wali_kelas"`
		Keterangan string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	namaClean := strings.TrimSpace(body.NamaKelas)
	_, err := database.DB.Exec(`
		UPDATE kelas SET nama_kelas = ?, tingkat = ?, wali_kelas = ?, keterangan = ?
		WHERE id = ?
	`, namaClean, body.Tingkat, body.WaliKelas, body.Keterangan, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal memperbarui kelas: Nama mungkin bentrok"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data kelas berhasil diperbarui"})
}

// DeleteKelas menghapus kelas
func DeleteKelas(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	// Cek apakah ada siswa di kelas ini
	var namaKelas string
	err := database.DB.QueryRow(`SELECT nama_kelas FROM kelas WHERE id = ?`, id).Scan(&namaKelas)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas tidak ditemukan"})
		return
	}

	var count int
	database.DB.QueryRow(`SELECT COUNT(*) FROM siswa WHERE kelas = ?`, namaKelas).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kelas tidak bisa dihapus karena masih ada siswa di dalamnya"})
		return
	}

	database.DB.Exec(`DELETE FROM kelas WHERE id = ?`, id)
	c.JSON(http.StatusOK, gin.H{"message": "Kelas berhasil dihapus"})
}
