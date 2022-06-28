package ui

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawTiled draws the given image to fill the provided width and height.
func DrawTiled(screen *ebiten.Image, image *ebiten.Image, op *ebiten.DrawImageOptions, width, height int) {
	bgWidth := image.Bounds().Dx()
	bgHeight := image.Bounds().Dy()

	rows := math.Ceil(float64(width)/float64(bgWidth)) + 1
	cols := math.Ceil(float64(height)/float64(bgHeight)) + 1

	for y := 0.0; y < rows; y++ {
		for x := 0.0; x < cols; x++ {
			bgOp := &ebiten.DrawImageOptions{}
			bgOp.GeoM.Concat(op.GeoM)
			bgOp.GeoM.Translate(x*float64(bgWidth), y*float64(bgHeight))
			screen.DrawImage(image, bgOp)
		}
	}

}
