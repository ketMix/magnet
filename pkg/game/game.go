package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebijam22/pkg/data"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// For now...
var (
	// Our internal screen width and height.
	cellWidth, cellHeight     int
	screenWidth, screenHeight int
	normalFace, boldFace      font.Face
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

	// Setup audio context.
	audio.NewContext(44100)

	// Load our global fonts.
	d, err := data.ReadFile("fonts/OpenSansPX.ttf")
	if err != nil {
		return err
	}
	tt, err := opentype.Parse(d)
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
	d, err = data.ReadFile("fonts/OpenSansPXBold.ttf")
	if err != nil {
		return err
	}
	tt, err = opentype.Parse(d)
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

	// Load configurations
	err = data.LoadConfigurations()
	if err != nil {
		return err
	}

	// Load data
	err = data.LoadData()

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
	if inpututil.IsKeyJustReleased(ebiten.KeyF) || inpututil.IsKeyJustReleased(ebiten.KeyF11) || (inpututil.IsKeyJustReleased(ebiten.KeyEnter) && ebiten.IsKeyPressed(ebiten.KeyAlt)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
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
	if g.state != nil {
		if err := s.Dispose(); err != nil {
			panic(err)
		}
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	g.state = s
	return nil
}
