package models

import (
	"fmt"
	"time"

	"example.com/tictactoe/board"
	"example.com/tictactoe/db"
)

type game struct {
	GameId       int64
	UserId       int64       `binding:"required"`
	CurrentState state       `binding:"required"`
	Board        board.Board `binding:"required"`
	DateTime     time.Time   `binding:"required"`
}

type state struct {
	computerFirst bool
	usershape     int64
}

func NewGame(userId int64, computerFirst bool, userShape int64) (game, error) {
	query := "INSERT INTO games (user_id, state, board, date) VALUES (?, ?, ?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return game{}, err
	}
	defer stmt.Close()

	var gameState state = state{
		computerFirst: computerFirst,
		usershape:     userShape,
	}
	newBoard := board.NewBoard()
	currentTime := time.Now()

	result, err := stmt.Exec(userId, gameState, newBoard, currentTime)

	if err != nil {
		fmt.Println(err)
		return game{}, err
	}

	gameId, err := result.LastInsertId()

	var newGame game = game{
		GameId:       gameId,
		UserId:       userId,
		CurrentState: gameState,
		Board:        newBoard,
		DateTime:     currentTime,
	}

	return newGame, err
}
