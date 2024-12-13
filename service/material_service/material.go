package material_service

import (
	"github.com/vgbhj/SKAT/models"
)

// AddMaterial добавляет материал в базу данных
func AddMaterial(material models.Material, fileBytes []byte) error {
	mat := map[string]interface{}{
		"id":      material.ID,
		"name":    material.Name,
		"desc":    material.Desc,
		"file":    fileBytes,
		"user_id": material.UserID,
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
