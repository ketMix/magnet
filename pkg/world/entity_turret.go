package world

import (
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type TurretEntity struct {
	BaseEntity
	target Entity
	// owner ActorEntity // ???
	headAnimation Animation
	cost          int
}

func NewTurretEntity(config data.EntityConfig) *TurretEntity {
	return &TurretEntity{
		BaseEntity: BaseEntity{
			animation: Animation{
				images: config.Images,
			},
			physics: PhysicsObject{
				polarity: config.Polarity,
			},
			turret: Turret{
				damage:      config.Damage,
				speed:       config.ProjecticleSpeed,
				rate:        config.AttackRate,
				attackRange: config.AttackRange,
			},
		},
		headAnimation: Animation{
			images: config.HeadImages,
		},
		cost: config.Points,
	}
}

func (e *TurretEntity) Update(world *World) (request Request, err error) {
	// Tick the turret
	e.turret.Tick(world.Speed)

	// Attempt to acquire target
	e.AcquireTarget(world)

	// Rotate head to face target (or just SPEEN)
	if e.target == nil {
		e.headAnimation.rotation += 0.02 * world.Speed
	} else {
		e.headAnimation.rotation = math.Atan2(e.physics.Y-e.target.Physics().Y, e.physics.X-e.target.Physics().X)
	}

	// Make request to fire if we have target and can fire
	if e.target != nil && !e.target.Trashed() && e.turret.CanFire(world.Speed) {
		px, py := e.physics.X, e.physics.Y
		tx, ty := e.target.Physics().X, e.target.Physics().Y

		vX, vY := GetDirection(px, py, tx, ty)

		vX *= world.Speed
		vY *= world.Speed

		projectile := &ProjecticleEntity{
			BaseEntity: BaseEntity{
				physics: PhysicsObject{
					vX:       vX * e.turret.speed,
					vY:       vY * e.turret.speed,
					polarity: e.physics.polarity,
				},
			},
			lifetime: 500,
			damage:   e.turret.damage,
		}

		request = SpawnProjecticleRequest{
			x:          px,
			y:          py,
			projectile: projectile,
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

	headColor := ebiten.ColorM{}
	headColor.Scale(data.GetPolarityColorScale(e.physics.polarity))

	// Draw da head
	op.GeoM.Translate(0, -5)
	for i := float64(0); i < 3; i++ {
		darken := .25 + i - (i / 3)
		headOp := &ebiten.DrawImageOptions{}
		headOp.GeoM.Concat(op.GeoM)
		headOp.GeoM.Translate(0, -i)
		headOp.ColorM.Scale(darken, darken, darken, 1)
		headOp.ColorM.Concat(headColor)
		e.headAnimation.Draw(screen, headOp)
	}
}

// Finds the closest entity within attack radius and sets the current target if found
func (e *TurretEntity) AcquireTarget(world *World) {
	// Collect our entities within our attack radius.
	entities := ObjectsWithinRadius(world.enemies, e.physics.X, e.physics.Y, e.turret.attackRange)

	// Sort from closest to further. This is a bit inefficient but I don't care.
	sort.Slice(entities, func(i, j int) bool {
		a := GetMagnitude(GetDistanceVector(e.physics.X, e.physics.Y, entities[i].Physics().X, entities[i].Physics().Y))
		b := GetMagnitude(GetDistanceVector(e.physics.X, e.physics.Y, entities[j].Physics().X, entities[j].Physics().Y))
		return a < b
	})

	// Set it to the first entry, as it should be the closest.
	if len(entities) > 0 {
		e.target = entities[0]
	} else {
		e.target = nil
	}
}
