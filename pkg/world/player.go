package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

// Player represents a player that controls an entity. It handles input and makes the entity dance.
type Player struct {
	// entity is the player-controlled entity.
	Entity Entity
	// I suppose the toolbelt should be here.
	Toolbelt Toolbelt
}

func NewPlayer() *Player {
	return &Player{
		Toolbelt: Toolbelt{
			items: []*ToolbeltItem{
				{kind: ToolGun, key: ebiten.Key1},
				{kind: ToolTurret, key: ebiten.Key2, polarity: data.NegativePolarity},
				{kind: ToolWall, key: ebiten.Key3},
				{kind: ToolDestroy, key: ebiten.Key4},
			},
		},
	}
}

// Arglebargle
func (p *Player) Update(w *World) error {
	// FIXME: This should be only be called when the window is changed.
	p.Toolbelt.Position()

	// Increment turret tick
	p.Entity.Turret().Tick()

	// Handle our toolbelt first.
	if req := p.Toolbelt.Update(); req != nil {
		return nil
	}
	if p.Entity != nil {
		var action EntityAction
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			// Right-click to delete.
			cx, cy := w.GetCursorPosition()
			tx, ty := w.GetClosestCellPosition(cx, cy)
			action = &EntityActionMove{
				x:        float64(tx)*float64(data.CellWidth) + float64(data.CellWidth)/2,
				y:        float64(ty+1)*float64(data.CellHeight) + float64(data.CellHeight)/2,
				distance: 8,
				// We wrap the place action as a move action's next step.
				next: &EntityActionPlace{
					x:    tx,
					y:    ty,
					kind: ToolDestroy,
				},
			}
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// Send turret placement request at the cell closest to the mouse.
			cx, cy := w.GetCursorPosition()
			tx, ty := w.GetClosestCellPosition(cx, cy)

			switch p.Toolbelt.activeItem.kind {
			case ToolTurret:
				fallthrough
			case ToolWall:
				fallthrough
			case ToolDestroy:
				action = &EntityActionMove{
					x:        float64(tx)*float64(data.CellWidth) + float64(data.CellWidth)/2,
					y:        float64(ty+1)*float64(data.CellHeight) + float64(data.CellHeight)/2,
					distance: 8,
					// We wrap the place action as a move action's next step.
					next: &EntityActionPlace{
						x:        tx,
						y:        ty,
						kind:     p.Toolbelt.activeItem.kind,
						polarity: p.Toolbelt.activeItem.polarity,
					},
				}
			case ToolGun:
				// Check if we can fire
				if p.Entity.Turret().CanFire() {
					action = &EntityActionShoot{
						targetX:  float64(cx),
						targetY:  float64(cy),
						polarity: p.Toolbelt.activeItem.polarity,
					}
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
				x:        p.Entity.Physics().X + x,
				y:        p.Entity.Physics().Y + y,
				distance: 0.5,
			}
		}
		if action != nil && (p.Entity.Action() == nil || p.Entity.Action().Replaceable()) {
			// TODO: Add a "chainable" action field that will instead add a new action as the next action in the deepest nested next action.
			p.Entity.SetAction(action)
		}
	}

	return nil
}
