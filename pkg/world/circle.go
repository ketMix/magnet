package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// circles is our naughty cache of circle images.
var circles map[int]*ebiten.Image = make(map[int]*ebiten.Image)

// getCircleImage returns an image with the outline of a circle. It generates one and caches it for future use if it does not exist.
func getCircleImage(radius int) *ebiten.Image {
	if c, ok := circles[radius]; ok {
		return c
	}

	// Bresenham, my old friend.
	r := radius * 2
	img := ebiten.NewImage(int(r)+1, int(r)+1)
	drawEllipsePoints := func(cx, cy, x, y int) {
		img.Set(cx+x, cy+y, color.RGBA{127, 127, 127, 255})
		img.Set(cx-x, cy+y, color.RGBA{127, 127, 127, 255})
		img.Set(cx+x, cy-y, color.RGBA{127, 127, 127, 255})
		img.Set(cx-x, cy-y, color.RGBA{127, 127, 127, 255})
	}

	d := 5 - 4*radius
	x := 0
	y := radius

	deltaA := (-2*radius + 5) * 4
	deltaB := 3 * 4

	drawEllipsePoints(radius, radius, x, y)

	for x <= y {
		drawEllipsePoints(radius, radius, x, y)
		drawEllipsePoints(radius, radius, y, x)

		x++

		if d > 0 {
			d += deltaA
			y--
			deltaA += 4 * 4
			deltaB += 2 * 2
		} else {
			d += deltaB
			deltaA += 2 * 4
			deltaB += 2 * 4
		}
	}

	circles[radius] = img
	return img
}

func drawCircle(screen *ebiten.Image, screenOp *ebiten.DrawImageOptions, radius int, r, g, b, a float64) {
	c := getCircleImage(radius)
	cop := &ebiten.DrawImageOptions{}
	cop.ColorM.Scale(r, g, b, a)
	cop.GeoM.Concat(screenOp.GeoM)
	cop.GeoM.Translate(-float64(c.Bounds().Dx())/2, -float64(c.Bounds().Dy())/2)
	screen.DrawImage(c, cop)
}
