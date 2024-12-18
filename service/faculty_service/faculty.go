package faculty_service

import (
	"github.com/vgbhj/SKAT/models"
)

// AddFaculty добавляет новый факультет
func AddFaculty(name string, universityID int) error {
	return models.AddFaculty(name, universityID)
}

// GetFaculty получает факультет по его ID
func GetFaculty(id int) (*models.Faculty, error) {
	faculty, err := models.GetFaculty(id)
	if err != nil {
		return nil, err
	}
	return faculty, nil
}

// GetFacultyIDByName получает ID факультета по его имени
func GetFacultyIDByName(name string) (int, error) {
	facultyID, err := models.GetFacultyIDByName(name)
	if err != nil {
		return 0, err
	}
	return facultyID, nil
}
