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

	// Add players here...?
	s.game.players = append(s.game.players, world.NewPlayer())
	s.game.players[0].Local = true
	// Add other player!
	if s.game.net.Active() {
		s.game.players = append(s.game.players, world.NewPlayer())
	}

	// Build the level.
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
	// Remove players. Should this be moved to a preplay state? Something between menu and travel for setting up players.
	s.game.players = make([]*world.Player, 0)
	return nil
}

func (s *PlayState) Update() error {
	// World mode handling.
	switch s.world.Mode {
	case world.LossMode:
		// TODO: Show hit "R" to restart or something. Also maybe stats.
		if s.game.net.Hosting() || !s.game.net.Active() {
			s.game.SetState(&TravelState{
				game:        s.game,
				targetLevel: "001", // FIXME: replace with current level!
			})
		}
	case world.VictoryMode:
		// TODO: Show end game stats, if possible! Then some sort of "hit okay" to travel button/key.
		if s.game.net.Hosting() || !s.game.net.Active() {
			s.game.SetState(&TravelState{
				game:        s.game,
				targetLevel: "001", // FIXME: replace with next level!
			})
		}
	}

	// If we're the host/solo and we hit R, restart the level. If we're the client, send a request.
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		if s.game.net.Hosting() || !s.game.net.Active() {
			s.game.SetState(&TravelState{
				game:        s.game,
				targetLevel: "001", // ???
			})
		} else {
			s.game.net.Send(net.TravelMessage{})
		}
		return nil
	}

	// Handle our network updates.
	for _, msg := range s.game.net.Messages() {
		switch msg := msg.(type) {
		case net.TravelMessage:
			if !s.game.net.Hosting() {
				s.game.SetState(&TravelState{
					game:        s.game,
					targetLevel: msg.Destination,
					restarting:  true,
				})
			} else {
				s.AddMessage(Message{
					content: fmt.Sprintf("%s wants to restart! Hit 'r' to conform.", s.game.net.OtherName),
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
			if s.game.net.Active() {
				s.game.net.Send(action)
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

	// Draw current game mode overlay...?
	if s.world.Mode == world.BuildMode {
		bounds := text.BoundString(normalFace, "build mode")
		text.Draw(
			screen,
			"build mode",
			normalFace,
			8,
			bounds.Dy()+8,
			color.White,
		)
		// Hmm.
		msg := "hit <spacebar> to start combat waves"
		bounds = text.BoundString(normalFace, msg)
		text.Draw(
			screen,
			msg,
			normalFace,
			world.ScreenWidth/2-bounds.Dx()/2,
			world.ScreenHeight-50,
			color.White,
		)
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
