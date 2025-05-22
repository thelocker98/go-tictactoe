package models

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/tictactoe/board"
	"example.com/tictactoe/db"
)

type Game struct {
	GameId         int64
	UserOwnerId    int64       `binding:"required"`
	UserOwnerShape int64       `binding:"required"`
	UserOwnerTurn  bool        `binding:"required"`
	UserPlayerId   int64       `binding:"required"`
	Board          board.Board `binding:"required"`
	DateTime       time.Time   `binding:"required"`
}

func NewGame(userOwnerId int64, userOwnerShape int64, userOwnerTurnFirst bool, userPlayerId int64) (Game, error) {
	query := "INSERT INTO games (user_owner_id, user_owner_shape, user_owner_turn_first, user_player_id, board, date) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return Game{}, err
	}
	defer stmt.Close()

	newBoard := board.NewBoard()
	currentTime := time.Now()

	jsonBoardState, _ := json.Marshal(newBoard)

	result, err := stmt.Exec(userOwnerId, userOwnerShape, userOwnerTurnFirst, userPlayerId, jsonBoardState, currentTime)

	if err != nil {
		fmt.Println(err)
		return Game{}, err
	}

	gameId, err := result.LastInsertId()

	var newGame Game = Game{
		GameId:         gameId,
		UserOwnerId:    userOwnerId,
		UserOwnerShape: userOwnerShape,
		UserOwnerTurn:  userOwnerTurnFirst,
		UserPlayerId:   userPlayerId,
		Board:          newBoard,
		DateTime:       currentTime,
	}

	return newGame, err
}

func GetGameById(gameId int64) (*Game, error) {
	query := `SELECT * FROM games WHERE id = ?`
	row := db.DB.QueryRow(query, gameId)

	var game Game
	var byteBoard []byte

	err := row.Scan(&game.GameId, &game.UserOwnerId, &game.UserOwnerShape, &game.UserOwnerTurn, &game.UserPlayerId, &byteBoard, &game.DateTime)

	var board board.Board
	json.Unmarshal(byteBoard, &board)
	game.Board = board

	if err != nil {
		return nil, err
	}

	return &game, nil
}
