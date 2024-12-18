package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/SKAT/db"
	"github.com/vgbhj/SKAT/models"
	"github.com/vgbhj/SKAT/service/subject_service"
)

// @Summary AddSubject
// @Description Add a new subject
// @Tags subjects
// @Accept json
// @Produce json
// @Param subject body models.Subject true "Subject"
// @Success 200 {object} models.SuccessResponse "Subject added successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/subject [post]
func AddSubject(c *gin.Context) {
	var subject models.Subject

	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid input", Details: err.Error()})
		return
	}

	// Вызов функции сервиса для добавления предмета
	if err := subject_service.AddSubject(subject.Name); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not save subject to database", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Subject added successfully"})
}

// @Summary GetSubject
// @Description Retrieve a single subject by its ID
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Success 200 {object} models.Subject "Subject details"
// @Failure 400 {object} models.ErrorResponse "Invalid subject ID"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve subject"
// @Router /api/subject/{id} [get]
func GetSubject(c *gin.Context) {
	idStr := c.Param("id") // Получаем ID из параметров URL

	// Преобразование строки в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid subject ID"})
		return
	}

	subject, err := subject_service.GetSubject(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve subject"})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// @Summary GetSubjectIDByName
// @Description Retrieve a subject ID by its name
// @Tags subjects
// @Accept json
// @Produce json
// @Param name path string true "Subject Name"
// @Success 200 {object} models.SuccessResponse "Subject ID"
// @Failure 400 {object} models.ErrorResponse "Invalid subject name"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve subject ID"
// @Router /api/subject/name/{name} [get]
func GetSubjectIDByName(c *gin.Context) {
	name := c.Param("name") // Получаем имя предмета из параметров URL

	subjectID, err := subject_service.GetSubjectIDByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve subject ID"})
		return
	}

	c.JSON(http.StatusOK, strconv.Itoa(subjectID))
}

// @Summary GetSubjects
// @Description Retrieve all subjects
// @Tags subjects
// @Accept json
// @Produce json
// @Success 200 {array} models.Subject "List of subjects"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve subjects"
// @Router /api/subjects [get]
func GetSubjects(c *gin.Context) {
	db := db.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM subject")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve subjects", Details: err.Error()})
		return
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not scan subject", Details: err.Error()})
			return
		}
		subjects = append(subjects, subject)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Error occurred during rows iteration", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, subjects)
}
