package game

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// For now...
var (
	// Our internal screen width and height.
	cellWidth, cellHeight     int
	screenWidth, screenHeight int
	normalFace, boldFace      font.Face

	// Images for drawing lines.
	emptyImage    *ebiten.Image
	emptySubImage *ebiten.Image

	// Images
	turretNegativeImage                *ebiten.Image
	turretPositiveImage                *ebiten.Image
	spawnerImage, spawnerShadowImage   *ebiten.Image
	toolSlotImage, toolSlotActiveImage *ebiten.Image
	toolDestroyImage                   *ebiten.Image
	toolGunImage                       *ebiten.Image
	projecticlePositiveImage           *ebiten.Image
	projecticleNegativeImage           *ebiten.Image
	projecticleNeutralImage            *ebiten.Image

	// SFX
	turretPlaceSound *Sound
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

	// Load configurations
	err = LoadConfigurations()
	if err != nil {
		return err
	}

	//
	emptyImage = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	emptyImage.Fill(color.White)

	// IMAGES //
	if img, err := readImage("turret-negative.png"); err == nil {
		turretNegativeImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("turret-positive.png"); err == nil {
		turretPositiveImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("spawner.png"); err == nil {
		spawnerImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("spawner-shadow.png"); err == nil {
		spawnerShadowImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Tools
	if img, err := readImage("toolslot.png"); err == nil {
		toolSlotImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("toolslot-active.png"); err == nil {
		toolSlotActiveImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("tool-destroy.png"); err == nil {
		toolDestroyImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("tool-gun.png"); err == nil {
		toolGunImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// Projecticles
	if img, err := readImage("projecticle-positive.png"); err == nil {
		projecticlePositiveImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("projecticle-negative.png"); err == nil {
		projecticleNegativeImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("projecticle-neutral.png"); err == nil {
		projecticleNeutralImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// SOUNDS //
	if snd, err := readSound("turret-place.ogg"); err == nil {
		turretPlaceSound = snd
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
