package game

import "github.com/hajimehoshi/ebiten/v2"

type PlayState struct {
	game *Game
}

func (s *PlayState) Init() error {
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
