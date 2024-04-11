package handlers

import (
	"errors"
	"example/buddyseller-api/db/datastore"
	"example/buddyseller-api/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type PostgresUserHandler struct {
	DS *datastore.Queries
}

func (h *PostgresUserHandler) GetUsers(ctx *gin.Context) {
	users, err := h.DS.ListUsers(ctx)
	if err != nil {
		operationErr := &operationError{Entity: "users", Operation: OperationGet, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": operationErr.Msg()})
		log.Println(operationErr.Error())
		return
	}

	if users == nil {
		ctx.JSON(http.StatusOK, gin.H{"data": make([]interface{}, 0)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *PostgresUserHandler) GetUserById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	user, err := h.DS.GetUser(ctx, int32(id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("No user found with id '%v'.", id),
			})
			log.Println(err.Error())
			return
		}

		operationErr := &operationError{Entity: "user", Operation: OperationGet, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *PostgresUserHandler) CreateUser(ctx *gin.Context) {

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

	hashedPassword, err := utils.HashPassword(userParams.Password)

	if err != nil {
		msg := "Error hashing password."
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": msg,
		})
		log.Println(msg)
		return
	}

	userParams.Password = hashedPassword

	pk, err := h.DS.CreateUser(ctx, userParams)

	if err != nil {
		operationErr := &operationError{Entity: "user", Operation: OperationSave, origin: err}
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

func (h *PostgresUserHandler) UpdateUser(ctx *gin.Context) {
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
	rowsAffected, err := h.DS.UpdateUser(ctx, userParams)

	if err != nil {
		operationErr := &operationError{Entity: "user", Operation: OperationUpdate, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	if rowsAffected == 0 {
		msg := fmt.Sprintf("No user found with id '%v'.", userId)
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": msg,
		})
		log.Println(msg)
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

func (h *PostgresUserHandler) UpdatePassword(ctx *gin.Context) {
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
	rowsAffected, err := h.DS.UpdatePassword(ctx, userParams)

	if err != nil {
		operationErr := &operationError{Entity: "user", Operation: OperationUpdate, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	if rowsAffected == 0 {
		msg := fmt.Sprintf("No user found with id '%v'.", userId)
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": msg,
		})
		log.Println(msg)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *PostgresUserHandler) DeleteUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		notValidParamErr := &notValidParamError{Param: "id", Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": notValidParamErr.Msg()})
		log.Println(notValidParamErr.Error())
		return
	}

	rowsAffected, err := h.DS.DeleteUser(ctx, int32(userId))

	if err != nil {
		operationErr := &operationError{Entity: "user", Operation: OperationDelete, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	if rowsAffected == 0 {
		msg := fmt.Sprintf("No user found with id '%v'.", userId)
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": msg,
		})
		log.Println(msg)
		return
	}

	ctx.Status(http.StatusNoContent)
}
