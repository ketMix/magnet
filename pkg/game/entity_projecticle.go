package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ProjecticleEntity struct {
	BaseEntity
	elapsed   int
	lifetime  int
	animation Animation
}

func NewProjecticleEntity() *ProjecticleEntity {
	return &ProjecticleEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{},
		},
		lifetime: 500, // Make the default lifetime 500 ticks. This should be set to a value that makes sense for the projectile's speed so it remains alive for however long it needs to.
		animation: Animation{
			images: []*ebiten.Image{
				projecticlePositiveImage,
				projecticleNegativeImage,
				projecticleNeutralImage,
			},
		},
	}
}

func (e *ProjecticleEntity) Update(world *World) (request Request, err error) {
	e.elapsed++

	// Set the current animation frame index manually.
	switch e.physics.polarity {
	case PositivePolarity:
		e.animation.index = 0
	case NegativePolarity:
		e.animation.index = 1
	case NeutralPolarity:
		e.animation.index = 2
	}

	// If our projecticle is magnetic, we need to potentially update projecticle vector
	if e.physics.polarity != NeutralPolarity {
		// Grab set of physics objects from entities where projecticle collides with magnet radius
		// For each collision
		//  - get magnetic vector
		//  - add to initial vector
		for _, entity := range world.entities {
			if entity.IsCollided(e.BaseEntity) {
				e.Trash()
			}
			if entity.IsWithinMagneticField(e.BaseEntity) {
				mX, mY := entity.Physics().GetMagneticVector(e.physics)
				e.physics.vX = e.physics.vX + mX
				e.physics.vY = e.physics.vY + mY
			}
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

	e.animation.Draw(screen, op)
}
