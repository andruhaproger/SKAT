package models

import (
	"fmt"

	"github.com/vgbhj/SKAT/db"
)

type University struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null" example:"Harvard University"`
}

func AddUniversity(name string) error {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	_, err := dbConn.Exec("INSERT INTO university (name) VALUES ($1)", name)
	if err != nil {
		return err
	}
	return nil
}

func GetUniversity(id int) (*University, error) {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	var university University
	err := dbConn.QueryRow("SELECT id, name FROM university WHERE id = $1", id).Scan(&university.ID, &university.Name)
	if err != nil {
		fmt.Println("Error retrieving university:", err)
		return nil, err
	}
	return &university, nil
}

func GetUniversityIDByName(name string) (int, error) {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	var universityID int
	err := dbConn.QueryRow("SELECT id FROM university WHERE name = $1", name).Scan(&universityID)
	if err != nil {
		fmt.Println("Error retrieving university ID by name:", err)
		return 0, err
	}
	return universityID, nil
}
