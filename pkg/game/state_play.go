package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
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
	// Draw level text centered at top of screen for now.
	bounds := text.BoundString(boldFace, s.level.title)
	centeredX := screenWidth/2 - bounds.Min.X - bounds.Dx()/2
	text.Draw(screen, s.level.title, boldFace, centeredX, bounds.Dy()+1, color.White)
}
