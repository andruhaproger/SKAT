package main

import (
	"log"

	"github.com/vgbhj/SKAT/config"
	"github.com/vgbhj/SKAT/db"
)

func init() {
	config.LoadEnvs()
	db.ConnectDB() // from database.go
}

func main() {
	db := db.ConnectDB()
	// close the db connection
	defer db.Close()

	// Миграция для таблицы users
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`)
	if err != nil {
		log.Fatal("Migration for users failed:", err)
	}

	// Миграция для таблицы material
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS material (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		file BYTEA NOT NULL,
		user_id INTEGER,
		upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
	);
	`)
	if err != nil {
		log.Fatal("Migration for material failed:", err)
	}

	log.Println("Migrations completed successfully.")
}
