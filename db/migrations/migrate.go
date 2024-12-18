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

	// Удаление таблиц в обратном порядке их создания
	tables := []string{
		"material",
		"favorite",
		"log",
		"users",
		"access",
		"year",
		"subject",
		"faculty",
		"university",
	}

	for _, table := range tables {
		_, err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE;")
		if err != nil {
			log.Fatalf("Failed to drop table %s: %v", table, err)
		}
		log.Printf("Table %s dropped successfully.", table)
	}

	log.Println("All tables dropped successfully.")

	// Миграция для таблицы university
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS university (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	`)
	if err != nil {
		log.Fatal("Migration for university failed:", err)
	}

	// Добавление одного значения в таблицу university
	_, err = db.Exec(`
	INSERT INTO university (name) 
	VALUES ('Example University')
	ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		log.Fatal("Inserting university failed:", err)
	}

	// Миграция для таблицы faculty
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS faculty (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		university_id INTEGER REFERENCES university(id)
	);
	`)
	if err != nil {
		log.Fatal("Migration for faculty failed:", err)
	}

	// Добавление одного значения в таблицу faculty
	_, err = db.Exec(`
	INSERT INTO faculty (name, university_id) 
	VALUES ('Example Faculty', 1)
	ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		log.Fatal("Inserting faculty failed:", err)
	}

	// Миграция для таблицы subject
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS subject (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255)
	);
	`)
	if err != nil {
		log.Fatal("Migration for subject failed:", err)
	}

	// Добавление одного значения в таблицу subject
	_, err = db.Exec(`
	INSERT INTO subject (name) 
	VALUES ('Example Subject')
	ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		log.Fatal("Inserting subject failed:", err)
	}

	// Миграция для таблицы year
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS year (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	`)
	if err != nil {
		log.Fatal("Migration for year failed:", err)
	}

	// Добавление одного значения в таблицу year
	_, err = db.Exec(`
	INSERT INTO year (name) 
	VALUES ('2024')
	ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		log.Fatal("Inserting year failed:", err)
	}

	// Миграция для таблицы access
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS access (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT
	);
	`)
	if err != nil {
		log.Fatal("Migration for access failed:", err)
	}

	// Миграция для таблицы users
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		university_id INTEGER REFERENCES university(id),
		access INTEGER
	);
	`)
	if err != nil {
		log.Fatal("Migration for users failed:", err)
	}

	// Добавление пользователя admin
	_, err = db.Exec(`
	INSERT INTO users (username, password, first_name, last_name, university_id, access) 
	VALUES 
	('admin', '$2a$10$YgzepzPAE0OZWr9P6mQVu.Ind9xcSN/DGCfOiVT8XClxWjWLWbfpa', 'Admin', 'Admin', 1, 1)
	ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		log.Fatal("Inserting admin user failed:", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS material (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		file_url VARCHAR(255) NOT NULL, -- URL или путь к файлу в MinIO
		user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
		faculty_id INTEGER REFERENCES faculty(id),
		subject_id INTEGER REFERENCES subject(id),
		year_id INTEGER REFERENCES year(id),
		upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		rating FLOAT DEFAULT 0,
		university_id INTEGER REFERENCES university(id)
	);
	`)

	if err != nil {
		log.Fatal("Migration for material failed:", err)
	}

	// Миграция для таблицы favorite
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS favorite (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		material_id INTEGER REFERENCES material(id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		log.Fatal("Migration for favorite failed:", err)
	}

	// Миграция для таблицы log
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS log (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		action VARCHAR(255) NOT NULL,
		target_id INTEGER,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		log.Fatal("Migration for log failed:", err)
	}

	log.Println("Migrations completed successfully.")
}
