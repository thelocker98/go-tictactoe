package board

import (
	"fmt"
	"time"

	"example.com/tictactoe/player"
)

func (b *Board) Play(p *player.Player, loc int64) error {
	if b.Board[loc] == 0 {
		b.Board[loc] = p.Shape
		return nil
	}

	return fmt.Errorf("location %d is already taken", loc)
}

func (b *Board) UserPlay(p *player.Player) error {
	fmt.Printf("%s's turn to play!\n", p.Name)
	b.PrintBoard()
	fmt.Print("\nWhere do you want to play? [0-8]: ")

	var loc int64 = 0
	fmt.Scan(&loc)

	if loc < 0 || loc > 8 {
		return fmt.Errorf("%s's input is not valid: %d", p.Name, loc)
	}

	return b.Play(p, loc)
}

func (b *Board) ComputerPlay(p *player.Player) error {
	fmt.Printf("%s's turn to play!\n", p.Name)
	b.PrintBoard()
	fmt.Println("Thinking...")
	time.Sleep(time.Second * 2)

	move := FindBestMove(b, *p)

	b.Play(p, move)
	return nil
}
