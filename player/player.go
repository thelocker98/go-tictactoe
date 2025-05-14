package player

import (
	"fmt"
	"strings"
)

type Player struct {
	Name    string
	GoFirst bool
	Shape   int64 `binding:"required"`
}

func New(name string, goFirst bool, shape string) (*Player, error) {
	shape = strings.ToUpper(shape)
	var structShape int64

	if shape == "O" {
		structShape = -1
	} else if shape == "X" {
		structShape = 1
	} else {
		return nil, fmt.Errorf("input shape not valid: %s", shape)
	}

	return &Player{
		Name:    name,
		GoFirst: goFirst,
		Shape:   structShape,
	}, nil
}

func NewUser() (*Player, error) {
	var prompt string
	var player Player

	fmt.Println("Welcome to TicTacToe")

	fmt.Print("What is your name?: ")
	fmt.Scan(&player.Name)

	fmt.Print("Do you want to go first? [y/n]: ")
	fmt.Scan(&prompt)

	if strings.ToLower(prompt) == "n" {
		player.GoFirst = false
	} else if strings.ToLower(prompt) == "y" {
		player.GoFirst = true
	} else {
		return &player, fmt.Errorf("input shape not valid: %s", prompt)
	}

	fmt.Print("Do you want to be X or O? [x/o]: ")
	fmt.Scan(&prompt)
	fmt.Print("\n\n")

	if strings.ToUpper(prompt) == "X" {
		player.Shape = 1
	} else if strings.ToUpper(prompt) == "O" {
		player.Shape = -1
	} else {
		return &player, fmt.Errorf("input shape not valid: %s", prompt)
	}

	return &player, nil
}

func (p Player) InvertShape() Player {
	p.Shape *= -1
	return p
}

func (p *Player) EndGame() {
	if p.Shape != 0 {
		fmt.Printf("%s is the winner!\n", p.Name)
		fmt.Println("Great Game!")
	} else {
		fmt.Println("It's a tie!")
		fmt.Println("Great Game!")
	}
}
