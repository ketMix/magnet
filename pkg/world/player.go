package world

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebijam22/pkg/data"
)

// Player represents a player that controls an entity. It handles input and makes the entity dance.
type Player struct {
	//
	Local bool
	// entity is the player-controlled entity.
	Entity Entity
	// I suppose the toolbelt should be here.
	Toolbelt Toolbelt
	// ReadyForWave means the players are done building and ready to start the waves.
	ReadyForWave bool
	// Name is acquired from the initial connection name.
	Name string
	//
	HoveringPlacement     bool
	HoveringPlace         EntityActionPlace
	HoverColumn, HoverRow int // X and Y hover coordinate in terms of columns/rows
	// Current points the player has.
	Points int
}

func NewPlayer() *Player {
	// Hehehe
	items := []*ToolbeltItem{
		{tool: ToolGun, key: ebiten.Key1},
	}

	// Collect our toolbelt items.
	var toolbeltItems []data.EntityConfig
	for _, v := range data.TurretConfigs {
		toolbeltItems = append(toolbeltItems, v)
	}

	// Sort them.
	sort.SliceStable(toolbeltItems, func(i, j int) bool {
		return toolbeltItems[i].ToolbeltOrder < toolbeltItems[j].ToolbeltOrder
	})

	i := 2
	for _, v := range toolbeltItems {
		items = append(items, &ToolbeltItem{
			tool: ToolTurret, key: ebiten.Key0 + ebiten.Key(i), polarity: data.NegativePolarity, kind: v,
		})
		i++
	}

	items = append(items, &ToolbeltItem{tool: ToolWall, key: ebiten.Key0 + ebiten.Key(i)})
	i++
	items = append(items, &ToolbeltItem{tool: ToolDestroy, key: ebiten.Key0 + ebiten.Key(i)})

	return &Player{
		Toolbelt: Toolbelt{
			items: items,
		},
	}
}

// Arglebargle
func (p *Player) Update(w *World) (EntityAction, error) {
	// Increment turret tick
	if p.Entity != nil {
		p.Entity.Turret().Tick(w.Speed)
	}

	// Just bail if this player is not a local entity.
	if !p.Local {
		return nil, nil
	}

	// FIXME: This should be only be called when the window is changed.
	p.Toolbelt.Position()

	// Do _not_ handle inputs if the window is not focused.
	if !ebiten.IsFocused() {
		return nil, nil
	}

	// Handle our toolbelt first.
	if req := p.Toolbelt.Update(); req != nil {
		return nil, nil
	}

	cx, cy := w.GetCursorPosition()
	tx, ty := w.GetClosestCellPosition(cx, cy)
	p.HoverColumn = tx
	p.HoverRow = ty

	if p.Entity != nil {
		var action EntityAction
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// TODO: Show placement preview
			p.HoveringPlacement = true
			p.HoveringPlace = EntityActionPlace{
				X:        tx,
				Y:        ty,
				Kind:     p.Toolbelt.activeItem.kind.Title,
				Tool:     p.Toolbelt.activeItem.tool,
				Polarity: p.Toolbelt.activeItem.polarity,
			}
		} else if p.HoveringPlacement {
			p.HoveringPlacement = false
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			// Right-click to delete.
			cx, cy := w.GetCursorPosition()
			tx, ty := w.GetClosestCellPosition(cx, cy)
			action = &EntityActionMove{
				X:        float64(tx)*float64(data.CellWidth) + float64(data.CellWidth)/2,
				Y:        float64(ty+1)*float64(data.CellHeight) + float64(data.CellHeight)/2,
				Distance: 8,
				// We wrap the place action as a move action's next step.
				Next: &EntityActionPlace{
					X:    tx,
					Y:    ty,
					Tool: ToolDestroy,
				},
			}
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && p.Toolbelt.activeItem.tool == ToolGun {
			if p.Toolbelt.activeItem.polarity == data.NeutralPolarity {
				p.Entity.Turret().rate = p.Entity.Turret().defaultRate * 2
			} else {
				p.Entity.Turret().rate = p.Entity.Turret().defaultRate
			}
			// Check if we can fire
			cx, cy := w.GetCursorPosition()
			if p.Entity.Turret().CanFire(w.Speed) {
				action = &EntityActionShoot{
					TargetX:  float64(cx),
					TargetY:  float64(cy),
					Polarity: p.Toolbelt.activeItem.polarity,
				}
			}

		} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			// Send turret placement request at the cell closest to the mouse.
			cx, cy := w.GetCursorPosition()
			tx, ty := w.GetClosestCellPosition(cx, cy)

			switch p.Toolbelt.activeItem.tool {
			case ToolTurret:
				fallthrough
			case ToolWall:
				fallthrough
			case ToolDestroy:
				action = &EntityActionMove{
					X:        float64(tx)*float64(data.CellWidth) + float64(data.CellWidth)/2,
					Y:        float64(ty)*float64(data.CellHeight) + float64(data.CellHeight)/2,
					Distance: 8,
					// We wrap the place action as a move action's next step.
					Next: &EntityActionPlace{
						X:        tx,
						Y:        ty,
						Kind:     p.Toolbelt.activeItem.kind.Title,
						Tool:     p.Toolbelt.activeItem.tool,
						Polarity: p.Toolbelt.activeItem.polarity,
					},
				}
			}
		} else if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyD) {
			// Sloppy/lazy keyboard movement. FIXME: We should probably abstract this out to a keybinds system where a slice of keys can be matched to make a "command". This command would automatically be added to some sort of current commands queue that would then be used to generate the appropriate player->entity action.
			x := 0.0
			y := 0.0
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				x--
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				x++
			}
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				y--
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				y++
			}

			action = &EntityActionMove{
				X:        p.Entity.Physics().X + x,
				Y:        p.Entity.Physics().Y + y,
				Distance: 0.5,
			}
		}
		if action != nil && (p.Entity.Action() == nil || p.Entity.Action().Replaceable()) {
			// TODO: Add a "chainable" action field that will instead add a new action as the next action in the deepest nested next action.
			//p.Entity.SetAction(action)
			return action, nil
		}
	}

	return nil, nil
}
