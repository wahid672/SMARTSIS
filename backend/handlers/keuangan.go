package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"smartsis-backend/database"
)

// ==========================================
// 1. RINGKASAN KAS (CASHFLOW DASHBOARD)
// ==========================================
func GetCashflowSummary(c *gin.Context) {
	var totalMasuk, totalKeluar float64

	database.DB.QueryRow(`SELECT COALESCE(SUM(nominal), 0) FROM transaksi_kas WHERE tipe = 'masuk'`).Scan(&totalMasuk)
	database.DB.QueryRow(`SELECT COALESCE(SUM(nominal), 0) FROM transaksi_kas WHERE tipe = 'keluar'`).Scan(&totalKeluar)

	var totalTagihanLunas, totalTagihanBelum float64
	database.DB.QueryRow(`SELECT COALESCE(SUM(terbayar), 0) FROM tagihan_siswa WHERE status = 'lunas'`).Scan(&totalTagihanLunas)
	database.DB.QueryRow(`SELECT COALESCE(SUM(nominal - terbayar), 0) FROM tagihan_siswa WHERE status = 'belum_lunas'`).Scan(&totalTagihanBelum)

	c.JSON(http.StatusOK, gin.H{
		"total_masuk":           totalMasuk,
		"total_keluar":          totalKeluar,
		"saldo_kas":             totalMasuk - totalKeluar,
		"total_tagihan_lunas":   totalTagihanLunas,
		"total_tagihan_tunggak": totalTagihanBelum,
	})
}

// ==========================================
// 2. POS KEUANGAN & BUKU KAS UMUM
// ==========================================

// GetPosKeuangan - Mengambil pos keuangan lengkap dengan akumulasi transaksi & saldo
func GetPosKeuangan(c *gin.Context) {
	query := `
		SELECT 
			p.id, p.nama_pos, p.tipe, COALESCE(p.keterangan, ''),
			COALESCE(SUM(CASE WHEN t.tipe = 'masuk' THEN t.nominal ELSE 0 END), 0) as total_masuk,
			COALESCE(SUM(CASE WHEN t.tipe = 'keluar' THEN t.nominal ELSE 0 END), 0) as total_keluar
		FROM pos_keuangan p
		LEFT JOIN transaksi_kas t ON t.pos_id = p.id
		GROUP BY p.id, p.nama_pos, p.tipe, p.keterangan
		ORDER BY p.tipe ASC, p.nama_pos ASC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pos keuangan: " + err.Error()})
		return
	}
	defer rows.Close()

	type PosItem struct {
		ID          int     `json:"id"`
		NamaPos     string  `json:"nama_pos"`
		Tipe        string  `json:"tipe"`
		Keterangan  string  `json:"keterangan"`
		TotalMasuk  float64 `json:"total_masuk"`
		TotalKeluar float64 `json:"total_keluar"`
		Saldo       float64 `json:"saldo"`
	}

	var list []PosItem
	for rows.Next() {
		var p PosItem
		if err := rows.Scan(&p.ID, &p.NamaPos, &p.Tipe, &p.Keterangan, &p.TotalMasuk, &p.TotalKeluar); err == nil {
			if p.Tipe == "masuk" {
				p.Saldo = p.TotalMasuk
			} else {
				p.Saldo = -p.TotalKeluar
			}
			list = append(list, p)
		}
	}
	if list == nil {
		list = []PosItem{}
	}
	c.JSON(http.StatusOK, list)
}

func CreatePosKeuangan(c *gin.Context) {
	var body struct {
		NamaPos    string `json:"nama_pos" binding:"required"`
		Tipe       string `json:"tipe" binding:"required"`
		Keterangan string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama pos dan tipe (masuk/keluar) wajib diisi"})
		return
	}

	nama := strings.TrimSpace(body.NamaPos)
	var count int
	_ = database.DB.QueryRow(`SELECT COUNT(*) FROM pos_keuangan WHERE LOWER(nama_pos) = LOWER(?) AND tipe = ?`, nama, body.Tipe).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pos keuangan dengan nama dan tipe tersebut sudah ada"})
		return
	}

	res, err := database.DB.Exec(`INSERT INTO pos_keuangan (nama_pos, tipe, keterangan) VALUES (?, ?, ?)`,
		nama, body.Tipe, body.Keterangan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah pos keuangan: " + err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Pos keuangan berhasil dibuat", "id": id})
}

func UpdatePosKeuangan(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		NamaPos    string `json:"nama_pos" binding:"required"`
		Tipe       string `json:"tipe" binding:"required"`
		Keterangan string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap"})
		return
	}

	nama := strings.TrimSpace(body.NamaPos)
	var count int
	_ = database.DB.QueryRow(`SELECT COUNT(*) FROM pos_keuangan WHERE LOWER(nama_pos) = LOWER(?) AND tipe = ? AND id != ?`, nama, body.Tipe, id).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama pos keuangan sudah digunakan pada pos lain"})
		return
	}

	_, err := database.DB.Exec(`UPDATE pos_keuangan SET nama_pos = ?, tipe = ?, keterangan = ? WHERE id = ?`,
		nama, body.Tipe, body.Keterangan, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pos keuangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pos keuangan berhasil diperbarui"})
}

func DeletePosKeuangan(c *gin.Context) {
	id := c.Param("id")

	var transCount int
	_ = database.DB.QueryRow(`SELECT COUNT(*) FROM transaksi_kas WHERE pos_id = ?`, id).Scan(&transCount)
	if transCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pos keuangan tidak dapat dihapus karena sudah memiliki riwayat transaksi buku kas."})
		return
	}

	var tagihanCount int
	_ = database.DB.QueryRow(`SELECT COUNT(*) FROM pos_tagihan WHERE pos_id = ?`, id).Scan(&tagihanCount)
	if tagihanCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pos keuangan tidak dapat dihapus karena terhubung dengan Pos Tagihan Siswa."})
		return
	}

	_, err := database.DB.Exec(`DELETE FROM pos_keuangan WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pos keuangan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pos keuangan berhasil dihapus"})
}

func GetTransaksiKas(c *gin.Context) {
	posID := c.Query("pos_id")
	tipe := c.Query("tipe")

	query := `
		SELECT t.id, t.pos_id, p.nama_pos, t.kode_transaksi, t.tanggal, t.tipe, t.nominal, COALESCE(t.keterangan, ''), t.created_at
		FROM transaksi_kas t
		JOIN pos_keuangan p ON p.id = t.pos_id
		WHERE 1=1
	`
	var args []interface{}
	if posID != "" {
		query += ` AND t.pos_id = ?`
		args = append(args, posID)
	}
	if tipe != "" {
		query += ` AND t.tipe = ?`
		args = append(args, tipe)
	}

	query += ` ORDER BY t.tanggal DESC, t.id DESC LIMIT 150`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil transaksi kas: " + err.Error()})
		return
	}
	defer rows.Close()

	type TransaksiItem struct {
		ID            int     `json:"id"`
		PosID         int     `json:"pos_id"`
		NamaPos       string  `json:"nama_pos"`
		KodeTransaksi string  `json:"kode_transaksi"`
		Tanggal       string  `json:"tanggal"`
		Tipe          string  `json:"tipe"`
		Nominal       float64 `json:"nominal"`
		Keterangan    string  `json:"keterangan"`
		CreatedAt     string  `json:"created_at"`
	}

	var list []TransaksiItem
	for rows.Next() {
		var t TransaksiItem
		if err := rows.Scan(&t.ID, &t.PosID, &t.NamaPos, &t.KodeTransaksi, &t.Tanggal, &t.Tipe, &t.Nominal, &t.Keterangan, &t.CreatedAt); err == nil {
			list = append(list, t)
		}
	}
	if list == nil {
		list = []TransaksiItem{}
	}
	c.JSON(http.StatusOK, list)
}

func CreateTransaksiKas(c *gin.Context) {
	var body struct {
		PosID      int     `json:"pos_id" binding:"required"`
		Tanggal    string  `json:"tanggal" binding:"required"`
		Tipe       string  `json:"tipe" binding:"required"`
		Nominal    float64 `json:"nominal" binding:"required"`
		Keterangan string  `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lengkapi pos, tanggal, tipe, dan nominal transaksi"})
		return
	}

	kode := fmt.Sprintf("KAS-%s-%d", strings.ToUpper(body.Tipe), time.Now().Unix())
	userID := c.MustGet("user_id")

	_, err := database.DB.Exec(`
		INSERT INTO transaksi_kas (pos_id, kode_transaksi, tanggal, tipe, nominal, keterangan, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, body.PosID, kode, body.Tanggal, body.Tipe, body.Nominal, body.Keterangan, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat transaksi kas: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaksi kas berhasil disimpan", "kode": kode})
}

// ==========================================
// 3. MASTER POS TAGIHAN SISWA
// ==========================================
func GetPosTagihan(c *gin.Context) {
	query := `
		SELECT 
			pt.id, pt.pos_id, COALESCE(pk.nama_pos, ''), pt.tahun_ajaran_id,
			COALESCE(ta.nama_tahun, '') as nama_tahun, COALESCE(ta.semester, '') as semester,
			pt.nama_tagihan, COALESCE(pt.tipe_bayar, 'bulanan'), pt.tipe_periode, pt.nominal, COALESCE(pt.keterangan, '')
		FROM pos_tagihan pt
		LEFT JOIN pos_keuangan pk ON pk.id = pt.pos_id
		LEFT JOIN tahun_ajaran ta ON ta.id = pt.tahun_ajaran_id
		ORDER BY pt.id DESC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pos tagihan: " + err.Error()})
		return
	}
	defer rows.Close()

	type PosTagihanItem struct {
		ID            int     `json:"id"`
		PosID         int     `json:"pos_id"`
		NamaPos       string  `json:"nama_pos"`
		TahunAjaranID int     `json:"tahun_ajaran_id"`
		NamaTahun     string  `json:"nama_tahun"`
		Semester      string  `json:"semester"`
		NamaTagihan   string  `json:"nama_tagihan"`
		TipeBayar     string  `json:"tipe_bayar"`
		TipePeriode   string  `json:"tipe_periode"`
		Nominal       float64 `json:"nominal"`
		Keterangan    string  `json:"keterangan"`
	}

	var list []PosTagihanItem
	for rows.Next() {
		var it PosTagihanItem
		if err := rows.Scan(&it.ID, &it.PosID, &it.NamaPos, &it.TahunAjaranID, &it.NamaTahun, &it.Semester, &it.NamaTagihan, &it.TipeBayar, &it.TipePeriode, &it.Nominal, &it.Keterangan); err == nil {
			list = append(list, it)
		}
	}
	if list == nil {
		list = []PosTagihanItem{}
	}
	c.JSON(http.StatusOK, list)
}

func CreatePosTagihan(c *gin.Context) {
	var body struct {
		PosID         int     `json:"pos_id" binding:"required"`
		TahunAjaranID int     `json:"tahun_ajaran_id" binding:"required"`
		NamaTagihan   string  `json:"nama_tagihan" binding:"required"`
		TipeBayar     string  `json:"tipe_bayar" binding:"required"` // 'bulanan' / 'bebas'
		Nominal       float64 `json:"nominal"`
		Keterangan    string  `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lengkapi pos, tahun ajaran, nama pembayaran, dan tipe bayar"})
		return
	}

	tipePeriode := "bulanan"
	if body.TipeBayar == "bebas" {
		tipePeriode = "tahunan"
	}

	res, err := database.DB.Exec(`
		INSERT INTO pos_tagihan (pos_id, tahun_ajaran_id, nama_tagihan, tipe_bayar, tipe_periode, nominal, keterangan)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, body.PosID, body.TahunAjaranID, strings.TrimSpace(body.NamaTagihan), body.TipeBayar, tipePeriode, body.Nominal, body.Keterangan)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat jenis tagihan: " + err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Jenis tagihan pembayaran berhasil ditambahkan", "id": id})
}

func UpdatePosTagihan(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		PosID         int     `json:"pos_id" binding:"required"`
		TahunAjaranID int     `json:"tahun_ajaran_id" binding:"required"`
		NamaTagihan   string  `json:"nama_tagihan" binding:"required"`
		TipeBayar     string  `json:"tipe_bayar" binding:"required"`
		Nominal       float64 `json:"nominal"`
		Keterangan    string  `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap"})
		return
	}

	tipePeriode := "bulanan"
	if body.TipeBayar == "bebas" {
		tipePeriode = "tahunan"
	}

	_, err := database.DB.Exec(`
		UPDATE pos_tagihan
		SET pos_id = ?, tahun_ajaran_id = ?, nama_tagihan = ?, tipe_bayar = ?, tipe_periode = ?, nominal = ?, keterangan = ?
		WHERE id = ?
	`, body.PosID, body.TahunAjaranID, strings.TrimSpace(body.NamaTagihan), body.TipeBayar, tipePeriode, body.Nominal, body.Keterangan, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pos tagihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pos tagihan berhasil diperbarui"})
}

func DeletePosTagihan(c *gin.Context) {
	id := c.Param("id")

	var tagihanCount int
	_ = database.DB.QueryRow(`SELECT COUNT(*) FROM tagihan_siswa WHERE pos_tagihan_id = ?`, id).Scan(&tagihanCount)
	if tagihanCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pos tagihan tidak dapat dihapus karena sudah memiliki data tagihan siswa."})
		return
	}

	_, err := database.DB.Exec(`DELETE FROM pos_tagihan WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pos tagihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pos tagihan berhasil dihapus"})
}

// ==========================================
// 4. SETTING TARIF TAGIHAN KELAS & SISWA
// ==========================================

// GetTarifKelas - Mengambil daftar siswa dan tagihan pada pos tertentu di suatu kelas
func GetTarifKelas(c *gin.Context) {
	posTagihanIDStr := c.Query("pos_tagihan_id")
	kelas := c.Query("kelas")

	if posTagihanIDStr == "" || kelas == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pos_tagihan_id dan kelas wajib diisi"})
		return
	}

	posTagihanID, _ := strconv.Atoi(posTagihanIDStr)

	// Info Pos Tagihan
	var posInfo struct {
		ID            int     `json:"id"`
		NamaTagihan   string  `json:"nama_tagihan"`
		TipeBayar     string  `json:"tipe_bayar"`
		NamaPos       string  `json:"nama_pos"`
		TahunAjaranID int     `json:"tahun_ajaran_id"`
		NamaTahun     string  `json:"nama_tahun"`
		Nominal       float64 `json:"nominal"`
	}
	err := database.DB.QueryRow(`
		SELECT pt.id, pt.nama_tagihan, COALESCE(pt.tipe_bayar, 'bulanan'), pk.nama_pos, pt.tahun_ajaran_id, ta.nama_tahun, pt.nominal
		FROM pos_tagihan pt
		JOIN pos_keuangan pk ON pk.id = pt.pos_id
		JOIN tahun_ajaran ta ON ta.id = pt.tahun_ajaran_id
		WHERE pt.id = ?
	`, posTagihanID).Scan(&posInfo.ID, &posInfo.NamaTagihan, &posInfo.TipeBayar, &posInfo.NamaPos, &posInfo.TahunAjaranID, &posInfo.NamaTahun, &posInfo.Nominal)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pos tagihan tidak ditemukan"})
		return
	}

	// Ambil Siswa di kelas tersebut
	rows, err := database.DB.Query(`SELECT id, nis, nama_lengkap, kelas FROM siswa WHERE kelas = ? ORDER BY nama_lengkap ASC`, kelas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa: " + err.Error()})
		return
	}
	defer rows.Close()

	type BulanItem struct {
		ID       int     `json:"id"`
		Bulan    string  `json:"bulan"`
		Nominal  float64 `json:"nominal"`
		Terbayar float64 `json:"terbayar"`
		Status   string  `json:"status"`
	}

	type SiswaTarifRow struct {
		SiswaID      int         `json:"siswa_id"`
		NIS          string      `json:"nis"`
		NamaLengkap  string      `json:"nama_lengkap"`
		Kelas        string      `json:"kelas"`
		TotalTagihan float64     `json:"total_tagihan"`
		BulanList    []BulanItem `json:"bulan_list"`
	}

	var students []SiswaTarifRow
	for rows.Next() {
		var s SiswaTarifRow
		if err := rows.Scan(&s.SiswaID, &s.NIS, &s.NamaLengkap, &s.Kelas); err == nil {
			// Query rincian tagihan per bulan untuk siswa ini
			tRows, tErr := database.DB.Query(`
				SELECT id, COALESCE(bulan, periode), nominal, terbayar, status
				FROM tagihan_siswa
				WHERE siswa_id = ? AND pos_tagihan_id = ?
				ORDER BY id ASC
			`, s.SiswaID, posTagihanID)

			s.BulanList = []BulanItem{}
			if tErr == nil {
				for tRows.Next() {
					var b BulanItem
					if err := tRows.Scan(&b.ID, &b.Bulan, &b.Nominal, &b.Terbayar, &b.Status); err == nil {
						s.TotalTagihan += b.Nominal
						s.BulanList = append(s.BulanList, b)
					}
				}
				tRows.Close()
			}
			students = append(students, s)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"pos_tagihan": posInfo,
		"kelas":       kelas,
		"students":    students,
	})
}

// SaveTarifKelasBatch - Menerapkan tarif tagihan 12 bulan (atau bebas) untuk seluruh siswa sekelas
func SaveTarifKelasBatch(c *gin.Context) {
	var body struct {
		PosTagihanID  int                `json:"pos_tagihan_id" binding:"required"`
		TahunAjaranID int                `json:"tahun_ajaran_id" binding:"required"`
		Kelas         string             `json:"kelas" binding:"required"`
		TipeBayar     string             `json:"tipe_bayar" binding:"required"`
		NominalBebas  float64            `json:"nominal_bebas"`
		BulanTarif    map[string]float64 `json:"bulan_tarif"` // {"Juli": 350000, "Agustus": 350000, ...}
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tarif tidak valid: " + err.Error()})
		return
	}

	// Ambil semua siswa di kelas ini
	rows, err := database.DB.Query(`SELECT id FROM siswa WHERE kelas = ?`, body.Kelas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar siswa kelas: " + err.Error()})
		return
	}
	defer rows.Close()

	var siswaIDs []int
	for rows.Next() {
		var sid int
		if rows.Scan(&sid) == nil {
			siswaIDs = append(siswaIDs, sid)
		}
	}

	if len(siswaIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada siswa yang terdaftar di kelas " + body.Kelas})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}
	defer tx.Rollback()

	daftarBulan := []string{"Juli", "Agustus", "September", "Oktober", "November", "Desember", "Januari", "Februari", "Maret", "April", "Mei", "Juni"}

	for _, sid := range siswaIDs {
		if body.TipeBayar == "bulanan" {
			for _, b := range daftarBulan {
				nominal := body.BulanTarif[b]
				// Cek apakah sudah ada tagihan untuk siswa + pos + bulan
				var existingID int
				var existingStatus string
				err := tx.QueryRow(`
					SELECT id, status FROM tagihan_siswa 
					WHERE siswa_id = ? AND pos_tagihan_id = ? AND (bulan = ? OR periode = ?)
				`, sid, body.PosTagihanID, b, b).Scan(&existingID, &existingStatus)

				if err == sql.ErrNoRows {
					// Insert baru
					_, _ = tx.Exec(`
						INSERT INTO tagihan_siswa (siswa_id, pos_tagihan_id, tahun_ajaran_id, periode, bulan, nominal, terbayar, status)
						VALUES (?, ?, ?, ?, ?, ?, 0, 'belum_lunas')
					`, sid, body.PosTagihanID, body.TahunAjaranID, b, b, nominal)
				} else if existingStatus != "lunas" {
					// Update jika belum lunas
					_, _ = tx.Exec(`
						UPDATE tagihan_siswa SET nominal = ?, bulan = ? WHERE id = ?
					`, nominal, b, existingID)
				}
			}
		} else {
			// Tipe Bebas
			var existingID int
			var existingStatus string
			err := tx.QueryRow(`
				SELECT id, status FROM tagihan_siswa 
				WHERE siswa_id = ? AND pos_tagihan_id = ?
			`, sid, body.PosTagihanID).Scan(&existingID, &existingStatus)

			if err == sql.ErrNoRows {
				_, _ = tx.Exec(`
					INSERT INTO tagihan_siswa (siswa_id, pos_tagihan_id, tahun_ajaran_id, periode, bulan, nominal, terbayar, status)
					VALUES (?, ?, ?, 'Tagihan Bebas', 'Tagihan Bebas', ?, 0, 'belum_lunas')
				`, sid, body.PosTagihanID, body.TahunAjaranID, body.NominalBebas)
			} else if existingStatus != "lunas" {
				_, _ = tx.Exec(`
					UPDATE tagihan_siswa SET nominal = ? WHERE id = ?
				`, body.NominalBebas, existingID)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan tarif ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Tarif pembayaran berhasil diterapkan untuk %d siswa di kelas %s", len(siswaIDs), body.Kelas),
	})
}

// DeleteTarifKelas - Menghapus tagihan yang belum lunas untuk satu kelas pada pos tertentu
func DeleteTarifKelas(c *gin.Context) {
	posTagihanID := c.Query("pos_tagihan_id")
	kelas := c.Query("kelas")

	if posTagihanID == "" || kelas == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pos_tagihan_id dan kelas wajib diisi"})
		return
	}

	res, err := database.DB.Exec(`
		DELETE FROM tagihan_siswa
		WHERE pos_tagihan_id = ? AND status != 'lunas' AND siswa_id IN (
			SELECT id FROM siswa WHERE kelas = ?
		)
	`, posTagihanID, kelas)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus tagihan kelas: " + err.Error()})
		return
	}

	affected, _ := res.RowsAffected()
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Berhasil menghapus %d data tagihan yang belum lunas untuk kelas %s", affected, kelas),
	})
}

// ==========================================
// 5. TRANSAKSI PEMBAYARAN KASIR & WALI
// ==========================================
func GetAllTagihan(c *gin.Context) {
	siswaIDStr := c.Query("siswa_id")
	posTagihanIDStr := c.Query("pos_tagihan_id")
	status := c.Query("status")
	kelas := c.Query("kelas")

	userType := c.GetString("user_type")
	userIDVal, _ := c.Get("user_id")

	baseQuery := `
		SELECT ts.id, ts.siswa_id, s.nama_lengkap, s.nis, s.kelas,
		       pt.id as pos_tagihan_id, pt.nama_tagihan, COALESCE(pt.tipe_bayar, 'bulanan'),
		       ta.nama_tahun, ta.semester,
		       COALESCE(ts.bulan, ts.periode), ts.nominal, ts.terbayar, ts.status, COALESCE(ts.tanggal_bayar, '')
		FROM tagihan_siswa ts
		JOIN siswa s ON s.id = ts.siswa_id
		JOIN pos_tagihan pt ON pt.id = ts.pos_tagihan_id
		JOIN tahun_ajaran ta ON ta.id = ts.tahun_ajaran_id
		WHERE 1=1
	`
	var args []interface{}

	// Jika akun wali, pastikan hanya melihat siswa miliknya
	if userType == "wali" {
		baseQuery += ` AND ts.siswa_id IN (SELECT siswa_id FROM wali_siswa WHERE wali_user_id = ?)`
		args = append(args, userIDVal)
	}

	if siswaIDStr != "" {
		sID, _ := strconv.Atoi(siswaIDStr)
		baseQuery += ` AND ts.siswa_id = ?`
		args = append(args, sID)
	}

	if posTagihanIDStr != "" {
		pID, _ := strconv.Atoi(posTagihanIDStr)
		baseQuery += ` AND ts.pos_tagihan_id = ?`
		args = append(args, pID)
	}

	if status != "" {
		baseQuery += ` AND ts.status = ?`
		args = append(args, status)
	}

	if kelas != "" {
		baseQuery += ` AND s.kelas = ?`
		args = append(args, kelas)
	}

	baseQuery += ` ORDER BY ts.status ASC, s.nama_lengkap ASC, ts.id ASC LIMIT 250`

	rows, err := database.DB.Query(baseQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tagihan: " + err.Error()})
		return
	}
	defer rows.Close()

	type TagihanItem struct {
		ID           int     `json:"id"`
		SiswaID      int     `json:"siswa_id"`
		NamaSiswa    string  `json:"nama_siswa"`
		NIS          string  `json:"nis"`
		Kelas        string  `json:"kelas"`
		PosTagihanID int     `json:"pos_tagihan_id"`
		NamaTagihan  string  `json:"nama_tagihan"`
		TipeBayar    string  `json:"tipe_bayar"`
		NamaTahun    string  `json:"nama_tahun"`
		Semester     string  `json:"semester"`
		Periode      string  `json:"periode"`
		Nominal      float64 `json:"nominal"`
		Terbayar     float64 `json:"terbayar"`
		Status       string  `json:"status"`
		TanggalBayar string  `json:"tanggal_bayar"`
	}

	var list []TagihanItem
	for rows.Next() {
		var item TagihanItem
		if err := rows.Scan(&item.ID, &item.SiswaID, &item.NamaSiswa, &item.NIS, &item.Kelas,
			&item.PosTagihanID, &item.NamaTagihan, &item.TipeBayar,
			&item.NamaTahun, &item.Semester,
			&item.Periode, &item.Nominal, &item.Terbayar, &item.Status, &item.TanggalBayar); err == nil {
			list = append(list, item)
		}
	}
	if list == nil {
		list = []TagihanItem{}
	}
	c.JSON(http.StatusOK, list)
}

// BayarTagihan memproses pembayaran SPP / tagihan dan otomatis membuat catatan Kas Masuk
func BayarTagihan(c *gin.Context) {
	idStr := c.Param("id")
	tagihanID, _ := strconv.Atoi(idStr)

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}
	defer tx.Rollback()

	var siswaID, posKasID int
	var nominal float64
	var namaTagihan, periode, namaSiswa string

	err = tx.QueryRow(`
		SELECT ts.siswa_id, ts.nominal, pt.nama_tagihan, COALESCE(ts.bulan, ts.periode), s.nama_lengkap, pt.pos_id
		FROM tagihan_siswa ts
		JOIN pos_tagihan pt ON pt.id = ts.pos_tagihan_id
		JOIN siswa s ON s.id = ts.siswa_id
		WHERE ts.id = ?
	`, tagihanID).Scan(&siswaID, &nominal, &namaTagihan, &periode, &namaSiswa, &posKasID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tagihan tidak ditemukan"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// Update status tagihan lunas
	_, err = tx.Exec(`
		UPDATE tagihan_siswa SET status = 'lunas', terbayar = nominal, tanggal_bayar = ?
		WHERE id = ?
	`, now, tagihanID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status tagihan"})
		return
	}

	if posKasID <= 0 {
		posKasID = 1
	}

	// Otomatis masukkan ke Transaksi Kas Masuk
	kodeKas := fmt.Sprintf("BYR-%d-%d", tagihanID, time.Now().Unix())
	ket := fmt.Sprintf("Pembayaran %s (%s) an. %s", namaTagihan, periode, namaSiswa)
	today := time.Now().Format("2006-01-02")
	userID := c.MustGet("user_id")

	_, _ = tx.Exec(`
		INSERT INTO transaksi_kas (pos_id, kode_transaksi, tanggal, tipe, nominal, keterangan, created_by)
		VALUES (?, ?, ?, 'masuk', ?, ?, ?)
	`, posKasID, kodeKas, today, nominal, ket, userID)

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran tagihan berhasil diproses & dicatat di Buku Kas Masuk"})
}
