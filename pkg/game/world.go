package game

import (
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// World is a struct for our cells and entities.
type World struct {
	game             *Game // Ewwww
	width, height    int
	cells            [][]LiveCell
	entities         []Entity
	currentTileset   TileSet
	cameraX, cameraY float64
}

// BuildFromLevel builds the world's cells and entities from a given base level.
func (w *World) BuildFromLevel(level Level) error {
	tileset := level.tileset
	if tileset == "" {
		tileset = "nature"
	}
	ts, err := loadTileSet(tileset)
	if err != nil {
		return err
	}
	w.currentTileset = ts

	w.width = 0
	w.height = 0
	w.cells = make([][]LiveCell, 0)
	for y, r := range level.cells {
		w.height++
		w.cells = append(w.cells, make([]LiveCell, 0))
		for x, c := range r {
			// Create any entities that should be there.
			if c.kind == PlayerCell {
				var target *Player
				// Only add it if we actually need to add a player.
				for _, p := range w.game.players {
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
					w.PlaceEntityInCell(e, x, y)
				}
			} else if c.kind == SouthSpawnCell || c.kind == NorthSpawnCell {
				e := NewSpawnerEntity()
				w.PlaceEntityInCell(e, x, y)
			} else if c.kind == EnemyPositiveCell {
				e := NewEnemyEntity(PositivePolarity)
				w.PlaceEntityInCell(e, x, y)
			} else if c.kind == EnemyNegativeCell {
				e := NewEnemyEntity(NegativePolarity)
				w.PlaceEntityInCell(e, x, y)
			}
			// Create the cell.
			cell := LiveCell{
				kind: c.kind,
				alt:  c.alt,
			}
			if c.kind != NoneCell && c.kind != PlayerCell {
				cell.blocked = true
			}
			w.cells[y] = append(w.cells[y], cell)
			if w.width < x+1 {
				w.width = x + 1
			}
		}
	}
	return nil
}

// Update updates the world.
func (w *World) Update() error {
	// TODO: Process physics

	// For now we're effectively recreating the entities slice per update, so as to allow for entity update followed by entity deletion.
	t := w.entities[:0]
	var requests []Request
	for _, e := range w.entities {
		if request, err := e.Update(w); err != nil {
			return err
		} else if request != nil {
			requests = append(requests, request)
		}
		if !e.Trashed() {
			t = append(t, e)
		}
	}
	w.entities = t

	// Iterate through our requests.
	for _, r := range requests {
		switch r := r.(type) {
		case UseToolRequest:
			if r.kind == ToolTurret {
				// TODO: Check if location is valid.
				c := w.GetCell(r.x, r.y)
				if c != nil {
					if c.IsOpen() {
						e := NewTurretEntity()
						w.PlaceEntityInCell(e, r.x, r.y)
						turretPlaceSound.Play(1)
						c.entity = e
					}
				}
				// TODO: Mark cell as blocked.
			} else if r.kind == ToolDestroy {
				c := w.GetCell(r.x, r.y)
				if c != nil {
					if c.entity != nil {
						c.entity.Trash()
						c.entity = nil
					}
				}
			}
		case SpawnProjecticleRequest:
			e := NewProjecticleEntity()
			e.physics.vX = r.vX
			e.physics.vY = r.vY
			e.physics.polarity = r.polarity
			w.PlaceEntityAt(e, r.x, r.y)
		}
	}

	return nil
}

// Draw draws the world, wow.
func (w *World) Draw(screen *ebiten.Image) {
	// Get our camera position.
	screenOp := &ebiten.DrawImageOptions{}
	// FIXME: Base this on some sort of player lookup or a global self reference.
	if w.game.players[0].entity != nil {
		w.cameraX = -w.game.players[0].entity.Physics().X + float64(screenWidth)/2
		w.cameraY = -w.game.players[0].entity.Physics().Y + float64(screenHeight)/2
	}
	screenOp.GeoM.Translate(
		w.cameraX,
		w.cameraY,
	)

	// Draw the map.
	for y, r := range w.cells {
		for x, c := range r {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Concat(screenOp.GeoM)
			op.GeoM.Translate(float64(x*cellWidth), float64(y*cellHeight))
			if c.kind == BlockedCell {
				// Don't mind my magic numbers.
				op.GeoM.Translate(0, -11)
				screen.DrawImage(w.currentTileset.blockedImage, op)
			} else if c.kind == EmptyCell {
				// nada
			} else {
				if c.alt {
					screen.DrawImage(w.currentTileset.openImage2, op)
				} else {
					screen.DrawImage(w.currentTileset.openImage, op)
				}
			}
		}
	}

	// Check for any special pending renders, such as move target or pending turret location.
	for _, p := range w.game.players {
		if p.entity != nil {
			if p.entity.Action() != nil && p.entity.Action().Next() != nil {
				switch a := p.entity.Action().Next().(type) {
				case *EntityActionPlace:
					if a.kind == ToolTurret {
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
	}

	// Make a sorted list of our entities to render.
	sortedEntities := make([]Entity, len(w.entities))
	copy(sortedEntities, w.entities)
	sort.SliceStable(sortedEntities, func(i, j int) bool {
		return sortedEntities[i].Physics().Y < sortedEntities[j].Physics().Y
	})
	// Draw our entities.
	for _, e := range sortedEntities {
		e.Draw(screen, screenOp)
	}
}

/** ENTITIES **/

// PlaceEntity places the entity into the world, aligned by cell and centered within a cell.
func (w *World) PlaceEntityInCell(e Entity, x, y int) {
	w.PlaceEntityAt(e, float64(x*cellWidth+cellWidth/2), float64(y*cellHeight+cellHeight/2))
}

// PlaceEntityAt places the entity into the world at the given specific coordinates.
func (w *World) PlaceEntityAt(e Entity, x, y float64) {
	e.Physics().X = x
	e.Physics().Y = y
	w.entities = append(w.entities, e)
}

/** CELLS **/

func (w *World) GetCell(x, y int) *LiveCell {
	if x < 0 || x >= w.width || y < 0 || y >= w.height {
		return nil
	}
	return &w.cells[y][x]
}

func (w *World) GetClosestCellPosition(x, y int) (int, int) {
	tx, ty := math.Floor(float64(x)/float64(cellWidth)), math.Floor(float64(y)/float64(cellHeight))
	return int(tx), int(ty)
}

// LiveCell is a position in a live level.
type LiveCell struct {
	entity  Entity
	blocked bool
	kind    CellKind // Same as Level
	alt     bool     // Same as Level
}

// IsOpen does what you think it does.
func (c *LiveCell) IsOpen() bool {
	return c.entity == nil && c.blocked == false
}
