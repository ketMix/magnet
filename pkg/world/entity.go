package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/goro/pathing"
)

type Entity interface {
	Physics() *PhysicsObject
	Turret() *Turret
	Animation() *Animation
	Trashed() bool
	Trash()
	Update(world *World) (Request, error)
	Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) // Eh, might as well allow the entities to draw themselves.
	Action() EntityAction
	SetAction(a EntityAction)
	IsCollided(t Entity) bool
	IsWithinMagneticField(t Entity) bool
	// Why not.
	CanPathfind() bool
	SetSteps(p []pathing.Step)
	NetID() int
	SetNetID(nid int)
	IsProjectile() bool
}

type BaseEntity struct {
	physics   PhysicsObject
	trashed   bool
	action    EntityAction
	turret    Turret
	health    int
	maxHealth int
	animation Animation
	netID     int
}

func (e *BaseEntity) Physics() *PhysicsObject {
	return &e.physics
}

func (e *BaseEntity) Turret() *Turret {
	return &e.turret
}

func (e *BaseEntity) Animation() *Animation {
	return &e.animation
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

func (e *BaseEntity) CanPathfind() bool {
	return false
}

func (e *BaseEntity) SetSteps(s []pathing.Step) {
}

func (e *BaseEntity) NetID() int {
	return e.netID
}

func (e *BaseEntity) SetNetID(nid int) {
	e.netID = nid
}

func (e *BaseEntity) IsProjectile() bool {
	return false
}

// Check whether or not the provided entity collides
// Should probably use entity sprites or add hitboxes to physics object
func (e *BaseEntity) IsCollided(t Entity) bool {
	x, y := e.physics.X+e.physics.OffsetX, e.physics.Y+e.physics.OffsetY
	tx, ty := t.Physics().X+t.Physics().OffsetX, t.Physics().Y+t.Physics().OffsetY
	return IsWithinRadius(x, y, tx, ty, t.Physics().radius)
}

// Check whether or not the provided entity is within magnetic field
func (e *BaseEntity) IsWithinMagneticField(t Entity) bool {
	// Can't be within field if ya aren't got the POWER
	if !e.physics.magnetic {
		return false
	}

	// Should probably extend this from center or edge of entity's sprite/hitbox
	x, y := e.physics.X, e.physics.Y
	tx, ty := t.Physics().X, t.Physics().Y
	return IsWithinRadius(x, y, tx, ty, e.physics.magnetRadius)
}
