package ai

import (
	"example.com/tictactoe/board"
)

func score(b *board.Board, shape int64) (bool, int64) {
	win, winner := b.CheckWin()
	if win && winner != 0 {
		if winner == shape {
			return win, 10
		} else {
			return win, -10
		}
	}
	return win, 0
}

func FindBestMove(b *board.Board, shape int64) int64 {
	var maxScore int64 = -20
	var loc int64 = -1

	for _, i := range board.RandomizeBoardIndex() {
		if b.Board[i] == 0 {
			b.Board[i] = shape
			thisScore := -minmax(*b, shape*-1)

			if thisScore > maxScore {
				maxScore = thisScore
				loc = i
			}

			b.Board[i] = 0
		}
	}
	return loc
}

func minmax(b board.Board, shape int64) int64 {
	win, scoreVal := score(&b, shape)
	if win {
		return scoreVal
	}

	var maxScore int64 = -20
	var loc int64 = -1

	for _, i := range board.RandomizeBoardIndex() {
		if b.Board[i] == 0 {
			b.Board[i] = shape
			thisScore := -minmax(b, shape*-1)
			if thisScore > maxScore {
				maxScore = thisScore
				loc = i
			}

			b.Board[i] = 0
		}
	}

	if loc == -1 {
		return 0
	} else {
		return maxScore
	}
}
