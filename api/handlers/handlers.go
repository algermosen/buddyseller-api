package handlers

import "github.com/gin-gonic/gin"

type UserHandler interface {
	GetUsers(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type ProductHandler interface {
	GetProducts(ctx *gin.Context)
	GetProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	CreateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
}

type OrderHandler interface {
	PlaceOrder(ctx *gin.Context)
	GetOrders(ctx *gin.Context)
	GetOrderById(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
	UpdateStatus(ctx *gin.Context)
}

type SessionHandler interface {
	Login(ctx *gin.Context)
}
