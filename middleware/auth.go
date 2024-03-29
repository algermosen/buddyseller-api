package middleware

import (
	"net/http"

	"example/buddyseller-api/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	var err error
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized", "error": "Empty token"})
		return
	}

	userId, err := utils.VerifyToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized", "error": err.Error()})
		return
	}

	c.Set("userId", userId)
	c.Next()
}
