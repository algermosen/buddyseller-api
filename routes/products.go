package routes

import (
	"net/http"
	"strconv"

	"example/buddyseller-api/models"

	"github.com/gin-gonic/gin"
)

func getProducts(ctx *gin.Context) {
	products, err := models.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch products. Try again later", "error": err.Error()})
	}

	ctx.JSON(http.StatusOK, products)
}

func getProductById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Not recognized as a valid Id"})
		return
	}

	product, err := models.GetProductById(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func createProduct(ctx *gin.Context) {

	var product models.Product
	err := ctx.ShouldBindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not match all the fields",
		})

		return
	}

	product.Save()

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "GOOD JOB!",
		"data":    product,
	})
}

func updateProduct(ctx *gin.Context) {
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

func deleteProduct(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = models.DeleteProduct(productId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
