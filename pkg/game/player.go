package game

// Player represents a player that controls an entity. It handles input and makes the entity dance.
type Player struct {
	// entity is the player-controlled entity.
	entity Entity
}

func (p *Player) Update() error {
	return nil
}
