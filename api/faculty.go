package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/SKAT/db"
	"github.com/vgbhj/SKAT/models"
	"github.com/vgbhj/SKAT/service/faculty_service"
)

// @Summary AddFaculty
// @Description Add a new faculty
// @Tags faculties
// @Accept json
// @Produce json
// @Param faculty body models.Faculty true "Faculty"
// @Success 200 {object} models.SuccessResponse "Faculty added successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/faculty [post]
func AddFaculty(c *gin.Context) {
	var faculty models.Faculty

	if err := c.ShouldBindJSON(&faculty); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid input", Details: err.Error()})
		return
	}

	// Вызов функции сервиса для добавления факультета
	if err := faculty_service.AddFaculty(faculty.Name, faculty.UniversityID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not save faculty to database", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Faculty added successfully"})
}

// @Summary GetFaculty
// @Description Retrieve a single faculty by its ID
// @Tags faculties
// @Accept json
// @Produce json
// @Param id path int true "Faculty ID"
// @Success 200 {object} models.Faculty "Faculty details"
// @Failure 400 {object} models.ErrorResponse "Invalid faculty ID"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve faculty"
// @Router /api/faculty/{id} [get]
func GetFaculty(c *gin.Context) {
	idStr := c.Param("id") // Получаем ID из параметров URL

	// Преобразование строки в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid faculty ID"})
		return
	}

	faculty, err := faculty_service.GetFaculty(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve faculty"})
		return
	}

	c.JSON(http.StatusOK, faculty)
}

// @Summary GetFacultyIDByName
// @Description Retrieve faculty ID by its name
// @Tags faculties
// @Accept json
// @Produce json
// @Param name path string true "Faculty Name"
// @Success 200 {integer} int "Faculty ID"
// @Failure 400 {object} models.ErrorResponse "Invalid faculty name"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve faculty ID"
// @Router /api/faculty/name/{name} [get]
func GetFacultyIDByName(c *gin.Context) {
	name := c.Param("name") // Получаем имя факультета из параметров URL

	facultyID, err := faculty_service.GetFacultyIDByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve faculty ID"})
		return
	}

	c.JSON(http.StatusOK, facultyID)
}

// @Summary GetFaculties
// @Description Retrieve all faculties
// @Tags faculties
// @Accept json
// @Produce json
// @Success 200 {array} models.Faculty "List of faculties"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve faculties"
// @Router /api/faculties [get]
func GetFaculties(c *gin.Context) {
	db := db.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, university_id FROM faculty")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve faculties", Details: err.Error()})
		return
	}
	defer rows.Close()

	var faculties []models.Faculty
	for rows.Next() {
		var faculty models.Faculty
		if err := rows.Scan(&faculty.ID, &faculty.Name, &faculty.UniversityID); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not scan faculty", Details: err.Error()})
			return
		}
		faculties = append(faculties, faculty)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Error occurred during rows iteration", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, faculties)
}
