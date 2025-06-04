package models

import (
	"encoding/json"
	"errors"
	"time"

	"example.com/tictactoe/board"
	"example.com/tictactoe/db"
)

type Game struct {
	GameId         int64
	UserOwnerId    int64       `binding:"required"`
	UserOwnerShape int64       `binding:"required"`
	CurrentTurn    int64       `binding:"required"`
	UserPlayerId   int64       `binding:"required"`
	Status         string      `binding:"required"`
	Board          board.Board `binding:"required"`
	DateTime       time.Time   `binding:"required"`
}

func NewGame(userOwnerId int64, userOwnerShape int64, currentTurn int64, userPlayerId int64) (Game, error) {
	status := "PENDING"
	if !(userOwnerShape == 1 || userOwnerShape == -1) {
		return Game{}, errors.New("invalid userOwnerShape")
	}

	query := "INSERT INTO games (user_owner_id, user_owner_shape, current_turn, user_player_id, status, board, date) VALUES (?, ?, ?, ?, ?, ?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return Game{}, err
	}
	defer stmt.Close()

	newBoard := board.NewBoard()
	currentTime := time.Now()

	jsonBoardState, _ := json.Marshal(newBoard)

	result, err := stmt.Exec(userOwnerId, userOwnerShape, currentTurn, userPlayerId, status, jsonBoardState, currentTime)

	if err != nil {
		return Game{}, err
	}

	gameId, err := result.LastInsertId()

	var newGame Game = Game{
		GameId:         gameId,
		UserOwnerId:    userOwnerId,
		UserOwnerShape: userOwnerShape,
		CurrentTurn:    currentTurn,
		UserPlayerId:   userPlayerId,
		Board:          newBoard,
		DateTime:       currentTime,
	}

	return newGame, err
}

func UpdateGame(game *Game) error {
	query := `
UPDATE games
SET user_owner_id = ?,
    user_owner_shape = ?,
    current_turn = ?,
    user_player_id = ?,
    status = ?,
    board = ?,
    date = ?
WHERE id = ?`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	currentTime := time.Now()

	jsonBoardState, _ := json.Marshal(game.Board)

	_, err = stmt.Exec(game.UserOwnerId, game.UserOwnerShape, game.CurrentTurn, game.UserPlayerId, game.Status, jsonBoardState, currentTime, game.GameId)

	return err
}

func DeleteGame(game *Game) error {
	query := `DELETE FROM games WHERE id = ?`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(game.GameId)

	return err
}

func GetGameById(gameId int64) (*Game, error) {
	query := `SELECT * FROM games WHERE id = ?`
	row := db.DB.QueryRow(query, gameId)

	var game Game
	var byteBoard []byte

	err := row.Scan(&game.GameId, &game.UserOwnerId, &game.UserOwnerShape, &game.CurrentTurn, &game.UserPlayerId, &game.Status, &byteBoard, &game.DateTime)

	var board board.Board
	json.Unmarshal(byteBoard, &board)
	game.Board = board

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func GetGameByUserId(userId int64) ([]Game, error) {
	query := `SELECT * FROM games WHERE user_owner_id = ? or user_player_id = ? ORDER BY date DESC;`
	rows, err := db.DB.Query(query, userId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game

	for rows.Next() {
		var game Game
		var byteBoard []byte

		err = rows.Scan(&game.GameId, &game.UserOwnerId, &game.UserOwnerShape, &game.CurrentTurn, &game.UserPlayerId, &game.Status, &byteBoard, &game.DateTime)
		if err != nil {
			return nil, err
		}

		var board board.Board
		json.Unmarshal(byteBoard, &board)
		game.Board = board

		games = append(games, game)
	}
	return games, nil
}
