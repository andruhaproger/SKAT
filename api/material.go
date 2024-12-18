package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vgbhj/SKAT/models"
	"github.com/vgbhj/SKAT/service/material_service"
)

var minioClient *minio.Client

func init() {
	var err error
	minioClient, err = minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("youraccesskey", "yoursecretkey", ""),
		Secure: false,
	})
	if err != nil {
		fmt.Println("Error initializing MinIO client:", err)
	}
}

// AddMaterialHandler обрабатывает HTTP-запрос на добавление материала
func AddMaterial(c *gin.Context) {
	// получение пользователя с мидлваре
	userId, exists := c.Get("currentUserId")
	if !exists {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "User ID not found"})
		return
	}

	// Преобразование userId в int
	userIdInt, ok := userId.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "User ID is not of type int"})
		return
	}

	// Получаем ID фака по его имени
	facultyId, err := models.GetFacultyIDByName(c.PostForm("faculty_name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "faculty not found", Details: err.Error()})
		return
	}

	// Получаем ID предмета для события
	subjectID, err := models.GetSubjectIDByName(c.PostForm("subject_name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "subject not found", Details: err.Error()})
		return
	}

	// Получаем ID года для события
	yearID, err := models.GetSubjectIDByName(c.PostForm("year_name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "year not found", Details: err.Error()})
		return
	}

	// Получаем ID года для события
	universityID, err := models.GetUniversityIDByName(c.PostForm("university_name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "university not found", Details: err.Error()})
		return
	}

	// Чтение файла из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	data, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer data.Close()

	bucketName := "mybucket"
	objectName := file.Filename

	// Создание бакета, если он не существует
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			fmt.Printf("Bucket %s already exists.\n", bucketName)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bucket: " + err.Error()})
			return
		}
	}

	// Загрузка файла в MinIO
	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, data, file.Size, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to MinIO: " + err.Error()})
		return
	}

	// Сохранение метаданных в базе данных
	fileURL := fmt.Sprintf("%s/%s", bucketName, objectName) // Путь к файлу в MinIO
	material := models.Material{
		Name:         c.PostForm("name"),
		Desc:         c.PostForm("description"),
		FileURL:      fileURL,
		UserID:       userIdInt,
		FacultyID:    &facultyId,
		SubjectID:    &subjectID,
		YearID:       &yearID,
		UniversityID: &universityID,
	}

	if err := material_service.AddMaterial(material); err != nil {
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

	// Получаем файл из MinIO
	object, err := minioClient.GetObject(context.Background(), "mybucket", material.FileURL, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file from MinIO", "details": err.Error()})
		return
	}
	defer object.Close()

	// Читаем содержимое файла
	fileData, err := io.ReadAll(object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file from MinIO", "details": err.Error()})
		return
	}

	// Возвращаем файл в ответе
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", material.FileURL))
	c.Data(http.StatusOK, "application/octet-stream", fileData)
}
