package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/net"
	"github.com/kettek/ebijam22/pkg/world"
)

// For now...
var (
// Our internal screen width and height.
)

// Game is our ebiten engine interface compliant type.
type Game struct {
	Options data.Options
	// Our game current game state.
	state State
	//
	net net.Connection
	//
	players []*world.Player
}

// Init is used to set up all initial game structures.
func (g *Game) Init() (err error) {
	// Set our cell width and height.
	data.CellWidth = 16
	data.CellHeight = 11

	// Use nearest-neighbor for scaling.
	ebiten.SetScreenFilterEnabled(false)

	// Size our screen.
	ebiten.SetWindowSize(1280, 720)

	// Setup audio context.
	audio.NewContext(48000)

	// Load configurations
	err = data.LoadConfigurations()
	if err != nil {
		return err
	}

	// Load data
	err = data.LoadData()

	// Mute sound/music if flag exists.
	if g.Options.NoMusic {
		data.BGM.Muted = true
	}
	if g.Options.NoSound {
		data.SFX.Muted = true
	}

	// FIXME: Don't manually network connect here. This should be handled in some intermediate state, like "preplay" or a lobby.
	if g.Options.Host != "" || g.Options.Join != "" || g.Options.Await || g.Options.Search != "" {
		g.net = net.NewConnection(g.Options.Name)
		if g.Options.Host != "" {
			if err := g.net.AwaitDirect(g.Options.Host, ""); err != nil {
				panic(err)
			}
		} else if g.Options.Join != "" {
			if err := g.net.AwaitDirect("", g.Options.Join); err != nil {
				panic(err)
			}
		} else if g.Options.Await {
			if err := g.net.AwaitHandshake(g.Options.Handshaker, "", ""); err != nil {
				panic(err)
			}
		} else if g.Options.Search != "" {
			if err := g.net.AwaitHandshake(g.Options.Handshaker, "", g.Options.Search); err != nil {
				panic(err)
			}
		}
		go g.net.Loop()
	}

	// Set our initial menu state.
	if err := g.SetState(&MenuState{
		game: g,
	}); err != nil {
		return err
	}

	return
}

func (g *Game) Players() []*world.Player {
	return g.players
}

func (g *Game) Net() *net.Connection {
	return &g.net
}

// Update updates, how about that.
func (g *Game) Update() error {
	// Call update on our BGM to ensure it's playing
	data.BGM.Update()

	if inpututil.IsKeyJustReleased(ebiten.KeyF) || inpututil.IsKeyJustReleased(ebiten.KeyF11) || (inpututil.IsKeyJustReleased(ebiten.KeyEnter) && ebiten.IsKeyPressed(ebiten.KeyAlt)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return g.state.Update()
}

// Draw draws to the given ebiten screen buffer image.
func (g *Game) Draw(screen *ebiten.Image) {
	g.state.Draw(screen)

	// This should be handled differently than directly drawing network state here.
	if g.net.Active() {
		var img *ebiten.Image
		if g.net.Connected() {
			img, _ = data.GetImage("online.png")
		} else {
			img, _ = data.GetImage("offline.png")
		}
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(
				float64(world.ScreenWidth)-float64(img.Bounds().Dx())-8,
				float64(world.ScreenHeight-img.Bounds().Dy()-8),
			)
			screen.DrawImage(img, op)
		}
		// This is _bad_
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			float64(world.ScreenWidth)-float64(img.Bounds().Dx())-8,
			float64(img.Bounds().Dy()-8),
		)
		for i := range g.players {
			name := g.net.Name
			if i == 1 {
				name = g.net.OtherName
			}
			imgs := data.Player2Init.Images
			if g.net.Hosting() {
				if i == 0 {
					imgs = data.PlayerInit.Images
				}
			} else {
				if i == 1 {
					imgs = data.PlayerInit.Images
				}
			}
			screen.DrawImage(imgs[0], op)

			bounds := text.BoundString(data.NormalFace, name)
			text.Draw(screen, name, data.NormalFace, int(op.GeoM.Element(0, 2))-bounds.Dx()-4, int(op.GeoM.Element(1, 2))+8, color.White)

			op.GeoM.Translate(0, float64(imgs[0].Bounds().Dy())+8)
		}
	} else {
	}
}

// Layout sets up "virtual" screen dimensions in contrast to the window's dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return world.ScreenWidth, world.ScreenHeight
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
