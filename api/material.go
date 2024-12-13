package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/SKAT/models"
	"github.com/vgbhj/SKAT/service/material_service"
)

// AddMaterialHandler обрабатывает HTTP-запрос на добавление материала
func AddMaterial(c *gin.Context) {
	var material models.Material
	material.Name = c.PostForm("name")
	material.Desc = c.PostForm("description")

	idStr := c.PostForm("id")
	userIDStr := c.PostForm("user_id")

	// Преобразование строки в int
	var err error
	if material.ID, err = strconv.Atoi(idStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if material.UserID, err = strconv.Atoi(userIDStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	// Получение остальных полей
	material.Name = c.PostForm("name")
	material.Desc = c.PostForm("description")

	// Чтение файла из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Открытие файла
	data, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer data.Close()

	// Чтение содержимого файла
	fileBytes, err := ioutil.ReadAll(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
		return
	}

	// Вызов функции сервиса для добавления материала
	if err := material_service.AddMaterial(material, fileBytes); err != nil {
		// Выводим текст ошибки в лог (опционально)
		fmt.Println("Error adding material:", err)

		// Возвращаем ошибку в ответе
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save material to database", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material added successfully"})
}

// GetMaterialHandler обрабатывает HTTP-запрос на получение материала
func GetMaterial(c *gin.Context) {
	idStr := c.Param("id") // Получаем ID из параметров URL

	// Преобразуем строку в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	material, err := material_service.GetMaterial(id) // Вызов функции сервиса для получения материала
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve material"})
		return
	}

	c.JSON(http.StatusOK, material)
}
