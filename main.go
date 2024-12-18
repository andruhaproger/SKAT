package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vgbhj/SKAT/api"
	"github.com/vgbhj/SKAT/config"
	"github.com/vgbhj/SKAT/db"
	_ "github.com/vgbhj/SKAT/docs"
)

func init() {
	config.LoadEnvs()
	db.ConnectDB()

}

func main() {
	router := gin.Default()

	// Настройка маршрута для Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/auth/signup", api.CreateUser)
	router.POST("/auth/login", api.Login)
	router.POST("/material", api.AddMaterial)
	router.GET("/material/:id", api.GetMaterial)
	// router.GET("/user/profile", middlewares.CheckAuth, controllers.GetUserProfile)
	router.Run()
}
