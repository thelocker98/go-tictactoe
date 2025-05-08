package game

import (
	"example.com/tictactoe/player"
)

func GetWinner(winner int, user1 *player.Player, user2 *player.Player) *player.Player {
	if winner == user1.Shape {
		return user1
	} else if winner == user2.Shape {
		return user2
	}

	return &player.Player{Name: "Tie", Shape: 0}
}
