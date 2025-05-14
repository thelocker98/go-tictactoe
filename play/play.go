package play

import (
	"fmt"
	"time"

	"example.com/tictactoe/ai"
	"example.com/tictactoe/board"
	"example.com/tictactoe/player"
)

func play(b *board.Board, p *player.Player, loc int64) error {
	if b.Board[loc] == 0 {
		b.Board[loc] = p.Shape
		return nil
	}

	return fmt.Errorf("location %d is already taken", loc)
}

func UserPlay(b *board.Board, p *player.Player) error {
	fmt.Printf("%s's turn to play!\n", p.Name)
	b.PrintBoard()
	fmt.Print("\nWhere do you want to play? [0-8]: ")

	var loc int64 = 0
	fmt.Scan(&loc)

	if loc < 0 || loc > 8 {
		return fmt.Errorf("%s's input is not valid: %d", p.Name, loc)
	}

	return play(b, p, loc)
}

func ComputerPlay(b *board.Board, p *player.Player) error {
	fmt.Printf("%s's turn to play!\n", p.Name)
	b.PrintBoard()
	fmt.Println("Thinking...")
	time.Sleep(time.Second * 2)

	move := ai.FindBestMove(b, *p)

	play(b, p, move)
	return nil
}

func GetWinner(winner int64, user1 *player.Player, user2 *player.Player) *player.Player {
	if winner == user1.Shape {
		return user1
	} else if winner == user2.Shape {
		return user2
	}

	return &player.Player{Name: "Tie", Shape: 0}
}
