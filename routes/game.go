package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/tictactoe/models"
	"github.com/gin-gonic/gin"
)

func rejectGame(context *gin.Context) {
	userId := context.GetInt64("userId")

	gameId, err1 := strconv.ParseInt(context.Param("id"), 10, 64)
	game, err2 := models.GetGameById(gameId)

	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "game does not exist"})
		return
	}

	if game.UserPlayerId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "you do not have permission to access this game"})
		return
	}
	if game.Status == "PENDING" {
		game.Status = "DENYED"

		err := models.UpdateGame(game)
		fmt.Println(game, err)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "error writing database"})
			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Game Rejected Successfuly"})
		return
	}

	context.JSON(http.StatusBadRequest, gin.H{"message": "game is in the wrong state"})
}

func acceptGame(context *gin.Context) {
	userId := context.GetInt64("userId")

	gameId, err1 := strconv.ParseInt(context.Param("id"), 10, 64)
	game, err2 := models.GetGameById(gameId)

	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "game does not exist"})
		return
	}

	if game.UserPlayerId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "you do not have permission to access this game"})
		return
	}
	if game.Status == "PENDING" {
		game.Status = "ACCEPT"

		err := models.UpdateGame(game)
		fmt.Println(game, err)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "error writing database"})
			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Game Accept Successfuly"})
		return
	}

	context.JSON(http.StatusBadRequest, gin.H{"message": "game is in the wrong state"})
}
