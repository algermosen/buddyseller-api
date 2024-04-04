package handlers

import (
	"example/buddyseller-api/db/datastore"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostgresOrderHandler struct {
	ds *datastore.Queries
}

func (handler *PostgresOrderHandler) PlaceOrder(ctx *gin.Context) {}
func (handler *PostgresOrderHandler) GetOrders(ctx *gin.Context) {
	orders, err := handler.ds.ListOrders(ctx)
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

	order, err := handler.ds.GetOrder(ctx, int32(id))

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

	err = handler.ds.UpdateOrderStatus(ctx, updateStatusParams)

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

	err = handler.ds.UpdateOrderStatus(ctx, updateStatusParams)

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
