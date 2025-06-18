package player

import (
	"fmt"
)

type Player struct {
	Name    string
	GoFirst bool
	Shape   int64 `binding:"required"`
}

func New(name string, goFirst bool, shape int64) (*Player, error) {

	if shape != 1 && shape != -1 {
		return nil, fmt.Errorf("input shape not valid: %d", shape)
	}

	return &Player{
		Name:    name,
		GoFirst: goFirst,
		Shape:   shape,
	}, nil
}

func (p Player) InvertShape() Player {
	p.Shape *= -1
	return p
}
