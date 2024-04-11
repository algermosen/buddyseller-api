package handlers

import (
	"errors"
	"example/buddyseller-api/api/dtos"
	"example/buddyseller-api/db/datastore"
	"example/buddyseller-api/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type PostgresSessionHandler struct {
	DS *datastore.Queries
}

func (h *PostgresSessionHandler) Login(ctx *gin.Context) {
	var user dtos.UserLoginDto
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		jsonErr := &jsonBindingError{Err: err}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": jsonErr.Msg(),
		})
		log.Println(jsonErr.Error())
		return
	}

	userCredentials, err := h.DS.GetUserCredentials(ctx, user.Code)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Could not validate user.",
			})
			log.Printf("Code '%v' not found.\n", user.Code)
			return
		}

		operationErr := &operationError{Entity: "user", Operation: OperationGet, origin: err}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": operationErr.Msg(),
		})
		log.Println(operationErr.Error())
		return
	}

	passwordIsValid := utils.CheckPassword(user.Password, userCredentials.Password)

	if !passwordIsValid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Could not validate user."})
		log.Println("Could not validate user.")
		return
	}

	token, err := utils.GenerateToken(userCredentials.Email, userCredentials.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token."})
		log.Println(err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Logged.", "token": token})
}
