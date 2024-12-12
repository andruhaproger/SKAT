package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/SKAT/db"
	"github.com/vgbhj/SKAT/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	db.DB.Where("username=?", authInput.Username).Find(&userFound)

	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: authInput.Username,
		Password: string(passwordHash),
	}

	db.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}
