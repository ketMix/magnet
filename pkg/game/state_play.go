package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayState struct {
	game  *Game
	level Level
}

func (s *PlayState) Init() error {
	fmt.Println("TODO: generate live collision map/entities from", s.level)
	return nil
}

func (s *PlayState) Dispose() error {
	return nil
}

func (s *PlayState) Update() error {
	return nil
}

func (s *PlayState) Draw(screen *ebiten.Image) {
}
