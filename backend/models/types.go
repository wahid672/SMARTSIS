package models

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type MenuPermission struct {
	MenuID    int    `json:"menu_id"`
	Code      string `json:"code"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	Path      string `json:"path"`
	SortOrder int    `json:"sort_order"`
	CanView   bool   `json:"can_view"`
	CanCreate bool   `json:"can_create"`
	CanEdit   bool   `json:"can_edit"`
	CanDelete bool   `json:"can_delete"`
}

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	RoleLabel   string `json:"role_label"`
	UserType    string `json:"user_type"` // 'pegawai' atau 'wali'
}

type SiswaSummary struct {
	ID           int    `json:"id"`
	NIS          string `json:"nis"`
	NISN         string `json:"nisn"`
	NamaLengkap  string `json:"nama_lengkap"`
	JenisKelamin string `json:"jenis_kelamin"`
	Kelas        string `json:"kelas"`
	Hubungan     string `json:"hubungan"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Email atau No HP
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string           `json:"token"`
	User        User             `json:"user"`
	Menus       []MenuPermission `json:"menus"`
	DaftarSiswa []SiswaSummary   `json:"daftar_siswa,omitempty"` // Terisi jika akun wali
}
