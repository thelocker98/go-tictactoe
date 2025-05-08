package board

import (
	"fmt"

	"example.com/tictactoe/player"
)

func score(b *Board, p *player.Player) (bool, int64) {
	win, winner := b.CheckWin()
	if win && winner != 0 {
		if winner == p.Shape {
			return win, 10
		} else {
			return win, -10
		}
	}
	return win, 0
}

func FindBestMove(b *Board, p player.Player) int64 {
	var maxScore int64 = -20
	var loc int64 = -1

	for _, i := range RandomizeBoardIndex() {
		if b.Board[i] == 0 {
			b.Board[i] = p.Shape
			thisScore := -b.minmax(p.InvertShape())
			//_, thisScore := score(b, &p)

			if thisScore > maxScore {
				maxScore = thisScore
				loc = i
			}

			b.Board[i] = 0
		}
	}
	fmt.Println("best move", loc)
	return loc
}

func (b Board) minmax(p player.Player) int64 {
	win, scoreVal := score(&b, &p)
	if win {
		return scoreVal
	}

	var maxScore int64 = -20
	var loc int64 = -1

	for _, i := range RandomizeBoardIndex() {
		if b.Board[i] == 0 {
			b.Board[i] = p.Shape
			thisScore := -b.minmax(p.InvertShape())
			//_, thisScore := score(&b, &p)

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
