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

// @Summary AddMaterial
// @Description Upload a material and save its metadata
// @Tags materials
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param name formData string true "Material Name"
// @Param description formData string true "Material Description"
// @Param faculty_name formData string true "Faculty Name"
// @Param subject_name formData string true "Subject Name"
// @Param year_name formData string true "Year Name"
// @Param university_name formData string true "University Name"
// @Param file formData file true "Material File"
// @Success 200 {object} models.SuccessResponse "Material added successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/material [post]
func AddMaterial(c *gin.Context) {
	// получение пользователя с мидлваре
	userId, exists := c.Get("currentUserId")
	if !exists {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "User ID not found"})
		return
	}

	// Преобразование userId в int
	var userIdInt int
	switch v := userId.(type) {
	case int:
		userIdInt = v
	case string:
		var err error
		userIdInt, err = strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid User ID format"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "User ID is not of a valid type"})
		return
	}

	facultyId, err := models.GetFacultyIDByName(c.PostForm("faculty_name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "faculty not found", Details: err.Error()})
		return
	}

	subjectID, err := models.GetSubjectIDByName(c.PostForm("subject_name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "subject not found", Details: err.Error()})
		return
	}

	yearID, err := models.GetSubjectIDByName(c.PostForm("year_name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "year not found", Details: err.Error()})
		return
	}

	universityID, err := models.GetUniversityIDByName(c.PostForm("university_name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "university not found", Details: err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "File is required"})
		return
	}

	data, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not open file"})
		return
	}
	defer data.Close()

	bucketName := "mybucket"
	objectName := file.Filename

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			fmt.Printf("Bucket %s already exists.\n", bucketName)
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create bucket", Details: err.Error()})
			return
		}
	}

	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, data, file.Size, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to upload file to MinIO", Details: err.Error()})
		return
	}

	fileURL := fmt.Sprintf("%s/%s", bucketName, objectName)
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
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not save material to database", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Material added successfully"})
}

// @Summary GetMaterial
// @Description Retrieve material details and file
// @Tags materials
// @Accept json
// @Produce octet-stream
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Material ID"
// @Success 200 {file} string "Material File"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/material/{id} [get]
func GetMaterial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid material ID"})
		return
	}

	material, err := material_service.GetMaterial(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve material"})
		return
	}

	object, err := minioClient.GetObject(context.Background(), "mybucket", material.FileURL, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve file from MinIO", Details: err.Error()})
		return
	}
	defer object.Close()

	fileData, err := io.ReadAll(object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not read file from MinIO", Details: err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", material.FileURL))
	c.Data(http.StatusOK, "application/octet-stream", fileData)
}
