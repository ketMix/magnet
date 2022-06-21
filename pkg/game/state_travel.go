package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type TravelState struct {
	game        *Game
	done        bool
	targetLevel string
	loadedLevel data.Level
}

func (s *TravelState) Init() (err error) {
	s.loadedLevel, err = data.NewLevel(s.targetLevel)
	if err != nil {
		return err
	}

	return nil
}

func (s *TravelState) Dispose() error {
	return nil
}

func (s *TravelState) Update() error {
	s.game.SetState(&PlayState{
		game:  s.game,
		level: s.loadedLevel,
	})
	return nil
}

func (s *TravelState) Draw(screen *ebiten.Image) {
}
