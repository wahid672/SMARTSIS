package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"smartsis-backend/database"
	"smartsis-backend/models"
)

// GetAllPegawai mengambil daftar pegawai dan staf
func GetAllPegawai(c *gin.Context) {
	query := `
		SELECT u.id, u.name, COALESCE(u.email, ''), COALESCE(u.phone_number, ''),
		       u.role_id, r.name, r.label, u.user_type, u.is_active
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.user_type = 'pegawai'
		ORDER BY u.id ASC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai"})
		return
	}
	defer rows.Close()

	type PegawaiItem struct {
		models.User
		IsActive bool `json:"is_active"`
	}

	var list []PegawaiItem
	for rows.Next() {
		var p PegawaiItem
		var active int
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.PhoneNumber,
			&p.RoleID, &p.RoleName, &p.RoleLabel, &p.UserType, &active); err == nil {
			p.IsActive = active == 1
			list = append(list, p)
		}
	}

	c.JSON(http.StatusOK, list)
}

// CreatePegawai menambah akun pegawai/guru baru
func CreatePegawai(c *gin.Context) {
	var body struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password" binding:"required"`
		RoleID      int    `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, email, password, dan role wajib diisi"})
		return
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
		return
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(body.Email))
	cleanPhone := strings.TrimSpace(body.PhoneNumber)

	res, err := database.DB.Exec(`
		INSERT INTO users (name, email, phone_number, password_hash, role_id, user_type, is_active)
		VALUES (?, ?, ?, ?, ?, 'pegawai', 1)
	`, body.Name, cleanEmail, cleanPhone, string(hashPass), body.RoleID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau nomor telepon sudah terdaftar pada akun lain"})
		return
	}

	newID, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Pegawai berhasil ditambahkan", "id": newID})
}

// UpdatePegawai mengubah profil atau role pegawai
func UpdatePegawai(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := strconv.Atoi(idStr)

	var body struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"` // opsional
		RoleID      int    `json:"role_id" binding:"required"`
		IsActive    bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(body.Email))
	activeInt := 0
	if body.IsActive {
		activeInt = 1
	}

	if body.Password != "" {
		hashPass, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		_, err := database.DB.Exec(`
			UPDATE users SET name = ?, email = ?, phone_number = ?, password_hash = ?, role_id = ?, is_active = ?
			WHERE id = ? AND user_type = 'pegawai'
		`, body.Name, cleanEmail, body.PhoneNumber, string(hashPass), body.RoleID, activeInt, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal memperbarui: Email/No HP mungkin bentrok"})
			return
		}
	} else {
		_, err := database.DB.Exec(`
			UPDATE users SET name = ?, email = ?, phone_number = ?, role_id = ?, is_active = ?
			WHERE id = ? AND user_type = 'pegawai'
		`, body.Name, cleanEmail, body.PhoneNumber, body.RoleID, activeInt, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal memperbarui: Email/No HP mungkin bentrok"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data pegawai berhasil diperbarui"})
}

// DeletePegawai menghapus akun pegawai
func DeletePegawai(c *gin.Context) {
	idStr := c.Param("id")
	userID, _ := strconv.Atoi(idStr)

	if userID == 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Administrator Utama tidak boleh dihapus"})
		return
	}

	res, err := database.DB.Exec(`DELETE FROM users WHERE id = ? AND user_type = 'pegawai'`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pegawai"})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pegawai berhasil dihapus"})
}
