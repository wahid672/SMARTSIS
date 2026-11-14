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
		"total_masuk":          totalMasuk,
		"total_keluar":         totalKeluar,
		"saldo_kas":            totalMasuk - totalKeluar,
		"total_tagihan_lunas":  totalTagihanLunas,
		"total_tagihan_tunggak": totalTagihanBelum,
	})
}

// ==========================================
// 2. POS KEUANGAN & BUKU KAS UMUM
// ==========================================
func GetPosKeuangan(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT id, nama_pos, tipe, COALESCE(keterangan, '') FROM pos_keuangan ORDER BY tipe ASC, nama_pos ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pos keuangan"})
		return
	}
	defer rows.Close()

	type PosItem struct {
		ID         int    `json:"id"`
		NamaPos    string `json:"nama_pos"`
		Tipe       string `json:"tipe"`
		Keterangan string `json:"keterangan"`
	}

	var list []PosItem
	for rows.Next() {
		var p PosItem
		if err := rows.Scan(&p.ID, &p.NamaPos, &p.Tipe, &p.Keterangan); err == nil {
			list = append(list, p)
		}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama pos dan tipe wajib diisi"})
		return
	}

	res, err := database.DB.Exec(`INSERT INTO pos_keuangan (nama_pos, tipe, keterangan) VALUES (?, ?, ?)`,
		strings.TrimSpace(body.NamaPos), body.Tipe, body.Keterangan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah pos keuangan"})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Pos keuangan berhasil dibuat", "id": id})
}

func GetTransaksiKas(c *gin.Context) {
	query := `
		SELECT t.id, t.pos_id, p.nama_pos, t.kode_transaksi, t.tanggal, t.tipe, t.nominal, COALESCE(t.keterangan, ''), t.created_at
		FROM transaksi_kas t
		JOIN pos_keuangan p ON p.id = t.pos_id
		ORDER BY t.tanggal DESC, t.id DESC
		LIMIT 100
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil transaksi kas"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat transaksi kas"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaksi kas berhasil disimpan", "kode": kode})
}

// ==========================================
// 3. TAGIHAN SISWA / SPP (BULANAN/SEMESTER/TAHUNAN)
// ==========================================
func GetAllTagihan(c *gin.Context) {
	siswaIDStr := c.Query("siswa_id")
	var rows *sql.Rows
	var err error

	baseQuery := `
		SELECT ts.id, ts.siswa_id, s.nama_lengkap, s.nis, s.kelas,
		       pt.id as pos_tagihan_id, pt.nama_tagihan, pt.tipe_periode,
		       ta.nama_tahun, ta.semester,
		       ts.periode, ts.nominal, ts.terbayar, ts.status, COALESCE(ts.tanggal_bayar, '')
		FROM tagihan_siswa ts
		JOIN siswa s ON s.id = ts.siswa_id
		JOIN pos_tagihan pt ON pt.id = ts.pos_tagihan_id
		JOIN tahun_ajaran ta ON ta.id = ts.tahun_ajaran_id
	`

	if siswaIDStr != "" {
		sID, _ := strconv.Atoi(siswaIDStr)
		rows, err = database.DB.Query(baseQuery+` WHERE ts.siswa_id = ? ORDER BY ts.id DESC`, sID)
	} else {
		rows, err = database.DB.Query(baseQuery + ` ORDER BY ts.status ASC, s.nama_lengkap ASC, ts.id DESC LIMIT 150`)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tagihan"})
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
		TipePeriode  string  `json:"tipe_periode"`
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
			&item.PosTagihanID, &item.NamaTagihan, &item.TipePeriode,
			&item.NamaTahun, &item.Semester,
			&item.Periode, &item.Nominal, &item.Terbayar, &item.Status, &item.TanggalBayar); err == nil {
			list = append(list, item)
		}
	}
	c.JSON(http.StatusOK, list)
}

// BayarTagihan memproses pembayaran SPP dan otomatis membuat catatan Kas Masuk
func BayarTagihan(c *gin.Context) {
	idStr := c.Param("id")
	tagihanID, _ := strconv.Atoi(idStr)

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	var siswaID int
	var nominal float64
	var namaTagihan, periode, namaSiswa string

	err = tx.QueryRow(`
		SELECT ts.siswa_id, ts.nominal, pt.nama_tagihan, ts.periode, s.nama_lengkap
		FROM tagihan_siswa ts
		JOIN pos_tagihan pt ON pt.id = ts.pos_tagihan_id
		JOIN siswa s ON s.id = ts.siswa_id
		WHERE ts.id = ?
	`, tagihanID).Scan(&siswaID, &nominal, &namaTagihan, &periode, &namaSiswa)

	if err != nil {
		tx.Rollback()
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
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status tagihan"})
		return
	}

	// Otomatis masukkan ke Transaksi Kas Masuk
	kodeKas := fmt.Sprintf("SPP-%d-%d", tagihanID, time.Now().Unix())
	ket := fmt.Sprintf("Pembayaran %s (%s) an. %s", namaTagihan, periode, namaSiswa)
	today := time.Now().Format("2006-01-02")
	userID := c.MustGet("user_id")

	// Ambil atau gunakan pos penerimaan SPP (ID 1)
	_, _ = tx.Exec(`
		INSERT INTO transaksi_kas (pos_id, kode_transaksi, tanggal, tipe, nominal, keterangan, created_by)
		VALUES (1, ?, ?, 'masuk', ?, ?, ?)
	`, kodeKas, today, nominal, ket, userID)

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran tagihan berhasil diproses & dicatat di Kas Masuk"})
}
