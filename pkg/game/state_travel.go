package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/net"
	"github.com/kettek/ebijam22/pkg/world"
)

type TravelState struct {
	game        *Game
	done        bool
	targetLevel string
	loadedLevel data.Level
	ready       bool
	restarting  bool
	//
	magnetImage *ebiten.Image
	magnetSpin  float64
}

func (s *TravelState) Init() (err error) {
	// Load some assets. This will be abstracted elsewhere.
	if img, err := data.ReadImage("magnet.png"); err == nil {
		ebiten.SetWindowIcon([]image.Image{img})
		s.magnetImage = ebiten.NewImageFromImage(img)
	} else {
		panic(err)
	}

	// TODO: Probably should keep looping until we a receive a "OkayToTravel" message or something.
	// See if we need to handle networked level loading logic.
	if s.game.net.Active() {
		// If we're hosting, send the required travel to other client.
		if s.game.net.Hosting() {
			s.game.net.Send(net.TravelMessage{
				Destination: s.targetLevel,
			})
			if err := s.LoadLevel(); err != nil {
				return err
			}
		} else if s.restarting {
			if err := s.LoadLevel(); err != nil {
				return err
			}
		}
	} else {
		if err := s.LoadLevel(); err != nil {
			return err
		}
	}

	return nil
}

func (s *TravelState) Dispose() error {
	return nil
}

func (s *TravelState) LoadLevel() (err error) {
	s.loadedLevel, err = data.NewLevel(s.targetLevel)
	s.ready = true
	return err
}

func (s *TravelState) Update() error {
	// Always spin that magnet.
	s.magnetSpin += math.Pi / 180 * 4

	// If we're marked as ready, let's go.
	if s.ready {
		s.game.SetState(&PlayState{
			game:  s.game,
			level: s.loadedLevel,
		})
		return nil
	}

	// If we're connected and not hosting, wait for a travel message.
	if s.game.net.Connected() && !s.game.net.Hosting() {
		for _, msg := range s.game.net.Messages() {
			switch m := msg.(type) {
			case net.TravelMessage:
				s.targetLevel = m.Destination
				if err := s.LoadLevel(); err != nil {
					return err
				}
			default:
				break
			}
		}
	}

	return nil
}

func (s *TravelState) Draw(screen *ebiten.Image) {
	// Rotate our magnet about its center.
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(s.magnetImage.Bounds().Dx())/2, -float64(s.magnetImage.Bounds().Dy())/2)
	op.GeoM.Rotate(s.magnetSpin)

	// Render it at the center of the screen.
	op.GeoM.Translate(float64(world.ScreenWidth/2), float64(world.ScreenHeight/2))
	screen.DrawImage(s.magnetImage, &op)
}
