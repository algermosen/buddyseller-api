package routes

import (
	"example/buddyseller-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	privateRoutes := server.Group("/")
	{
		privateRoutes.Use(middleware.Authenticate)
	}

	server.GET("/products", getProducts)
}
