package routes

import (
	"example/buddyseller-api/dtos"
	"example/buddyseller-api/models"
	"example/buddyseller-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	var userDto dtos.UserLoginDto
	err := c.ShouldBindJSON(&userDto)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid object", "error": err.Error()})
		return
	}

	user, err := models.ValidateUserCredentials(&userDto)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Could not validate user", "error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Login successful", "token": token})
}
