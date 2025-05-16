package routes

import (
	"fmt"
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

func getBestMove(c *gin.Context) {
	var p player.Player
	var b board.Board
	var webBoard webBoard

	err := c.ShouldBindJSON(&webBoard)

	if err != nil || len(webBoard.Board) != 9 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse board"})
		return
	}

	p.Shape = int64(webBoard.Shape)
	b.Board = [9]int64(webBoard.Board)

	move := ai.FindBestMove(&b, p)

	if move != -1 {
		b.Board[move] = p.Shape
	}

	win, winner := b.CheckWin()

	c.JSON(http.StatusCreated, gin.H{"win": win, "winner": winner, "move": move, "board": b.Board})
	createGame()
}

func createGame() {
	currentgame, _ := models.NewGame(1, true, -1)
	fmt.Println(currentgame)

	models.GetGameById(currentgame.GameId)
}
