package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ProjecticleEntity struct {
	BaseEntity
	polarity Polarity
}

func NewProjecticleEntity() *ProjecticleEntity {
	return &ProjecticleEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{},
		},
	}
}

func (e *ProjecticleEntity) Update() (request Request, err error) {
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

	// // Do we need to destroy this object?
	// var offScreenX = (e.physics.X > float64(screenWidth) || e.physics.X < 0)
	// var offScreenY = (e.physics.Y > float64(screenWidth) || e.physics.Y < 0)
	// if offScreenX || offScreenY {
	// 	e.Trash()
	// }

	return request, nil
}

func (e *ProjecticleEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	var image *ebiten.Image
	switch e.polarity {
	case POSITIVE:
		image = projecticlePositiveImage
	case NEGATIVE:
		image = projecticleNegativeImage
	case NEUTRAL:
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
