package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyEntity struct {
	BaseEntity
	animation Animation
}

func NewEnemyEntity(polarity Polarity) *EnemyEntity {
	images := []*ebiten.Image{
		enemyPositive1Image,
		enemyPositive2Image,
	}
	if polarity == NegativePolarity {
		images = []*ebiten.Image{
			enemyNegative1Image,
			enemyNegative2Image,
		}
	}
	return &EnemyEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{
				polarity:       polarity,
				magnetic:       true,
				magnetStrength: 1,
				magnetRadius:   100,
			},
		},
		animation: Animation{
			images:    images,
			frameTime: 60,
			speed:     0.25,
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
