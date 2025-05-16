package models

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/tictactoe/board"
	"example.com/tictactoe/db"
)

type Game struct {
	GameId       int64
	UserId       int64       `binding:"required"`
	CurrentState State       `binding:"required"`
	Board        board.Board `binding:"required"`
	DateTime     time.Time   `binding:"required"`
}

type State struct {
	ComputerFirst bool  `json:"computerFirst"`
	UserShape     int64 `json:"userShape"`
}

func NewGame(userId int64, computerFirst bool, userShape int64) (Game, error) {
	query := "INSERT INTO games (user_id, state, board, date) VALUES (?, ?, ?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return Game{}, err
	}
	defer stmt.Close()

	var gameState State = State{
		ComputerFirst: computerFirst,
		UserShape:     userShape,
	}

	newBoard := board.NewBoard()
	newBoard.Board[4] = 1
	currentTime := time.Now()

	jsonBoardState, _ := json.Marshal(newBoard)
	jsonGameState, _ := json.Marshal(gameState)

	result, err := stmt.Exec(userId, jsonGameState, jsonBoardState, currentTime)

	if err != nil {
		fmt.Println(err)
		return Game{}, err
	}

	gameId, err := result.LastInsertId()

	var newGame Game = Game{
		GameId:       gameId,
		UserId:       userId,
		CurrentState: gameState,
		Board:        newBoard,
		DateTime:     currentTime,
	}

	return newGame, err
}

func GetGameById(gameId int64) (*Game, error) {
	query := `SELECT * FROM games WHERE id = ?`
	row := db.DB.QueryRow(query, gameId)

	var game Game
	var byteState []byte
	var byteBoard []byte

	err := row.Scan(&game.GameId, &game.UserId, &byteState, &byteBoard, &game.DateTime)

	var state State
	json.Unmarshal(byteState, &state)

	var board board.Board
	json.Unmarshal(byteBoard, &board)

	game.CurrentState = state
	game.Board = board

	fmt.Println(game)

	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	return &game, nil
}
