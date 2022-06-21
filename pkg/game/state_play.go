package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/world"
)

type PlayState struct {
	game  *Game
	level data.Level
	world world.World
}

func (s *PlayState) Init() error {
	s.world.Game = s.game // Eww
	if err := s.world.BuildFromLevel(s.level); err != nil {
		return err
	}
	return nil
}

func (s *PlayState) Dispose() error {
	// Remove player entity reference.
	for _, p := range s.game.players {
		p.Entity = nil
	}
	return nil
}

func (s *PlayState) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		s.game.SetState(&TravelState{
			game:        s.game,
			targetLevel: "001", // ???
		})
		return nil
	}
	// Update our players.
	for _, p := range s.game.players {
		if err := p.Update(&s.world); err != nil {
			return err
		}
	}

	// Update our world.
	if err := s.world.Update(); err != nil {
		return err
	}

	return nil
}

func (s *PlayState) Draw(screen *ebiten.Image) {
	// Draw our world.
	s.world.Draw(screen)

	// Draw level text centered at top of screen for now.
	bounds := text.BoundString(boldFace, s.level.Title)
	centeredX := world.ScreenWidth/2 - bounds.Min.X - bounds.Dx()/2
	text.Draw(screen, s.level.Title, boldFace, centeredX, bounds.Dy()+1, color.White)

	// Draw our player's belt!
	s.game.players[0].Toolbelt.Draw(screen)
}
