package world

import (
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kettek/ebijam22/pkg/data"
)

type TurretBeamEntity struct {
	TurretEntity
	beamTick int
}

func NewTurretBeamEntity(config data.EntityConfig) *TurretBeamEntity {
	return &TurretBeamEntity{
		TurretEntity: *NewTurretEntity(config),
	}
}

func (e *TurretBeamEntity) Update(world *World) (request Request, err error) {
	if e.locked {
		e.lockedTicker++
		if e.lockedTicker%60 < 30 {
			e.headAnimation.rotation += 0.02 * world.Speed
		} else {
			e.headAnimation.rotation -= 0.02 * world.Speed
		}
		return

	}
	// Tick the turret
	e.turret.Tick(world.Speed)

	// Rotate head to face target (or just SPEEN)
	if e.target == nil || e.target.Trashed() {
		e.headAnimation.rotation += 0.02 * world.Speed
		e.AcquireTarget(world)
	} else {
		e.headAnimation.rotation = math.Atan2(e.physics.Y-e.target.Physics().Y, e.physics.X-e.target.Physics().X)
	}

	// Make request to fire if we have target and can fire
	if e.target != nil && !e.target.Trashed() {
		e.beamTick++
		// Damage?
		if e.turret.CanFire(world.Speed) {
			if e2, ok := e.target.(*EnemyEntity); ok {
				data.SFX.Play("turret-beam.ogg")
				e2.health -= e.turret.damage
			}
		}
	}

	return request, nil
}

// Finds the closest entity of same polarity within attack radius and sets the current target if found.
func (e *TurretBeamEntity) AcquireTarget(world *World) {
	// Collect our entities within our attack radius.
	entities := ObjectsWithinRadius(world.enemies, e.physics.X, e.physics.Y, e.turret.attackRange)

	// Always target our own polarity.
	if e.physics.polarity == data.NegativePolarity {
		entities = ObjectsWithPolarity(entities, data.PositivePolarity)
	} else {
		entities = ObjectsWithPolarity(entities, data.NegativePolarity)
	}

	// Sort from furthest to closest, with a priority for low health targets. This is a bit inefficient but I don't care.
	sort.Slice(entities, func(i, j int) bool {
		a := float64(entities[i].health)
		b := float64(entities[j].health)
		return a < b
	})

	// Set it to the first entry, as it should be the closest.
	if len(entities) > 0 {
		e.target = entities[0]
	} else {
		e.target = nil
	}
}

func (e *TurretBeamEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	e.TurretEntity.Draw(screen, screenOp)

	if e.target != nil && !e.target.Trashed() {
		x := screenOp.GeoM.Element(0, 2)
		y := screenOp.GeoM.Element(1, 2) - 6 // -6 for the head offset
		c := data.GetPolarityColor(e.physics.polarity)
		c.A = uint8(100 + math.Sin(float64(e.beamTick))*255)
		// Draw that beam.
		ebitenutil.DrawLine(screen, x+e.physics.X, y+e.physics.Y, x+e.target.Physics().X, y+e.target.Physics().Y+4, c)
	}

}
