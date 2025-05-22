package board

import (
	"fmt"
	"math/rand"
)

type Board struct {
	Board [9]int64 `binding:"required"`
}

func NewBoard() Board {
	return Board{
		Board: [9]int64{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func (b *Board) PrintBoard() {
	for c := 0; c < 3; c++ {
		fmt.Print("-------------\n")
		for r := 0; r < 3; r++ {
			loc := c*3 + r
			if b.Board[loc] != 0 {
				result := map[int64]string{-1: "O", 1: "X"}[b.Board[loc]]
				fmt.Print("| ", result, " ")
			} else {
				fmt.Print("| ", loc, " ")
			}
		}
		fmt.Print("|\n")
	}
	fmt.Print("-------------\n")
}

func (b *Board) CheckWin() (bool, int64) {
	win := false
	var winner int64 = 0

	var i int64
	for i = 0; i < 3; i++ {
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

func RandomizeBoardIndex() [9]int64 {
	numbers := [9]int64{0, 1, 2, 3, 4, 5, 6, 7, 8}

	// Shuffle the slice using Fisher-Yates algorithm
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	return numbers
}
