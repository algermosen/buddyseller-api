package middleware

import (
	"log"
	"net/http"
	"time"

	"example/buddyseller-api/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	var err error
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Must provide a valid token."})
		return
	}

	claims, err := utils.VerifyToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not valid token."})
		log.Printf("error validating token: %v", err.Error())
		return
	}

	exp, err := claims.GetExpirationTime()

	if err != nil || exp.Before(time.Now()) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not valid token."})
		log.Printf("token expired '%v'", token)
		return
	}

	c.Set("userId", claims["user_id"])
	log.Printf("%v", c.Keys)
	c.Next()
}
