package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vgbhj/SKAT/db"
	"github.com/vgbhj/SKAT/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthInput структура для входных данных аутентификации
type AuthInput struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// SuccessResponse структура для успешного ответа
type SuccessResponse struct {
	Message string `json:"message" example:"User created successfully"`
}

// ErrorResponse структура для ошибок
type ErrorResponse struct {
	Error   string `json:"error" example:"Description of the error"`
	Details string `json:"details,omitempty" example:"Optional detailed error message"`
}

// TokenResponse структура для успешного ответа с токеном
type TokenResponse struct {
	Token string `json:"token" example:"your_jwt_token_here"`
}

// @Summary CreateUser
// @Description Create a new user with a username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body AuthInput true "User details"
// @Success 200 {object} SuccessResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Username already exists or invalid input"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/signup [post]
func CreateUser(c *gin.Context) {
	var authInput AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input", Details: err.Error()})
		return
	}

	db := db.ConnectDB()
	defer db.Close()

	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", authInput.Username).Scan(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Username already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Password hashing failed", Details: err.Error()})
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password, access_level) VALUES ($1, $2, 1)", authInput.Username, string(passwordHash))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not create user", Details: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Message: "User created successfully"})
}

// @Summary Login
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body AuthInput true "User credentials"
// @Success 200 {object} TokenResponse "JWT token"
// @Failure 400 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/login [post]
func Login(c *gin.Context) {
	var authInput AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input", Details: err.Error()})
		return
	}

	dbql := db.ConnectDB()
	defer dbql.Close()

	var userFound models.User
	err := dbql.QueryRow("SELECT id, password FROM users WHERE username = $1", authInput.Username).Scan(&userFound.ID, &userFound.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found"})
		return
	}

	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid password"})
		return
	}

	// Generate JWT token
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token", Details: err.Error()})
		return
	}

	err = db.RedisClient.Set(db.Ctx, "auth_token:"+token, userFound.ID, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	// Используем TokenResponse для возврата токена
	c.JSON(http.StatusOK, TokenResponse{Token: token})
}
