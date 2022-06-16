package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// For now...
var (
	// Our internal screen width and height.
	cellWidth, cellHeight     int
	screenWidth, screenHeight int
	normalFace, boldFace      font.Face
	playerImage               *ebiten.Image
	grassImage                *ebiten.Image
)

// Game is our ebiten engine interface compliant type.
type Game struct {
	// Our game current game state.
	state State
	//
	players []*Player
}

// Init is used to set up all initial game structures.
func (g *Game) Init() (err error) {
	// Default to 640x360 for now.
	screenWidth = 640
	screenHeight = 360

	// Set our cell width and height.
	cellWidth = 16
	cellHeight = 11

	// Use nearest-neighbor for scaling.
	ebiten.SetScreenFilterEnabled(false)

	// Size our screen.
	ebiten.SetWindowSize(1280, 720)

	// Load our global fonts.
	data, err := readFile("fonts/OpenSansPX.ttf")
	if err != nil {
		return err
	}
	tt, err := opentype.Parse(data)
	if err != nil {
		return err
	}
	if normalFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	}); err != nil {
		return err
	}
	data, err = readFile("fonts/OpenSansPXBold.ttf")
	if err != nil {
		return err
	}
	tt, err = opentype.Parse(data)
	if err != nil {
		return err
	}
	if boldFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	}); err != nil {
		return err
	}

	// Load some images.
	if img, err := readImage("player.png"); err == nil {
		playerImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}
	if img, err := readImage("grass.png"); err == nil {
		grassImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Set our initial menu state.
	if err := g.SetState(&MenuState{
		game: g,
	}); err != nil {
		return err
	}

	return
}

// Update updates, how about that.
func (g *Game) Update() error {
	return g.state.Update()
}

// Draw draws to the given ebiten screen buffer image.
func (g *Game) Draw(screen *ebiten.Image) {
	g.state.Draw(screen)
}

// Layout sets up "virtual" screen dimensions in contrast to the window's dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) SetState(s State) error {
	if err := s.Init(); err != nil {
		panic(err)
	}
	g.state = s
	return nil
}
