package game

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/world"
)

type MenuState struct {
	game        *Game
	magnetImage *ebiten.Image
	titleImage  *ebiten.Image
	title       string
	magnetSpin  float64
	buttons     []data.Button
	titleFadeIn int
}

func (s *MenuState) Init() error {
	// Title Text
	s.title = "ebijam 2022"

	// Load some assets. This will be abstracted elsewhere.
	if img, err := data.ReadImage("/ui/magnet.png"); err == nil {
		ebiten.SetWindowIcon([]image.Image{img})
		s.magnetImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Load our title image.
	if img, err := data.ReadImage("/ui/title.png"); err == nil {
		s.titleImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Set Main Menu Buttons
	x := world.ScreenWidth / 2
	y := int(float64(world.ScreenHeight) / 1.25)
	startGameButton := data.NewButton(
		x,
		y,
		"Solo Game",
		func() {
			s.game.SetState(&SoloMenuState{
				game: s.game,
			})
		},
	)
	y += startGameButton.Image().Bounds().Dy() * 2
	networkButton := data.NewButton(
		x,
		y,
		"Network Game",
		func() {
			s.game.SetState(&NetworkMenuState{
				game: s.game,
			})
		},
	)
	s.buttons = []data.Button{
		*startGameButton,
		*networkButton,
	}
	// Start the tunes
	data.BGM.Set("menu.ogg")
	return nil
}

func (s *MenuState) Dispose() error {
	return nil
}

func (s *MenuState) Update() error {
	// Are we skippin menu
	if s.game.Options.NoMenu {
		s.StartGame()
	}

	// Spin at 4 degrees per update.
	s.magnetSpin += math.Pi / 180 * 4

	// Update buttons
	for _, button := range s.buttons {
		button.Update()
	}

	s.titleFadeIn++

	return nil
}

func (s *MenuState) Draw(screen *ebiten.Image) {
	// Let's first draw that background.
	titleOp := &ebiten.DrawImageOptions{}
	// Darken it a lil.
	d := math.Max(0.5, 1.0-float64(s.titleFadeIn)/120.0)
	titleOp.ColorM.Scale(d, d, d, 1.0)
	screen.DrawImage(s.titleImage, titleOp)

	// Rotate our magnet about its center.
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(s.magnetImage.Bounds().Dx())/2, -float64(s.magnetImage.Bounds().Dy())/2)
	op.GeoM.Rotate(s.magnetSpin)

	// Render it at the center of the screen.
	op.GeoM.Translate(float64(world.ScreenWidth)/2, float64(world.ScreenHeight/2))
	screen.DrawImage(s.magnetImage, &op)

	// Draw our title
	data.DrawStaticText(
		s.title,
		data.BoldFace,
		world.ScreenWidth/3,
		world.ScreenHeight/2,
		color.White,
		screen,
		true,
	)

	// Draw game buttons
	for _, button := range s.buttons {
		button.Draw(screen, &ebiten.DrawImageOptions{})
	}
}

func (s *MenuState) StartGame() {
	s.game.SetState(&TravelState{
		game:        s.game,
		targetLevel: s.game.Options.Map,
	})
}
