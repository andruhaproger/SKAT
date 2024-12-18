package models

import (
	"fmt"

	"github.com/vgbhj/SKAT/db"
)

type Faculty struct {
	// Используяется только в GET
	ID           int    `json:"id" gorm:"primary_key"`
	Name         string `json:"name" gorm:"not null" example:"Engineering"`
	UniversityID int    `json:"university_id" gorm:"not null"`
}

// Faculty Operations
func AddFaculty(name string, universityID int) error {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	_, err := dbConn.Exec("INSERT INTO faculty (name, university_id) VALUES ($1, $2)", name, universityID)
	if err != nil {
		return err
	}
	return nil
}

func GetFaculty(id int) (*Faculty, error) {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	var faculty Faculty
	err := dbConn.QueryRow("SELECT id, name, university_id FROM faculty WHERE id = $1", id).Scan(&faculty.ID, &faculty.Name, &faculty.UniversityID)
	if err != nil {
		fmt.Println("Error retrieving faculty:", err)
		return nil, err
	}
	return &faculty, nil
}

// GetFacultyIDByName получает ID факультета по его имени
func GetFacultyIDByName(name string) (int, error) {
	dbConn := db.ConnectDB() // Получаем соединение с базой данных
	defer dbConn.Close()

	var facultyID int

	err := dbConn.QueryRow("SELECT id FROM faculty WHERE name = $1", name).Scan(&facultyID)
	if err != nil {
		// Выводим ошибку в лог
		fmt.Println("Error retrieving faculty ID by name:", err)

		// Возвращаем ошибку
		return 0, err
	}

	return facultyID, nil
}
