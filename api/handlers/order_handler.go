package handlers

import (
	"context"
	"example/buddyseller-api/api/dtos"
	"example/buddyseller-api/db"
	"example/buddyseller-api/db/datastore"
	"example/buddyseller-api/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostgresOrderHandler struct {
	DS *datastore.Queries
}

func (handler *PostgresOrderHandler) PlaceOrder(ctx *gin.Context) {
	var newOrder dtos.NewOrderDto
	err := ctx.ShouldBindJSON(&newOrder)

	if err != nil {
		jsonErr := &jsonBindingError{Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": jsonErr.Msg(),
		})
		log.Println(jsonErr.Error())
		return
	}

	err = handler.placeNewOrder(ctx, &newOrder)

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
	})
}

func (handler *PostgresOrderHandler) placeNewOrder(ctx context.Context, newOrder *dtos.NewOrderDto) error {
	productIDs := make([]int32, len(newOrder.Items))
	mappedPrices := make(map[int32]int32)

	// Extract ProductIDs from the array of OrderItemDto
	for i, orderItem := range newOrder.Items {
		productIDs[i] = orderItem.ProductID
		mappedPrices[orderItem.ProductID] = orderItem.Quantity
	}

	priceRows, err := handler.DS.ListProductPrices(ctx, productIDs)

	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	dstx := handler.DS.WithTx(tx)

	const tax = 0.18
	var totalAmount float64 = 0.0

	for _, priceRow := range priceRows {
		totalAmount = totalAmount + utils.NumericToFloat(priceRow.Price)*float64(mappedPrices[priceRow.ID])
	}

	orderParams := datastore.CreateOrderParams{
		ClientName:  utils.StringToText(newOrder.ClientName),
		ClientEmail: utils.StringToText(newOrder.ClientEmail),
		ClientPhone: utils.StringToText(newOrder.ClientPhone),
		Note:        utils.StringToText(newOrder.Note),
		TotalAmount: utils.FloatToNumeric(totalAmount),
		Tax:         utils.FloatToNumeric(totalAmount * tax),
	}

	pk, err := dstx.CreateOrder(ctx, orderParams)

	if err != nil {
		return err
	}

	orderItems := make([]datastore.CreateOrderItemsParams, len(priceRows))

	for _, priceRow := range priceRows {
		orderItems = append(orderItems, datastore.CreateOrderItemsParams{
			UnitPrice: priceRow.Price,
			ProductID: priceRow.ID,
			Quantity:  mappedPrices[priceRow.ID],
			OrderID:   pk,
		})
	}

	_, err = dstx.CreateOrderItems(ctx, orderItems)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (handler *PostgresOrderHandler) GetOrders(ctx *gin.Context) {
	orders, err := handler.DS.ListOrders(ctx)
	if err != nil {
		operationErr := &operationError{Entity: "orders", Operation: OperationGet, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": operationErr.Msg()})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (handler *PostgresOrderHandler) GetOrderById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	order, err := handler.DS.GetOrder(ctx, int32(id))

	if err != nil {
		operationErr := &operationError{Entity: "order", Operation: OperationGet, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (handler *PostgresOrderHandler) CancelOrder(ctx *gin.Context) {
	orderId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	var updateStatusParams = datastore.UpdateOrderStatusParams{
		ID:     int32(orderId),
		Status: datastore.OrderStatusCancelled,
	}

	err = handler.DS.UpdateOrderStatus(ctx, updateStatusParams)

	if err != nil {
		operationErr := &operationError{Entity: "order", Operation: OperationUpdate, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *PostgresOrderHandler) UpdateStatus(ctx *gin.Context) {
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
		updateStatusParams.Status = status
	} else {
		notValidParamErr := &notValidParamError{Param: "status", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	err = handler.DS.UpdateOrderStatus(ctx, updateStatusParams)

	if err != nil {
		operationErr := &operationError{Entity: "order", Operation: OperationUpdate, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.Status(http.StatusOK)
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
