package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"smartsis-backend/database"
	"smartsis-backend/models"
)

var JWTSecret = []byte("smartsis-super-secret-key-2026")

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	identifier := strings.TrimSpace(req.Identifier)

	// Deteksi pintar: apakah email (ada @) atau nomor HP/WA
	var query string
	if strings.Contains(identifier, "@") {
		query = `SELECT u.id, u.name, COALESCE(u.email, ''), COALESCE(u.phone_number, ''), 
		                u.password_hash, u.role_id, r.name, r.label, u.user_type, u.is_active
		         FROM users u
		         JOIN roles r ON r.id = u.role_id
		         WHERE LOWER(u.email) = LOWER(?)`
	} else {
		// Bersihkan format no hp (misal +62 atau spasi)
		cleanPhone := strings.ReplaceAll(identifier, " ", "")
		cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
		if strings.HasPrefix(cleanPhone, "+62") {
			cleanPhone = "0" + cleanPhone[3:]
		} else if strings.HasPrefix(cleanPhone, "62") {
			cleanPhone = "0" + cleanPhone[2:]
		}

		query = `SELECT u.id, u.name, COALESCE(u.email, ''), COALESCE(u.phone_number, ''), 
		                u.password_hash, u.role_id, r.name, r.label, u.user_type, u.is_active
		         FROM users u
		         JOIN roles r ON r.id = u.role_id
		         WHERE u.phone_number = ? OR u.phone_number = ?`
		identifier = cleanPhone
	}

	var user models.User
	var passwordHash string
	var isActive int

	var row *sql.Row
	if strings.Contains(req.Identifier, "@") {
		row = database.DB.QueryRow(query, identifier)
	} else {
		row = database.DB.QueryRow(query, identifier, req.Identifier)
	}

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber,
		&passwordHash, &user.RoleID, &user.RoleName, &user.RoleLabel, &user.UserType, &isActive)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email / No HP atau kata sandi salah"})
		return
	}

	if isActive == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda dinonaktifkan oleh administrator"})
		return
	}

	// Cek kata sandi
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email / No HP atau kata sandi salah"})
		return
	}

	// Buat token JWT
	tokenClaims := jwt.MapClaims{
		"user_id":   user.ID,
		"role_id":   user.RoleID,
		"user_type": user.UserType,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 hari
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan token"})
		return
	}

	// Ambil hak akses menu dinamis berdasarkan role_id
	menus, _ := GetMenusByRole(user.RoleID)

	// Jika tipe user adalah 'wali', ambil juga daftar anak/siswanya
	var daftarSiswa []models.SiswaSummary
	if user.UserType == "wali" {
		daftarSiswa, _ = GetSiswaByWali(user.ID)
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:       tokenString,
		User:        user,
		Menus:       menus,
		DaftarSiswa: daftarSiswa,
	})
}

func GetMenusByRole(roleID int) ([]models.MenuPermission, error) {
	query := `
		SELECT m.id, m.code, m.title, m.icon, m.path, COALESCE(m.category, 'LAINNYA'), m.sort_order,
		       rm.can_view, rm.can_create, rm.can_edit, rm.can_delete
		FROM menus m
		JOIN role_menus rm ON rm.menu_id = m.id
		WHERE rm.role_id = ? AND rm.can_view = 1 AND m.is_active = 1
		ORDER BY m.sort_order ASC
	`
	rows, err := database.DB.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.MenuPermission
	for rows.Next() {
		var item models.MenuPermission
		var canView, canCreate, canEdit, canDelete int
		if err := rows.Scan(&item.MenuID, &item.Code, &item.Title, &item.Icon, &item.Path, &item.Category, &item.SortOrder,
			&canView, &canCreate, &canEdit, &canDelete); err == nil {
			item.CanView = canView == 1
			item.CanCreate = canCreate == 1
			item.CanEdit = canEdit == 1
			item.CanDelete = canDelete == 1
			list = append(list, item)
		}
	}
	return list, nil
}

func GetSiswaByWali(waliUserID int) ([]models.SiswaSummary, error) {
	query := `
		SELECT s.id, s.nis, COALESCE(s.nisn, ''), s.nama_lengkap, s.jenis_kelamin, s.kelas, ws.hubungan
		FROM siswa s
		JOIN wali_siswa ws ON ws.siswa_id = s.id
		WHERE ws.wali_user_id = ? AND s.is_active = 1
	`
	rows, err := database.DB.Query(query, waliUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SiswaSummary
	for rows.Next() {
		var s models.SiswaSummary
		if err := rows.Scan(&s.ID, &s.NIS, &s.NISN, &s.NamaLengkap, &s.JenisKelamin, &s.Kelas, &s.Hubungan); err == nil {
			list = append(list, s)
		}
	}
	return list, nil
}
