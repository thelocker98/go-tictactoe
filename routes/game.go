package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/tictactoe/models"
	"github.com/gin-gonic/gin"
)

type webCreateGame struct {
	Opponent  int64 `json:"opponent" binding:"required"`
	UserFirst int64 `json:"user_first" binding:"required"`
	Shape     int64 `json:"shape" binding:"required"`
}

func createGame(context *gin.Context) {
	var webCreateGame webCreateGame
	userId := context.GetInt64("userId")

	err := context.ShouldBindJSON(&webCreateGame)

	if err != nil || (webCreateGame.UserFirst != userId && webCreateGame.UserFirst != webCreateGame.Opponent) || userId == webCreateGame.Opponent {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"})
		return
	}

	newGame, err := models.NewGame(context.GetInt64("userId"), webCreateGame.Shape, webCreateGame.UserFirst, webCreateGame.Opponent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"gameId": newGame.GameId})
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
		game.Status = "ACCEPTED"

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

func deleteGame(context *gin.Context) {
	userId := context.GetInt64("userId")

	gameId, err1 := strconv.ParseInt(context.Param("id"), 10, 64)
	game, err2 := models.GetGameById(gameId)

	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "game does not exist"})
		return
	}

	if game.UserOwnerId != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "you do not have permission to access this game"})
		return
	}
	if game.Status == "DENYED" || game.Status == "PENDING" {

		err := models.DeleteGame(game)
		fmt.Println(game, err)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "error writing database"})
			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Game Deleted Successfuly"})
		return
	}

	context.JSON(http.StatusBadRequest, gin.H{"message": "game is in the wrong state"})
}
