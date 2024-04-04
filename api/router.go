package api

import (
	"example/buddyseller-api/api/handlers"

	"github.com/gin-gonic/gin"
)

func RouterSetup(
	userHandler handlers.UserHandler,
	productHandler handlers.ProductHandler,
	sessionHandler handlers.SessionHandler,
	orderHandler handlers.OrderHandler,
) *gin.Engine {
	r := gin.Default()

	userGroup := r.Group("/users")
	{
		userGroup.GET("", userHandler.GetUsers)
		userGroup.GET("/:id", userHandler.GetUserById)
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.PATCH("/:id", userHandler.UpdateUser)
		userGroup.PATCH("/password", userHandler.UpdatePassword)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
	}

	orderGroup := r.Group("/orders")
	{
		orderGroup.POST("/orders", orderHandler.PlaceOrder)
		orderGroup.GET("/orders", orderHandler.GetOrders)
		orderGroup.GET("/orders/:id", orderHandler.GetOrderById)
		orderGroup.DELETE("/orders/:id", orderHandler.CancelOrder)
		orderGroup.PATCH("/orders/:id/:status", orderHandler.UpdateStatus)
	}

	productGroup := r.Group("")
	{
		productGroup.GET("/products", productHandler.GetProducts)
		productGroup.GET("/products/:identifier", productHandler.GetProduct)
		productGroup.POST("/products", productHandler.CreateProduct)
		productGroup.PATCH("/products/:id", productHandler.UpdateProduct)
		productGroup.DELETE("/products/:id", productHandler.DeleteProduct)
	}

	r.POST("/login", sessionHandler.Login)

	return r
}
