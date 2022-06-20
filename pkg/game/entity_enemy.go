package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/goro/pathing"
)

type EnemyEntity struct {
	BaseEntity
	steps []pathing.Step
	speed float64
}

func NewEnemyEntity(config EntityConfig) *EnemyEntity {
	return &EnemyEntity{
		speed: config.speed,
		BaseEntity: BaseEntity{
			animation: Animation{
				images:    config.images,
				frameTime: 60,
				speed:     0.25,
			},
			health: config.health,
			physics: PhysicsObject{
				polarity:       config.polarity,
				magnetic:       config.magnetic,
				magnetStrength: config.magnetStrength,
				magnetRadius:   config.magnetRadius,
				radius:         config.radius,
			},
		},
	}
}

func (e *EnemyEntity) Update(world *World) (request Request, err error) {
	if e.health <= 0 {
		e.Trash()
		return
	}
	// Update animation.
	e.animation.Update()

	// Attempt to move along path to player's core
	if len(e.steps) != 0 {
		tx := float64(e.steps[0].X()*cellWidth + cellWidth/2)
		ty := float64(e.steps[0].Y()*cellHeight + cellHeight/2)
		if math.Abs(e.physics.X-float64(e.steps[0].X()*cellWidth+cellWidth/2)) < 1 && math.Abs(e.physics.Y-float64(e.steps[0].Y()*cellHeight+cellHeight/2)) < 1 {
			//
			e.steps = e.steps[1:]
		} else {
			r := math.Atan2(e.physics.Y-ty, e.physics.X-tx)
			x := math.Cos(r) * e.speed
			y := math.Sin(r) * e.speed

			e.physics.X -= x
			e.physics.Y -= y
		}

		// TODO: move towards step[0], then remove it when near its center. If the last one is to be removed, then we have reached the core.
	} else {
		// No mo steppes
		e.Trash()
	}

	return request, nil
}

func (e *EnemyEntity) CanPathfind() bool {
	return true
}

func (e *EnemyEntity) SetSteps(s []pathing.Step) {
	e.steps = s
}

func (e *EnemyEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	// Draw animation.
	e.animation.Draw(screen, op)
}
