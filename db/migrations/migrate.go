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
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

}
