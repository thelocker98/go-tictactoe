package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/tictactoe/ai"
	"example.com/tictactoe/models"
	"github.com/gin-gonic/gin"
)

type webBoard struct {
	Board []int64 `json:"board" binding:"required"`
	Shape int64   `json:"shape" binding:"required"`
}

type webGame struct {
	Board        [9]int64
	YourUserName string
	YourTurn     bool
	YourShape    int64
	Win          bool
	Winner       int64
}

type webPlay struct {
	Win    bool
	Winner int64
	Board  [9]int64
}

func playMove(context *gin.Context) {
	type move struct {
		GameId int64  `json:"gameid" binding:"required"`
		Move   *int64 `json:"move" binding:"required"`
	}
	var currentMove move

	userId := context.GetInt64("userId")
	err1 := context.ShouldBindJSON(&currentMove)

	game, err3 := models.GetGameById(currentMove.GameId)

	if err1 != nil || err3 != nil || !(*currentMove.Move >= 0 && *currentMove.Move <= 8) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse value"})
		return
	}

	// chech game state
	win, _ := game.Board.CheckWin()
	if game.Status != "ACCEPTED" || win || game.CurrentTurn != userId {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Wrong Game State"})
		return
	}

	// make play
	shape := game.UserOwnerShape
	if userId != game.UserOwnerId {
		shape = shape * -1
	}

	err := game.Board.Play(shape, *currentMove.Move)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid play"})
		return
	}

	// update database
	if game.CurrentTurn == game.UserOwnerId {
		game.CurrentTurn = game.UserPlayerId
	} else {
		game.CurrentTurn = game.UserOwnerId
	}

	err = models.UpdateGame(game)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "error updating game"})
		return
	}

	if game.CurrentTurn == 1 {
		err = ai.ComputerPlayMove(game.GameId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"}) // Computer failed to make a move
			return
		}
	}

	var responsePlay webPlay
	responsePlay.Win, responsePlay.Winner = game.Board.CheckWin()
	responsePlay.Board = game.Board.Board
	context.JSON(http.StatusOK, gin.H{"game": responsePlay})
}

func getBoardLayout(context *gin.Context) {
	userId := context.GetInt64("userId")

	gameId, err1 := strconv.ParseInt(context.Param("id"), 10, 64)
	game, err2 := models.GetGameById(gameId)

	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "game does not exist"})
		return
	}

	if game.CurrentTurn == 1 {
		err := ai.ComputerPlayMove(game.GameId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"}) // Computer failed to make a move
			return
		}

		// Pull updated game
		game, _ = models.GetGameById(gameId)
		fmt.Print("game update", game)
	}

	if (userId != game.UserOwnerId && userId != game.UserPlayerId) || game.Status != "ACCEPTED" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "do not have access to this game"})
		return
	}

	user, _ := models.GetUserById(userId)
	shape := game.UserOwnerShape
	if userId != game.UserOwnerId {
		shape = shape * -1
	}
	turn := false
	if userId == game.CurrentTurn {
		turn = true
	}

	webGame := webGame{
		Board:        game.Board.Board,
		YourUserName: user.UserName,
		YourShape:    shape,
		YourTurn:     turn,
	}
	webGame.Win, webGame.Winner = game.Board.CheckWin()
	context.JSON(http.StatusOK, gin.H{"game": webGame})
}
