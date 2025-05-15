package routes

import (
	"net/http"

	"example.com/tictactoe/models"
	"example.com/tictactoe/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse JSON"})
		return
	}

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not save user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user})
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse JSON"})
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})

}
