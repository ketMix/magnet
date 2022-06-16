package game

import "github.com/hajimehoshi/ebiten/v2"

// Player represents a player that controls an entity. It handles input and makes the entity dance.
type Player struct {
	// entity is the player-controlled entity.
	entity Entity
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
		}
		if action != nil && (p.entity.Action() == nil || p.entity.Action().Replaceable()) {
			p.entity.SetAction(action)
		}
	}

	return nil
}
