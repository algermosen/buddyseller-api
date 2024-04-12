package api

import (
	"example/buddyseller-api/api/handlers"
	"example/buddyseller-api/api/middleware"

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
		userGroup.POST("", handlers.UserHandler.CreateUser)
		userGroup.PATCH("/:id", handlers.UserHandler.UpdateUser)
		userGroup.PATCH("/password", handlers.UserHandler.UpdatePassword)
		userGroup.DELETE("/:id", handlers.UserHandler.DeleteUser)
	}

	orderGroup := r.Group("/orders")
	{
		orderGroup.Use(middleware.Authenticate)
		orderGroup.POST("", handlers.OrderHandler.PlaceOrder)
		orderGroup.GET("", handlers.OrderHandler.GetOrders)
		orderGroup.GET("/:id", handlers.OrderHandler.GetOrderById)
		orderGroup.DELETE("/:id", handlers.OrderHandler.CancelOrder)
		orderGroup.PATCH("/:id/:status", handlers.OrderHandler.UpdateStatus)
	}

	productGroup := r.Group("/products")
	{
		productGroup.GET("", handlers.ProductHandler.GetProducts)
		productGroup.GET("/:identifier", handlers.ProductHandler.GetProduct)
		productGroup.POST("", handlers.ProductHandler.CreateProduct)
		productGroup.PATCH("/:id", handlers.ProductHandler.UpdateProduct)
		productGroup.DELETE("/:id", handlers.ProductHandler.DeleteProduct)
	}

	r.POST("/login", handlers.SessionHandler.Login)

	return r
}
