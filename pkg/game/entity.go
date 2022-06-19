package game

import "github.com/hajimehoshi/ebiten/v2"

type Entity interface {
	Physics() *PhysicsObject
	Trashed() bool
	Trash()
	Update(world *World) (Request, error)
	Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) // Eh, might as well allow the entities to draw themselves.
	Action() EntityAction
	SetAction(a EntityAction)
	IsCollided(t BaseEntity) bool
	IsWithinMagneticField(t BaseEntity) bool
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

// Check whether or not the provided entity collides
// Should probably use entity sprites or add hitboxes to physics object
func (e *BaseEntity) IsCollided(t BaseEntity) bool {
	hitboxRadius := 0.0 // set this here for now for testing
	x, y := e.physics.X, e.physics.Y
	tx, ty := t.physics.X, t.physics.Y
	return IsWithinRadius(x, y, tx, ty, hitboxRadius)
}

// Check whether or not the provided entity is within magnetic field
func (e *BaseEntity) IsWithinMagneticField(t BaseEntity) bool {
	// Can't be within field if ya aren't got the POWER
	if !e.physics.magnetic {
		return false
	}

	// Should probably extend this from center or edge of entity's sprite/hitbox
	x, y := e.physics.X, e.physics.Y
	tx, ty := t.physics.X, t.physics.Y
	return IsWithinRadius(x, y, tx, ty, e.physics.magnetRadius)
}
