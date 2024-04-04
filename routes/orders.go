package routes

import (
	"example/buddyseller-api/dtos"
	"example/buddyseller-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getOrders(ctx *gin.Context) {
	orders, err := models.GetAllOrders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch orders. Try again later", "error": err.Error()})
	}

	ctx.JSON(http.StatusOK, orders)
}

func getOrderById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Not recognized as a valid Id"})
		return
	}

	order, err := models.GetOrderById(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func placeOrder(ctx *gin.Context) {

	var order dtos.NewOrderDto
	err := ctx.ShouldBindJSON(&order)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not match all the fields",
		})

		return
	}

	err = models.PlaceOrder(order)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error Saving product",
			"details": err.Error(),
			"caller":  "product.Save()",
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "GOOD JOB!",
		"data":    order,
	})
}

func updateStatus(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var updatedProduct models.Product
	err = ctx.ShouldBindJSON(&updatedProduct)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedProduct.ID = productId
	err = updatedProduct.Update()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
}

func cancelOrder(ctx *gin.Context) {
	orderId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = models.CancelOrder(orderId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
