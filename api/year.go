package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/SKAT/db"
	"github.com/vgbhj/SKAT/models"
	"github.com/vgbhj/SKAT/service/year_service"
)

// @Summary AddYear
// @Description Add a new year
// @Tags years
// @Accept json
// @Produce json
// @Param year body models.Year true "Year"
// @Success 200 {object} models.SuccessResponse "Year added successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/year [post]
// AddYear обрабатывает HTTP-запрос на добавление года
func AddYear(c *gin.Context) {
	var year models.Year

	if err := c.ShouldBindJSON(&year); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid input", Details: err.Error()})
		return
	}

	// Вызов функции сервиса для добавления года
	if err := year_service.AddYear(year.Name); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not save year to database", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Year added successfully"})
}

// @Summary GetYear
// @Description Retrieve a single year by its ID
// @Tags years
// @Accept json
// @Produce json
// @Param id path int true "Year ID"
// @Success 200 {object} models.Year "Year details"
// @Failure 400 {object} models.ErrorResponse "Invalid year ID"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve year"
// @Router /api/year/{id} [get]
func GetYear(c *gin.Context) {
	idStr := c.Param("id") // Получаем ID из параметров URL

	// Преобразование строки в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid year ID"})
		return
	}

	year, err := year_service.GetYear(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve year"})
		return
	}

	c.JSON(http.StatusOK, year)
}

// @Summary GetYearIDByName
// @Description Retrieve a year ID by its name
// @Tags years
// @Accept json
// @Produce json
// @Param name path string true "Year Name"
// @Success 200 {integer} int "Year ID"
// @Failure 400 {object} models.ErrorResponse "Invalid year name"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve year ID"
// @Router /api/year/name/{name} [get]
func GetYearIDByName(c *gin.Context) {
	name := c.Param("name") // Получаем имя года из параметров URL

	yearID, err := year_service.GetYearIDByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve year ID"})
		return
	}

	c.JSON(http.StatusOK, yearID)
}

// @Summary GetYears
// @Description Retrieve all years
// @Tags years
// @Accept json
// @Produce json
// @Success 200 {array} models.Year "List of years"
// @Failure 500 {object} models.ErrorResponse "Could not retrieve years"
// @Router /api/years [get]
func GetYears(c *gin.Context) {
	db := db.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM year")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not retrieve years", Details: err.Error()})
		return
	}
	defer rows.Close()

	var years []models.Year
	for rows.Next() {
		var year models.Year
		if err := rows.Scan(&year.ID, &year.Name); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Could not scan year", Details: err.Error()})
			return
		}
		years = append(years, year)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Error occurred during rows iteration", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, years)
}
