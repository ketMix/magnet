package game

import (
	"image/color"
	"math"
	"sort"

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
					e := NewActorEntity(target)
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
	e.Physics().X = float64(x*cellWidth + cellWidth/2)
	e.Physics().Y = float64(y*cellHeight + cellHeight/2)
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
	var requests []Request
	for _, e := range s.entities {
		if request, err := e.Update(); err != nil {
			panic(err)
		} else if request != nil {
			requests = append(requests, request)
		}
		if !e.Trashed() {
			t = append(t, e)
		}
	}
	s.entities = t

	// Iterate through our requests.
	for _, r := range requests {
		switch r := r.(type) {
		case SpawnTurretRequest:
			// TODO: Check if location is still valid.
			e := NewTurretEntity()
			s.PlaceEntity(e, r.x, r.y)
			// TODO: Mark cell as blocked.
		}
	}

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

	// Check for any special pending renders, such as move target or pending turret location.
	for _, p := range s.game.players {
		if p.entity != nil {
			if p.entity.Action() != nil && p.entity.Action().Next() != nil {
				switch a := p.entity.Action().Next().(type) {
				case *EntityActionPlace:
					op := &ebiten.DrawImageOptions{}
					op.ColorM.Scale(1, 1, 1, 0.5)
					op.GeoM.Concat(screenOp.GeoM)
					op.GeoM.Translate(float64(a.x*cellWidth)+float64(cellWidth/2), float64(a.y*cellHeight)+float64(cellHeight/2))
					// Draw from center.
					op.GeoM.Translate(
						-float64(turretBaseImage.Bounds().Dx())/2,
						-float64(turretBaseImage.Bounds().Dy())/2,
					)
					screen.DrawImage(turretBaseImage, op)
				}
			}
		}
	}

	// Make a sorted list of our entities to render.
	sortedEntities := make([]Entity, len(s.entities))
	copy(sortedEntities, s.entities)
	sort.SliceStable(sortedEntities, func(i, j int) bool {
		return sortedEntities[i].Physics().Y < sortedEntities[j].Physics().Y
	})
	// Draw our entities.
	for _, e := range sortedEntities {
		e.Draw(screen, screenOp)
	}

	// Draw level text centered at top of screen for now.
	bounds := text.BoundString(boldFace, s.level.title)
	centeredX := screenWidth/2 - bounds.Min.X - bounds.Dx()/2
	text.Draw(screen, s.level.title, boldFace, centeredX, bounds.Dy()+1, color.White)
}

// getCursorPosition returns the cursor position relative to the map.
func (s *PlayState) getCursorPosition() (x, y int) {
	x, y = ebiten.CursorPosition()
	x -= int(s.cameraX)
	y -= int(s.cameraY)
	return x, y
}

// getClosestCellPosition returns the closest cell position to the passed x and y coords.
func (s *PlayState) getClosestCellPosition(x, y int) (int, int) {
	tx, ty := math.Floor(float64(x)/float64(cellWidth)), math.Floor(float64(y)/float64(cellHeight))
	return int(tx), int(ty)
}
