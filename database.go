package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func initDatabase() error {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/ka_ping_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return err
	}

	// Test database connection
	if err = db.Ping(); err != nil {
		return err
	}

	// Create tables if they don't exist
	createTables()
	return nil
}

func createTables() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS devices (
		id INT AUTO_INCREMENT PRIMARY KEY,
		uuid VARCHAR(36) UNIQUE NOT NULL,
		hostname VARCHAR(255) NOT NULL,
		os VARCHAR(100) NOT NULL,
		mac VARCHAR(17) NOT NULL,
		public_ip VARCHAR(45) NOT NULL,
		country VARCHAR(100),
		region VARCHAR(100),
		city VARCHAR(100),
		latitude VARCHAR(20),
		longitude VARCHAR(20),
		last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_uuid (uuid),
		INDEX idx_last_seen (last_seen)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
	log.Println("Database tables created/verified successfully")
}

func closeDatabase() {
	if db != nil {
		db.Close()
	}
}
