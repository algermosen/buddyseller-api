package handlers

import (
	"example/buddyseller-api/db/datastore"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostgresUserHandler struct {
	DS *datastore.Queries
}

func (handler *PostgresUserHandler) GetUsers(ctx *gin.Context) {
	users, err := handler.DS.ListUsers(ctx)
	if err != nil {
		operationErr := &operationError{Entity: "users", Operation: OperationGet, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": operationErr.Msg()})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (handler *PostgresUserHandler) GetUserById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	user, err := handler.DS.GetUser(ctx, int32(id))

	if err != nil {
		operationErr := &operationError{Entity: "user", Operation: OperationGet, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *PostgresUserHandler) CreateUser(ctx *gin.Context) {

	var userParams datastore.CreateUserParams
	err := ctx.ShouldBindJSON(&userParams)

	if err != nil {
		jsonErr := &jsonBindingError{Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": jsonErr.Msg(),
		})
		log.Println(jsonErr.Error())
		return
	}

	pk, err := handler.DS.CreateUser(ctx, userParams)

	if err != nil {
		operationErr := &operationError{Entity: "user", Operation: OperationSave, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created succesfully.",
		"data": map[string]any{
			"ID":    pk,
			"Name":  userParams.Name,
			"Code":  userParams.Code,
			"Email": userParams.Email,
		},
	})
}

func (handler *PostgresUserHandler) UpdateUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	var userParams datastore.UpdateUserParams
	err = ctx.ShouldBindJSON(&userParams)

	if err != nil {
		jsonErr := &jsonBindingError{Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": jsonErr.Msg(),
		})
		log.Println(jsonErr.Error())
		return
	}

	userParams.ID = int32(userId)
	err = handler.DS.UpdateUser(ctx, userParams)

	if err != nil {
		operationErr := &operationError{Entity: "user", Operation: OperationUpdate, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated succesfully.",
		"data": map[string]any{
			"ID":    userParams.ID,
			"Name":  userParams.Name,
			"Code":  userParams.Code,
			"Email": userParams.Email,
		},
	})
}

func (handler *PostgresUserHandler) UpdatePassword(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	var userParams datastore.UpdatePasswordParams
	err = ctx.ShouldBindJSON(&userParams)

	if err != nil {
		jsonErr := &jsonBindingError{Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": jsonErr.Msg(),
		})
		log.Println(jsonErr.Error())
		return
	}

	userParams.ID = int32(userId)
	err = handler.DS.UpdatePassword(ctx, userParams)

	if err != nil {
		operationErr := &operationError{Entity: "user", Operation: OperationUpdate, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *PostgresUserHandler) DeleteUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	err = handler.DS.DeleteUser(ctx, int32(userId))

	if err != nil {
		operationErr := &operationError{Entity: "user", Operation: OperationDelete, Err: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}
