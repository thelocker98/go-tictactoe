package ai

import (
	"errors"
	"fmt"

	"example.com/tictactoe/models"
	"example.com/tictactoe/player"
)

func ComputerGameAccept(gameId int64) error {
	// Get game by id
	game, err := models.GetGameById(gameId)

	// Check that the computer is a part of this game and that the game was retrived from the database correctly
	if game.UserPlayerId != 1 || err != nil {
		return errors.New("game does not have correct permissions")
	}

	// Make sure that the game is in the pending state
	if game.Status == "PENDING" {
		game.Status = "ACCEPTED"

		err := models.UpdateGame(game)
		if err != nil {
			return errors.New("failed to update game status")
		}
		return nil
	}
	return errors.New("game in wrong state")
}

func ComputerPlayMove(gameId int64) error {
	// Get game using ID
	game, err := models.GetGameById(gameId)
	var p player.Player

	// Check it is computers turn and make sure their was no error getting the game by id
	if game.CurrentTurn != 1 || err != nil {
		return errors.New("not computers turn")
	}

	// Check that the game is not already won. If the game is already won that do not play and skip turn
	win, _ := game.Board.CheckWin()
	if win {
		return nil
	}

	// create the computer player struct and then find the best move
	p.Shape = game.UserOwnerShape * -1
	move := FindBestMove(&game.Board, p)

	// verify that the move is not an error and then update the game with the new move
	if move != -1 {
		game.Board.Board[move] = p.Shape

		game.CurrentTurn = game.UserOwnerId

		fmt.Println(game.Board)

		err = models.UpdateGame(game)
		if err != nil {
			return errors.New("could not update game")
		}

		return nil
	}

	// alert the system that their was an error finding a valid move for the computer
	return errors.New("invalid move by computer")
}
