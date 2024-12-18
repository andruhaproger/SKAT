package models

import (
	"fmt"

	"github.com/vgbhj/SKAT/db"
)

type Material struct {
	// Только для гет запроса
	ID           int    `json:"id" gorm:"primary_key"`
	Name         string `json:"name" gorm:"not null"`
	Desc         string `json:"description"`
	FileURL      string `json:"file_url" gorm:"not null"`
	UserID       int    `json:"user_id"`
	FacultyID    *int   `json:"faculty_id"`
	SubjectID    *int   `json:"subject_id"`
	YearID       *int   `json:"year_id"`
	UniversityID *int   `json:"university_id"`
	UploadDate   string `json:"upload_date"`
}

func AddMaterial(data map[string]interface{}) error {
	db := db.ConnectDB()
	defer db.Close()

	material := Material{
		Name:         data["name"].(string),
		Desc:         data["desc"].(string),
		FileURL:      data["file_url"].(string),
		UserID:       (data["user_id"].(int)),
		FacultyID:    data["faculty_id"].(*int),
		SubjectID:    data["subject_id"].(*int),
		YearID:       data["year_id"].(*int),
		UniversityID: data["university_id"].(*int),
	}

	_, err := db.Exec(`
		INSERT INTO material (name, description, file_url, user_id, faculty_id, subject_id, year_id, university_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		material.Name, material.Desc, material.FileURL, material.UserID, material.FacultyID, material.SubjectID, material.YearID, material.UniversityID)
	if err != nil {
		return err
	}

	return nil
}

func GetMaterial(id int) (*Material, error) {
	db := db.ConnectDB()
	defer db.Close()

	var material Material

	err := db.QueryRow(`
		SELECT id, name, description, file_url, user_id, faculty_id, subject_id, year_id, university_id, upload_date 
		FROM material WHERE id = $1`, id).Scan(
		&material.ID, &material.Name, &material.Desc, &material.FileURL, &material.UserID,
		&material.FacultyID, &material.SubjectID, &material.YearID, &material.UniversityID,
		&material.UploadDate)

	if err != nil {
		// Выводим ошибку в лог
		fmt.Println("Error retrieving material:", err)

		// Возвращаем ошибку
		return nil, err
	}

	return &material, nil
}
