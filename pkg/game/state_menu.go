package game

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/data/assets/lang"
	"github.com/kettek/ebijam22/pkg/world"
)

type MenuState struct {
	game             *Game
	ebitenImage      *ebiten.Image
	magnetImage      *ebiten.Image
	backgroundImage  *ebiten.Image
	title            string
	titleImage       *ebiten.Image
	magnetSpin       float64
	buttons          []*data.Button
	backgroundFadeIn int
	shouldQuit       bool
}

func (s *MenuState) Init() error {
	// Title Text
	s.title = "" // forever hidden

	// Load some assets. This will be abstracted elsewhere.
	if img, err := data.ReadImage("/ui/magnet.png"); err == nil {
		ebiten.SetWindowIcon([]image.Image{img})
		s.magnetImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Load our title image.
	if img, err := data.ReadImage("/ui/logo.png"); err == nil {
		s.titleImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Load our background image.
	if img, err := data.ReadImage("/ui/title.png"); err == nil {
		s.backgroundImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Load our ebiten image.
	if img, err := data.ReadImage("/ui/ebitengine-gamejam.png"); err == nil {
		s.ebitenImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Load our flag images.
	var jaFlagImage *ebiten.Image
	var usFlagImage *ebiten.Image

	if img, err := data.ReadImage("/ui/ja.png"); err == nil {
		jaFlagImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}
	if img, err := data.ReadImage("/ui/us.png"); err == nil {
		usFlagImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Set Main Menu Buttons
	x := world.ScreenWidth / 2
	y := int(float64(world.ScreenHeight) / 2)
	startGameButton := data.NewButton(
		x,
		y,
		lang.SoloGame,
		func() {
			s.game.SetState(&SoloMenuState{
				game: s.game,
			})
		},
	)
	startGameButton.Hover = true
	y += startGameButton.Image().Bounds().Dy() * 3
	networkButton := data.NewButton(
		x,
		y,
		lang.NetworkGame,
		func() {
			s.game.SetState(&NetworkMenuState{
				game: s.game,
			})
		},
	)
	networkButton.Hover = true
	y += networkButton.Image().Bounds().Dy() * 3
	exitButton := data.NewButton(
		x,
		y,
		lang.Exit,
		func() {
			s.shouldQuit = true
		},
	)
	exitButton.Hover = true

	// Use buttons for credits.
	x = world.ScreenWidth - world.ScreenWidth/6
	y = world.ScreenHeight / 8
	credits1aButton := data.NewButton(
		x,
		y,
		lang.KettekContrib,
		func() {
			// open github/website
		},
	)
	y += int(float64(credits1aButton.Image().Bounds().Dy()) * 1.5)
	credits1bButton := data.NewButton(
		x,
		y,
		"kettek",
		func() {
			OpenFile("https://kettek.net")
		},
	)
	credits1bButton.Bold = true
	credits1bButton.Underline = true
	credits1bButton.Hover = true

	y += credits1bButton.Image().Bounds().Dy() * 3
	credits2aButton := data.NewButton(
		x,
		y,
		lang.LiqMixContrib,
		func() {
			s.game.SetState(&MusicMenuState{
				game: s.game,
			})
		},
	)
	y += int(float64(credits2aButton.Image().Bounds().Dy()) * 1.5)
	credits2bButton := data.NewButton(
		x,
		y,
		"liqMix",
		func() {
			OpenFile("https://liq.mx")
		},
	)
	credits2bButton.Bold = true
	credits2bButton.Underline = true
	credits2bButton.Hover = true

	y += credits2bButton.Image().Bounds().Dy() * 3
	credits3aButton := data.NewButton(
		x,
		y,
		lang.AmaruukContrib,
		func() {
		},
	)
	y += int(float64(credits3aButton.Image().Bounds().Dy()) * 1.5)
	credits3bButton := data.NewButton(
		x,
		y,
		"Amaruuk",
		func() {
			OpenFile("https://birdtooth.studio")
		},
	)
	credits3bButton.Bold = true
	credits3bButton.Underline = true
	credits3bButton.Hover = true

	// Create our flag lang clickables, set them in bottom left corner
	x = int(float64(world.ScreenWidth) * 0.05)
	y = int(float64(world.ScreenHeight) * 0.95)
	jaLangButton := data.NewImageButton(
		x, y,
		jaFlagImage,
		func() {
			data.ChangeLang(data.Japanese)
		},
	)
	jaLangButton.Hover = true

	x += jaFlagImage.Bounds().Dx() + 5
	usLangButton := data.NewImageButton(
		x, y,
		usFlagImage,
		func() {
			data.ChangeLang(data.English)
		},
	)
	usLangButton.Hover = true

	s.buttons = []*data.Button{
		startGameButton,
		networkButton,
		exitButton,
		credits1aButton,
		credits1bButton,
		credits2aButton,
		credits2bButton,
		credits3aButton,
		credits3bButton,
		jaLangButton,
		usLangButton,
	}

	// Start the tunes
	data.BGM.Set("menu.ogg")
	return nil
}

func (s *MenuState) Dispose() error {
	return nil
}

func (s *MenuState) Update() error {
	// It seems the idiomatic way to quit an ebiten program is to return an error...?
	if s.shouldQuit {
		return NoError{}
	}
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

	s.backgroundFadeIn++
	return nil
}

func (s *MenuState) Draw(screen *ebiten.Image) {
	// Let's first draw that background.
	titleOp := &ebiten.DrawImageOptions{}
	// Darken it a lil.
	d := math.Max(0.5, 1.0-float64(s.backgroundFadeIn)/120.0)
	titleOp.ColorM.Scale(d, d, d, 1.0)
	screen.DrawImage(s.backgroundImage, titleOp)

	// Rotate our magnet about its center.
	op := ebiten.DrawImageOptions{}
	/*op.GeoM.Translate(-float64(s.magnetImage.Bounds().Dx())/2, -float64(s.magnetImage.Bounds().Dy())/2)
	op.GeoM.Rotate(s.magnetSpin)

	// Render it at the center of the screen.
	op.GeoM.Translate(float64(world.ScreenWidth)/2, float64(world.ScreenHeight/2))
	screen.DrawImage(s.magnetImage, &op)*/

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

	// Draw our real title.
	top := &ebiten.DrawImageOptions{}
	top.GeoM.Translate(float64(world.ScreenWidth)/2-float64(s.titleImage.Bounds().Dx())/2, 16)
	screen.DrawImage(s.titleImage, top)

	// Draw left-hand ebiten
	op = ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(world.ScreenWidth)/7, float64(world.ScreenHeight)/4)
	op.GeoM.Translate(-float64(s.ebitenImage.Bounds().Dx())/2, -float64(s.ebitenImage.Bounds().Dy())/2)
	screen.DrawImage(s.ebitenImage, &op)

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
