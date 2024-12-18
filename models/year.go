package models

import (
	"fmt"

	"github.com/vgbhj/SKAT/db"
)

type Year struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null" example:"2024"`
}

// Year Operations
func AddYear(name string) error {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	_, err := dbConn.Exec("INSERT INTO year (name) VALUES ($1)", name)
	if err != nil {
		return err
	}
	return nil
}

func GetYear(id int) (*Year, error) {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	var year Year
	err := dbConn.QueryRow("SELECT id, name FROM year WHERE id = $1", id).Scan(&year.ID, &year.Name)
	if err != nil {
		fmt.Println("Error retrieving year:", err)
		return nil, err
	}
	return &year, nil
}

func GetYearIDByName(name string) (int, error) {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	var yearID int
	err := dbConn.QueryRow("SELECT id FROM year WHERE name = $1", name).Scan(&yearID)
	if err != nil {
		fmt.Println("Error retrieving year ID by name:", err)
		return 0, err
	}
	return yearID, nil
}
