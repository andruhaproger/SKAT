package models

import (
	"fmt"

	"github.com/vgbhj/SKAT/db"
)

type Material struct {
	ID     int    `json:"id" gorm:"primary_key"`
	Name   string `json:"name" gorm:"not null"`
	Desc   string `json:"description"`
	File   []byte `json:"-"`
	UserID int    `json:"user_id"`
}

func AddMaterial(data map[string]interface{}) error {
	db := db.ConnectDB()
	defer db.Close()

	material := Material{
		ID:     (data["id"].(int)),
		Name:   data["name"].(string),
		Desc:   data["desc"].(string),
		File:   data["file"].([]byte),
		UserID: (data["user_id"].(int)),
	}
	_, err := db.Exec("INSERT INTO material (name, description, file, user_id) VALUES ($1, $2, $3, $4)",
		material.Name, material.Desc, material.File, material.UserID)
	if err != nil {
		return err
	}

	return nil
}

func GetMaterial(id int) (*Material, error) {
	db := db.ConnectDB()
	defer db.Close()

	var material Material

	err := db.QueryRow("SELECT id, name, description, file, user_id FROM material WHERE id = $1", id).Scan(
		&material.ID, &material.Name, &material.Desc, &material.File, &material.UserID)

	if err != nil {
		// Выводим ошибку в лог
		fmt.Println("Error retrieving material:", err)

		// Возвращаем ошибку
		return nil, err
	}

	return &material, nil
}
