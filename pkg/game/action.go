package game

// EntityAction is the interface to represent an entity's desired action.
type EntityAction interface {
	// Replaceble returns if the action can be immediately replaced by another.
	Replaceable() bool
	// Complete returns if the action is completed. If it is, then it is removed by the entity.
	Complete() bool
}

// EntityActionMove represents an action that should move the entity towards a target location.
type EntityActionMove struct {
	// x and y represents the target position to move to.
	x, y float64
	// relative represents if the movement is considered as relative to the entity's current position. Should this even be a thing?
	relative bool
	//
	complete bool
}

func (a *EntityActionMove) Replaceable() bool {
	return true
}

func (a *EntityActionMove) Complete() bool {
	return a.complete
}
