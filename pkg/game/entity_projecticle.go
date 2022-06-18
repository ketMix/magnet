package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ProjecticleEntity struct {
	BaseEntity
	polarity Polarity
	elapsed  int
	lifetime int
}

func NewProjecticleEntity() *ProjecticleEntity {
	return &ProjecticleEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{},
		},
		lifetime: 500, // Make the default lifetime 500 ticks. This should be set to a value that makes sense for the projectile's speed so it remains alive for however long it needs to.
	}
}

func (e *ProjecticleEntity) Update() (request Request, err error) {
	e.elapsed++
	// Grab initial vector
	// Grab set of physics objects from entities where projecticle collides with magnet radius
	// For each collision
	//  - get vector direction of projecticle to object
	//  - get magnet factor by multiplying magnet strength by distance // 2 (magnetic effect drops with inverse square relationship)
	//  - add to initial vector
	// Update projecticle's position by resulting vector

	// For now, just continue in straight line
	e.physics.X += e.physics.vX
	e.physics.Y += e.physics.vY

	// NOTE: We could use an offscreen oob check, but that would be based on the map width/height, which we don't want here, as it would involve passing either those dimensions on construction or having the world as a field on this entity. So, we're just using a lifetime tick counter.
	if e.elapsed >= e.lifetime {
		e.Trash()
	}

	return request, nil
}

func (e *ProjecticleEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	var image *ebiten.Image
	switch e.polarity {
	case PositivePolarity:
		image = projecticlePositiveImage
	case NegativePolarity:
		image = projecticleNegativeImage
	case NeutralPolarity:
		image = projecticleNeutralImage
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	// Draw from center.
	op.GeoM.Translate(
		-float64(image.Bounds().Dx())/2,
		-float64(image.Bounds().Dy())/2,
	)
	screen.DrawImage(image, op)
}
