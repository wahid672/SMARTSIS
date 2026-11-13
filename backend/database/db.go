package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() *sql.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./siakad.db"
	}

	// SQLite dengan mode WAL & timeout agar aman dan cepat untuk konkurensi
	dsn := "file:" + dbPath + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)"

	var err error
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka database SQLite: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Gagal ping database: %v", err)
	}

	log.Printf("Database SQLite terkoneksi: %s", dbPath)
	MigrateSchema()
	SeedInitialData()

	return DB
}
