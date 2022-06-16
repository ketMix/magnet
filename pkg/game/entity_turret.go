package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type TurretEntity struct {
	BaseEntity
	// owner ActorEntity // ???
}

func NewTurretEntity() *TurretEntity {
	return &TurretEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{},
		},
	}
}

func (e *TurretEntity) Update() (request Request, err error) {
	// A mystery.
	return request, nil
}

func (e *TurretEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)
	// Draw from center.
	op.GeoM.Translate(
		-float64(turretBaseImage.Bounds().Dx())/2,
		-float64(turretBaseImage.Bounds().Dy())/2,
	)
	screen.DrawImage(turretBaseImage, op)
}
