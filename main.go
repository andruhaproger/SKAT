package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vgbhj/SKAT/api"
	"github.com/vgbhj/SKAT/config"
	"github.com/vgbhj/SKAT/db"
)

func init() {
	config.LoadEnvs()
	db.ConnectDB()

}

func main() {
	router := gin.Default()

	router.POST("/auth/signup", api.CreateUser)
	router.POST("/auth/login", api.Login)
	// router.GET("/user/profile", middlewares.CheckAuth, controllers.GetUserProfile)
	router.Run()
}
