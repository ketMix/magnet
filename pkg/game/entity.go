package game

import "github.com/hajimehoshi/ebiten/v2"

type Entity interface {
	Physics() *PhysicsObject
	Trashed() bool
	Trash()
	Update() (Request, error)
	Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) // Eh, might as well allow the entities to draw themselves.
	Action() EntityAction
	SetAction(a EntityAction)
}

type BaseEntity struct {
	physics PhysicsObject
	trashed bool
	action  EntityAction
}

func (e *BaseEntity) Physics() *PhysicsObject {
	return &e.physics
}

func (e *BaseEntity) Trashed() bool {
	return e.trashed
}

func (e *BaseEntity) Trash() {
	e.trashed = true
}

func (e *BaseEntity) Action() EntityAction {
	return e.action
}

func (e *BaseEntity) SetAction(a EntityAction) {
	e.action = a
}
