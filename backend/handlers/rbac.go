package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"smartsis-backend/database"
	"smartsis-backend/models"
)

// GetAllRoles mengambil seluruh daftar role
func GetAllRoles(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.name, r.label, COALESCE(r.description, ''), COUNT(u.id) as user_count
		FROM roles r
		LEFT JOIN users u ON u.role_id = r.id
		GROUP BY r.id
		ORDER BY r.id ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar role"})
		return
	}
	defer rows.Close()

	type RoleWithCount struct {
		models.Role
		UserCount int `json:"user_count"`
	}

	var list []RoleWithCount
	for rows.Next() {
		var r RoleWithCount
		if err := rows.Scan(&r.ID, &r.Name, &r.Label, &r.Description, &r.UserCount); err == nil {
			list = append(list, r)
		}
	}
	c.JSON(http.StatusOK, list)
}

// CreateRole membuat role baru
func CreateRole(c *gin.Context) {
	var body struct {
		Name        string `json:"name" binding:"required"`
		Label       string `json:"label" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan label role wajib diisi"})
		return
	}

	codeName := strings.ToLower(strings.ReplaceAll(body.Name, " ", "_"))

	res, err := database.DB.Exec(`INSERT INTO roles (name, label, description) VALUES (?, ?, ?)`,
		codeName, body.Label, body.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama role sudah ada atau tidak valid"})
		return
	}

	roleID, _ := res.LastInsertId()

	// Default beri akses dashboard & profil
	database.DB.Exec(`INSERT OR IGNORE INTO role_menus (role_id, menu_id, can_view)
		SELECT ?, id, 1 FROM menus WHERE code IN ('dashboard', 'profil')`, roleID)

	c.JSON(http.StatusCreated, gin.H{"message": "Role baru berhasil ditambahkan", "id": roleID})
}

// DeleteRole menghapus role (role ID 1/Super Admin tidak boleh dihapus)
func DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	roleID, _ := strconv.Atoi(idStr)

	if roleID <= 3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Role bawaan sistem tidak boleh dihapus"})
		return
	}

	// Cek apakah ada pegawai yang memakai role ini
	var count int
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role_id = ?`, roleID).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak bisa dihapus karena masih digunakan oleh pegawai"})
		return
	}

	database.DB.Exec(`DELETE FROM role_menus WHERE role_id = ?`, roleID)
	database.DB.Exec(`DELETE FROM roles WHERE id = ?`, roleID)

	c.JSON(http.StatusOK, gin.H{"message": "Role berhasil dihapus"})
}

// GetRolePermissions mengambil semua menu master beserta status centang izin untuk role tertentu
func GetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	roleID, _ := strconv.Atoi(idStr)

	query := `
		SELECT m.id, m.code, m.title, m.icon, m.path, COALESCE(m.category, 'LAINNYA'), m.sort_order,
		       COALESCE(rm.can_view, 0) as can_view,
		       COALESCE(rm.can_create, 0) as can_create,
		       COALESCE(rm.can_edit, 0) as can_edit,
		       COALESCE(rm.can_delete, 0) as can_delete
		FROM menus m
		LEFT JOIN role_menus rm ON rm.menu_id = m.id AND rm.role_id = ?
		WHERE m.is_active = 1
		ORDER BY m.sort_order ASC
	`

	rows, err := database.DB.Query(query, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil hak akses menu"})
		return
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

	c.JSON(http.StatusOK, list)
}

// UpdateRolePermissions menyimpan matriks centang hak akses menu untuk role tertentu
func UpdateRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	roleID, _ := strconv.Atoi(idStr)

	var permissions []struct {
		MenuID    int  `json:"menu_id"`
		CanView   bool `json:"can_view"`
		CanCreate bool `json:"can_create"`
		CanEdit   bool `json:"can_edit"`
		CanDelete bool `json:"can_delete"`
	}

	if err := c.ShouldBindJSON(&permissions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data hak akses tidak valid"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
		return
	}

	// Hapus hak akses lama untuk role ini
	_, _ = tx.Exec(`DELETE FROM role_menus WHERE role_id = ?`, roleID)

	stmt, err := tx.Prepare(`
		INSERT INTO role_menus (role_id, menu_id, can_view, can_create, can_edit, can_delete)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyiapkan query"})
		return
	}
	defer stmt.Close()

	for _, p := range permissions {
		// Hanya simpan jika minimal can_view aktif
		if p.CanView {
			view := 1
			create := 0
			edit := 0
			del := 0
			if p.CanCreate {
				create = 1
			}
			if p.CanEdit {
				edit = 1
			}
			if p.CanDelete {
				del = 1
			}
			_, _ = stmt.Exec(roleID, p.MenuID, view, create, edit, del)
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan hak akses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hak akses menu berhasil diperbarui"})
}
