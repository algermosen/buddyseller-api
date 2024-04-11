package handlers

import (
	"errors"
	"example/buddyseller-api/db/datastore"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type PostgresProductHandler struct {
	DS *datastore.Queries
}

func (h *PostgresProductHandler) GetProducts(ctx *gin.Context) {
	products, err := h.DS.ListProducts(ctx)

	if err != nil {
		operationErr := &operationError{Entity: "products", Operation: OperationGet, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": operationErr.Msg()})
		log.Println(operationErr.Error())
		return
	}

	if products == nil {
		ctx.JSON(http.StatusOK, gin.H{"data": make([]interface{}, 0)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *PostgresProductHandler) GetProduct(ctx *gin.Context) {
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
		product, err = h.DS.GetProductBySku(ctx, sku)
	} else {
		log.Println("Getting product with id...")
		product, err = h.DS.GetProductById(ctx, int32(id))
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("No product found with identifier '%v'.", identifier),
			})
			log.Println(err.Error())
			return
		}

		operationErr := &operationError{Entity: "product", Operation: OperationGet, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (h *PostgresProductHandler) UpdateProduct(ctx *gin.Context) {
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
	rowsAffected, err := h.DS.UpdateProduct(ctx, productParams)

	if err != nil {
		operationErr := &operationError{Entity: "product", Operation: OperationUpdate, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	if rowsAffected == 0 {
		msg := fmt.Sprintf("No product found with identifier '%v'.", productId)
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": msg,
		})
		log.Println(msg)
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

func (h *PostgresProductHandler) CreateProduct(ctx *gin.Context) {
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

	pk, err := h.DS.CreateProduct(ctx, productParams)

	if err != nil {
		operationErr := &operationError{Entity: "product", Operation: OperationSave, origin: err}
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

func (h *PostgresProductHandler) DeleteProduct(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	rowsAffected, err := h.DS.DeleteProduct(ctx, int32(productId))

	if err != nil {
		operationErr := &operationError{Entity: "product", Operation: OperationGet, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	if rowsAffected == 0 {
		msg := fmt.Sprintf("No product found with identifier '%v'.", productId)
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": msg,
		})
		log.Println(msg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product deleted succesfully.",
	})
}
