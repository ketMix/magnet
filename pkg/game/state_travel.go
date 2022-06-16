package game

import "github.com/hajimehoshi/ebiten/v2"

type TravelState struct {
	game *Game
	done bool
}

func (s *TravelState) Init() error {
	return nil
}

func (s *TravelState) Dispose() error {
	return nil
}

func (s *TravelState) Update() error {
	s.game.SetState(&PlayState{})
	return nil
}

func (s *TravelState) Draw(screen *ebiten.Image) {
}
