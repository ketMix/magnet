package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// For now...
var (
	// Our internal screen width and height.
	screenWidth, screenHeight int
)

// Game is our ebiten engine interface compliant type.
type Game struct {
	// Our game current game state.
	state State
}

// Init is used to set up all initial game structures.
func (g *Game) Init() (err error) {
	// Default to 640x360 for now.
	screenWidth = 640
	screenHeight = 360

	// Use nearest-neighbor for scaling.
	ebiten.SetScreenFilterEnabled(false)

	// Size our screen.
	ebiten.SetWindowSize(1280, 720)

	// Set our initial menu state.
	if err := g.SetState(&MenuState{
		game: g,
	}); err != nil {
		panic(err)
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
