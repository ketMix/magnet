package world

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/data/ui"
	"github.com/kettek/ebijam22/pkg/net"
	"github.com/kettek/goro/pathing"
)

// FIXME: This is a sin.
var ScreenWidth int = 640
var ScreenHeight int = 360

// World is a struct for our cells and entities.
type World struct {
	Game             Game // Ewwww x2
	Mode             WorldMode
	width, height    int
	cells            [][]LiveCell
	entities         []Entity
	netIDs           int
	trashedIDs       []int // This is a slice of trashed IDs for the current wave. This is used to ensure entities are not created if they're marked as trashed. This can happen due to out of order arrival of packets.
	spawners         []*SpawnerEntity
	enemies          []*EnemyEntity
	actors           []*ActorEntity
	currentTileset   data.TileSet
	CameraX, CameraY float64
	cameraShakeTimer int
	path             pathing.Path
	// Our waves, acquired from BuildFromLevel.
	waves       []*data.Wave
	CurrentWave int
	MaxWave     int
	// Might as well store the core's position.
	cores        []*CoreEntity
	coreX, coreY int
	// Overall game speed
	Speed float64
	//
	backgroundTimer int
	backgroundImage *ebiten.Image
	backgroundIndex int
	hasNextLevel    bool
}

// BuildFromLevel builds the world's cells and entities from a given base level.
func (w *World) BuildFromLevel(level data.Level) error {
	w.hasNextLevel = level.Next != ""
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
				// Add all players to the same spot. We _could_ adjust level parsing to have "n" and "s" for players.
				// Only add it if we actually need to add a player.
				for i, p := range w.Game.Players() {
					if i > 0 && !w.Game.Net().Active() {
						// Ignore players beyond 0 if we have no net.
						continue
					}
					if p.Entity == nil {
						c := data.PlayerInit
						xoffset := 0
						if i == 0 {
							if w.Game.Net().Active() && !w.Game.Net().Hosting() {
								c = data.Player2Init
								xoffset = 1
							}
						} else if i == 1 {
							if w.Game.Net().Active() && w.Game.Net().Hosting() {
								c = data.Player2Init
								xoffset = 1
							}
						}

						fmt.Println("Adding player entity", i, c.Title)
						e := NewActorEntity(p, c)
						// Tie 'em together.
						e.player = p
						p.Entity = e
						w.actors = append(w.actors, e)
						// And place.
						w.PlaceEntityInCell(e, x+xoffset, y)
					}
				}
			} else if c.Kind == data.SouthSpawnCell {
				e := NewSpawnerEntity(data.NegativePolarity)
				w.PlaceEntityInCell(e, x, y)
				w.spawners = append(w.spawners, e)
			} else if c.Kind == data.NorthSpawnCell {
				e := NewSpawnerEntity(data.PositivePolarity)
				w.PlaceEntityInCell(e, x, y)
				w.spawners = append(w.spawners, e)
			} else if c.Kind == data.EnemyPositiveCell {
				e := NewEnemyEntity(data.EnemyConfigs["walker-positive"])
				w.PlaceEntityInCell(e, x, y)
			} else if c.Kind == data.EnemyNegativeCell {
				e := NewEnemyEntity(data.EnemyConfigs["walker-negative"])
				w.PlaceEntityInCell(e, x, y)
			} else if c.Kind == data.CoreCell {
				e := NewCoreEntity(data.CoreConfig)
				w.PlaceEntityInCell(e, x, y)
				w.coreX = x
				w.coreY = y
				e.id = len(w.cores)
				// Do we want more than 1 core...?
				w.cores = append(w.cores, e)
			}
			// Create the cell.
			cell := LiveCell{
				kind:     c.Kind,
				polarity: c.Polarity,
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

	// Clone our waves list from the level.
	for _, wave := range level.Waves {
		w.waves = append(w.waves, wave.Clone())
	}

	w.SetWaves()

	// Set our player points/orbs.
	for _, pl := range w.Game.Players() {
		pl.Points = 0
	}
	w.SplitPoints(level.Points)
	if w.Game.Net().Hosting() {
		w.SendPlayerPoints()
	}

	w.UpdatePathing()
	return nil
}

func (w *World) ProcessNetMessage(msg net.Message) error {
	if w.Game.Net().Hosting() {
		switch msg := msg.(type) {
		case EntityActionMove:
			w.Game.Players()[1].Entity.SetAction(&msg)
		case EntityActionShoot:
			// let th' boy shoot
			w.Game.Players()[1].Entity.SetAction(&msg)
		case UseToolRequest:
			w.ProcessRequest(msg)
		}
	} else {
		switch msg := msg.(type) {
		case DamageCoreRequest:
			w.DamageCore(msg)
		case PlaySoundRequest:
			data.SFX.Play(msg.Sound)
		case PointsSync:
			w.SyncPoints(msg)
		case EntityPropertySync:
			w.SyncEntity(msg)
		case EntityActionMove:
			w.Game.Players()[1].Entity.SetAction(&msg)
		case SpawnEnemyRequest:
			w.SpawnEnemyEntity(msg)
		case SpawnOrbRequest:
			w.SpawnOrbEntity(msg)
		case CollectOrbRequest:
			w.CollectOrb(msg)
		case SpawnProjecticleRequest:
			w.SpawnProjecticleEntity(msg)
		case TrashEntityRequest:
			w.ProcessRequest(msg)
		case UseToolRequest:
			w.HandleToolRequest(msg)
		case BuildMode:
			w.SetMode(&msg)
		case WaveMode:
			w.SetMode(&msg)
		case VictoryMode:
			w.SetMode(&msg)
		case LossMode:
			w.SetMode(&msg)
		case PostGameMode:
			w.SetMode(&msg)
		default:
			fmt.Printf("unhandled net %+v\n", msg)
		}
	}
	return nil
}

func (w *World) ProcessRequest(r Request) {
	switch r := r.(type) {
	case MultiRequest:
		for _, rq := range r.Requests {
			w.ProcessRequest(rq)
		}
	case EntityPropertySync:
		if w.Game.Net().Hosting() {
			w.Game.Net().Send(r)
		}
	case DamageCoreRequest:
		if !w.Game.Net().Active() || w.Game.Net().Hosting() {
			w.DamageCore(r)
			if w.Game.Net().Hosting() {
				w.Game.Net().SendReliable(r)
			}
		}
	case UseToolRequest:
		// NOTE: Technically a client could just send this request and we won't do any distance checking.
		// Disallow tool use during wave mode.
		if _, ok := w.Mode.(*WaveMode); ok {
			return
		}
		// Deny clients from directly processing tool request. Sorry, lil buckaroos.
		if w.Game.Net().Active() && !w.Game.Net().Hosting() {
			// Send it to the overlord.
			w.Game.Net().SendReliable(r)
			return
		}
		if r.Tool == ToolTurret {
			if c := w.GetCell(r.X, r.Y); c != nil {
				if w.IsPlacementValid(r.X, r.Y) && c.IsOpen() {
					if r.local {
						r.Owner = w.Game.Players()[0].Name
					} else {
						r.Owner = w.Game.Players()[1].Name
					}

					pl := w.Game.GetPlayerByName(r.Owner)
					config := data.TurretConfigs[r.Kind]
					if pl.Points >= config.Points {
						e := w.HandleToolRequest(r)
						if e != nil {
							pl.Points -= config.Points
							w.SendPlayerPoints()
							// Let the client know to make our turret.
							if w.Game.Net().Hosting() {
								r.NetID = e.NetID()
								w.Game.Net().SendReliable(r)
							}
						}
					} else {
						if !r.local {
							w.Game.Net().SendReliable(PlaySoundRequest{
								Sound: "denied.ogg",
							})
						} else {
							data.SFX.Play("denied.ogg")
						}
					}
				} else {
					if !r.local {
						w.Game.Net().SendReliable(PlaySoundRequest{
							Sound: "denied.ogg",
						})
					} else {
						data.SFX.Play("denied.ogg")
					}
				}
			}
		} else if r.Tool == ToolDestroy {
			if r.local {
				r.Owner = w.Game.Players()[0].Name
			} else {
				r.Owner = w.Game.Players()[1].Name
			}
			w.HandleToolRequest(r)
			if w.Game.Net().Hosting() {
				w.Game.Net().SendReliable(r)
			}
		} else if r.Tool == ToolWall {
			if r.local {
				r.Owner = w.Game.Players()[0].Name
			} else {
				r.Owner = w.Game.Players()[1].Name
			}
			c := w.GetCell(r.X, r.Y)
			if c != nil {
				if w.IsPlacementValid(r.X, r.Y) && c.IsOpen() {
					pl := w.Game.GetPlayerByName(r.Owner)
					if pl.Points >= 3 {
						e := w.HandleToolRequest(r)
						if e != nil {

							pl.Points -= 3
							w.SendPlayerPoints()

							if w.Game.Net().Hosting() {
								r.NetID = e.NetID()
								w.Game.Net().SendReliable(r)
							}
						}
					} else {
						if !r.local {
							w.Game.Net().SendReliable(PlaySoundRequest{
								Sound: "denied.ogg",
							})
						} else {
							data.SFX.Play("denied.ogg")
						}
					}
				} else {
					if !r.local {
						w.Game.Net().SendReliable(PlaySoundRequest{
							Sound: "denied.ogg",
						})
					} else {
						data.SFX.Play("denied.ogg")
					}
				}
			}
		}
	case SpawnProjecticleRequest:
		if !w.Game.Net().Active() || w.Game.Net().Hosting() {
			e := w.SpawnProjecticleEntity(r)
			if w.Game.Net().Active() && w.Game.Net().Hosting() {
				w.Game.Net().SendReliable(SpawnProjecticleRequest{
					X:        r.X,
					Y:        r.Y,
					VX:       r.VX,
					VY:       r.VY,
					Polarity: r.Polarity,
					Damage:   r.Damage,
					NetID:    e.netID,
				})
			}
		}
	case SpawnEnemyRequest:
		if !w.Game.Net().Active() || w.Game.Net().Hosting() {
			e := w.SpawnEnemyEntity(r)
			// Hmm.
			if w.Game.Net().Active() && w.Game.Net().Hosting() {
				w.Game.Net().SendReliable(SpawnEnemyRequest{
					X:        r.X,
					Y:        r.Y,
					Polarity: r.Polarity,
					Kind:     r.Kind,
					NetID:    e.netID,
				})
			}
		}
	case SpawnOrbRequest:
		if !w.Game.Net().Active() || w.Game.Net().Hosting() {
			e := w.SpawnOrbEntity(r)
			// Hmm.
			if w.Game.Net().Active() && w.Game.Net().Hosting() {
				r.NetID = e.NetID()
				w.Game.Net().SendReliable(r)
			}
		}
	case CollectOrbRequest:
		// Only handle orb requests if we're the server or solo.
		if !w.Game.Net().Active() || w.Game.Net().Hosting() {
			w.CollectOrb(r)
			// Send it to client so they can play a sound if they collected it.
			if w.Game.Net().Hosting() {
				w.Game.Net().Send(r)
			}
		}
	case TrashEntityRequest:
		// Trash entities if we are local or host.
		if !w.Game.Net().Active() || w.Game.Net().Hosting() {
			r.entity.Trash()
			for i, e := range w.enemies {
				if e == r.entity {
					w.enemies = append(w.enemies[:i], w.enemies[i+1:]...)
					break
				}
			}
			if w.Game.Net().Hosting() {
				w.Game.Net().SendReliable(r)
			}
		} else if w.Game.Net().Active() {
			if !r.local {
				// Mark it as trashed.
				w.trashedIDs = append(w.trashedIDs, r.NetID)
				for _, e := range w.entities {
					if e.NetID() == r.NetID {
						e.Trash()
						for i, e2 := range w.enemies {
							if e2 == e {
								w.enemies = append(w.enemies[:i], w.enemies[i+1:]...)
								break
							}
						}
						break
					}
				}
			}
		}
	}
}

// ???
func (w *World) HandleToolRequest(r UseToolRequest) Entity {
	pl := w.Game.GetPlayerByName(r.Owner)
	if pl == nil {
		panic("PLAYER NIL")
	}

	if r.Tool == ToolTurret {
		config := data.TurretConfigs[r.Kind]
		// This is kind of stupid.
		var e Entity
		if config.AttackType == "beam" {
			te := NewTurretBeamEntity(config)
			te.owner = r.Owner

			// Hmm... this feels kind of gross.
			if r.Owner != "" {
				if pl := w.Game.GetPlayerByName(r.Owner); pl != nil {
					te.colorMultiplier = pl.Entity.(*ActorEntity).colorMultiplier
				}
			}

			e = te
		} else {
			te := NewTurretEntity(config)
			te.owner = r.Owner

			// Hmm... this feels kind of gross.
			if r.Owner != "" {
				if pl := w.Game.GetPlayerByName(r.Owner); pl != nil {
					te.colorMultiplier = pl.Entity.(*ActorEntity).colorMultiplier
				}
			}

			e = te
		}

		if w.Game.Net().Hosting() {
			e.SetNetID(w.GetNextNetID())
		} else {
			if r.NetID != 0 {
				for _, trashedID := range w.trashedIDs {
					if trashedID == r.NetID {
						// Oh! we don't want to construct this, as it has already been trashed.
						return nil
					}
				}
				e.SetNetID(r.NetID)
			}
		}
		e.Physics().polarity = r.Polarity
		w.PlaceEntityInCell(e, r.X, r.Y)
		data.SFX.Play("turret-place.ogg")

		if c := w.GetCell(r.X, r.Y); c != nil {
			c.entity = e
		}
		w.UpdatePathing()

		return e
	} else if r.Tool == ToolDestroy {
		c := w.GetCell(r.X, r.Y)
		if c != nil {
			if c.entity != nil {
				var points int
				ownerName := r.Owner
				switch e := c.entity.(type) {
				case *TurretEntity:
					ownerName = e.owner
					points = e.cost
				case *TurretBeamEntity:
					ownerName = e.owner
					points = e.cost
				case *WallEntity:
					ownerName = e.owner
					points = 3
				}
				if ownerName != r.Owner {
					if r.Owner == w.Game.Players()[0].Name {
						// TODO: Show some sort of "that isn't yours!" message on screen.
						fmt.Printf("that is %s's, not yours!\n", ownerName)
						data.SFX.Play("denied.ogg")
					}
				} else {

					if !w.Game.Net().Active() || w.Game.Net().Hosting() {
						pl.Points += points
					}
					w.SendPlayerPoints()

					c.entity.Trash()
					c.entity = nil
					w.UpdatePathing()
				}
			}
		}
	} else if r.Tool == ToolWall {
		e := NewWallEntity()
		e.owner = r.Owner
		w.PlaceEntityInCell(e, r.X, r.Y)
		data.SFX.Play("turret-place.ogg")

		if w.Game.Net().Hosting() {
			e.netID = w.GetNextNetID()
		} else {
			e.netID = r.NetID
		}

		// Hmm... this feels kind of gross.
		if r.Owner != "" {
			if pl := w.Game.GetPlayerByName(r.Owner); pl != nil {
				e.colorMultiplier = pl.Entity.(*ActorEntity).colorMultiplier
			}
		}

		if c := w.GetCell(r.X, r.Y); c != nil {
			c.entity = e
		}
		w.UpdatePathing()

		return e
	}

	return nil
}

func (w *World) SpawnEnemyEntity(r SpawnEnemyRequest) *EnemyEntity {
	enemyConfig := data.EnemyConfigs[r.Kind]
	e := NewEnemyEntity(enemyConfig)
	if w.Game.Net().Hosting() {
		e.netID = w.GetNextNetID()
	} else {
		if r.NetID != 0 {
			for _, trashedID := range w.trashedIDs {
				if trashedID == r.NetID {
					// Oh! we don't want to construct this, as it has already been trashed.
					return nil
				}
			}
			e.netID = r.NetID
		}
	}
	e.physics.polarity = r.Polarity
	w.enemies = append(w.enemies, e)
	w.PlaceEntityAt(e, r.X, r.Y)
	w.UpdatePathing()

	return e
}

func (w *World) SpawnOrbEntity(r SpawnOrbRequest) *OrbEntity {
	e := NewOrbEntity(r.Worth)
	if w.Game.Net().Hosting() {
		e.netID = w.GetNextNetID()
	} else {
		if r.NetID != 0 {
			for _, trashedID := range w.trashedIDs {
				if trashedID == r.NetID {
					// Oh! we don't want to construct this, as it has already been trashed.
					return nil
				}
			}
			e.netID = r.NetID
		}
	}
	w.PlaceEntityAt(e, r.X, r.Y)

	return e
}

func (w *World) CollectOrb(r CollectOrbRequest) {
	if !w.Game.Net().Active() || w.Game.Net().Hosting() {
		w.SplitPoints(r.Worth)
		w.SendPlayerPoints()
	}
	if r.Collector == w.Game.Players()[0].Name {
		s := data.SFX.Play("pop.ogg")
		if !data.SFX.Muted {
			if r.Worth <= 10 {
				s.SetVolume(0.5)
			} else if r.Worth <= 15 {
				s.SetVolume(1)
			} else {
				s.SetVolume(1.5)
			}
		}
	}
}

func (w *World) DamageCore(r DamageCoreRequest) {
	for _, c := range w.cores {
		if c.id == r.ID {
			c.health -= r.Damage
			data.SFX.Play("core-damage.ogg")
			w.cameraShakeTimer = 30
			if c.health <= 0 && !c.destroyed {
				c.destroyed = true
				data.SFX.Play("loss-hit.ogg")
			}
			return
		}
	}
}

func (w *World) SpawnProjecticleEntity(r SpawnProjecticleRequest) *ProjecticleEntity {
	e := NewProjecticleEntity()
	e.physics.vX = r.VX
	e.physics.vY = r.VY
	if w.Game.Net().Hosting() {
		e.netID = w.GetNextNetID()
	} else {
		if r.NetID != 0 {
			for _, trashedID := range w.trashedIDs {
				if trashedID == r.NetID {
					// Oh! we don't want to construct this, as it has already been trashed.
					return nil
				}
			}
			e.netID = r.NetID
		}
	}
	e.physics.polarity = r.Polarity
	e.damage = r.Damage
	w.PlaceEntityAt(e, r.X, r.Y)

	// We should probably control how loud we shot web in accordance to the source -- as in, if it is the other player's shot, it should be quieter.
	data.SFX.Play("shot.ogg")

	return e
}

func (w *World) SyncEntity(r EntityPropertySync) {
	for _, e := range w.entities {
		if e.NetID() == r.NetID {
			e.Physics().X = r.X
			e.Physics().Y = r.Y
			switch e := e.(type) {
			case *EnemyEntity:
				e.health = r.Health
			}
			break
		}
	}
}

func (w *World) SplitPoints(value int) {
	// Get it's split value.
	worth := math.Max(1, math.Floor(float64(value)/float64(len(w.Game.Players()))))
	for _, pl := range w.Game.Players() {
		pl.Points += int(worth)
	}
}

func (w *World) SendPlayerPoints() {
	if w.Game.Net().Hosting() {
		m := PointsSync{
			Points: make(map[string]int),
		}
		// Generate a points sync message.
		for _, pl := range w.Game.Players() {
			m.Points[pl.Name] = pl.Points
		}
		w.Game.Net().SendReliable(m)
	}
}

// SyncPoints synchronizes the players' points to match the provided data.
func (w *World) SyncPoints(r PointsSync) {
	for name, value := range r.Points {
		if pl := w.Game.GetPlayerByName(name); pl != nil {
			pl.Points = value
		}
	}
}

// SetMode sets the current game mode to the one indicated. If we're a client and the mode set is local, this does nothing.
func (w *World) SetMode(m WorldMode) {
	// If it is a locally generated mode and we're a client, do nothing.
	if w.Game.Net().Active() && !w.Game.Net().Hosting() && m.Local() {
		return
	}
	// Otherwise update the mode.
	w.Mode = m

	m.Init(w)

	if _, ok := w.Mode.(*WaveMode); ok {
		// Clear out our trashed IDs history on wave start.
		w.trashedIDs = make([]int, 0)
	}

	// Also send the network message if we're the host.
	if w.Game.Net().Active() && w.Game.Net().Hosting() && m.Local() {
		w.Game.Net().SendReliable(m)
	}
}

// Update updates the world.
func (w *World) Update() error {
	if w.cameraShakeTimer > 0 {
		w.cameraShakeTimer--
	}
	// Silly background processing.
	if len(w.currentTileset.BackgroundImages) > 0 {
		w.backgroundTimer++
		if w.backgroundTimer >= 30 {
			w.backgroundTimer = 0
			w.backgroundIndex++
		}
		if w.backgroundIndex >= len(w.currentTileset.BackgroundImages) {
			w.backgroundIndex = 0
		}
		w.backgroundImage = w.currentTileset.BackgroundImages[w.backgroundIndex]
	}

	// TODO: Move this elsewhere
	if inpututil.IsKeyJustPressed(ebiten.KeyAlt) {
		for _, e := range w.entities {
			if e, ok := e.(*TurretEntity); ok {
				e.showRange = true
			}
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyAlt) {
		for _, e := range w.entities {
			if e, ok := e.(*TurretEntity); ok {
				e.showRange = false
			}
		}
	}

	// TODO: Add delay between mode switching!
	if nextMode, _ := w.Mode.Update(w); nextMode != nil {
		w.SetMode(nextMode)
	}

	// Update our entities and get any requests.
	var requests []Request
	for _, e := range w.entities {
		if request, err := e.Update(w); err != nil {
			return err
		} else if request != nil {
			requests = append(requests, request)
		}
	}

	// Iterate through our requests.
	for _, r := range requests {
		w.ProcessRequest(r)
	}

	// Clean up any destroyed entities.
	t := w.entities[:0]
	for _, e := range w.entities {
		if !e.Trashed() {
			t = append(t, e)
		}
	}
	w.entities = t

	return nil
}

// Draw draws the world, wow.
func (w *World) Draw(screen *ebiten.Image) {
	// Get our camera position.
	screenOp := &ebiten.DrawImageOptions{}

	// FIXME: Base this on some sort of player lookup or a global self reference.
	if w.Game.Players()[0].Entity != nil {
		w.CameraX = -w.Game.Players()[0].Entity.Physics().X + float64(ScreenWidth)/2
		w.CameraY = -w.Game.Players()[0].Entity.Physics().Y + float64(ScreenHeight)/2
	}

	// Shake the camera if the timer is set.
	if w.cameraShakeTimer > 0 {
		w.CameraX += math.Sin(float64(w.cameraShakeTimer)*2) / 2
		w.CameraY += math.Cos(float64(w.cameraShakeTimer)*2) / 2
	}

	screenOp.GeoM.Translate(
		w.CameraX,
		w.CameraY,
	)

	// Draw any mode overlays
	w.Mode.Draw(w, screen)

	// Draw da background.
	if w.backgroundImage != nil {
		bgOp := &ebiten.DrawImageOptions{}
		width := ScreenWidth * 2
		height := ScreenHeight * 2
		bgOp.GeoM.Translate(-w.CameraX/float64(width/32), -w.CameraY/float64(height/32))
		ui.DrawTiled(screen, w.backgroundImage, bgOp, width, height)
	}

	// Draw the map.
	for y, r := range w.cells {
		for x, c := range r {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Concat(screenOp.GeoM)
			op.GeoM.Translate(float64(x*data.CellWidth), float64(y*data.CellHeight))
			if c.kind == data.BlockedCell {
				// Don't mind my magic numbers.
				op.GeoM.Translate(0, -11)
				screen.DrawImage(w.currentTileset.BlockedImage, op)
			} else if c.kind == data.EmptyCell {
				// nada
			} else {
				if c.polarity == data.PositivePolarity {
					screen.DrawImage(w.currentTileset.OpenPositiveImage, op)
				} else if c.polarity == data.NegativePolarity {
					screen.DrawImage(w.currentTileset.OpenNegativeImage, op)
				} else {
					screen.DrawImage(w.currentTileset.OpenNeutralImage, op)
				}
			}
		}
	}

	// Check for any special pending renders, such as move target or pending turret location.
	for _, p := range w.Game.Players() {
		if p.Entity != nil {
			if p.HoveringPlacement {
				if p.HoveringPlace.Tool == ToolTurret || p.HoveringPlace.Tool == ToolWall {
					image := GetToolImage(p.HoveringPlace.Tool, p.HoveringPlace.Kind)
					op := &ebiten.DrawImageOptions{}
					op.ColorM.Scale(data.GetPolarityColorScale(p.HoveringPlace.Polarity))
					op.ColorM.Scale(1, 1, 1, 0.5)
					op.GeoM.Concat(screenOp.GeoM)
					op.GeoM.Translate(float64(p.HoveringPlace.X*data.CellWidth)+float64(data.CellWidth/2), float64(p.HoveringPlace.Y*data.CellHeight)+float64(data.CellHeight/2))
					// Draw from center.
					op.GeoM.Translate(
						-float64(image.Bounds().Dx())/2,
						-float64(image.Bounds().Dy())/2,
					)
					screen.DrawImage(image, op)
				}
			}
			if p.Entity.Action() != nil && p.Entity.Action().GetNext() != nil {
				switch a := p.Entity.Action().GetNext().(type) {
				case *EntityActionPlace:
					// Draw transparent version of tool for placement
					if a.Tool == ToolTurret || a.Tool == ToolWall {
						image := GetToolImage(a.Tool, a.Kind)
						op := &ebiten.DrawImageOptions{}
						op.ColorM.Scale(data.GetPolarityColorScale(a.Polarity))
						op.ColorM.Scale(1, 1, 1, 0.5)
						op.GeoM.Concat(screenOp.GeoM)
						op.GeoM.Translate(float64(a.X*data.CellWidth)+float64(data.CellWidth/2), float64(a.Y*data.CellHeight)+float64(data.CellHeight/2))
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
		a := sortedEntities[i]
		b := sortedEntities[j]
		if a.IsProjectile() && !b.IsProjectile() {
			return false
		} else if !a.IsProjectile() && b.IsProjectile() {
			return true
		}
		return a.Physics().Y < b.Physics().Y
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

// ArePlayersReady returns true if all players are ready to start.
func (w *World) ArePlayersReady() bool {
	playersCount := len(w.Game.Players())
	for _, p := range w.Game.Players() {
		if p.ReadyForWave {
			playersCount--
		}
	}
	return playersCount == 0
}

// AreCoresDead returns true if all cores are dead.
func (w *World) AreCoresDead() bool {
	coreCount := len(w.cores)
	for _, c := range w.cores {
		if c.health <= 0 {
			coreCount--
		}
	}
	return coreCount < 1
}

/** WAVES **/
func (w *World) AreSpawnersHolding() bool {
	spawnerCount := len(w.spawners)

	for _, s := range w.spawners {
		if s.heldWave {
			spawnerCount--
		}
	}

	return spawnerCount == 0
}

func (w *World) AreEnemiesDead() bool {
	return len(w.enemies) == 0
}

func (w *World) AreWavesComplete() bool {
	spawnerCount := len(w.spawners)

	for _, s := range w.spawners {
		if s.wave == nil {
			spawnerCount--
		}
	}

	return spawnerCount == 0 && w.AreEnemiesDead()
}

func (w *World) SetWaves() {
	for i, wave := range w.waves {
		if i >= len(w.spawners) {
			// Ignore, wave definition is beyond our actual spawner count.
			continue
		}
		w.spawners[i].wave = wave
		// Set the spawn elapsed to the first spawn's spawn rate so as to ensure immediate spawning.
		if wave.Spawns != nil {
			w.spawners[i].spawnElapsed = float64(wave.Spawns.Spawnrate)
		}
		// Make sure to hold until build is done.
		w.spawners[i].heldWave = true
		// Okay, this is a bit stupid.
		count := 1
		for wave := wave.Next; wave != nil; wave = wave.Next {
			count++
		}
		if w.MaxWave < count {
			w.MaxWave = count
		}
	}
}

/** PATHING **/
func (w *World) UpdatePathing() {
	// Hmm.
	for _, e := range w.enemies {
		w.UpdateEntityPathing(e)
	}
	for _, e := range w.spawners {
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
	w.PlaceEntityAt(e, float64(x*data.CellWidth+data.CellWidth/2), float64(y*data.CellHeight+data.CellHeight/2))
}

// PlaceEntityAt places the entity into the world at the given specific coordinates.
func (w *World) PlaceEntityAt(e Entity, x, y float64) {
	e.Physics().X = x
	e.Physics().Y = y
	w.entities = append(w.entities, e)
}

// ObjectsWithinRadius is a generic function that can apply to a slice of any object that has a Physics() *PhysicsObject method.
func ObjectsWithinRadius[K interface{ Physics() *PhysicsObject }](l []K, x, y, radius float64) []K {
	var results []K
	for _, target := range l {
		if IsWithinRadius(x, y, target.Physics().X, target.Physics().Y, radius) {
			results = append(results, target)
		}
	}
	return results
}

func ObjectsWithPolarity[K interface{ Physics() *PhysicsObject }](l []K, p data.Polarity) []K {
	var results []K
	for _, target := range l {
		if target.Physics().polarity == p {
			results = append(results, target)
		}
	}
	return results
}

func ObjectsNearest[K interface{ Physics() *PhysicsObject }](l []K, x, y float64) []K {
	var results []K

	results = make([]K, len(l))
	copy(results, l)

	sort.Slice(results, func(i, j int) bool {
		a := GetMagnitude(GetDistanceVector(x, y, results[i].Physics().X, results[i].Physics().Y))
		b := GetMagnitude(GetDistanceVector(x, y, results[j].Physics().X, results[j].Physics().Y))
		return a < b
	})

	return results
}

/** CELLS **/

func (w *World) GetCell(x, y int) *LiveCell {
	if x < 0 || x >= w.width || y < 0 || y >= w.height {
		return nil
	}
	return &w.cells[y][x]
}

func (w *World) GetClosestCellPosition(x, y int) (int, int) {
	tx, ty := math.Floor(float64(x)/float64(data.CellWidth)), math.Floor(float64(y)/float64(data.CellHeight))
	return int(tx), int(ty)
}

// GetCursorPosition returns the cursor position relative to the map.
func (w *World) GetCursorPosition() (x, y int) {
	x, y = ebiten.CursorPosition()
	x -= int(w.CameraX)
	y -= int(w.CameraY)
	return x, y
}

func (w *World) GetNextNetID() int {
	w.netIDs++
	return w.netIDs
}

// LiveCell is a position in a live level.
type LiveCell struct {
	entity   Entity
	blocked  bool
	kind     data.CellKind // Same as Level
	polarity data.Polarity
}

// IsOpen does what you think it does.
func (c *LiveCell) IsOpen() bool {
	return c.kind == data.NorthSpawnCell || c.kind == data.SouthSpawnCell || (c.entity == nil && c.blocked == false)
}

// This is _not_ the place for this.
func DrawWaves(w *World, screen *ebiten.Image, spawnerOp *ebiten.DrawImageOptions) {
	for _, spawner := range w.spawners {
		yAdjust := 0
		if spawner.wave != nil {
			x := 0
			for spawnList := spawner.wave.Spawns; spawnList != nil; spawnList = spawnList.Next {
				t := fmt.Sprintf("%d", spawnList.Count)
				bounds := text.BoundString(data.NormalFace, t)
				text.Draw(screen, t, data.NormalFace, x+int(spawnerOp.GeoM.Element(0, 2)), int(spawnerOp.GeoM.Element(1, 2))+bounds.Dy(), color.White)
				x += bounds.Dx() + 3
				// Draw th' dude(s).
				eop := ebiten.DrawImageOptions{}
				eop.ColorM.Scale(data.GetPolarityColorScale(spawner.physics.polarity))
				eop.GeoM.Concat(spawnerOp.GeoM)
				eop.GeoM.Translate(float64(x), 0)
				for _, kind := range spawnList.Kinds {
					if enemy, ok := data.EnemyConfigs[kind]; ok {
						ebitenutil.DrawRect(screen, eop.GeoM.Element(0, 2)-1, eop.GeoM.Element(1, 2)-1, float64(enemy.WalkImages[0].Bounds().Dx())+1, float64(enemy.WalkImages[0].Bounds().Dy())+1, color.RGBA{128, 128, 128, 64})

						if enemy.WalkImages[0].Bounds().Dy() >= yAdjust {
							yAdjust = enemy.WalkImages[0].Bounds().Dy()
						}

						// Actually draw the non-shameful image.
						screen.DrawImage(enemy.WalkImages[0], &eop)

						x += enemy.WalkImages[0].Bounds().Dx() + 2
						eop.GeoM.Translate(float64(enemy.WalkImages[0].Bounds().Dx())+2, 0)
					}
				}
			}
		}
		spawnerOp.GeoM.Translate(0, float64(yAdjust)+2)
	}

}
