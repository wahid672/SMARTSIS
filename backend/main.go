package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"smartsis-backend/database"
	"smartsis-backend/handlers"
	"smartsis-backend/middleware"
)

func main() {
	// Inisialisasi Database SQLite
	database.InitDB()

	// Gin Engine
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// API Routing
	api := r.Group("/api")
	{
		// Public Auth
		api.POST("/auth/login", handlers.Login)

		// Protected Route
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/user/profile", func(c *gin.Context) {
				userID := c.MustGet("user_id")
				c.JSON(http.StatusOK, gin.H{
					"message": "Profil berhasil diakses",
					"user_id": userID,
				})
			})

			// Endpoint mengambil menu dinamis (bisa dipanggil kapan saja)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("SMARTSIS Backend menyala di port :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
