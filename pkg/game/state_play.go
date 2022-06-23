package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebijam22/pkg/data"
	"github.com/kettek/ebijam22/pkg/net"
	"github.com/kettek/ebijam22/pkg/world"
)

type PlayState struct {
	game     *Game
	level    data.Level
	world    world.World
	messages []Message
}

func (s *PlayState) Init() error {
	s.world.Game = s.game // Eww
	if err := s.world.BuildFromLevel(s.level); err != nil {
		return err
	}
	return nil
}

func (s *PlayState) Dispose() error {
	// Remove player entity reference.
	for _, p := range s.game.players {
		p.Entity = nil
	}
	return nil
}

func (s *PlayState) Update() error {
	// If we're the host/solo and we hit R, restart the level. If we're the client, send a request.
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		if s.game.Net.Hosting() || !s.game.Net.Active() {
			s.game.SetState(&TravelState{
				game:        s.game,
				targetLevel: "001", // ???
			})
		} else {
			s.game.Net.Send(net.TravelMessage{})
		}
		return nil
	}

	// Handle our network updates.
	for _, msg := range s.game.Net.Messages() {
		switch msg := msg.(type) {
		case net.TravelMessage:
			if !s.game.Net.Hosting() {
				s.game.SetState(&TravelState{
					game:        s.game,
					targetLevel: msg.Destination,
					restarting:  true,
				})
			} else {
				s.AddMessage(Message{
					content: fmt.Sprintf("%s wants to restart! Hit 'r' to conform.", s.game.Net.OtherName),
				})
			}
		default:
			// Send unhandled messages to the world.
			s.world.ProcessNetMessage(msg)
		}
	}

	// Process our local messages.
	t := s.messages[:0]
	for _, m := range s.messages {
		m.lifetime++
		if m.lifetime < m.deathtime {
			t = append(t, m)
		}
	}
	s.messages = t

	// Update our players.
	for _, p := range s.game.players {
		action, err := p.Update(&s.world)
		if err != nil {
			return err
		}
		if action != nil {
			// If our net is active, send our desired action to the other.
			if s.game.Net.Active() {
				s.game.Net.Send(action)
			}
			p.Entity.SetAction(action)
		}
	}

	// Update our world.
	if err := s.world.Update(); err != nil {
		return err
	}

	return nil
}

func (s *PlayState) Draw(screen *ebiten.Image) {
	// Draw our world.
	s.world.Draw(screen)

	// Draw level text centered at top of screen for now.
	bounds := text.BoundString(boldFace, s.level.Title)
	centeredX := world.ScreenWidth/2 - bounds.Min.X - bounds.Dx()/2
	text.Draw(screen, s.level.Title, boldFace, centeredX, bounds.Dy()+1, color.White)

	// Draw our messages from most recent to oldest, bottom to top.
	mx := 8
	my := world.ScreenHeight - 40
	for i := len(s.messages) - 1; i >= 0; i-- {
		m := s.messages[i]
		bounds := text.BoundString(normalFace, m.content)

		d := float64(m.lifetime) / float64(m.deathtime)
		c := color.RGBA{
			255,
			255,
			255,
			255 - uint8(255*d),
		}

		text.Draw(
			screen,
			m.content,
			normalFace,
			mx,
			my,
			c,
		)
		my -= bounds.Dy() + 2
	}

	// Draw our player's belt!
	s.game.players[0].Toolbelt.Draw(screen)
}

func (s *PlayState) AddMessage(m Message) {
	if m.deathtime <= 0 || m.deathtime >= 1000 {
		m.deathtime = 300
	}
	s.messages = append(s.messages, m)
}
