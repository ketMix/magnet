package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// For now...
var (
	magnetImage *ebiten.Image
	magnetSpin  float64
)

// Game is our ebiten engine interface compliant type.
type Game struct {
	Width, Height int
}

// Init is used to set up all initial game structures.
func (g *Game) Init() (err error) {
	// Default to 640x360 for now.
	g.Width = 640
	g.Height = 360

	// Use nearest-neighbor for scaling.
	ebiten.SetScreenFilterEnabled(false)

	// Size our screen.
	ebiten.SetWindowSize(1280, 720)

	// Load some assets. This will be abstracted elsewhere.
	img, _ := readImage("magnet.png")
	magnetImage = ebiten.NewImageFromImage(img)

	return
}

// Update updates, how about that.
func (g *Game) Update() error {
	// Spin at 4 degrees per update.
	magnetSpin += math.Pi / 180 * 4
	return nil
}

// Draw draws to the given ebiten screen buffer image.
func (g *Game) Draw(screen *ebiten.Image) {
	// Rotate our magnet about its center.
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(magnetImage.Bounds().Dx())/2, -float64(magnetImage.Bounds().Dy())/2)
	op.GeoM.Rotate(magnetSpin)

	// Render it at the center of the screen.
	op.GeoM.Translate(float64(g.Width/2), float64(g.Height/2))
	screen.DrawImage(magnetImage, &op)
}

// Layout sets up "virtual" screen dimensions in contrast to the window's dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}
