package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type TurretEntity struct {
	BaseEntity
	attackRadius float64
	turret       Turret
	target       Entity
	// owner ActorEntity // ???
}

func NewTurretEntity() *TurretEntity {
	return &TurretEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{
				polarity: NeutralPolarity,
			},
		},
		turret: Turret{
			speed: 1,
			rate:  1,
		},
		attackRadius: 100,
	}
}

func (e *TurretEntity) Update(world *World) (request Request, err error) {
	// Tick the turret
	e.turret.Tick()

	// Attempt to acquire target
	e.AcquireTarget(&world.entities)

	// Make request to fire if we have target and can fire
	if e.target != nil && e.turret.CanFire() {
		px, py := e.physics.X, e.physics.Y
		tx, ty := e.target.Physics().X, e.target.Physics().Y

		vX, vY := GetDirection(px, py, tx, ty)
		request = SpawnProjecticleRequest{
			x:        px,
			y:        py,
			vX:       vX * e.turret.speed,
			vY:       vY * e.turret.speed,
			polarity: e.physics.polarity,
		}
	}

	return request, nil
}

func (e *TurretEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)
	// Draw from center.
	op.GeoM.Translate(
		-float64(turretBaseImage.Bounds().Dx())/2,
		-float64(turretBaseImage.Bounds().Dy())/2,
	)
	screen.DrawImage(turretBaseImage, op)
}

// Finds the closest entity within attack radius and sets the current target if found
// !! Iterates through world entity list, could probably be optimized !!
func (e *TurretEntity) AcquireTarget(entities *[]Entity) {
	// hmm...
	minDistance := 100000.0
	var target Entity

	x, y := e.physics.X, e.physics.Y
	for _, entity := range *entities {
		switch entity.(type) {
		case *EnemyEntity:
			tx, ty := entity.Physics().X, entity.Physics().Y
			if IsWithinRadius(x, y, tx, ty, e.attackRadius) {
				magnitude := GetMagnitude(GetDistanceVector(x, y, tx, ty))
				if magnitude < minDistance {
					minDistance = magnitude
					target = entity
				}
			}
		}
	}
	e.target = target
}
