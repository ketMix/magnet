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
	game          *Game
	levelDataName string
	level         data.Level
	world         world.World
	messages      []Message
	clickables    []ClickableUI
}

func (s *PlayState) Init() error {
	s.world.Game = s.game // Eww
	s.world.Speed = s.game.Options.Speed

	// Add players here...?
	s.game.players = append(s.game.players, world.NewPlayer())
	s.game.players[0].Local = true
	// Add other player!
	if s.game.net.Active() {
		s.game.players = append(s.game.players, world.NewPlayer())
		// Set player names if networked.
		s.game.players[0].Name = s.game.net.Name
		s.game.players[1].Name = s.game.net.OtherName
	}

	// Build the level.
	if err := s.world.BuildFromLevel(s.level); err != nil {
		return err
	}

	s.world.Mode = &world.PreGameMode{}
	s.clickables = []ClickableUI{
		NewBGMIcon(),
		NewSFXIcon(),
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
	// Update the clickables
	for _, c := range s.clickables {
		c.Update()
	}

	// World mode handling.
	switch s.world.Mode.(type) {
	case *world.LossMode:
		// TODO: Show hit "R" to restart or something. Also maybe stats.
		if s.game.net.Hosting() || !s.game.net.Active() {
			s.game.SetState(&TravelState{
				game:        s.game,
				targetLevel: s.levelDataName,
			})
		}
	case *world.VictoryMode:
		// TODO: Show end game stats, if possible! Then some sort of "hit okay" to travel button/key.
		if s.level.Next != "" {
			if s.game.net.Hosting() || !s.game.net.Active() {
				fmt.Println("TRAVELING TO NEXT")
				s.game.SetState(&TravelState{
					game:        s.game,
					targetLevel: s.level.Next,
				})
			}
		} else {
			// We beat the video game.
		}
	}

	// If we're the host/solo and we hit R, restart the level. If we're the client, send a request.
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		if s.game.net.Hosting() || !s.game.net.Active() {
			s.game.SetState(&TravelState{
				game:        s.game,
				targetLevel: s.levelDataName, // ???
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
					fromLive:    true,
				})
			} else {
				s.AddMessage(Message{
					content: fmt.Sprintf("%s wants to restart! Hit 'r' to conform.", s.game.net.OtherName),
				})
			}
		case world.StartModeRequest:
			s.game.players[1].ReadyForWave = true
			if !s.world.ArePlayersReady() {
				s.AddMessage(Message{
					content: fmt.Sprintf("%s wants to start! Hit '<spacebar>' to conform.", s.game.net.OtherName),
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
	bounds := text.BoundString(data.BoldFace, s.level.Title)
	centeredX := world.ScreenWidth/2 - bounds.Min.X - bounds.Dx()/2
	text.Draw(screen, s.level.Title, data.BoldFace, centeredX, bounds.Dy()+1, color.White)

	// Draw our messages from most recent to oldest, bottom to top.
	mx := 8
	my := world.ScreenHeight - 40
	for i := len(s.messages) - 1; i >= 0; i-- {
		m := s.messages[i]
		bounds := text.BoundString(data.NormalFace, m.content)

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
			data.NormalFace,
			mx,
			my,
			c,
		)
		my -= bounds.Dy() + 2
	}

	// Draw mode.
	s.world.Mode.Draw(screen)

	// Draw the waves and current points.
	mx = 8
	my = 8
	{
		t := fmt.Sprintf("wave: %d/%d points: %d", s.world.CurrentWave, s.world.MaxWave, s.world.Points)
		bounds := text.BoundString(data.NormalFace, t)
		text.Draw(
			screen,
			t,
			data.NormalFace,
			mx,
			my+bounds.Dy(),
			color.White,
		)
		mx += bounds.Dx()
	}

	// Draw our clickables
	offset := 16
	if s.clickables != nil {
		for i, c := range s.clickables {
			c.SetPos(float64(mx+(i+1)*offset), float64(12))
			c.Draw(screen, &ebiten.DrawImageOptions{})
		}
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
