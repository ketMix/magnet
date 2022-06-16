package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerEntity struct {
	BaseEntity
	player *Player
}

func NewPlayerEntity(player *Player) *PlayerEntity {
	return &PlayerEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{},
		},
		player: player,
	}
}

func (e *PlayerEntity) Update() error {
	// React to pending commands?
	return nil
}

func (e *PlayerEntity) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)
	// Draw from center.
	op.GeoM.Translate(
		-float64(playerImage.Bounds().Dx())/2,
		-float64(playerImage.Bounds().Dy())/2,
	)
	screen.DrawImage(playerImage, op)
}
