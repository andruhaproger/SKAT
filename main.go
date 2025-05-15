package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vgbhj/SKAT/api"
	"github.com/vgbhj/SKAT/config"
	"github.com/vgbhj/SKAT/db"
	_ "github.com/vgbhj/SKAT/docs"
	"github.com/vgbhj/SKAT/middleware"
)

func init() {
	config.LoadEnvs()
	db.ConnectDB()
	db.ConnectRedis()
}

func main() {
	router := gin.Default()

	// Настройка маршрута для Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/api/signup", api.CreateUser)
	router.POST("/api/login", api.Login)

	router.POST("/api/material", middleware.CheckAuth, api.AddMaterial)
	router.GET("/api/material/:id", api.GetMaterial)
	router.GET("/api/materials", api.GetMaterials)

	router.POST("/api/faculty", api.AddFaculty)
	router.GET("/api/faculty/:id", api.GetFaculty)
	router.GET("/api/faculty/name/:name", api.GetFacultyIDByName)
	router.GET("/api/faculties", api.GetFaculties)

	router.POST("/api/university", api.AddUniversity)
	router.GET("/api/university/:id", api.GetUniversity)
	router.GET("/api/university/name/:name", api.GetUniversityIDByName)
	router.GET("/api/universities", api.GetUniversities)

	router.POST("/api/subject", api.AddSubject)
	router.GET("/api/subject/:id", api.GetSubject)
	router.GET("/api/subject/name/:name", api.GetFacultyIDByName)
	router.GET("/api/subjects", api.GetSubjects)

	router.POST("/api/year", api.AddYear)
	router.GET("/api/year/:id", api.GetYear)
	router.GET("/api/year/name/:name", api.GetYearIDByName)
	router.GET("/api/years", api.GetYears)
	router.Run()
}
