package routes

import (
	"errors"
	"fmt"

	"example.com/tictactoe/ai"
	"example.com/tictactoe/models"
)

type webGame struct {
	Board        [9]int64 `json:"board"`
	YourUserName string   `json:"player"`
	YourTurn     bool     `json:"yourTurn"`
	YourShape    int64    `json:"yourShape"`
	Win          bool     `json:"win"`
	Winner       int64    `json:"winner"`
}

func getBoardLayout(userId int64, gameId int64) (webGame, error) {
	game, err := models.GetGameById(gameId)

	if err != nil {
		return webGame{}, err
	}

	if game.CurrentTurn == 1 {
		err := ai.ComputerPlayMove(game.GameId)
		if err != nil {
			return webGame{}, err
		}

		// Pull updated game
		game, _ = models.GetGameById(gameId)
		fmt.Print("game update", game)
	}

	if (userId != game.UserOwnerId && userId != game.UserPlayerId) || game.Status != "ACCEPTED" {
		return webGame{}, err
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

	return webGame, nil
}

func playMove(userId int64, WSgamedata WSdataIn) (webGame, error) {

	//userId := context.GetInt64("userId")

	game, err := models.GetGameById(WSgamedata.GameId)

	if err != nil || !(WSgamedata.Move >= 0 && WSgamedata.Move <= 8) {
		return webGame{}, err
	}

	// chech game state
	win, _ := game.Board.CheckWin()
	if game.Status != "ACCEPTED" || win || game.CurrentTurn != userId {
		return webGame{}, err
	}

	// make play
	shape := game.UserOwnerShape
	if userId != game.UserOwnerId {
		shape = shape * -1
	}

	err = game.Board.Play(shape, WSgamedata.Move)
	if err != nil {
		return webGame{}, errors.New(fmt.Sprint("error invalid move", err)) // error invalid move
	}

	// update database
	if game.CurrentTurn == game.UserOwnerId {
		game.CurrentTurn = game.UserPlayerId
	} else {
		game.CurrentTurn = game.UserOwnerId
	}

	err = models.UpdateGame(game)
	if err != nil {
		return webGame{}, errors.New(fmt.Sprint("error updating game", err))
	}

	if game.CurrentTurn == 1 {
		err = ai.ComputerPlayMove(game.GameId)
		if err != nil {
			return webGame{}, errors.New(fmt.Sprint("could not parse", err))
		}
	}

	var responsePlay webGame
	responsePlay.Win, responsePlay.Winner = game.Board.CheckWin()
	responsePlay.Board = game.Board.Board

	return responsePlay, nil
}
