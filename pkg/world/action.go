package world

import "github.com/kettek/ebijam22/pkg/data"

// EntityAction is the interface to represent an entity's desired action.
type EntityAction interface {
	// Replaceble returns if the action can be immediately replaced by another.
	Replaceable() bool
	// Complete returns if the action is completed. If it is, then it is removed by the entity.
	Complete() bool
	Next() EntityAction
}

// EntityActionMove represents an action that should move the entity towards a target location.
type EntityActionMove struct {
	// x and y represents the target position to move to.
	x, y float64
	// distance represents the distance from the target that should be considered valid.
	distance float64
	// relative represents if the movement is considered as relative to the entity's current position. Should this even be a thing?
	relative bool
	//
	complete bool
	//
	next EntityAction
}

func (a *EntityActionMove) Replaceable() bool {
	return true
}

func (a *EntityActionMove) Complete() bool {
	return a.complete
}

func (a *EntityActionMove) Next() EntityAction {
	return a.next
}

type EntityActionPlace struct {
	// x and y are cell positions to place at.
	x, y     int
	complete bool
	kind     ToolKind
	polarity data.Polarity
}

func (a *EntityActionPlace) Replaceable() bool {
	return true
}

func (a *EntityActionPlace) Complete() bool {
	return a.complete
}

func (a *EntityActionPlace) Next() EntityAction {
	return nil
}

type EntityActionShoot struct {
	// targetX and targetY represent the position to fire at.
	targetX, targetY float64
	polarity         data.Polarity
	complete         bool
	next             EntityAction
}

func (a *EntityActionShoot) Replaceable() bool {
	return false
}

func (a *EntityActionShoot) Complete() bool {
	return a.complete
}

func (a *EntityActionShoot) Next() EntityAction {
	return a.next
}
