package api

import (
	"example/buddyseller-api/api/handlers"

	"github.com/gin-gonic/gin"
)

type RouterHandlers struct {
	UserHandler    handlers.UserHandler
	ProductHandler handlers.ProductHandler
	SessionHandler handlers.SessionHandler
	OrderHandler   handlers.OrderHandler
}

func RouterSetup(handlers *RouterHandlers) *gin.Engine {
	r := gin.Default()

	userGroup := r.Group("/users")
	{
		userGroup.GET("", handlers.UserHandler.GetUsers)
		userGroup.GET("/:id", handlers.UserHandler.GetUserById)
		userGroup.POST("/", handlers.UserHandler.CreateUser)
		userGroup.PATCH("/:id", handlers.UserHandler.UpdateUser)
		userGroup.PATCH("/password", handlers.UserHandler.UpdatePassword)
		userGroup.DELETE("/:id", handlers.UserHandler.DeleteUser)
	}

	orderGroup := r.Group("/orders")
	{
		orderGroup.POST("/orders", handlers.OrderHandler.PlaceOrder)
		orderGroup.GET("/orders", handlers.OrderHandler.GetOrders)
		orderGroup.GET("/orders/:id", handlers.OrderHandler.GetOrderById)
		orderGroup.DELETE("/orders/:id", handlers.OrderHandler.CancelOrder)
		orderGroup.PATCH("/orders/:id/:status", handlers.OrderHandler.UpdateStatus)
	}

	productGroup := r.Group("")
	{
		productGroup.GET("/products", handlers.ProductHandler.GetProducts)
		productGroup.GET("/products/:identifier", handlers.ProductHandler.GetProduct)
		productGroup.POST("/products", handlers.ProductHandler.CreateProduct)
		productGroup.PATCH("/products/:id", handlers.ProductHandler.UpdateProduct)
		productGroup.DELETE("/products/:id", handlers.ProductHandler.DeleteProduct)
	}

	r.POST("/login", handlers.SessionHandler.Login)

	return r
}
