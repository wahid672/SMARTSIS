package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"smartsis-backend/database"
	"smartsis-backend/handlers"
	"smartsis-backend/middleware"
)

// Embed folder hasil build frontend Ionic Vue 3 ke dalam binary executable
//
//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	// Inisialisasi Database SQLite
	database.InitDB()

	// Gin Engine
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// 1. API Routing
	api := r.Group("/api")
	{
		// Public Auth (Smart Login: Email pegawai vs No HP wali)
		api.POST("/auth/login", handlers.Login)

		// Public Settings (Agar halaman login bisa menampilkan nama lembaga & logo)
		api.GET("/settings", handlers.GetSettings)

		// Protected Route
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// RBAC Management
			protected.GET("/admin/roles", handlers.GetAllRoles)
			protected.POST("/admin/roles", handlers.CreateRole)
			protected.DELETE("/admin/roles/:id", handlers.DeleteRole)
			protected.GET("/admin/roles/:id/permissions", handlers.GetRolePermissions)
			protected.POST("/admin/roles/:id/permissions", handlers.UpdateRolePermissions)

			// Pegawai Management
			protected.GET("/admin/pegawai", handlers.GetAllPegawai)
			protected.POST("/admin/pegawai", handlers.CreatePegawai)
			protected.PUT("/admin/pegawai/:id", handlers.UpdatePegawai)
			protected.DELETE("/admin/pegawai/:id", handlers.DeletePegawai)

			// Kelas Management
			protected.GET("/admin/kelas", handlers.GetAllKelas)
			protected.POST("/admin/kelas", handlers.CreateKelas)
			protected.PUT("/admin/kelas/:id", handlers.UpdateKelas)
			protected.DELETE("/admin/kelas/:id", handlers.DeleteKelas)

			// Siswa & Wali Management
			protected.GET("/admin/siswa", handlers.GetAllSiswa)
			protected.POST("/admin/siswa", handlers.CreateSiswa)
			protected.PUT("/admin/siswa/:id", handlers.UpdateSiswa)
			protected.DELETE("/admin/siswa/:id", handlers.DeleteSiswa)

			// Master Tahun Ajaran
			protected.GET("/admin/tahun-ajaran", handlers.GetAllTahunAjaran)
			protected.POST("/admin/tahun-ajaran", handlers.CreateTahunAjaran)
			protected.PUT("/admin/tahun-ajaran/:id/activate", handlers.SetActiveTahunAjaran)
			protected.DELETE("/admin/tahun-ajaran/:id", handlers.DeleteTahunAjaran)

			// Modul Keuangan & SPP
			protected.GET("/keuangan/summary", handlers.GetCashflowSummary)
			protected.GET("/keuangan/pos", handlers.GetPosKeuangan)
			protected.POST("/keuangan/pos", handlers.CreatePosKeuangan)
			protected.GET("/keuangan/kas", handlers.GetTransaksiKas)
			protected.POST("/keuangan/kas", handlers.CreateTransaksiKas)
			protected.GET("/keuangan/tagihan", handlers.GetAllTagihan)
			protected.POST("/keuangan/tagihan/:id/bayar", handlers.BayarTagihan)

			// Modul Presensi & Absensi Santri
			protected.GET("/presensi/kelas", handlers.GetPresensiByKelasDanTanggal)
			protected.POST("/presensi/batch", handlers.SavePresensiBatch)
			protected.GET("/presensi/wali", handlers.GetRekapPresensiWali)

			// Update Settings (Hanya bisa diakses setelah login)
			protected.POST("/settings", handlers.UpdateSettings)

			protected.GET("/user/profile", func(c *gin.Context) {
				userID := c.MustGet("user_id")
				c.JSON(http.StatusOK, gin.H{
					"message": "Profil berhasil diakses",
					"user_id": userID,
				})
			})

			// Endpoint mengambil menu dinamis
			protected.GET("/user/menus", func(c *gin.Context) {
				roleID := int(c.MustGet("role_id").(float64))
				menus, err := handlers.GetMenusByRole(roleID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil menu"})
					return
				}
				c.JSON(http.StatusOK, menus)
			})
		}
	}

	// 2. Serve Frontend PWA dari Single Binary (go:embed)
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err == nil {
		fileServer := http.FileServer(http.FS(distFS))

		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			// Jika request API yang 404, tetap kembalikan JSON
			if strings.HasPrefix(path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint API tidak ditemukan"})
				return
			}

			// Cek apakah file fisik ada di static build
			f, err := distFS.Open(strings.TrimPrefix(path, "/"))
			if err == nil {
				f.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}

			// Jika bukan file fisik (misal rute SPA seperti /tabs/dashboard, /login), kembalikan index.html
			indexFile, err := distFS.Open("index.html")
			if err == nil {
				defer indexFile.Close()
				stat, _ := indexFile.Stat()
				modTime := time.Now()
				if stat != nil {
					modTime = stat.ModTime()
				}
				http.ServeContent(c.Writer, c.Request, "index.html", modTime, indexFile.(io.ReadSeeker))
				return
			}

			c.String(http.StatusNotFound, "File static tidak ditemukan")
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("SMARTSIS Backend & Frontend PWA menyala di port :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
