package world

import (
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type TurretEntity struct {
	BaseEntity
	target          Entity
	colorMultiplier [3]float64 // Color multiplier, passed in when in multiplayer.
	owner           string
	headAnimation   Animation
	cost            int
	showRange       bool
	polarizer       bool
	locked          bool // locked is used to lock the entity when the mode changes to a loss.
	lockedTicker    int
}

func NewTurretEntity(config data.EntityConfig) *TurretEntity {
	polarizer := false
	if config.AttackType == "polarizer" {
		polarizer = true
	}
	return &TurretEntity{
		polarizer: polarizer,
		BaseEntity: BaseEntity{
			animation: Animation{
				images: config.Images,
			},
			physics: PhysicsObject{
				polarity:       config.Polarity,
				magnetic:       config.Magnetic,
				magnetStrength: config.MagnetStrength,
				magnetRadius:   config.MagnetRadius,
				radius:         config.Radius,
			},
			turret: Turret{
				damage:         config.Damage,
				speed:          config.ProjecticleSpeed,
				rate:           config.AttackRate,
				attackRange:    config.AttackRange,
				projecticleNum: config.ProjecticleNum,
			},
		},
		headAnimation: Animation{
			images: config.HeadImages,
		},
		colorMultiplier: [3]float64{1, 1, 1},
		cost:            config.Points,
	}
}

func (e *TurretEntity) Update(world *World) (request Request, err error) {
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

	// Attempt to acquire target
	e.AcquireTarget(world)

	// Rotate head to face target (or just SPEEN)
	if e.target == nil || e.target.Trashed() {
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

		const spreadArc = 45.0
		var vectors = SplitVectorByDegree(spreadArc, vX, vY, e.turret.projecticleNum)
		var projecticleRequests MultiRequest
		for _, v := range vectors {
			request := SpawnProjecticleRequest{
				X:        px,
				Y:        py,
				VX:       v.vX * e.turret.speed,
				VY:       v.vY * e.turret.speed,
				Polarity: e.physics.polarity,
				Damage:   e.turret.damage,
			}
			projecticleRequests.Requests = append(projecticleRequests.Requests, request)
		}
		request = projecticleRequests
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

	op.ColorM.Scale(e.colorMultiplier[0], e.colorMultiplier[1], e.colorMultiplier[2], 1)

	DrawTurret(screen, op, e.animation, e.headAnimation, e.physics.polarity)

	// This is temporary, as all things in life are.
	if e.showRange {
		r, g, b, a := data.GetPolarityColorScale(e.physics.polarity)
		drawCircle(screen, op, int(e.turret.attackRange), r, g, b, a)
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

func DrawTurret(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions, bodyAnimation Animation, headAnimation Animation, polarity data.Polarity) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.ColorM.Concat(screenOp.ColorM)

	bodyAnimation.Draw(screen, op)

	headColor := ebiten.ColorM{}
	headColor.Concat(screenOp.ColorM)
	headColor.Scale(data.GetPolarityColorScale(polarity))

	// Draw da head
	op.GeoM.Translate(0, -5)
	for i := float64(0); i < 3; i++ {
		darken := .25 + i - (i / 3)
		headOp := &ebiten.DrawImageOptions{}
		headOp.GeoM.Concat(op.GeoM)
		headOp.GeoM.Translate(0, -i)
		headOp.ColorM.Scale(darken, darken, darken, 1)
		headOp.ColorM.Concat(headColor)
		headAnimation.Draw(screen, headOp)
	}
}
