package routes

import (
	"fmt"
	"net/http"

	"example.com/tictactoe/ai"
	"example.com/tictactoe/board"
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
	fmt.Println(webBoard)

	if err != nil || len(webBoard.Board) != 9 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse board"})
		return
	}

	p.Shape = int64(webBoard.Shape)
	b.Board = [9]int64(webBoard.Board)

	b.PrintBoard()

	move := ai.FindBestMove(&b, p)

	b.Board[move] = p.Shape

	b.PrintBoard()

	win, winner := b.CheckWin()

	c.JSON(http.StatusCreated, gin.H{"win": win, "winner": winner, "move": move, "board": b.Board})

}
