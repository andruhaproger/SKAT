package models

import (
	"fmt"

	"github.com/vgbhj/SKAT/db"
)

type Subject struct {
	// Используяется только в GET
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" example:"Mathematics"`
}

func AddSubject(name string) error {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	_, err := dbConn.Exec("INSERT INTO subject (name) VALUES ($1)", name)
	if err != nil {
		return err
	}
	return nil
}

func GetSubject(id int) (*Subject, error) {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	var subject Subject
	err := dbConn.QueryRow("SELECT id, name FROM subject WHERE id = $1", id).Scan(&subject.ID, &subject.Name)
	if err != nil {
		fmt.Println("Error retrieving subject:", err)
		return nil, err
	}
	return &subject, nil
}

func GetSubjectIDByName(name string) (int, error) {
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	var subjectID int
	err := dbConn.QueryRow("SELECT id FROM subject WHERE name = $1", name).Scan(&subjectID)
	if err != nil {
		fmt.Println("Error retrieving subject ID by name:", err)
		return 0, err
	}
	return subjectID, nil
}
