package game

import "github.com/hajimehoshi/ebiten/v2"

type Entity interface {
	Physics() *PhysicsObject
	Trashed() bool
	Update() error
	Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) // Eh, might as well allow the entities to draw themselves.
}

type BaseEntity struct {
	physics PhysicsObject
	trashed bool
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
