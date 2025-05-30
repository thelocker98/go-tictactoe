package routes

import (
	"net/http"

	"example.com/tictactoe/ai"
	"example.com/tictactoe/board"
	"example.com/tictactoe/models"
	"example.com/tictactoe/player"
	"github.com/gin-gonic/gin"
)

type webBoard struct {
	Board []int64 `json:"board" binding:"required"`
	Shape int64   `json:"shape" binding:"required"`
}

type webCreateGame struct {
	Opponent int64 `json:"opponent" binding:"required"`
	First    *bool `json:"first" binding:"required"`
	Shape    int64 `json:"shape" binding:"required"`
}

func getBestMove(context *gin.Context) {
	var p player.Player
	var b board.Board
	var webBoard webBoard

	err := context.ShouldBindJSON(&webBoard)

	if err != nil || len(webBoard.Board) != 9 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse board"})
		return
	}

	p.Shape = int64(webBoard.Shape)
	b.Board = [9]int64(webBoard.Board)

	move := ai.FindBestMove(&b, p)

	if move != -1 {
		b.Board[move] = p.Shape
	}

	win, winner := b.CheckWin()

	context.JSON(http.StatusCreated, gin.H{"win": win, "winner": winner, "move": move, "board": b.Board})
}

func createGame(context *gin.Context) {
	var webCreateGame webCreateGame

	err := context.ShouldBindJSON(&webCreateGame)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"})
		return
	}

	newGame, err := models.NewGame(context.GetInt64("userId"), webCreateGame.Shape, *webCreateGame.First, webCreateGame.Opponent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"gameId": newGame.GameId})
}
