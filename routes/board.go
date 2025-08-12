package routes

import (
	"errors"
	"fmt"

	"gitea.locker98.com/locker98/go-tictactoe/ai"
	"gitea.locker98.com/locker98/go-tictactoe/models"
)

type webGame struct {
	Board  [9]int64 `json:"board"`
	User   int64    `json:"userId"`
	Turn   bool     `json:"userTurn"`
	Shape  int64    `json:"userShape"`
	Win    bool     `json:"userWin"`
	Winner int64    `json:"userWinner"`
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
	}

	if (userId != game.UserOwnerId && userId != game.UserPlayerId) || game.Status != "ACCEPTED" {
		return webGame{}, err
	}

	shape := game.UserOwnerShape
	if userId != game.UserOwnerId {
		shape = shape * -1
	}
	turn := false
	if userId == game.CurrentTurn {
		turn = true
	}

	webGame := webGame{
		Board: game.Board.Board,
		User:  userId,
		Shape: shape,
		Turn:  turn,
	}
	webGame.Win, webGame.Winner = game.Board.CheckWin()

	return webGame, nil
}

func playMove(userId int64, WSgamedata gameMoveData) (webGame, error) {
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

	turn := false
	if game.CurrentTurn == 1 {
		err = ai.ComputerPlayMove(game.GameId)
		if err != nil {
			return webGame{}, errors.New(fmt.Sprint("could not parse", err))
		}
		turn = true
		game, _ = models.GetGameById(game.GameId)
	}

	win, winner := game.Board.CheckWin()
	responsePlay := webGame{
		Board:  game.Board.Board,
		User:   userId,
		Turn:   turn,
		Shape:  shape,
		Win:    win,
		Winner: winner,
	}

	return responsePlay, nil
}
