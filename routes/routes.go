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
	server.GET("/products/:id", getProductById)
	server.GET("/products/sku/:sku", getProductBySku)
	server.POST("/products", createProduct)
	server.PATCH("/products/:id", updateProduct)
	server.DELETE("/products/:id", deleteProduct)

	server.GET("/users", getUsers)
	server.GET("/users/:id", getUserById)
	server.POST("/users", createUser)
	server.PATCH("/users/:id", updateUser)
	server.PATCH("/users/password", updateUser) // TODO: Use the updatePassword method
	server.DELETE("/users/:id", deleteUser)

	server.POST("/login", login)

	server.POST("/orders", placeOrder)
	server.GET("/orders", getOrders)
	server.GET("/orders/:id", getOrderById)
	server.DELETE("/orders/:id", cancelOrder)
	server.PATCH("/orders/:id/:status", updateStatus)
}
