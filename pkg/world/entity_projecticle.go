package world

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type TouchContainer struct {
	entity    Entity
	count     int
	touchTime time.Time
}

type ProjecticleEntity struct {
	BaseEntity
	elapsed         int
	lifetime        int
	damage          int
	touchedEntities []TouchContainer
}

func NewProjecticleEntity() *ProjecticleEntity {
	return &ProjecticleEntity{
		lifetime: 500, // Make the default lifetime 500 ticks. This should be set to a value that makes sense for the projectile's speed so it remains alive for however long it needs to.
		BaseEntity: BaseEntity{
			physics: PhysicsObject{
				radius: 4,
			},
		},
	}
}

func (e *ProjecticleEntity) Update(world *World) (request Request, err error) {
	e.elapsed++
	// Grab set of physics objects from entities where projecticle collides with magnet radius
	// For each collision
	//  - get magnetic vector
	//  - add to initial vector

	getTouch := func(entity Entity) *TouchContainer {
		for i, c := range e.touchedEntities {
			if c.entity == entity {
				return &e.touchedEntities[i]
			}
		}
		return nil
	}

	for _, entity := range world.entities {
		switch entity := entity.(type) {
		case *EnemyEntity:
			if e.IsCollided(entity) {
				entity.health -= e.damage
				e.Trash()
				break
			}
		case *TurretEntity:
			if entity.polarizer && e.IsCollided(entity) {
				t := getTouch(entity)
				if t == nil {
					e.physics.polarity = entity.physics.polarity
					e.touchedEntities = append(e.touchedEntities, TouchContainer{entity, 1, time.Now()})
				}
			} else if entity.reflector && e.IsCollided(entity) {
				t := getTouch(entity)
				if t == nil || time.Now().Sub(t.touchTime) >= time.Duration(50)*time.Millisecond {
					if math.Abs(e.physics.vX) > math.Abs(e.physics.vY) {
						e.physics.vX = -e.physics.vX
					} else {
						e.physics.vY = -e.physics.vY
					}
					if t != nil {
						t.count++
						t.touchTime = time.Now()
					} else {
						e.touchedEntities = append(e.touchedEntities, TouchContainer{entity, 1, time.Now()})
					}
				}
			}
		}

		// If our projecticle has polarity, we need to potentially update projecticle vector
		if e.physics.polarity != data.NeutralPolarity && entity.IsWithinMagneticField(e) {
			mX, mY := entity.Physics().GetMagneticVector(e.physics)
			e.physics.vX = (e.physics.vX + mX) * world.Speed
			e.physics.vY = (e.physics.vY + mY) * world.Speed
		}
	}

	// Update projecticle's position by resulting vector
	e.physics.X += e.physics.vX
	e.physics.Y += e.physics.vY

	// NOTE: We could use an offscreen oob check, but that would be based on the map width/height, which we don't want here, as it would involve passing either those dimensions on construction or having the world as a field on this entity. So, we're just using a lifetime tick counter.
	if e.elapsed >= e.lifetime {
		e.Trash()
	}
	return request, nil
}

func (e *ProjecticleEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	x1 := op.GeoM.Element(0, 2)
	y1 := op.GeoM.Element(1, 2)
	x2 := x1 + e.physics.vX
	y2 := y1 + e.physics.vY
	c := data.GetPolarityColor(e.physics.polarity)

	length := math.Hypot(x2-x1, y2-y1)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(2+length, 2)
	op2.GeoM.Rotate(math.Atan2(y2-y1, x2-x1))
	op2.GeoM.Translate(x1, y1)
	op2.ColorM.ScaleWithColor(c)
	// Filter must be 'nearest' filter (default).
	// Linear filtering would make edges blurred.
	screen.DrawImage(data.EmptySubImage, op2)
}

func (e *ProjecticleEntity) IsProjectile() bool {
	return true
}
