package routes

import (
	"net/http"
	"strconv"

	"example/buddyseller-api/models"

	"github.com/gin-gonic/gin"
)

func getUsers(ctx *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch users. Try again later",
			"error":   err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, users)
}

func getUserById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Not recognized as a valid Id"})
		return
	}

	user, err := models.GetUserById(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func createUser(ctx *gin.Context) {

	var user models.User
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not match all the fields",
		})

		return
	}

	err = user.Save()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error Saving User",
			"details": err.Error(),
			"caller":  "user.Save()",
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "GOOD JOB!",
		"data":    user,
	})
}

func updateUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var updatedUser models.User
	err = ctx.ShouldBindJSON(&updatedUser)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedUser.ID = userId
	err = updatedUser.Update()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func deleteUser(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = models.DeleteUser(productId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
