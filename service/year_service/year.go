package year_service

import (
	"github.com/vgbhj/SKAT/models"
)

// AddYear добавляет новый год в базу данных.
func AddYear(name string) error {
	return models.AddYear(name)
}

// GetYear получает год по его идентификатору.
func GetYear(id int) (*models.Year, error) {
	year, err := models.GetYear(id)
	if err != nil {
		return nil, err
	}
	return year, nil
}

// GetYearIDByName получает идентификатор года по его имени.
func GetYearIDByName(name string) (int, error) {
	yearID, err := models.GetYearIDByName(name)
	if err != nil {
		return 0, err
	}
	return yearID, nil
}
