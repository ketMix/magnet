package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/data/ui"
	"github.com/kettek/ebijam22/pkg/world"
)

type MusicMenuState struct {
	game        *Game
	title       string
	magnetImage *ebiten.Image
	magnetSpin  float64

	tiledBackgroundImages  []*ebiten.Image
	tiledBackgroundElapsed int
	tiledBackgroundIndex   int
	backgroundImage        *ebiten.Image

	buttons    []*data.Button
	animations []*world.Animation
}

func (s *MusicMenuState) Init() error {
	// Oh boy.
	t, err := data.LoadTileSet("magnet")
	if err != nil {
		return err
	}
	s.tiledBackgroundImages = t.BackgroundImages

	// Load our background image.
	if img, err := data.ReadImage("/ui/multiplayer.png"); err == nil {
		s.backgroundImage = ebiten.NewImageFromImage(img)
	} else {
		return err
	}

	// Title Text
	s.title = "Music Player"

	// Magnet Image
	if img, err := data.ReadImage("/ui/magnet.png"); err == nil {
		s.magnetImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	centeredX := world.ScreenWidth / 2

	// Create Buttons
	backButton := data.NewButton(
		15,
		10,
		"Back",
		func() {
			s.game.SetState(&MenuState{
				game: s.game,
			})
		},
	)
	backButton.Hover = true
	s.buttons = []*data.Button{
		backButton,
	}

	// Create music buttons
	offsetY := world.ScreenHeight / 2
	offset := 20
	tracks := data.BGM.GetAllTracks()
	for i := range tracks {
		trackName := tracks[i]
		trackButton := data.NewButton(
			centeredX,
			offsetY,
			data.FormatTrackName(trackName),
			func() {
				data.BGM.Set(trackName + ".ogg")
			},
		)
		trackButton.Hover = true
		s.buttons = append(s.buttons, trackButton)
		offsetY += offset
	}

	// Create the entities
	for i := range data.EnemyConfigs {
		config := data.EnemyConfigs[i]
		animation := world.NewAnimation(
			config.VictoryImages,
			30,
			1-config.Speed,
		)
		s.animations = append(s.animations, animation)
	}
	return nil
}

func (s *MusicMenuState) Dispose() error {
	return nil
}

func (s *MusicMenuState) Update() error {
	// Animate the background.
	s.tiledBackgroundElapsed++
	if s.tiledBackgroundElapsed >= 30 {
		s.tiledBackgroundElapsed = 0
		s.tiledBackgroundIndex++
		if s.tiledBackgroundIndex >= len(s.tiledBackgroundImages) {
			s.tiledBackgroundIndex = 0
		}
	}

	// Spin at 4 degrees per update.
	s.magnetSpin += math.Pi / 180 * 4
	if s.magnetSpin >= 2*math.Pi {
		s.magnetSpin = 0
	}
	// Update buttons
	for _, button := range s.buttons {
		button.Update()
	}

	// Update animations
	for _, animation := range s.animations {
		animation.Update()
	}
	return nil
}

func (s *MusicMenuState) Draw(screen *ebiten.Image) {
	// Draw our tiled background.
	bgOp := ebiten.DrawImageOptions{}
	ui.DrawTiled(screen, s.tiledBackgroundImages[s.tiledBackgroundIndex], &bgOp, world.ScreenWidth, world.ScreenHeight)

	// Draw our background.
	screenOp := &ebiten.DrawImageOptions{}
	screenOp.ColorM.Scale(0.5, 0.5, 0.5, 1)
	screen.DrawImage(s.backgroundImage, screenOp)

	centeredX := world.ScreenWidth / 2

	// Draw our title
	titleBounds := text.BoundString(data.BoldFace, s.title)
	data.DrawStaticText(
		s.title,
		data.BoldFace,
		centeredX,
		world.ScreenHeight/8,
		color.White,
		screen,
		true,
	)

	// Draw currently playing
	currentlyPlaying := "Currently Playing"
	currentTrack := data.FormatTrackName(data.BGM.GetCurrentTrack())
	textBounds := text.BoundString(data.NormalFace, currentlyPlaying)
	data.DrawStaticText(
		currentlyPlaying,
		data.NormalFace,
		centeredX,
		world.ScreenHeight/5,
		color.White,
		screen,
		true,
	)
	data.DrawStaticText(
		currentTrack,
		data.BoldFace,
		centeredX,
		world.ScreenHeight/5+textBounds.Dy()*2,
		color.White,
		screen,
		true,
	)

	// Rotate our magnet about its center.
	magnetOp := ebiten.DrawImageOptions{}
	magnetOp.GeoM.Translate(-float64(s.magnetImage.Bounds().Dx())/2, -float64(s.magnetImage.Bounds().Dy())/2)

	rightOp := ebiten.DrawImageOptions{}
	rightOp.GeoM.Concat(magnetOp.GeoM)
	rightOp.GeoM.Rotate(s.magnetSpin)
	rightOp.GeoM.Translate(float64(world.ScreenWidth/2)+float64(titleBounds.Dx())*0.7, float64(world.ScreenHeight/8))

	leftOp := ebiten.DrawImageOptions{}
	leftOp.GeoM.Concat(magnetOp.GeoM)
	leftOp.GeoM.Rotate(-s.magnetSpin)
	leftOp.GeoM.Translate(float64(world.ScreenWidth/2)-float64(titleBounds.Dx())*0.7, float64(world.ScreenHeight/8))

	// Render magnets on each side of title
	screen.DrawImage(s.magnetImage, &leftOp)
	screen.DrawImage(s.magnetImage, &rightOp)

	op := ebiten.DrawImageOptions{}

	// Draw game buttons
	for _, button := range s.buttons {
		if button.Text() == currentTrack {
			button.Underline = true
		} else {
			button.Underline = false
		}
		button.Draw(screen, &op)
	}
	animationOp := ebiten.DrawImageOptions{}
	animationOp.GeoM.Translate(float64(world.ScreenWidth)/2.5, float64(world.ScreenHeight)/2.5)
	offsetX := float64(world.ScreenWidth / 15 * (len(s.animations) / 4))
	// Draw animations
	for i := range s.animations {
		s.animations[i].Draw(screen, &animationOp)
		animationOp.GeoM.Translate(offsetX, 0)
	}
}
