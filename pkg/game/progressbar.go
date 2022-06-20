package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ProgressBar struct {
	image    *ebiten.Image
	progress float64
}

func NewProgressBar(width, height int, barColor color.RGBA) *ProgressBar {
	image := ebiten.NewImage(width, height)
	image.Fill(barColor)
	return &ProgressBar{
		image: image,
	}
}
func (pb *ProgressBar) Update() {
}

func (pb *ProgressBar) Draw(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1/pb.progress, 1)
	op.GeoM.Concat(screenOp.GeoM)
	screen.DrawImage(pb.image, op)
}
