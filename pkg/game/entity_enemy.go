package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyEntity struct {
	BaseEntity
}

func NewEnemyEntity(config EntityConfig) *EnemyEntity {
	return &EnemyEntity{
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
			},
		},
	}
}

func (e *EnemyEntity) Update(world *World) (request Request, err error) {
	// Update animation.
	e.animation.Update()

	// Attempt to move along path to player's core

	return request, nil
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
