package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type CoreEntity struct {
	BaseEntity
	healthBar *ProgressBar
	id        int // Simplified core ID that can be shared between clients, as it is based on map construction.
	destroyed bool
}

func NewCoreEntity(config data.EntityConfig) *CoreEntity {
	return &CoreEntity{
		BaseEntity: BaseEntity{
			health:    10,
			maxHealth: 10,
			physics: PhysicsObject{
				radius: 5,
			},
			animation: Animation{
				images:    config.Images,
				frameTime: 5,
				speed:     1,
			},
		},
		healthBar: NewProgressBar(
			10, 2, color.RGBA{32, 255, 32, 1},
		),
	}
}

func (e *CoreEntity) Update(world *World) (request Request, err error) {
	e.animation.Update()

	// Update healthbar
	e.healthBar.progress = float64(e.maxHealth) / float64(e.health)

	return request, nil
}

func (e *CoreEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	// Offset the obelisk/core.
	op.GeoM.Translate(0, -float64(e.animation.Image().Bounds().Dy()/3))

	e.animation.Draw(screen, op)

	// Draw from center.
	op.GeoM.Translate(
		-float64(e.animation.Image().Bounds().Dx())/2,
		0,
	)

	// Draw healthbar
	e.healthBar.Draw(screen, op)
}
