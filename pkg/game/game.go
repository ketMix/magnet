package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// For now...
var (
	// Our internal screen width and height.
	cellWidth, cellHeight              int
	screenWidth, screenHeight          int
	normalFace, boldFace               font.Face
	playerImage                        *ebiten.Image
	turretNegativeImage                *ebiten.Image
	turretPositiveImage                *ebiten.Image
	spawnerImage, spawnerShadowImage   *ebiten.Image
	toolSlotImage, toolSlotActiveImage *ebiten.Image
	toolDestroyImage                   *ebiten.Image
	toolGunImage                       *ebiten.Image
	projecticlePositiveImage           *ebiten.Image
	projecticleNegativeImage           *ebiten.Image
	projecticleNeutralImage            *ebiten.Image
	enemyPositive1Image                *ebiten.Image
	enemyPositive2Image                *ebiten.Image
	enemyNegative1Image                *ebiten.Image
	enemyNegative2Image                *ebiten.Image
	//
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

	// IMAGES //
	if img, err := readImage("player.png"); err == nil {
		playerImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

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

	// Enemies
	if img, err := readImage("enemy-positive-1.png"); err == nil {
		enemyPositive1Image = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("enemy-positive-2.png"); err == nil {
		enemyPositive2Image = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("enemy-negative-1.png"); err == nil {
		enemyNegative1Image = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	if img, err := readImage("enemy-negative-2.png"); err == nil {
		enemyNegative2Image = ebiten.NewImageFromImage(img)
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
	if err := s.Init(); err != nil {
		panic(err)
	}
	g.state = s
	return nil
}
