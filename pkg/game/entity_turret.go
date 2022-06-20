package game

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type TurretEntity struct {
	BaseEntity
	target Entity
	// owner ActorEntity // ???
}

func NewTurretEntity(config EntityConfig) *TurretEntity {
	return &TurretEntity{
		BaseEntity: BaseEntity{
			animation: Animation{
				images: config.images,
			},
			physics: PhysicsObject{
				polarity: config.polarity,
			},
			turret: Turret{
				damage:      config.damage,
				speed:       config.projecticleSpeed,
				rate:        config.attackRate,
				attackRange: config.attackRange,
			},
		},
	}
}

func (e *TurretEntity) Update(world *World) (request Request, err error) {
	// Tick the turret
	e.turret.Tick()

	// Attempt to acquire target
	e.AcquireTarget(world)

	// Make request to fire if we have target and can fire
	if e.target != nil && !e.target.Trashed() && e.turret.CanFire() {
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

	e.animation.Draw(screen, op)
}

// Finds the closest entity within attack radius and sets the current target if found
// !! Iterates through world entity list, could probably be optimized !!
func (e *TurretEntity) AcquireTarget(world *World) {
	// Collect our entities within our attack radius.
	entities := world.EntitiesWithinRadius(e.physics.X, e.physics.Y, e.turret.attackRange)

	// Filter out non-enemies.
	var filtered []Entity
	for _, entity := range entities {
		switch entity.(type) {
		case *EnemyEntity:
			filtered = append(filtered, entity)
		}
	}

	// Sort from closest to further. This is a bit inefficient but I don't care.
	sort.Slice(filtered, func(i, j int) bool {
		a := GetMagnitude(GetDistanceVector(e.physics.X, e.physics.Y, entities[i].Physics().X, entities[i].Physics().Y))
		b := GetMagnitude(GetDistanceVector(e.physics.X, e.physics.Y, entities[j].Physics().X, entities[j].Physics().Y))
		return a < b
	})

	// Set it to the first entry, as it should be the closest.
	if len(filtered) > 0 {
		e.target = filtered[0]
	} else {
		e.target = nil
	}
}
