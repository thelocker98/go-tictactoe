package board

import (
	"fmt"
)

type Board struct {
	Board [9]int
}

func NewBoard() Board {
	return Board{
		Board: [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func (b *Board) PrintBoard() {
	for c := 0; c < 3; c++ {
		fmt.Print("-------------\n")
		for r := 0; r < 3; r++ {
			loc := c*3 + r
			if b.Board[loc] != 0 {
				result := map[int]string{-1: "O", 1: "X"}[b.Board[loc]]
				fmt.Print("| ", result, " ")
			} else {
				fmt.Print("| ", loc, " ")
			}
		}
		fmt.Print("|\n")
	}
	fmt.Print("-------------\n")
}

func (b *Board) CheckWin() (bool, int) {
	win := false
	winner := 0

	for i := 0; i < 3; i++ {
		if b.Board[i*3] == b.Board[(i*3)+1] && b.Board[(i*3)+1] == b.Board[(i*3)+2] && b.Board[i*3] != 0 {
			win = true
			winner = b.Board[i*3]
			return win, winner
		}
		if b.Board[i] == b.Board[i+3] && b.Board[i+3] == b.Board[i+6] && b.Board[i] != 0 {
			win = true
			winner = b.Board[i]
			return win, winner
		}
	}

	if b.Board[0] == b.Board[4] && b.Board[4] == b.Board[8] && b.Board[4] != 0 {
		win = true
		winner = b.Board[4]
		return win, winner
	}
	if b.Board[2] == b.Board[4] && b.Board[4] == b.Board[6] && b.Board[4] != 0 {
		win = true
		winner = b.Board[4]
		return win, winner
	}

	for _, i := range b.Board {
		if i == 0 {
			return false, 0 // Not full and no winner yet
		}
	}

	return true, 0 // Tie
}

/*
X       X
  X   X
    X
  X   X
X       X



  OOOO
 O     O
O       O
 O     O
  OOOO


X | X | X
----------
X | X | X
----------
X | X | X

*/
