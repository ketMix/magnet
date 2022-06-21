package game

import (
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/goro/pathing"
)

// World is a struct for our cells and entities.
type World struct {
	game             *Game // Ewwww
	width, height    int
	cells            [][]LiveCell
	entities         []Entity
	currentTileset   data.TileSet
	cameraX, cameraY float64
	path             pathing.Path
	// Might as well store the core's position.
	coreX, coreY int
}

// BuildFromLevel builds the world's cells and entities from a given base level.
func (w *World) BuildFromLevel(level data.Level) error {
	tileset := level.Tileset
	if tileset == "" {
		tileset = "nature"
	}
	ts, err := data.LoadTileSet(tileset)
	if err != nil {
		return err
	}
	w.currentTileset = ts

	w.width = 0
	w.height = 0
	w.cells = make([][]LiveCell, 0)
	for y, r := range level.Cells {
		w.height++
		w.cells = append(w.cells, make([]LiveCell, 0))
		for x, c := range r {
			// Create any entities that should be there.
			if c.Kind == data.PlayerCell {
				var target *Player
				// Only add it if we actually need to add a player.
				for _, p := range w.game.players {
					if p.entity == nil {
						target = p
						break
					}
				}
				if target != nil {
					e := NewActorEntity(target, data.PlayerInit)
					// Tie 'em together.
					e.player = target
					target.entity = e
					// And place.
					w.PlaceEntityInCell(e, x, y)
				}
			} else if c.Kind == data.SouthSpawnCell {
				e := NewSpawnerEntity(data.NegativePolarity)
				w.PlaceEntityInCell(e, x, y)
			} else if c.Kind == data.NorthSpawnCell {
				e := NewSpawnerEntity(data.PositivePolarity)
				w.PlaceEntityInCell(e, x, y)
			} else if c.Kind == data.EnemyPositiveCell {
				e := NewEnemyEntity(data.EnemyConfigs["walker-positive"])
				w.PlaceEntityInCell(e, x, y)
			} else if c.Kind == data.EnemyNegativeCell {
				e := NewEnemyEntity(data.EnemyConfigs["walker-negative"])
				w.PlaceEntityInCell(e, x, y)
			} else if c.Kind == data.CoreCell {
				// Do we want more than 1 core...?
				e := NewCoreEntity(data.CoreConfig)
				w.PlaceEntityInCell(e, x, y)
				w.coreX = x
				w.coreY = y
			}
			// Create the cell.
			cell := LiveCell{
				kind: c.Kind,
				alt:  c.Alt,
			}
			if c.Kind == data.BlockedCell || c.Kind == data.EmptyCell {
				cell.blocked = true
			}
			w.cells[y] = append(w.cells[y], cell)
			if w.width < x+1 {
				w.width = x + 1
			}
		}
	}
	w.UpdatePathing()
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
				c := w.GetCell(r.x, r.y)
				if c != nil {
					if w.IsPlacementValid(r.x, r.y) && c.IsOpen() {
						e := NewTurretEntity(data.TurretConfigs["basic"])
						e.physics.polarity = r.polarity
						w.PlaceEntityInCell(e, r.x, r.y)
						turretPlaceSound.Play(1)
						c.entity = e
						w.UpdatePathing()
					}
				}
			} else if r.kind == ToolDestroy {
				c := w.GetCell(r.x, r.y)
				if c != nil {
					if c.entity != nil {
						c.entity.Trash()
						c.entity = nil
						w.UpdatePathing()
					}
				}
			} else if r.kind == ToolWall {
				c := w.GetCell(r.x, r.y)
				if c != nil {
					if w.IsPlacementValid(r.x, r.y) && c.IsOpen() {
						e := NewWallEntity()
						w.PlaceEntityInCell(e, r.x, r.y)
						turretPlaceSound.Play(1)
						c.entity = e
						w.UpdatePathing()
					}
				}
			}
		case SpawnProjecticleRequest:
			w.PlaceEntityAt(r.projectile, r.x, r.y)
		case SpawnEnemyRequest:
			e := NewEnemyEntity(r.enemyConfig)
			w.PlaceEntityAt(e, r.x, r.y)
			w.UpdatePathing()
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
			if c.kind == data.BlockedCell {
				// Don't mind my magic numbers.
				op.GeoM.Translate(0, -11)
				screen.DrawImage(w.currentTileset.BlockedImage, op)
			} else if c.kind == data.EmptyCell {
				// nada
			} else {
				if c.alt {
					screen.DrawImage(w.currentTileset.OpenImage2, op)
				} else {
					screen.DrawImage(w.currentTileset.OpenImage, op)
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
					// Draw transparent version of tool for placement
					if a.kind == ToolTurret || a.kind == ToolWall {
						image := GetToolKindImage(a.kind)
						op := &ebiten.DrawImageOptions{}
						op.ColorM.Scale(data.GetPolarityColorScale(a.polarity))
						op.ColorM.Scale(1, 1, 1, 0.5)
						op.GeoM.Concat(screenOp.GeoM)
						op.GeoM.Translate(float64(a.x*cellWidth)+float64(cellWidth/2), float64(a.y*cellHeight)+float64(cellHeight/2))
						// Draw from center.
						op.GeoM.Translate(
							-float64(image.Bounds().Dx())/2,
							-float64(image.Bounds().Dy())/2,
						)
						screen.DrawImage(image, op)
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

	// Pathing debug.
	/*for y, r := range w.cells {
		for x, c := range r {
			if c.IsOpen() {
				ebitenutil.DrawRect(screen, w.cameraX+float64(x*cellWidth+cellWidth/2), w.cameraY+float64(y*cellHeight+cellHeight/2), 2, 2, color.White)
			}
		}
	}*/
}

/** PATHING **/
func (w *World) UpdatePathing() {
	// Hmm.
	for _, e := range w.entities {
		w.UpdateEntityPathing(e)
	}
}

func (w *World) UpdateEntityPathing(e Entity) {
	if e.CanPathfind() {
		path := pathing.NewPathFromFunc(w.width, w.height, func(x, y int) uint32 {
			c := w.GetCell(x, y)
			if !c.IsOpen() {
				return pathing.MaximumCost
			}
			return 0
		}, pathing.AlgorithmAStar)
		cx, cy := w.GetClosestCellPosition(int(e.Physics().X), int(e.Physics().Y))
		steps := path.Compute(cx, cy, w.coreX, w.coreY)
		e.SetSteps(steps)
	}
}

func (w *World) IsPlacementValid(placeX, placeY int) bool {
	for _, e := range w.entities {
		path := pathing.NewPathFromFunc(w.width, w.height, func(x, y int) uint32 {
			c := w.GetCell(x, y)
			if !c.IsOpen() || (placeX == x && placeY == y) {
				return pathing.MaximumCost
			}
			return 1
		}, pathing.AlgorithmAStar)

		canPath := func(x1, y1, x2, y2 int) bool {
			steps := path.Compute(x1, y1, x2, y2)
			if x1 == x2 && y1 == y2 {
				return true
			}
			for _, s := range steps {
				if s.X() == x2 && s.Y() == y2 {
					return true
				}
			}
			return false
		}

		switch e.(type) {
		case *EnemyEntity:
			cx, cy := w.GetClosestCellPosition(int(e.Physics().X), int(e.Physics().Y))
			if !canPath(cx, cy, w.coreX, w.coreY) || (placeX == cx && placeY == cy) {
				return false
			}
		case *SpawnerEntity:
			cx, cy := w.GetClosestCellPosition(int(e.Physics().X), int(e.Physics().Y))
			if !canPath(cx, cy, w.coreX, w.coreY) || (placeX == cx && placeY == cy) {
				return false
			}
		}
	}
	return true
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

// EntitiesWithinRadius returns a slice of all entities within the radius of x, y
func (w *World) EntitiesWithinRadius(x, y float64, radius float64) []Entity {
	var entities []Entity
	for _, entity := range w.entities {
		if IsWithinRadius(x, y, entity.Physics().X, entity.Physics().Y, radius) {
			entities = append(entities, entity)
		}
	}
	return entities
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
	kind    data.CellKind // Same as Level
	alt     bool          // Same as Level
}

// IsOpen does what you think it does.
func (c *LiveCell) IsOpen() bool {
	return c.kind == data.NorthSpawnCell || c.kind == data.SouthSpawnCell || (c.entity == nil && c.blocked == false)
}
