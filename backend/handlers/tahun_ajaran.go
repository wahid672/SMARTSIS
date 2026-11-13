package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"smartsis-backend/database"
)

type TahunAjaran struct {
	ID       int    `json:"id"`
	Nama     string `json:"nama_tahun"`
	Semester string `json:"semester"`
	IsActive bool   `json:"is_active"`
}

// GetAllTahunAjaran mengambil semua daftar tahun ajaran
func GetAllTahunAjaran(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT id, nama_tahun, semester, is_active FROM tahun_ajaran ORDER BY id DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tahun ajaran"})
		return
	}
	defer rows.Close()

	var list []TahunAjaran
	for rows.Next() {
		var t TahunAjaran
		var active int
		if err := rows.Scan(&t.ID, &t.Nama, &t.Semester, &active); err == nil {
			t.IsActive = active == 1
			list = append(list, t)
		}
	}
	c.JSON(http.StatusOK, list)
}

// CreateTahunAjaran membuat tahun ajaran baru
func CreateTahunAjaran(c *gin.Context) {
	var body struct {
		Nama     string `json:"nama_tahun" binding:"required"`
		Semester string `json:"semester" binding:"required"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama tahun ajaran dan semester wajib diisi"})
		return
	}

	activeInt := 0
	if body.IsActive {
		activeInt = 1
		// Nonaktifkan tahun ajaran lainnya
		database.DB.Exec(`UPDATE tahun_ajaran SET is_active = 0`)
	}

	res, err := database.DB.Exec(`INSERT INTO tahun_ajaran (nama_tahun, semester, is_active) VALUES (?, ?, ?)`,
		strings.TrimSpace(body.Nama), body.Semester, activeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan tahun ajaran"})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Tahun ajaran berhasil dibuat", "id": id})
}

// SetActiveTahunAjaran mengaktifkan satu tahun ajaran sebagai acuan transaksi
func SetActiveTahunAjaran(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	database.DB.Exec(`UPDATE tahun_ajaran SET is_active = 0`)
	_, err := database.DB.Exec(`UPDATE tahun_ajaran SET is_active = 1 WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengaktifkan tahun ajaran"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tahun ajaran aktif berhasil diperbarui"})
}

// DeleteTahunAjaran menghapus tahun ajaran
func DeleteTahunAjaran(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var active int
	database.DB.QueryRow(`SELECT is_active FROM tahun_ajaran WHERE id = ?`, id).Scan(&active)
	if active == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun ajaran yang sedang aktif tidak boleh dihapus"})
		return
	}

	database.DB.Exec(`DELETE FROM tahun_ajaran WHERE id = ?`, id)
	c.JSON(http.StatusOK, gin.H{"message": "Tahun ajaran berhasil dihapus"})
}
