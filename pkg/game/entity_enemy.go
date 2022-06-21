package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/goro/pathing"
)

type EnemyEntity struct {
	BaseEntity
	steps     []pathing.Step
	healthBar *ProgressBar
	speed     float64
}

func NewEnemyEntity(config EntityConfig) *EnemyEntity {
	return &EnemyEntity{
		BaseEntity: BaseEntity{
			animation: Animation{
				images:    config.images,
				frameTime: 60,
				speed:     0.25,
			},
			health:    config.health,
			maxHealth: config.health,
			physics: PhysicsObject{
				polarity:       config.polarity,
				magnetic:       config.magnetic,
				magnetStrength: config.magnetStrength,
				magnetRadius:   config.magnetRadius,
				radius:         config.radius,
			},
		},
		healthBar: NewProgressBar(
			7,
			1,
			color.RGBA{255, 0, 0, 1},
		),
		speed: config.speed,
	}
}

func (e *EnemyEntity) Update(world *World) (request Request, err error) {
	if e.health <= 0 {
		e.Trash()
		return
	}
	// Update animation.
	e.animation.Update()

	// Update healthbar
	e.healthBar.progress = float64(e.maxHealth) / float64(e.health)

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

	// Draw healthbar if less than max health
	if e.health < e.maxHealth {
		// Center the health bar horizontally and position it at the bottom of our image.
		// NOTE: we're using the first image, as if we use whatever the current frame is, it might be of different dimensions, which would lead to the health bar changing position.
		op.GeoM.Translate(-float64(e.animation.images[0].Bounds().Dx())/2, float64(e.animation.images[0].Bounds().Dy()))
		e.healthBar.Draw(screen, op)
	}
}
