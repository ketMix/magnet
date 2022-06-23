package world

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/goro/pathing"
)

type EnemyEntity struct {
	BaseEntity
	steps     []pathing.Step
	healthBar *ProgressBar
	speed     float64
	lifetime  float64
}

func NewEnemyEntity(config data.EntityConfig) *EnemyEntity {
	return &EnemyEntity{
		BaseEntity: BaseEntity{
			animation: Animation{
				images:    config.Images,
				frameTime: 30,
				speed:     1 - config.Speed,
			},
			health:    config.Health,
			maxHealth: config.Health,
			physics: PhysicsObject{
				polarity:       config.Polarity,
				magnetic:       config.Magnetic,
				magnetStrength: config.MagnetStrength,
				magnetRadius:   config.MagnetRadius,
				radius:         config.Radius,
			},
		},
		healthBar: NewProgressBar(
			7,
			1,
			color.RGBA{255, 0, 0, 1},
		),
		speed: config.Speed,
	}
}

func (e *EnemyEntity) Update(world *World) (request Request, err error) {
	if e.health <= 0 {
		request = TrashEntityRequest{
			NetID:  e.netID,
			entity: e,
			local:  true,
		}
		return
	}
	e.lifetime++
	// Update animation.
	e.animation.Update()

	// Update healthbar
	e.healthBar.progress = float64(e.maxHealth) / float64(e.health)

	// Attempt to move along path to player's core
	if len(e.steps) != 0 {
		tx := float64(e.steps[0].X()*data.CellWidth + data.CellWidth/2)
		ty := float64(e.steps[0].Y()*data.CellHeight + data.CellHeight/2)
		if math.Abs(e.physics.X-float64(e.steps[0].X()*data.CellWidth+data.CellWidth/2)) < 1 && math.Abs(e.physics.Y-float64(e.steps[0].Y()*data.CellHeight+data.CellHeight/2)) < 1 {
			//
			e.steps = e.steps[1:]
		} else {
			r := math.Atan2(e.physics.Y-ty, e.physics.X-tx)
			x := math.Cos(r) * e.speed
			y := math.Sin(r) * e.speed

			e.physics.X -= x
			e.physics.Y -= y

			if x > 0 {
				e.animation.mirror = false
			} else {
				e.animation.mirror = true
			}
		}

		// TODO: move towards step[0], then remove it when near its center. If the last one is to be removed, then we have reached the core.
	} else {
		// No mo steppes
		request = TrashEntityRequest{
			NetID:  e.netID,
			entity: e,
			local:  true,
		}
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

	if e.lifetime < 10 {
		op.ColorM.Scale(1, 1, 1, e.lifetime/10)
	}

	if e.physics.polarity != data.NeutralPolarity {
		r, g, b, a := data.GetPolarityColorScale(e.physics.polarity)
		op.ColorM.Scale(r*2, g*2, b*2, a)
	}

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
