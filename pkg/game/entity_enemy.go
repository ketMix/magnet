package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyEntity struct {
	BaseEntity
	animationTick   float64
	animationSpeed  float64
	animationIndex  int
	animationImages []*ebiten.Image
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
		animationSpeed:  0.25,
		animationImages: images,
	}
}

func (e *EnemyEntity) Update(world *World) (request Request, err error) {
	// increment animation tick
	e.animationTick++

	// Attempt to move along path to player's core

	return request, nil
}

func (e *EnemyEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	// Progress animation index
	if e.animationTick >= (e.animationSpeed * 60) {
		e.animationTick = 0
		e.animationIndex += 1
		if e.animationIndex >= len(e.animationImages) {
			e.animationIndex = 0
		}
	}

	// Set image from animation index
	image := e.animationImages[e.animationIndex]

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	// Draw from center.
	// TODO: set enemy image on instatiation?
	op.GeoM.Translate(
		-float64(image.Bounds().Dx())/2,
		-float64(image.Bounds().Dy())/2,
	)
	screen.DrawImage(image, op)
}
