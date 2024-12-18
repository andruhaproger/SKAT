package material_service

import (
	"github.com/vgbhj/SKAT/models"
)

// AddMaterial добавляет материал в базу данных
func AddMaterial(material models.Material) error {
	mat := map[string]interface{}{
		"name":          material.Name,
		"desc":          material.Desc,
		"file_url":      material.FileURL,
		"user_id":       material.UserID,
		"faculty_id":    material.FacultyID,
		"subject_id":    material.SubjectID,
		"year_id":       material.YearID,
		"university_id": material.UniversityID,
	}
	return models.AddMaterial(mat)
}

// GetMaterial получает материал из базы данных по ID
func GetMaterial(id int) (*models.Material, error) {
	material, err := models.GetMaterial(id)
	if err != nil {
		return nil, err
	}
	return material, nil
}
