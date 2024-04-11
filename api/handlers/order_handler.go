package handlers

import (
	"context"
	"errors"
	"example/buddyseller-api/api/dtos"
	"example/buddyseller-api/db"
	"example/buddyseller-api/db/datastore"
	"example/buddyseller-api/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type PostgresOrderHandler struct {
	DS *datastore.Queries
}

func (h *PostgresOrderHandler) PlaceOrder(ctx *gin.Context) {
	var newOrder dtos.NewOrderDto
	userId := int32(ctx.Keys["userId"].(float64))

	if userId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User not valid.",
		})
		return
	}

	err := ctx.ShouldBindJSON(&newOrder)

	if err != nil {
		jsonErr := &jsonBindingError{Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": jsonErr.Msg(),
		})
		log.Println(jsonErr.Error())
		return
	}

	err = h.placeNewOrder(ctx, &newOrder, userId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error placing order.",
		})
		log.Println(err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "GOOD JOB!",
	})
}

func (h *PostgresOrderHandler) placeNewOrder(ctx context.Context, newOrder *dtos.NewOrderDto, userID int32) error {
	if len(newOrder.Items) == 0 {
		return fmt.Errorf("cannot create an order with no items")
	}

	productIDs := make([]int32, len(newOrder.Items))
	productQuantityIndex := make(map[int32]int32)

	// Extract ProductIDs from the array of OrderItemDto
	for i, orderItem := range newOrder.Items {
		productIDs[i] = orderItem.ProductID
		productQuantityIndex[orderItem.ProductID] = orderItem.Quantity
	}

	productRows, err := h.DS.ListProductsToOrder(ctx, productIDs)
	if err != nil {
		return fmt.Errorf("error listing product prices: \n%w", err)
	}

	var pRetrieved, pExpected = len(productRows), len(productIDs)
	if pRetrieved != pExpected {
		log.Printf("products tried to list: %v\n", productIDs)
		return fmt.Errorf("not all products could be retrieved: '%d' retrieved, '%d' expected", pRetrieved, pExpected)
	}

	const tax = 0.18
	var totalAmount float64 = 0.0

	for _, productRow := range productRows {
		qty := productQuantityIndex[productRow.ID]
		if qty > productRow.Stock {
			outOfStockErr := fmt.Errorf("product '%d' was requested %d units when only %d available", productRow.ID, qty, productRow.Stock)
			err = errors.Join(err, outOfStockErr)
			continue
		}
		totalAmount = totalAmount + utils.NumericToFloat(productRow.Price)*float64(qty)
	}

	if err != nil {
		return fmt.Errorf("invalid products: \n%v", err)
	}

	tx, err := db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed transaction initialization: \n%w", err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err == nil {
			log.Println("order transaction not completed")
		}
	}()

	dstx := h.DS.WithTx(tx)

	orderParams := datastore.CreateOrderParams{
		ClientName:  utils.StringToText(newOrder.ClientName),
		ClientEmail: utils.StringToText(newOrder.ClientEmail),
		ClientPhone: utils.StringToText(newOrder.ClientPhone),
		Note:        utils.StringToText(newOrder.Note),
		TotalAmount: utils.FloatToNumeric(totalAmount),
		Tax:         utils.FloatToNumeric(totalAmount * tax),
		UserID:      userID,
	}

	pk, err := dstx.CreateOrder(ctx, orderParams)
	if err != nil {
		return fmt.Errorf("order not created: %+v", orderParams)
	}

	orderItems := make([]datastore.CreateOrderItemsParams, len(productRows))

	for idx, priceRow := range productRows {
		qty := productQuantityIndex[priceRow.ID]

		product := datastore.UpdateStockParams{
			ID:    priceRow.ID,
			Stock: priceRow.Stock - qty,
		}

		err := dstx.UpdateStock(ctx, product)

		if err != nil {
			return fmt.Errorf("error updating stock of product '%d': %v", product.ID, err)
		}

		orderItems[idx] = datastore.CreateOrderItemsParams{
			UnitPrice: priceRow.Price,
			ProductID: priceRow.ID,
			Quantity:  qty,
			OrderID:   pk,
		}
	}

	_, err = dstx.CreateOrderItems(ctx, orderItems)
	if err != nil {
		log.Printf("price rows len: %d", len(productRows)) // Temporal
		log.Printf("order items tried to be created: %+v", orderItems)
		return fmt.Errorf("order items not created: %w", err)
	}

	return tx.Commit(ctx)
}

func (h *PostgresOrderHandler) GetOrders(ctx *gin.Context) {
	orders, err := h.DS.ListOrders(ctx)
	if err != nil {
		operationErr := &operationError{Entity: "orders", Operation: OperationGet, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": operationErr.Msg()})
		log.Println(operationErr.Error())
		return
	}

	if orders == nil {
		ctx.JSON(http.StatusOK, gin.H{"data": make([]interface{}, 0)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *PostgresOrderHandler) GetOrderById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	order, err := h.DS.GetOrder(ctx, int32(id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("No order found with id '%v'.", id),
			})
			log.Println(err.Error())
			return
		}

		operationErr := &operationError{Entity: "order", Operation: OperationGet, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	items, err := h.DS.GetOrderItemsDetail(ctx, int32(order.ID))

	if err != nil {
		operationErr := &operationError{Entity: "order_items", Operation: OperationGet, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  order,
		"items": items,
	})
}

func (h *PostgresOrderHandler) CancelOrder(ctx *gin.Context) {
	var cancelReason = struct{ CancelationReason string }{}
	err := ctx.ShouldBindJSON(&cancelReason)

	if err != nil {
		jsonErr := &jsonBindingError{Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": jsonErr.Msg(),
		})
		log.Println(jsonErr.Error())
		return
	}

	orderId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	var updateStatusParams = datastore.CancelOrderParams{
		ID:                 int32(orderId),
		CancellationReason: utils.StringToText(cancelReason.CancelationReason),
	}

	err = h.DS.CancelOrder(ctx, updateStatusParams)

	if err != nil {
		operationErr := &operationError{Entity: "order", Operation: OperationUpdate, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *PostgresOrderHandler) UpdateStatus(ctx *gin.Context) {
	orderId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	var updateStatusParams = datastore.UpdateOrderStatusParams{
		ID: int32(orderId),
	}

	if status, ok := isValidStatus(ctx.Param("status")); ok {
		updateStatusParams.Column2 = status
	} else {
		notValidParamErr := &notValidParamError{Param: "status", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	err = h.DS.UpdateOrderStatus(ctx, updateStatusParams)

	if err != nil {
		operationErr := &operationError{Entity: "order", Operation: OperationUpdate, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

func isValidStatus(value string) (datastore.OrderStatus, bool) {
	status := datastore.OrderStatus(value)
	switch status {
	case datastore.OrderStatusPending,
		datastore.OrderStatusCancelled,
		datastore.OrderStatusDelivered,
		datastore.OrderStatusShipped:
		return status, true
	default:
		return status, false
	}
}
