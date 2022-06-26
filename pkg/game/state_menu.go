package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/world"
)

type MenuState struct {
	game        *Game
	magnetImage *ebiten.Image
	magnetSpin  float64
}

func (s *MenuState) Init() error {
	// Load some assets. This will be abstracted elsewhere.
	if img, err := data.ReadImage("/ui/magnet.png"); err == nil {
		ebiten.SetWindowIcon([]image.Image{img})
		s.magnetImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	return nil
}

func (s *MenuState) Dispose() error {
	return nil
}

func (s *MenuState) Update() error {
	// Spin at 4 degrees per update.
	s.magnetSpin += math.Pi / 180 * 4

	// Travel for now.
	s.game.SetState(&TravelState{
		game:        s.game,
		targetLevel: s.game.Options.Map,
	})
	return nil
}

func (s *MenuState) Draw(screen *ebiten.Image) {
	// Rotate our magnet about its center.
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(s.magnetImage.Bounds().Dx())/2, -float64(s.magnetImage.Bounds().Dy())/2)
	op.GeoM.Rotate(s.magnetSpin)

	// Render it at the center of the screen.
	op.GeoM.Translate(float64(world.ScreenWidth/2), float64(world.ScreenHeight/2))
	screen.DrawImage(s.magnetImage, &op)
}
