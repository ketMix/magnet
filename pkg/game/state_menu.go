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
	title       string
	magnetSpin  float64
	buttons     []data.Button
	inputs      []data.TextInput
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

	// Set Main Menu Buttons
	startGameButton := data.NewButton(
		world.ScreenWidth/2,
		int(float64(world.ScreenHeight)/1.25),
		"Start Game",
		func() {
			s.StartGame()
		},
	)
	s.buttons = []data.Button{*startGameButton}

	playerNameInput := data.NewTextInput(
		"Player Name",
		int(float64(world.ScreenWidth)/1.25),
		int(float64(world.ScreenHeight)/4),
	)
	s.inputs = []data.TextInput{*playerNameInput}

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
	// Update inputs
	for _, input := range s.inputs {
		input.Update()
	}
	return nil
}

func (s *MenuState) Draw(screen *ebiten.Image) {
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
		button.Draw(screen, &op)
	}
	// Update inputs
	for _, input := range s.inputs {
		input.Draw(screen, &op)
	}
}

func (s *MenuState) StartGame() {
	s.game.SetState(&TravelState{
		game:        s.game,
		targetLevel: s.game.Options.Map,
	})
}
