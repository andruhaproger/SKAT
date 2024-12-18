package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/SKAT/db"
	"github.com/vgbhj/SKAT/models"
	"github.com/vgbhj/SKAT/service/university_service"
)

// @Summary AddUniversity
// @Description Add a new university
// @Tags universities
// @Accept json
// @Produce json
// @Param university body models.University true "University"
// @Success 200 {object} models.SuccessResponse "University added successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/university [post]
func AddUniversity(c *gin.Context) {
	var university models.University

	if err := c.ShouldBindJSON(&university); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid input", Details: err.Error()})
		return
	}

	// Вызов функции сервиса для добавления университета
	if err := university_service.AddUniversity(university.Name); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not save university to database", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "University added successfully"})
}

// @Summary GetUniversity
// @Description Retrieve a single university by its ID
// @Tags universities
// @Accept json
// @Produce json
// @Param id path int true "University ID"
// @Success 200 {object} models.University "University details"
// @Failure 400 {object} models.ErrorResponse "Invalid university ID"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve university"
// @Router /api/university/{id} [get]
func GetUniversity(c *gin.Context) {
	idStr := c.Param("id") // Получаем ID из параметров URL

	// Преобразование строки в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid university ID"})
		return
	}

	university, err := university_service.GetUniversity(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve university"})
		return
	}

	c.JSON(http.StatusOK, university)
}

// @Summary GetUniversityIDByName
// @Description Retrieve university ID by its name
// @Tags universities
// @Accept json
// @Produce json
// @Param name path string true "University Name"
// @Success 200 {object} models.SuccessResponse "University ID"
// @Failure 400 {object} models.ErrorResponse "Invalid university name"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve university ID"
// @Router /api/university/name/{name} [get]
func GetUniversityIDByName(c *gin.Context) {
	name := c.Param("name") // Получаем имя университета из параметров URL

	universityID, err := university_service.GetUniversityIDByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve university ID"})
		return
	}

	c.JSON(http.StatusOK, strconv.Itoa(universityID))
}

// @Summary GetUniversities
// @Description Retrieve all universities
// @Tags universities
// @Accept json
// @Produce json
// @Success 200 {array} models.University "List of universities"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve universities"
// @Router /api/universities [get]
func GetUniversities(c *gin.Context) {
	db := db.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM university")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve universities", Details: err.Error()})
		return
	}
	defer rows.Close()

	var universities []models.University
	for rows.Next() {
		var university models.University
		if err := rows.Scan(&university.ID, &university.Name); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not scan university", Details: err.Error()})
			return
		}
		universities = append(universities, university)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Error occurred during rows iteration", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, universities)
}
