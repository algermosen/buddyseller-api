package handlers

import (
	"example/buddyseller-api/db/datastore"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostgresProductHandler struct {
	ds *datastore.Queries
}

func (handler *PostgresProductHandler) GetProducts(ctx *gin.Context) {
	products, err := handler.ds.ListProducts(ctx)
	if err != nil {
		operationErr := &operationError{Entity: "products", Operation: OperationGet, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": operationErr.Msg()})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (handler *PostgresProductHandler) GetProduct(ctx *gin.Context) {
	identifier := ctx.Param("identifier")
	var (
		id      int64
		sku     string
		product datastore.Product
		err     error
	)

	id, err = strconv.ParseInt(identifier, 10, 32)
	if err != nil {
		sku = identifier
		log.Println("Getting product with sku...")
		product, err = handler.ds.GetProductBySku(ctx, sku)
	} else {
		log.Println("Getting product with id...")
		product, err = handler.ds.GetProductById(ctx, int32(id))
	}

	if err != nil {
		operationErr := &operationError{Entity: "product", Operation: OperationGet, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (handler *PostgresProductHandler) UpdateProduct(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	var productParams datastore.UpdateProductParams
	err = ctx.ShouldBindJSON(&productParams)

	if err != nil {
		jsonErr := &jsonBindingError{Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": jsonErr.Msg(),
		})
		log.Println(jsonErr.Error())
		return
	}

	productParams.ID = int32(productId)
	err = handler.ds.UpdateProduct(ctx, productParams)

	if err != nil {
		operationErr := &operationError{Entity: "product", Operation: OperationUpdate, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product updated succesfully.",
		"data": map[string]any{
			"ID":          productId,
			"Name":        productParams.Name,
			"Description": productParams.Description,
			"Sku":         productParams.Sku,
			"Price":       productParams.Price,
			"Stock":       productParams.Stock,
		},
	})
}

func (handler *PostgresProductHandler) CreateProduct(ctx *gin.Context) {
	var productParams datastore.CreateProductParams
	err := ctx.ShouldBindJSON(&productParams)

	if err != nil {
		jsonErr := &jsonBindingError{Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": jsonErr.Msg(),
		})
		log.Println(jsonErr.Error())
		return
	}

	pk, err := handler.ds.CreateProduct(ctx, productParams)

	if err != nil {
		operationErr := &operationError{Entity: "product", Operation: OperationSave, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Product created succesfully.",
		"data": map[string]any{
			"ID":          pk,
			"Name":        productParams.Name,
			"Description": productParams.Description,
			"Sku":         productParams.Sku,
			"Price":       productParams.Price,
			"Stock":       productParams.Stock,
		},
	})
}

func (handler *PostgresProductHandler) DeleteProduct(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	err = handler.ds.DeleteProduct(ctx, int32(productId))

	if err != nil {
		operationErr := &operationError{Entity: "product", Operation: OperationGet, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}
