package university_service

import (
	"github.com/vgbhj/SKAT/models"
)

// AddUniversity добавляет новый университет
func AddUniversity(name string) error {
	return models.AddUniversity(name)
}

// GetUniversity получает университет по его ID
func GetUniversity(id int) (*models.University, error) {
	university, err := models.GetUniversity(id)
	if err != nil {
		return nil, err
	}
	return university, nil
}

// GetUniversityIDByName получает ID университета по его имени
func GetUniversityIDByName(name string) (int, error) {
	universityID, err := models.GetUniversityIDByName(name)
	if err != nil {
		return 0, err
	}
	return universityID, nil
}
