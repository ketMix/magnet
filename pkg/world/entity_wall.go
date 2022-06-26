package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
)

type WallEntity struct {
	BaseEntity
	owner           string
	colorMultiplier [3]float64
}

func NewWallEntity() *WallEntity {
	wallImg, _ := data.GetImage("wall.png")
	return &WallEntity{
		BaseEntity: BaseEntity{
			physics: PhysicsObject{},
			animation: Animation{
				images: []*ebiten.Image{wallImg},
			},
		},
		colorMultiplier: [3]float64{1, 1, 1},
	}
}

func (e *WallEntity) Update(world *World) (request Request, err error) {
	return request, nil
}

func (e *WallEntity) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(screenOp.GeoM)
	op.GeoM.Translate(
		e.physics.X,
		e.physics.Y,
	)

	op.ColorM.Scale(e.colorMultiplier[0], e.colorMultiplier[1], e.colorMultiplier[2], 1)

	e.animation.Draw(screen, op)
}
