package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type PlayState struct {
	game             *Game
	level            Level
	cameraX, cameraY float64
	entities         []Entity
}

func (s *PlayState) Init() error {
	s.buildFromLevel()
	return nil
}

// buildLevel builds the world from the level field.
func (s *PlayState) buildFromLevel() {
	for y, r := range s.level.cells {
		for x, c := range r {
			switch c.kind {
			case BlockedCell:
			case CoreCell:
			case PlayerCell:
				// Add entity for player if one does not exist.
				var target *Player
				for _, p := range s.game.players {
					if p.entity == nil {
						target = p
						break
					}
				}
				if target != nil {
					e := NewPlayerEntity(target)
					// Tie 'em together.
					e.player = target
					target.entity = e
					// And place.
					s.PlaceEntity(e, x, y)
				}
			case SouthSpawnCell:
			case NorthSpawnCell:
			case PathCell:
			}
		}
	}
}

// PlaceEntity places the given entity into the game world, aligned by cell and centered within a cell.
func (s *PlayState) PlaceEntity(e Entity, x, y int) {
	e.Physics().X = float64(x)*float64(cellWidth) + float64(cellWidth)/2
	e.Physics().Y = float64(y)*float64(cellHeight) + float64(cellHeight)/2
	s.entities = append(s.entities, e)
}

func (s *PlayState) Dispose() error {
	// Delete current entities.
	return nil
}

func (s *PlayState) Update() error {
	// Update our players.
	for _, p := range s.game.players {
		if err := p.Update(s); err != nil {
			return err
		}
	}

	// For now we're effectively recreating the entities slice per update, so as to allow for entity update followed by entity deletion.
	t := s.entities[:0]
	for _, e := range s.entities {
		e.Update()
		if !e.Trashed() {
			t = append(t, e)
		}
	}
	s.entities = t

	return nil
}

func (s *PlayState) Draw(screen *ebiten.Image) {
	// Get our "camera" position.
	screenOp := &ebiten.DrawImageOptions{}
	// FIXME: Base this on some sort of player lookup or a global self reference.
	if s.game.players[0].entity != nil {
		s.cameraX = -s.game.players[0].entity.Physics().X + float64(screenWidth)/2
		s.cameraY = -s.game.players[0].entity.Physics().Y + float64(screenHeight)/2
	}
	screenOp.GeoM.Translate(
		s.cameraX,
		s.cameraY,
	)

	// Draw the map.
	for y, r := range s.level.cells {
		for x := range r {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Concat(screenOp.GeoM)
			op.GeoM.Translate(float64(x*cellWidth), float64(y*cellHeight))
			screen.DrawImage(grassImage, op)
		}
	}

	// Draw our entities.
	for _, e := range s.entities {
		e.Draw(screen, screenOp)
	}

	// Draw level text centered at top of screen for now.
	bounds := text.BoundString(boldFace, s.level.title)
	centeredX := screenWidth/2 - bounds.Min.X - bounds.Dx()/2
	text.Draw(screen, s.level.title, boldFace, centeredX, bounds.Dy()+1, color.White)
}
