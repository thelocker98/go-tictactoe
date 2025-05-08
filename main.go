package main

import (
	"fmt"

	"example.com/tictactoe/board"
	"example.com/tictactoe/game"
	"example.com/tictactoe/player"
)

func main() {
	clearScreen()
	fmt.Println("TicTacToe")
	currentboard := board.NewBoard()

	user, err := player.NewUser()
	if err != nil {
		panic(err)
	}

	computer_shape := map[int]string{-1: "O", 1: "X"}[user.Shape*-1]
	computer, _ := player.New("Computer", !user.GoFirst, computer_shape)

	win := false
	winner := 0
	user_first := user.GoFirst
	for !win {
		clearScreen()

		if user_first {
			err := currentboard.UserPlay(user)
			if err != nil {
				continue
			}

		} else {
			err := currentboard.ComputerPlay(computer)
			if err != nil {
				continue
			}
		}

		win, winner = currentboard.CheckWin()
		fmt.Println(win, winner)

		user_first = !user_first
	}
	clearScreen()

	winningUser := game.GetWinner(winner, user, computer)
	winningUser.EndGame()

	currentboard.PrintBoard()
}

func clearScreen() {
	//fmt.Println("\x1b[2J\x1b[H")
}
