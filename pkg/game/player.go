package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Player represents a player that controls an entity. It handles input and makes the entity dance.
type Player struct {
	// entity is the player-controlled entity.
	entity Entity
}

func NewPlayer() *Player {
	return &Player{}
}

// It's kind of weird to pass the play state, but oh well.
func (p *Player) Update(s *PlayState) error {
	if p.entity != nil {
		var action EntityAction
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			// TODO: Make target the center of the closest intersecting cell.
			action = &EntityActionMove{
				x: float64(x) - s.cameraX,
				y: float64(y) - s.cameraY,
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
				x: p.entity.Physics().X + x,
				y: p.entity.Physics().Y + y,
			}
		}
		if action != nil && (p.entity.Action() == nil || p.entity.Action().Replaceable()) {
			p.entity.SetAction(action)
		}
	}

	return nil
}
