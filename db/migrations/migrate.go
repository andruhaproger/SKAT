package main

import (
	"github.com/vgbhj/SKAT/config"
	"github.com/vgbhj/SKAT/db"
	"github.com/vgbhj/SKAT/models"
)

func init() {
	config.LoadEnvs()
	db.ConnectDB() // from database.go

}

func main() {
	db.DB.AutoMigrate(&models.User{})
}
