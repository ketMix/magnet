package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type CoreEntity struct {
	BaseEntity
}

func NewCoreEntity() *CoreEntity {
	return &CoreEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{},
		},
	}
}

func (e *CoreEntity) Update(world *World) (request Request, err error) {
	return request, nil
}

func (e *CoreEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	op.GeoM.Translate(
		-float64(spawnerImage.Bounds().Dx())/2,
		0,
	)

	// Draw from center.
	op.GeoM.Translate(
		-float64(spawnerImage.Bounds().Dx())/2,
		-float64(spawnerImage.Bounds().Dy())/2,
	)

	// TODO: Make an "active" mode that has an alternative image or an image underlay.
	screen.DrawImage(spawnerImage, op)
}
