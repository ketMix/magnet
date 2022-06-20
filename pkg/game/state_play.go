package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type PlayState struct {
	game  *Game
	level Level
	world World
}

func (s *PlayState) Init() error {
	s.world.game = s.game // Eww
	if err := s.world.BuildFromLevel(s.level); err != nil {
		return err
	}
	return nil
}

func (s *PlayState) Dispose() error {
	// Remove player entity reference.
	for _, p := range s.game.players {
		p.entity = nil
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
		if err := p.Update(s); err != nil {
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
	bounds := text.BoundString(boldFace, s.level.title)
	centeredX := screenWidth/2 - bounds.Min.X - bounds.Dx()/2
	text.Draw(screen, s.level.title, boldFace, centeredX, bounds.Dy()+1, color.White)

	// Draw our player's belt!
	s.game.players[0].toolbelt.Draw(screen)
}

// getCursorPosition returns the cursor position relative to the map.
func (s *PlayState) getCursorPosition() (x, y int) {
	x, y = ebiten.CursorPosition()
	x -= int(s.world.cameraX)
	y -= int(s.world.cameraY)
	return x, y
}
