package subject_service

import (
	"github.com/vgbhj/SKAT/models"
)

// AddSubject добавляет новый предмет
func AddSubject(name string) error {
	return models.AddSubject(name)
}

// GetSubject получает предмет по его ID
func GetSubject(id int) (*models.Subject, error) {
	subject, err := models.GetSubject(id)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

// GetSubjectIDByName получает ID предмета по его имени
func GetSubjectIDByName(name string) (int, error) {
	subjectID, err := models.GetSubjectIDByName(name)
	if err != nil {
		return 0, err
	}
	return subjectID, nil
}
