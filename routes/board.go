package routes

import (
	"net/http"
	"strconv"

	"example.com/tictactoe/models"
	"example.com/tictactoe/play"
	"example.com/tictactoe/player"
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
	userDB, err2 := models.GetUserById(userId)

	game, err3 := models.GetGameById(currentMove.GameId)

	if err1 != nil || err2 != nil || err3 != nil || !(*currentMove.Move >= 0 && *currentMove.Move <= 8) {
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
	var shape string
	shapeint := game.UserOwnerShape
	if userId != game.UserOwnerId {
		shapeint = shapeint * -1
	}
	if shapeint == -1 {
		shape = "O"
	} else if shapeint == 1 {
		shape = "X"
	}

	playerNew, err := player.New(userDB.UserName, false, shape)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create user"})
		return
	}

	err = play.Play(&game.Board, playerNew, *currentMove.Move)
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

/*
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
*/
